package model

// CheckinRecord 签到记录
type CheckinRecord struct {
	Base

	CustomerId            uint64 `json:"customer_id" gorm:"not null"` // 签到人ID
	HelpCheckinCustomerId uint64 `json:"help_checkin_customer_id"`    // 帮签人ID
	Day                   uint64 `json:"day" gorm:"not null"`         // 签到第几天
}
