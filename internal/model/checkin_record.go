package model

import "time"

// CheckinRecord 签到记录
type CheckinRecord struct {
	Base

	CustomerID            uint64    `json:"customer_id" gorm:"not null"`                     // 签到人ID
	HelpCheckinCustomerID uint64    `json:"help_checkin_customer_id" gorm:"default:0"`       // 帮签人ID
	Day                   uint64    `json:"day" gorm:"not null"`                             // 签到第几天
	NeedCheckinTime       time.Time `json:"need_checkin_time" gorm:"not null;type:datetime"` // 需要签到的日期
}

// CheckinRecordVO 记录客户签到
type CheckinRecordVO struct {
	AccessToken string `form:"access_token" json:"access_token" binding:"required" example:"客户token"`
	Day         uint64 `form:"customer_id" json:"customer_id" binding:"required" example:"客户ID"`
}

// CheckinRecordListResp 所有签到记录列表
type CheckinRecordListResp struct {
	CheckinRecord *CheckinRecord `json:"checkin_record"`
	Customer      *Customer      `json:"customer"`
}
