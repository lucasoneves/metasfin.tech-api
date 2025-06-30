package models

type Goal struct {
	GormModel
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetValue int    `json:"target_value"`
	Balance     int    `json:"balance"`
	UserID      uint   `json:"user_id"`
	Active      bool   `json:"active"`
	Completed   bool   `json:"completed"`
}
