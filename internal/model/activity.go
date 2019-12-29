package model

import "time"

// Activity 活动
type Activity struct {
	Base

	Name        string    `json:"name" gorm:"not null"`                     // 活动名称
	StartTime   time.Time `json:"start_time" gorm:"not null;type:datetime"` // 开始时间
	EndTime     time.Time `json:"end_time" gorm:"not null;type:datetime"`   // 结束时间
	Scope       string    `json:"scope"  gorm:"type:char(1);not null"`      // 参与方式
	PrizeAmount uint64    `json:"prize_amount" gorm:"not null"`             // 奖品数量
	Poster      string    `json:"poster" gorm:"not null"`                   // 活动海报
	PrizeNumber uint64    `json:"prize_number"`                             // 中奖号码
	PrizeIssued uint64    `json:"prize_issued" gorm:"-"`                    // 已发放的奖品数目
}

// ActivityVO 添加活动参数
type ActivityVO struct {
	Name        string    `json:"name" binding:"required"`         // 活动名称
	StartTime   time.Time `json:"start_time" binding:"required"`   // 开始时间
	EndTime     time.Time `json:"end_time" binding:"required"`     // 结束时间
	Scope       string    `json:"scope"  binding:"required"`       // 参与方式：O：无限制；R：必须在福利签领取过签到奖励
	PrizeAmount uint64    `json:"prize_amount" binding:"required"` // 奖品数量
	Poster      string    `json:"poster" binding:"required"`       // 活动海报
}

// ActivityListVO 查询活动列表参数
type ActivityListVO struct {
	Name     string `form:"name" json:"name"`
	PageNo   int    `form:"page_no" json:"page_no"`
	PageSize int    `form:"page_size" json:"page_size"`
}

// ActivityParticipantListVO 查询活动参与者列表参数
type ActivityParticipantListVO struct {
	ActivityID  uint64 `form:"activity_id" json:"activity_id"`
	Mobile      string `form:"mobile" json:"mobile"`
	IsSearchWin bool   `form:"is_search_win" json:"is_search_win"`
	PageNo      int    `form:"page_no" json:"page_no"`
	PageSize    int    `form:"page_size" json:"page_size"`
}
