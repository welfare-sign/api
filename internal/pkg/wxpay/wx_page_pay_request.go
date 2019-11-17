package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
)

type SignType string

const (
	SignType_MD5         SignType = "MD5"
	SignType_Hmac_SHA256 SignType = "HMAC-SHA256"
)

var (
	apiKey string //微信API密钥
)

func init() {
	apiKey = viper.GetString(config.KeyWXPayAPI)
}

type WxPagePayRequest struct {
	Values map[string]string
}

func (p *WxPagePayRequest) SetValue(key string, value string) {
	if p.Values == nil {
		p.Values = make(map[string]string)
	}
	p.Values[key] = value
}

func (p *WxPagePayRequest) DelValue(key string) {
	if p.Values == nil {
		return
	}

	delete(p.Values, key)
}

func (p *WxPagePayRequest) GetValue(key string) string {
	val, ok := p.Values[key]
	if !ok {
		return ""
	}
	return val
	//tp := reflect.TypeOf(val).Kind()
	//if tp == reflect.String {
	//	return val.(string)
	//} else if tp == reflect.Int {
	//	return strconv.Itoa(val.(int))
	//}
	//return ""
}

func (p *WxPagePayRequest) IsSet(key string) bool {
	val := p.GetValue(key)
	return len(val) > 0
}

func (p *WxPagePayRequest) MakeSign(signType SignType) string {

	//待签名字符串
	signStr := p.getSignString(apiKey)
	if signType == SignType_Hmac_SHA256 {
		return computeHmacSha256Hash(signStr, apiKey)
	} else {
		return computeMD5Hash(signStr)
	}
}

func (p *WxPagePayRequest) getSignString(apiKey string) string {
	var (
		i    = 0
		size = len(p.Values)
	)

	//第一步：把字典按Key的字母顺序排序
	sortedParams := make([]string, size)
	for key := range p.Values {
		sortedParams[i] = key
		i++
	}
	sort.Strings(sortedParams)

	//第二步：把所有参数名和参数值串在一起
	buffer := bytes.Buffer{}
	for _, key := range sortedParams {
		val := p.GetValue(key)
		if key == "sign" || len(val) == 0 {
			continue
		}
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString(val)
		buffer.WriteString("&")

	}

	//在string后加入API KEY
	buffer.WriteString("key=")
	buffer.WriteString(apiKey)

	return buffer.String()
}

func computeHmacSha256Hash(plaintext string, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(plaintext))
	hexStr := hex.EncodeToString(h.Sum(nil))

	return strings.ToUpper(hexStr)
}

func computeMD5Hash(plaintext string) string {
	hash := md5.New()
	hash.Write([]byte(plaintext))
	hexStr := hex.EncodeToString(hash.Sum(nil))

	return strings.ToUpper(hexStr)
}

func (p *WxPagePayRequest) CheckSign(signType SignType) bool {
	if !p.IsSet("sign") {
		return false
	}

	sign := p.GetValue("sign")
	//在计算参数签名
	calcSign := p.MakeSign(signType)

	return sign == calcSign
}

func (p *WxPagePayRequest) ToJson() string {
	str, err := json.Marshal(p.Values)
	if err != nil {
		return ""
	}

	return string(str)
}

func (p *WxPagePayRequest) ToXml() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("<xml>")
	for key, val := range p.Values {

		buffer.WriteString("<")
		buffer.WriteString(key)
		buffer.WriteString(">")

		//tp := reflect.TypeOf(val).Kind()
		//if tp == reflect.String {
		buffer.WriteString("<![CDATA[")
		buffer.WriteString(val)
		buffer.WriteString("]]>")
		//} else if tp == reflect.Int {
		//	buffer.WriteString(strconv.Itoa(val.(int)))
		//}

		buffer.WriteString("</")
		buffer.WriteString(key)
		buffer.WriteString(">")
	}
	buffer.WriteString("</xml>")

	return buffer.String()
}

func FromXml(body string) (*WxPagePayRequest, error) {
	var (
		params map[string]string
	)

	err := xml.Unmarshal([]byte(body), (*XmlMap)(&params))
	if err != nil {
		return nil, errors.WithMessage(err, "")
	}
	code, ok := params["result_code"]
	if !ok || code != "SUCCESS" {
		errDesc, _ := params["err_code_des"]
		return nil, errors.New(errDesc)
	}

	wxPayRequest := &WxPagePayRequest{
		Values: params,
	}

	//wxPayRequest.SetValue("return_code", unifiedResp.ReturnCode)
	//wxPayRequest.SetValue("return_msg", unifiedResp.ReturnMsg)
	//wxPayRequest.SetValue("appid", unifiedResp.AppId)
	//wxPayRequest.SetValue("mch_id", unifiedResp.MchId)
	//wxPayRequest.SetValue("nonce_str", unifiedResp.NonceStr)
	//wxPayRequest.SetValue("openid", unifiedResp.OpenId)
	//wxPayRequest.SetValue("sign", unifiedResp.Sign)
	//wxPayRequest.SetValue("result_code", unifiedResp.ResultCode)
	//wxPayRequest.SetValue("prepay_id", unifiedResp.PrepayId)
	//wxPayRequest.SetValue("trade_type", unifiedResp.TradeType)
	//wxPayRequest.SetValue("code_url", unifiedResp.CodeUrl)
	//wxPayRequest.SetValue("err_code", unifiedResp.ErrCode)
	//wxPayRequest.SetValue("err_code_des", unifiedResp.ErrCodeDes)
	//wxPayRequest.SetValue("bank_type", unifiedResp.BankType)
	//wxPayRequest.SetValue("cash_fee", unifiedResp.CashFee)
	//wxPayRequest.SetValue("fee_type", unifiedResp.FeeType)
	//wxPayRequest.SetValue("is_subscribe", unifiedResp.IsSubscribe)
	//wxPayRequest.SetValue("out_trade_no", unifiedResp.OutTradeNo)
	//wxPayRequest.SetValue("time_end", unifiedResp.TimeEnd)
	//wxPayRequest.SetValue("total_fee", unifiedResp.TotalFee)
	//wxPayRequest.SetValue("transaction_id", unifiedResp.TransactionId)

	flag := wxPayRequest.CheckSign(SignType_Hmac_SHA256)
	if !flag {
		return nil, errors.New("签名错误！")
	}

	return wxPayRequest, nil
}

func (p *WxPagePayRequest) Success() string {
	p.SetValue("return_code", "SUCCESS")
	p.SetValue("return_msg", "OK")

	return p.ToXml()
}
