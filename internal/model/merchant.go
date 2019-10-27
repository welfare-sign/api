package model

// Merchant 店铺
type Merchant struct {
	Base

	StoreName    string  `json:"store_name" gorm:"not null"`                            // 店名
	Address      string  `json:"address" gorm:"not null"`                               // 地址
	Lon          float64 `json:"lon" gorm:"not null"`                                   // 经度
	Lat          float64 `json:"lat" gorm:"not null"`                                   // 纬度
	CateringType string  `json:"catering_type" gorm:"type:varchar(100)"`                // 餐饮类型
	StoreAvatar  string  `json:"store_avatar" gorm:"not null"`                          // 店铺头像
	Poster       string  `json:"poster"`                                                // 商户海报
	ContactName  string  `json:"contact_name"`                                          // 联系人
	ContactPhone string  `json:"contact_phone" gorm:"unique;not null;type:varchar(50)"` // 联系人电话
	Received     uint64  `json:"received"`                                              // 已领礼品数量
	TotalReceive uint64  `json:"total_receive"`                                         // 该店礼品一共可领取总数
	CheckinDays  uint64  `json:"checkin_days"`                                          // 签到天数多少天可领取礼品
	CheckinNum   uint64  `json:"checkin_num"`                                           // 达到指定签到天数后，可领取的礼品数量
}

// MerchantVO 新增店铺参数
type MerchantVO struct {
	StoreName    string  `json:"store_name" validate:"required"`    // 店名
	Address      string  `json:"address" validate:"required"`       // 地址
	Lon          float64 `json:"lon" validate:"required"`           // 经度
	Lat          float64 `json:"lat" validate:"required"`           // 纬度
	CateringType string  `json:"catering_type"`                     // 餐饮类型
	StoreAvatar  string  `json:"store_avatar" validate:"required"`  // 店铺头像
	Poster       string  `json:"poster"`                            // 商户海报
	ContactName  string  `json:"contact_name"`                      // 联系人
	ContactPhone string  `json:"contact_phone" validate:"required"` // 联系人电话
	Received     uint64  `json:"received"`                          // 已领礼品数量
	TotalReceive uint64  `json:"total_receive" validate:"required"` // 该店礼品一共可领取总数
	CheckinDays  uint64  `json:"checkin_days"`                      // 签到天数多少天可领取礼品
	CheckinNum   uint64  `json:"checkin_num"`                       // 达到指定签到天数后，可领取的礼品数量
}

// MerchantListVO 获取店铺列表参数
type MerchantListVO struct {
	StoreName    string `form:"store_name" json:"store_name"`
	ContactName  string `form:"contact_name" json:"contact_name"`
	ContactPhone string `form:"contact_phone" json:"contact_phone"`
	PageNo       int    `form:"page_no" json:"page_no"`
	PageSize     int    `form:"page_size" json:"page_size"`
}

// MerchantLoginVO 商户登录参数
type MerchantLoginVO struct {
	ContactPhone string `json:"contact_phone"`
	Code         int    `json:"code"`
}
