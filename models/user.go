// models/user.go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
	// Adicione GoogleID string `json:"google_id" gorm:"uniqueIndex"` se quiser vincular
}
