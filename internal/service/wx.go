package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/util"
	"welfare-sign/internal/pkg/wsgin"
	"welfare-sign/internal/pkg/wxpay"
)

// GetWXConfig 获取微信配置
func (s *Service) GetWXConfig(ctx context.Context, url string) (*model.WXConfigResp, wsgin.APICode, error) {
	var c model.WXConfigResp
	appID := viper.GetString(config.KeyWxAppID)
	appSecret := viper.GetString(config.KeyWxAppSecret)
	c.Appid = appID
	c.Timestamp = time.Now().Unix()
	c.Noncestr = uuid.NewV4().String()

	// get jsapi_ticket
	ticket, err := s.dao.GetWXJSTicket()
	if err != nil {
		log.Warn(ctx, "GetWXConfig.GetWXJSTicket() error", zap.Error(err))
		return nil, apicode.ErrGetWXConfig, err
	}
	if ticket == "" {
		accessToken, err := s.dao.GetWXAccessToken()
		if err != nil {
			log.Warn(ctx, "GetWXConfig.GetWXAccessToken() error", zap.Error(err))
			return nil, apicode.ErrGetWXConfig, err
		}
		if accessToken == "" {
			// get access_token
			resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appId=%s&secret=%s", appID, appSecret))
			if err != nil {
				return nil, apicode.ErrGetWXConfig, err
			}
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, apicode.ErrGetWXConfig, err
			}
			accessToken = gjson.GetBytes(bytes, "access_token").String()
			if accessToken == "" {
				return nil, wsgin.APICodeSuccess, errors.New("获取access_token失败")
			}
			accessTokenExpire := time.Second * time.Duration(gjson.GetBytes(bytes, "expires_in").Int()-60)
			if err := s.dao.StoreWXAccessToken(accessToken, accessTokenExpire); err != nil {
				log.Warn(ctx, "GetWXConfig.StoreWXAccessToken() error", zap.Error(err))
				return nil, apicode.ErrGetWXConfig, err
			}
		}

		resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken))
		if err != nil {
			return nil, apicode.ErrGetWXConfig, err
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, apicode.ErrGetWXConfig, err
		}
		ticket = gjson.GetBytes(bytes, "ticket").String()
		if ticket == "" {
			return nil, wsgin.APICodeSuccess, errors.New("获取jsapi_ticket失败")
		}
		ticketExpire := time.Second * time.Duration(gjson.GetBytes(bytes, "expires_in").Int()-60)
		if err := s.dao.StoreWXJSTicket(ticket, ticketExpire); err != nil {
			log.Warn(ctx, "GetWXConfig.StoreWXJSTicket() error", zap.Error(err))
			return nil, apicode.ErrGetWXConfig, err
		}
	}
	sign := util.GetWXSignString(map[string]string{
		"noncestr":     c.Noncestr,
		"jsapi_ticket": ticket,
		"timestamp":    strconv.FormatInt(c.Timestamp, 10),
		"url":          url,
	})
	c.Signature = sign
	return &c, wsgin.APICodeSuccess, nil
}

// WXPay 微信支付
func (s *Service) WXPay(ctx context.Context, customerID uint64) (string, wsgin.APICode, error) {
	customer, _ := s.dao.FindCustomer(ctx, map[string]interface{}{
		"id":     customerID,
		"status": global.ActiveStatus,
	})
	if customer.ID == 0 {
		return "", apicode.ErrWXPay, errors.New("未查到用户信息")
	}

	req := prepareWxpayRequest(ctx, customer.OpenID)
	ret, err := wxpay.UnifiedOrder(req)
	if err != nil {
		return "", apicode.ErrWXPay, errors.WithMessage(err, "当前订单无法支付，请稍候再试")
	}
	//微信小程序支付prepay_id
	prepayId := ret.GetValue("prepay_id")

	req = miniWxpaySign(prepayId)
	return req.ToJson(), wsgin.APICodeSuccess, nil
}

func miniWxpaySign(prepayId string) *wxpay.WxPagePayRequest {
	request := wxpay.WxPagePayRequest{}
	request.SetValue("appId", viper.GetString(config.KeyWxAppID))
	request.SetValue("timeStamp", strconv.FormatInt(time.Now().Unix(), 10))
	request.SetValue("nonceStr", util.GenerateNonceStr(20))
	request.SetValue("package", "prepay_id="+prepayId)
	request.SetValue("signType", string(wxpay.SignType_Hmac_SHA256))
	request.SetValue("paySign", request.MakeSign(wxpay.SignType_Hmac_SHA256))
	request.SetValue("payFee", strconv.FormatFloat(viper.GetFloat64(config.KeyWXPayAmount)*100, 'f', 0, 64))

	request.DelValue("appId")
	return &request
}

