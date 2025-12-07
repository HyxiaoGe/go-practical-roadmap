package main

import (
	"log"

	"github.com/yourname/02-concurrency-worker/internal/app"
)

func main() {
	// 创建应用程序实例
	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// 在goroutine中启动应用
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Failed to run app: %v", err)
		}
	}()

	// 等待中断信号
	application.WaitForInterrupt()
}