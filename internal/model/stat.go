package model

// RegisterStat 统计注册用户结果
type RegisterStat struct {
	Date string `json:"date"` // 日期
	Num  uint64 `json:"num"`  // 注册用户数
}

// CheckinStat 统计用户签到结果
type CheckinStat struct {
	Date string `json:"date"` // 日期
	Num  uint64 `json:"num"`  // 签到次数
}
