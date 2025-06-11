package models

type Goal struct {
	GormModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Balance     int    `json:"balance"`
	UserID      uint   `json:"user_id"`
}
