package model

// IssueRecord 礼品发放记录
type IssueRecord struct {
	Base `json:"base,omitempty"`

	MerchantID   uint64   `json:"merchant_id,omitempty" gorm:"not null"` // 店铺ID
	CustomerID   uint64   `json:"customer_id,omitempty" gorm:"not null"` // 顾客ID
	TotalReceive uint64   `json:"total_receive,omitempty"`               // 在对应店铺可领取的总礼品数
	Received     uint64   `json:"received,omitempty"`                    // 在对应店铺已兑换的礼品数
	Merchant     Merchant `json:"merchant,omitempty" gorm:"-"`
	Customer     Customer `json:"customer,omitempty" gorm:"-"`
}
