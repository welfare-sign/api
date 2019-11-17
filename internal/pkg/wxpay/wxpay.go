package wxpay

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/util"
)

// 统一下单
func UnifiedOrder(inputObj *WxPagePayRequest) (*WxPagePayRequest, error) {
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	//检查必填参数
	if !inputObj.IsSet("out_trade_no") {
		return nil, errors.New("缺少统一支付接口必填参数out_trade_no！")
	}
	if !inputObj.IsSet("body") {
		return nil, errors.New("缺少统一支付接口必填参数body！")
	}
	if !inputObj.IsSet("total_fee") {
		return nil, errors.New("缺少统一支付接口必填参数total_fee！")
	}
	if !inputObj.IsSet("trade_type") {
		return nil, errors.New("缺少统一支付接口必填参数trade_type！")
	}
	if !inputObj.IsSet("notify_url") {
		return nil, errors.New("缺少统一支付接口必填参数notify_url！")
	}

	//关联参数
	if inputObj.GetValue("trade_type") == "NATIVE" && !inputObj.IsSet("product_id") {
		return nil, errors.New("统一支付接口中，缺少必填参数product_id！trade_type为NATIVE时，product_id为必填参数！")
	}

	inputObj.SetValue("appid", viper.GetString(config.KeyWxAppID))     //公众账号ID
	inputObj.SetValue("mch_id", viper.GetString(config.KeyWXPayMchID)) //商户号
	inputObj.SetValue("spbill_create_ip", util.GetIP())                //终端ip
	inputObj.SetValue("nonce_str", util.GenerateNonceStr(20))          //随机字符串
	inputObj.SetValue("sign_type", string(SignType_Hmac_SHA256))       //签名类型

	//签名
	inputObj.SetValue("sign", inputObj.MakeSign(SignType_Hmac_SHA256))
	xml := inputObj.ToXml()
	respXml, err := Post(xml, url)
	if err != nil { //调用微信支付接口时发生异常
		return nil, err
	}
	//解析微信支付返回结果
	payRequest, err := FromXml(respXml)
	if err != nil {
		return nil, err
	}
	return payRequest, nil
}

//微信支付订单查询
func OrderQuery(inputObj *WxPagePayRequest) (string, error) {
	url := "https://api.mch.weixin.qq.com/pay/orderquery"
	//检查必填参数
	if !inputObj.IsSet("out_trade_no") && !inputObj.IsSet("transaction_id") {
		return "", errors.New("订单查询接口中，out_trade_no、transaction_id至少填一个！")
	}
	inputObj.SetValue("appid", viper.GetString(config.KeyWxAppID))     //公众账号ID
	inputObj.SetValue("mch_id", viper.GetString(config.KeyWXPayMchID)) //商户号
	inputObj.SetValue("nonce_str", util.GenerateNonceStr(20))          //随机字符串
	inputObj.SetValue("sign_type", string(SignType_Hmac_SHA256))       //签名类型
	inputObj.SetValue("sign", inputObj.MakeSign(SignType_Hmac_SHA256)) //签名

	respXml, err := Post(inputObj.ToXml(), url)
	if err != nil {
		return "", errors.WithMessage(err, "查询微信支付订单状态异常！")
	}

	return respXml, nil
	//解析微信支付返回结果
	//payRequest, err := FromXml(respXml)

	/*
		<xml><return_code><![CDATA[SUCCESS]]></return_code>
		<return_msg><![CDATA[OK]]></return_msg>
		<appid><![CDATA[wx735916a7779a132e]]></appid>
		<mch_id><![CDATA[1540992911]]></mch_id>
		<device_info><![CDATA[]]></device_info>
		<nonce_str><![CDATA[W7glVX3Spw3u6JJr]]></nonce_str>
		<sign><![CDATA[66912C701CA498397D0F0B036C18ADBC9AD8A5E3A9A8FC437C2B317E0B516D6A]]></sign>
		<result_code><![CDATA[SUCCESS]]></result_code>
		<total_fee>2500</total_fee>
		<out_trade_no><![CDATA[1000172152]]></out_trade_no>
		<trade_state><![CDATA[NOTPAY]]></trade_state>
		<trade_state_desc><![CDATA[订单未支付]]></trade_state_desc>
		</xml>

		trade_state（交易状态）：SUCCESS—支付成功，REFUND—转入退款，OTPAY—未支付，CLOSED—已关闭，REVOKED—已撤销（付款码支付），	USERPAYING--用户支付中（付款码支付），PAYERROR--支付失败(其他原因，如银行返回失败)，支付状态机请见下单API页面
	*/

	//if err != nil {
	//	return nil, err
	//}
	//return payRequest, nil
}
