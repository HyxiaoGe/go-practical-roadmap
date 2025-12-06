package main

import (
	"fmt"
	"log"

	"go-practical-roadmap/01-web-api-template/internal/app"
)

func main() {
	fmt.Println("Go Web API Template Project")
	fmt.Println("This is a template project for building web APIs with Go, Gin, and Gorm.")

	// 创建应用实例
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// 在goroutine中启动应用
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Failed to run app: %v", err)
		}
	}()

	fmt.Println("Server is running on :8080")
	fmt.Println("Press Ctrl+C to stop")

	// 等待中断信号
	application.WaitForInterrupt()

	// 停止应用
	if err := application.Stop(); err != nil {
		log.Fatalf("Failed to stop app: %v", err)
	}

	fmt.Println("Server stopped")
}