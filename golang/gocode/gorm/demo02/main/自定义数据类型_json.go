package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Info struct {
	Status string `json:"status"`
	Addr   string `json:"addr"`
	Age    int    `json:"age"`
}

// Scan 从数据库中读取出来
func (i *Info) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, i)
	return err
}

// Value 存入数据库
func (i Info) Value() (driver.Value, error) {
	fmt.Printf("入库前  %#v,%T\n", i, i)
	return json.Marshal(i)
}

type AuthModel struct {
	ID   uint
	Name string
	Info Info `gorm:"type:string"`
}

func main() {
	//DB.AutoMigrate(&AuthModel{})

	//// 写数据
	//DB.Create(&AuthModel{
	//	Name: "张三",
	//	Info: Info{
	//		Status: "success",
	//		Addr:   "上海",
	//		Age:    18,
	//	},
	//})

	//	读取数据
	var authModel AuthModel
	DB.First(&authModel)
	fmt.Println(authModel)

}
