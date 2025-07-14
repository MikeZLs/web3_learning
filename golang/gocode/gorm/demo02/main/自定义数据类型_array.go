package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Array []string

// Scan 从数据库中读取出来
func (arr *Array) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, arr)
	return err
}

// Value 存入数据库
func (arr Array) Value() (driver.Value, error) {
	return json.Marshal(arr)
}

type HostModel struct {
	ID    uint   `json:"id"`
	IP    string `json:"ip"`
	Ports Array  `gorm:"type:string" json:"ports"`
}

func main() {
	//DB.AutoMigrate(&HostModel{})
	//
	//DB.Create(&HostModel{
	//	IP:    "192.168.200.21",
	//	Ports: []string{"80", "8080"},
	//})
	var host HostModel
	DB.Take(&host, 1)
	fmt.Println(host)
}
