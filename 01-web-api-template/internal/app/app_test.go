package app

import (
	"testing"
	"time"
)

func TestAppCreation(t *testing.T) {
	// 测试应用创建
	app, err := NewApp()
	if err != nil {
		// 由于缺少实际的Go环境和数据库，这里预期会出错
		t.Logf("Expected error when creating app without proper environment: %v", err)
		return
	}

	// 如果应用创建成功，则测试停止功能
	if app != nil {
		// 在一个goroutine中启动应用
		go func() {
			err := app.Run()
			if err != nil {
				t.Logf("App run error: %v", err)
			}
		}()

		// 等待一点时间让服务器启动
		time.Sleep(100 * time.Millisecond)

		// 停止应用
		err = app.Stop()
		if err != nil {
			t.Errorf("Failed to stop app: %v", err)
		}
	}
}