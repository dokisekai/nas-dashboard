package main

import (
	"log"
	"os"

	"nas-dashboard/internal/api"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"
	"nas-dashboard/internal/routes"
	"nas-dashboard/internal/sso"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	db, err := initDatabase()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Server will continue without database support")
	}

	r := gin.Default()

	// 初始化 API 处理器与 SSO Provider
	api.InitAPI(db)
	ssoServer := sso.NewSSOServer(db)

	// 全局中间件 + 模板
	r.Use(middleware.CORS())
	r.LoadHTMLGlob("templates/*.html")

	// 注册全部 HTTP 路由（按域分散到 internal/routes/*.go）
	routes.Setup(r, db, ssoServer)

	// 启动服务器：若存在证书则启用 HTTPS，否则降级为 HTTP
	const port = "8888"
	certFile := "certs/server.crt"
	keyFile := "certs/server.key"

	if _, err := os.Stat(certFile); err == nil {
		log.Printf("Server starting on https://0.0.0.0:%s (SSL Enabled)", port)
		if err := r.RunTLS("0.0.0.0:"+port, certFile, keyFile); err != nil {
			log.Fatal("Failed to start HTTPS server:", err)
		}
		return
	}

	log.Printf("Server starting on http://0.0.0.0:%s", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}

// initDatabase 初始化数据库（非阻塞，失败时继续运行）。
func initDatabase() (*gorm.DB, error) {
	log.Println("Initializing database...")

	cfg := database.LoadConfig()
	if err := database.Connect(cfg); err != nil {
		return nil, err
	}

	db := database.GetDB()

	if err := database.Migrate(); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
