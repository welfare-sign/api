package model

// WXPayRecord 微信支付流水记录
type WXPayRecord struct {
	Base

	OrderID         string `json:"order_id" gorm:"not null"`    // 内部订单号
	PayFee          uint64 `json:"pay_fee" gorm:"not null"`     // 支付金额，单位分
	TradeNo         string `json:"trade_no" gorm:"not null"`    // 微信内部订单号
	CustomerID      uint64 `json:"customer_id" gorm:"not null"` // 用户ID
	CompletePayTime string `json:"complete_pay_time"`           // 用户完成支付时间
	CheckinRecordID uint64 `json:"checkin_record_id"`           // 关联签到记录ID
}

// WXConfigResp 微信配置响应体
type WXConfigResp struct {
	Appid     string `json:"appid"`     // 公众号的唯一标识
	Timestamp int64  `json:"timestamp"` // 生成签名的时间戳
	Noncestr  string `json:"noncestr"`  // 生成签名的随机串
	Signature string `json:"signature"` // 签名
}
