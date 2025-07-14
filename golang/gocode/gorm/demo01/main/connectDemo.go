package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var mysqlLogger logger.Interface

func init() {

	mysqlLogger = logger.Default.LogMode(logger.Info)

	dsn := "root:root@tcp(127.0.0.1:3306)/test01?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//SkipDefaultTransaction: true,// 跳过默认事务
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   "f_", // 表名加前缀
		//	SingularTable: true, // 是否单数表名(后缀加不加字母 s)
		//	NoLowerCase:   true, // 不要小写转换
		//},
		//Logger: mysqlLogger,
	})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	//连接成功
	DB = db
}