func prepareWxpayRequest(ctx context.Context, openId string) *wxpay.WxPagePayRequest {
	//应付金额
	payPrice := viper.GetFloat64(config.KeyWXPayAmount)
	year, month, day := time.Now().Date()
	orderNo := strconv.Itoa(year) + strconv.Itoa(int(month)) + strconv.Itoa(day) + strconv.FormatInt(time.Now().Unix(), 10)
	request := wxpay.WxPagePayRequest{}
	request.SetValue("body", "支付5元补签")
	request.SetValue("out_trade_no", "J"+orderNo)
	request.SetValue("total_fee", strconv.FormatFloat(payPrice*100, 'f', 0, 64)) //分
	request.SetValue("trade_type", "JSAPI")
	request.SetValue("openid", openId)
	request.SetValue("notify_url", viper.GetString(config.KeyWXPayNotifyURL))

	return &request
}

// WxpayCallback 微信支付回调
func (s *Service) WxpayCallback(ctx context.Context, notifyData string) (wsgin.APICode, error) {
	ok, wxPayNotifyRequest := checkWeixinPayValidation(notifyData)
	if !ok {
		return wsgin.APICodeSuccess, nil
	}
	fmt.Println("WxpayCallback.wxPayNotifyRequest: ", wxPayNotifyRequest)
	soid := wxPayNotifyRequest.GetValue("out_trade_no") //订单号
	openid := wxPayNotifyRequest.GetValue("openid")
	completePayTime := wxPayNotifyRequest.GetValue("time_end") // 支付完成时间

	// 转换支付订单价格（分）
	payFee, _ := strconv.ParseUint(wxPayNotifyRequest.GetValue("total_fee"), 10, 64)

	// 微信支付订单交易号
	tradeNo := wxPayNotifyRequest.GetValue("transaction_id")

	return s.payOrderComplete(ctx, soid, uint64(payFee), tradeNo, openid, completePayTime)
}

func checkWeixinPayValidation(notifyData string) (bool, *wxpay.WxPagePayRequest) {
	payRequest, err := wxpay.FromXml(notifyData)
	if err != nil {
		return false, nil
	}

	//APPID
	if viper.GetString(config.KeyWxAppID) != payRequest.GetValue("appid") {
		return false, nil
	}
	//商户号
	if viper.GetString(config.KeyWXPayMchID) != payRequest.GetValue("mch_id") {
		return false, nil
	}
	return true, payRequest
}

func (s *Service) payOrderComplete(ctx context.Context, orderId string, payFee uint64, tradeNo, openid, completePayTime string) (wsgin.APICode, error) {
	//TODO: 更改签到记录, 校验金额
	customer, err := s.dao.FindCustomer(ctx, map[string]interface{}{
		"open_id": openid,
		"status":  global.ActiveStatus,
	})
	if err != nil {
		return apicode.ErrWXPayNotify, err
	}
	if customer.ID == 0 {
		return apicode.ErrWXPayNotify, errors.New("用户未找到")
	}

	if payFee != 500 {
		log.Warn(ctx, "payOrderComplete.payFee error", zap.Uint64("实际支付金额", payFee))
		return apicode.ErrWXPayNotify, errors.New("支付金额不正确")
	}

	unchecked, err := s.dao.GetUnchecked(ctx, customer.ID)
	if err != nil {
		return apicode.ErrWXPayNotify, err
	}
	if unchecked.ID == 0 {
		log.Warn(ctx, "payOrderComplete.GetUnchecked()", zap.String("notify: ", "当前用户没有需要补签的记录"))
		return wsgin.APICodeSuccess, nil
	}

	if err := s.dao.PayCheckin(ctx, unchecked.ID, customer.ID, &model.WXPayRecord{
		OrderID:         orderId,
		PayFee:          payFee,
		TradeNo:         tradeNo,
		CustomerID:      customer.ID,
		CompletePayTime: completePayTime,
	}); err != nil {
		return apicode.ErrWXPayNotify, err
	}

	return wsgin.APICodeSuccess, nil
}
