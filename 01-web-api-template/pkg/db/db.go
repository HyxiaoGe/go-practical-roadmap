package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB 数据库连接实例
var DB *gorm.DB

// Connect 连接数据库
func Connect(dsn string, driver string) error {
	var err error

	switch driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		// 默认使用sqlite
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return err
	}

	// 获取通用数据库对象
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}