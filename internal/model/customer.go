package model

import "time"

// Customer 顾客
type Customer struct {
	Base

	OpenID          string     `json:"open_id" gorm:"not null"`                 // 微信用户openid
	Nickname        string     `json:"nickname"`                                // 微信用户昵称
	Sex             int        `json:"sex"`                                     // 微信用户性别
	Country         string     `json:"country"`                                 // 微信用户所在国家
	Province        string     `json:"province"`                                // 微信用户所在市
	City            string     `json:"city"`                                    // 微信用户所在区
	Headimgurl      string     `json:"headimgurl"`                              // 微信用户头像
	Name            string     `json:"name" gorm:"not null"`                    // 称呼
	Mobile          string     `json:"mobile" gorm:"type:varchar(50);not null"` // 手机号
	LastCheckinTime *time.Time `json:"last_checkin_time" gorm:"type:datetime"`  // 最后一次签到时间
}

// CustomerListVO 查询顾客列表参数
type CustomerListVO struct {
	Name     string `form:"name" json:"name"`
	Mobile   string `form:"mobile" json:"mobile"`
	PageNo   int    `form:"page_no" json:"page_no"`
	PageSize int    `form:"page_size" json:"page_size"`
}

// CustomerLoginResp 客户登录之后返回的数据
type CustomerLoginResp struct {
	Customer          *Customer        `json:"customer"`            // 用户信息
	CheckinRecordList []*CheckinRecord `json:"checkin_record_list"` // 用户签到信息
}

// WxErrResp .
type WxErrResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// WxSuccessResp .
type WxSuccessResp struct {
	AccessToken  string `json:"access_token"`  // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn    int    `json:"expires_in"`    // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` // 用户刷新access_token,有效期为30天，当refresh_token失效之后，需要用户重新授权。
	OpenID       string `json:"openid"`        // 用户唯一标识
	Scope        string `json:"scope"`         // 用户授权的作用域，使用逗号（,）分隔
}

// WxUserResp .
type WxUserResp struct {
	OpenID     string `json:"openid"`     // 用户的唯一标识
	Nickname   string `json:"nickname"`   // 用户昵称
	Sex        int    `json:"sex"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province   string `json:"province"`   // 用户个人资料填写的省份
	City       string `json:"city"`       // 普通用户个人资料填写的城市
	Country    string `json:"country"`    // 国家
	Headimgurl string `json:"headimgurl"` // 用户头像
}
