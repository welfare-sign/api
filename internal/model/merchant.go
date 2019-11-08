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
	StoreName    string  `json:"store_name" binding:"required"`    // 店名
	Address      string  `json:"address" binding:"required"`       // 地址
	Lon          float64 `json:"lon" binding:"required"`           // 经度
	Lat          float64 `json:"lat" binding:"required"`           // 纬度
	CateringType string  `json:"catering_type"`                    // 餐饮类型
	StoreAvatar  string  `json:"store_avatar" binding:"required"`  // 店铺头像
	Poster       string  `json:"poster"`                           // 商户海报
	ContactName  string  `json:"contact_name"`                     // 联系人
	ContactPhone string  `json:"contact_phone" binding:"required"` // 联系人电话
	Received     uint64  `json:"-"`                                // 已领礼品数量
	TotalReceive uint64  `json:"total_receive" binding:"required"` // 该店礼品一共可领取总数
	CheckinDays  uint64  `json:"checkin_days"`                     // 签到天数多少天可领取礼品
	CheckinNum   uint64  `json:"checkin_num"`                      // 达到指定签到天数后，可领取的礼品数量
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
	ContactPhone string `json:"contact_phone" binding:"required" example:"手机号"`
	Code         string `json:"code" binding:"required" example:"验证码"`
}

// MerchantWriteOffVO 获取商户核销参数
type MerchantWriteOffVO struct {
	AccessToken string `form:"access_token" json:"access_token" binding:"required" example:"商户token"`
	CustomerID  uint64 `form:"customer_id" json:"customer_id" binding:"required" example:"客户ID"`
}

// MerchantWriteOffRespVO 商户核销响应
type MerchantWriteOffRespVO struct {
	Merchant    *Merchant    `json:"merchant"`
	Customer    *Customer    `json:"customer"`
	IssueRecord *IssueRecord `json:"issue_record"`
}

// MerchantExecWriteOffVO 商户执行核销参数
type MerchantExecWriteOffVO struct {
	MerchantID uint64
	CustomerID uint64
	Num        uint64
}

// NearMerchantVO 附近商家
type NearMerchantVO struct {
	Lon      float64 // 经度
	Lat      float64 // 维度
	Distince float64 // 距离多少公里内
	Num      int     // 返回数量
}
