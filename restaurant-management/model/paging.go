package model

type Paging struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}
