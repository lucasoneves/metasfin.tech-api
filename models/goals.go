package models

type Goal struct {
	GormModel
	Title       string  `json:"title"`
	Description string  `json:"description"`
	TargetValue float64 `json:"target_value"`
	Balance     float64 `json:"balance"`
	UserID      uint    `json:"user_id"`
	Active      bool    `json:"active"`
	Completed   bool    `json:"completed"`
}
