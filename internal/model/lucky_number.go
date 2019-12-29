package model

// LuckyNumberRecord 幸运数字记录
type LuckyNumberRecord struct {
	Base

	CustomerID  uint64 `json:"customer_id" gorm:"not null"`  // 签到人ID
	ActivityID  uint64 `json:"activity_id" gorm:"not null"`  // 活动ID
	Mobile      string `json:"mobile" gorm:"not null"`       // 参与人手机号
	LuckyNumber uint64 `json:"lucky_number" gorm:"not null"` // 用户填写的幸运数字
	Ranking     uint64 `json:"ranking"`                      // 本轮幸运数字排名

	Activity *Activity `json:"activity" gorm:"-"` // 活动信息
	Customer *Customer `json:"customer" gorm:"-"` // 用户信息
}
