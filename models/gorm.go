package models

import (
	"time"

	"gorm.io/gorm"
)

type GormModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"` // "omitempty" esconde o campo se for zero/nulo
}
