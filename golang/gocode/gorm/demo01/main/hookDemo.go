package main

import (
	"gorm.io/gorm"
)

// 在 GORM 中，Hooks（钩子函数） 是一种机制，允许你在 数据库操作之前或之后自动执行一些逻辑代码。
// 这就像是在生命周期的关键节点“挂上钩子”，GORM 会自动帮你触发这些钩子函数。

func (student *Student) BeforeCreate(db *gorm.DB) (err error) {
	email := "test@gmail.com"
	student.Email = &email
	return err
}

func (student *Student) AfterCreate(db *gorm.DB) (err error) { return err }

func (student *Student) BeforeDelete(db *gorm.DB) (err error) { return err }

func (student *Student) AfterDelete(db *gorm.DB) (err error) { return err }

func (student *Student) BeforeUpdate(db *gorm.DB) (err error) { return err } // 仅在调用 Update()/Updates() 时触发

func (student *Student) AfterUpdate(db *gorm.DB) (err error) { return err } // 仅在调用 Update()/Updates() 时触发

func (student *Student) BeforeSave(db *gorm.DB) (err error) { return err } // Create 和 Update 都会触发

func (student *Student) AfterSave(db *gorm.DB) (err error) { return err } // Create 和 Update 都会触发

func (student *Student) AfterFind(db *gorm.DB) (err error) { return err }
