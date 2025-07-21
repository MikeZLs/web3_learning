package database

// 该包将所有与数据库交互的初始化逻辑封装在一起，使代码更整洁。

import (
	"personal-blog/config"
	"personal-blog/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDatabase 初始化并返回一个GORM数据库连接实例(*gorm.DB)。
func InitDatabase() (*gorm.DB, error) {
	// 从配置模块获取数据库连接字符串
	dsn := config.GetDatabaseURL()
	// 使用GORM和MySQL驱动打开一个数据库连接池
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// AutoMigrate 使用GORM的自动迁移功能来同步数据库表结构。
// 它会根据在models包中定义的结构体，自动创建或更新数据库中的表。
// 这是一个非常方便的功能，尤其是在开发初期，模型结构经常变动时。
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
}
