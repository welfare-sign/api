package model

import "time"

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

// Base base model contain some basic attr
type Base struct {
	ID        uint64    `json:"id" gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;type:datetime"`
	UpdatedBy uint64    `json:"updated_by" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;type:datetime"`
	CreatedBy uint64    `json:"created_by" gorm:"not null"`
	Status    string    `json:"status" gorm:"type:char(1);not null"`
}

// SetDefaultAttr 设置默认的属性值
func (b *Base) SetDefaultAttr() {
	b.Status = "A"
	b.CreatedAt = time.Now()
	b.CreatedBy = 0
	b.UpdatedAt = time.Now()
	b.UpdatedBy = 0
}
