package models

type Model struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PageInfo struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
	Key   string `form:"key"`
}
