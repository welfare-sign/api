package model

// LuckyNumberRecord 幸运数字记录
type LuckyNumberRecord struct {
	Base

	CustomerID  uint64 `json:"customer_id" gorm:"not null"`  // 签到人ID
	LuckyNumber int64  `json:"lucky_number" gorm:"not null"` // 用户填写的幸运数字
	Ranking     uint64 `json:"ranking"`                      // 本轮幸运数字排名

	CompositeIndex *CompositeIndex `json:"composite_index" gorm:"-"` // 本轮内上证指数
}
