package model

// CompositeIndex 上证指数
type CompositeIndex struct {
	Base

	CompositeDate string  `json:"composite_date" gorm:"not null"` // 上证指数日期
	Points        float64 `json:"points" gorm:"not null"`         // 指数
}
