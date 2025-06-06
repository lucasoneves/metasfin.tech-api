package models

type Goal struct {
	GormModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	UserID      uint   `json:"user_id"`
}
