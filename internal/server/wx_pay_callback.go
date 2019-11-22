package server

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
	"welfare-sign/internal/pkg/wxpay"
)

// WxpayCallbackRequest .
type WxpayCallbackRequest struct {
	wsgin.BaseRequest

	Params   string `form:"params" json:"params"`
	Response gin.ResponseWriter
}

// New .
func (m *WxpayCallbackRequest) New() wsgin.Process {
	return &WxpayCallbackRequest{}
}

// Extract .
func (m *WxpayCallbackRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	m.Response = c.Writer
	defer c.Request.Body.Close()
	body, _ := ioutil.ReadAll(c.Request.Body)
	m.Params = string(body)
	return
}

// Exec 微信支付回调
// @Summary 用户支付
// @Description wx pay callback
// @Tags 微信
// @Accept json
// @Produce json
// @Param args body server.WxpayCallbackRequest true "参数"
// @Success 200 {string} string	"{"status":true}"
// @Router /wx/pay/notify [post]
func (m *WxpayCallbackRequest) Exec(ctx context.Context) interface{} {
	callback := wxpay.WxPagePayRequest{}
	//支付状态确认
	code, _ := svc.WxpayCallback(ctx, m.Params)
	if code == wsgin.APICodeSuccess {
		//通知微信
		respXml := callback.Success()
		m.Response.Write([]byte(respXml))
	}

	return nil
}
