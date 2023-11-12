package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        uint64    `gorm:"column:id;type:bigint;primary_key,AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	return nil
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()

	return nil
}
