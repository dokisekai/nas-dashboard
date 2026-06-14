package main

import (
	"fmt"
	"log"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"nas-dashboard/internal/service"
)

func main() {
	// 初始化数据库
	cfg := database.LoadConfig()
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 运行迁移
	if err := database.Migrate(); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// 获取数据库连接
	db := database.GetDB()

	// 检查是否已有用户
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Fatalf("Failed to check users: %v", err)
	}

	if userCount > 0 {
		fmt.Printf("数据库中已有 %d 个用户\n", userCount)

		// 显示现有用户信息
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			log.Fatalf("Failed to get users: %v", err)
		}

		fmt.Println("现有用户列表:")
		for _, user := range users {
			fmt.Printf("  - ID: %d, 用户名: %s, 邮箱: %s, 角色: %s, 状态: %v\n",
				user.ID, user.Username, user.Email, user.Role, user.IsActive)
		}

		// 创建admin用户
		fmt.Println("\n是否要创建admin用户? (y/n)")
		var answer string
		fmt.Scanln(&answer)

		if answer == "y" || answer == "Y" {
			// 创建认证服务
			authService := service.NewAuthService()

			// 创建默认管理员用户
			username := "admin"
			password := "admin123"
			email := "admin@nas.local"

			// 检查admin用户是否已存在
			var existingUser models.User
			if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
				fmt.Printf("❌ admin用户已存在 (ID: %d)\n", existingUser.ID)

				// 询问是否重置密码
				fmt.Println("是否要重置admin用户密码? (y/n)")
				var resetAnswer string
				fmt.Scanln(&resetAnswer)

				if resetAnswer == "y" || resetAnswer == "Y" {
					// 重置密码为 "admin123"
					newPassword := "admin123"
					hashedPassword, err := authService.HashPassword(newPassword)
					if err != nil {
						log.Fatalf("Failed to hash password: %v", err)
					}

					if err := db.Model(&existingUser).Update("password_hash", hashedPassword).Error; err != nil {
						log.Fatalf("Failed to update password: %v", err)
					}

					fmt.Printf("✅ admin用户密码已重置为: %s\n", newPassword)
					fmt.Printf("   登录地址: http://192.168.50.10:5174/login\n")
				}
				return
			}

			// 创建admin用户
			user, err := authService.CreateUser(username, password, email, "系统管理员", "admin")
			if err != nil {
				log.Fatalf("Failed to create admin user: %v", err)
			}

			fmt.Printf("✅ 成功创建管理员用户:\n")
			fmt.Printf("   用户名: %s\n", user.Username)
			fmt.Printf("   邮箱: %s\n", user.Email)
			fmt.Printf("   角色: %s\n", user.Role)
			fmt.Printf("\n请使用以下凭据登录:\n")
			fmt.Printf("   用户名: %s\n", username)
			fmt.Printf("   密码: %s\n", password)
			fmt.Printf("   登录地址: http://192.168.50.10:5174/login\n")
		}

		return
	}

	// 创建认证服务
	authService := service.NewAuthService()

	// 创建默认管理员用户
	username := "admin"
	password := "admin123"
	email := "admin@nas.local"

	// 创建管理员用户
	user, err := authService.CreateUser(username, password, email, "系统管理员", "admin")
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("✅ 成功创建管理员用户:\n")
	fmt.Printf("   用户名: %s\n", user.Username)
	fmt.Printf("   邮箱: %s\n", user.Email)
	fmt.Printf("   角色: %s\n", user.Role)
	fmt.Printf("\n请使用以下凭据登录:\n")
	fmt.Printf("   用户名: %s\n", username)
	fmt.Printf("   密码: %s\n", password)
	fmt.Printf("   登录地址: http://192.168.50.10:5174/login\n")
}