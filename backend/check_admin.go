package main

import (
	"fmt"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/service"
)

func main() {
	// 初始化数据库
	database.GetDB()

	authService := service.NewAuthService()

	// 检查admin用户是否存在
	user, err := authService.GetUserByUsername("admin")
	if err != nil {
		fmt.Println("Admin user not found")
		fmt.Println("Creating admin user...")

		newUser, err := authService.CreateUser("admin", "admin", "admin@nas.local", "Administrator", "admin")
		if err != nil {
			fmt.Printf("Error creating admin user: %v\n", err)
		} else {
			fmt.Printf("Admin user created successfully: %s (Password: admin)\n", newUser.Username)
		}
	} else {
		fmt.Printf("Admin user found: %s (Role: %s, Email: %s)\n", user.Username, user.Role, user.Email)
		fmt.Println("You can login with username: admin, password: admin")
	}
}
