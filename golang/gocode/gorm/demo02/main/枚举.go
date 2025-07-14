package main

import (
	"encoding/json"
	"fmt"
)

type Status int

const (
	Running Status = 1
	OffLine Status = 2
	Except  Status = 3
)

type Host struct {
	ID     uint   `json:"id"`
	Status Status `gorm:"size:8" json:"status"`
	IP     string `json:"ip"`
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Status) String() string {
	var str string
	switch s {
	case Running:
		str = "Running"
	case Except:
		str = "Except"
	case OffLine:
		str = "Status"
	}
	return str
}

func main() {
	//DB.AutoMigrate(&Host{})

	//DB.Create(&Host{
	//	IP:     "192.168.200.12",
	//	Status: Running,
	//})
	var host Host
	DB.Take(&host)
	fmt.Println(host)
	fmt.Printf("%#v,%T\n", host.Status, host.Status)
	data, _ := json.Marshal(host)
	fmt.Println(string(data))

}
