package model

// HelpCheckinMessage 补签消息
type HelpCheckinMessage struct {
	Base

	CustomerID      uint64 `json:"customer_id" gorm:"not null"` // 用户ID
	CheckinRecordID uint64 `json:"checkin_record_id"`           // 关联签到记录ID
	IsRead          string `json:"is_read" gorm:"type:char(1);not null"`
}
