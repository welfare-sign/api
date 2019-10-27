package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/pkg/util"
)

// getVerifyCode 发送验证码
// @Summary 发送验证码
// @Description send verify code
// @Tags 验证码
// @Accept json
// @Produce json
// @Param mobile query string false "手机号"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /verify_code [get]
func getVerifyCode(c *gin.Context) {
	mobile := c.Query("mobile")

	if mobile == "" || !util.IsMobile(mobile) {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrMobile,
			Message: apicode.MapZH[apicode.ErrMobile],
		})
		return
	}
	if err := svc.SendVerifyCode(c, mobile); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrSendSMS,
			Message: apicode.MapZH[apicode.ErrSendSMS],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse{
		Status:  true,
		Code:    apicode.Success,
		Message: apicode.MapZH[apicode.Success],
	})
}
