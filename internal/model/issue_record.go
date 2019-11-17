package model

// IssueRecord 礼品发放记录
type IssueRecord struct {
	Base

	MerchantID   uint64    `json:"merchant_id" gorm:"not null"` // 店铺ID
	CustomerID   uint64    `json:"customer_id" gorm:"not null"` // 顾客ID
	TotalReceive uint64    `json:"total_receive"`               // 在对应店铺可领取的总礼品数
	Received     uint64    `json:"received"`                    // 在对应店铺已兑换的礼品数
	Merchant     *Merchant `json:"merchant" gorm:"-"`
	Customer     *Customer `json:"customer" gorm:"-"`
}
