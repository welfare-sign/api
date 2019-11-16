package model

// WXConfigResp 微信配置响应体
type WXConfigResp struct {
	Appid     string `json:"appid"`     // 公众号的唯一标识
	Timestamp int64  `json:"timestamp"` // 生成签名的时间戳
	Noncestr  string `json:"noncestr"`  // 生成签名的随机串
	Signature string `json:"signature"` // 签名
}
