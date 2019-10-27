package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
)

// addMerchant 新增商户
// @Summary 新增商户
// @Description create new merchant
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body model.MerchantVO true "参数"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /merchants [post]
func addMerchant(c *gin.Context) {
	var vo model.MerchantVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   err.Error(),
		})
		return
	}
	if err := svc.AddMerchant(c, vo); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrModelCreate,
			Message: apicode.MapZH[apicode.ErrModelCreate],
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

// getMerchantList 获取商户列表
// @Summary 获取商户列表
// @Description get merchant list
// @Tags 商户
// @Accept json
// @Produce json
// @Param store_name query string false "商户名"
// @Param contact_name query string false "联系人"
// @Param contact_phone query string false "联系电话"
// @Param page_no query int false "页码" default(1)
// @Param page_size query int false "页数" default(10)
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /merchants [get]
func getMerchantList(c *gin.Context) {
	var vo model.MerchantListVO
	c.ShouldBindQuery(&vo)

	merchants, err := svc.GetMerchantList(c, vo)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrGetListData,
			Message: apicode.MapZH[apicode.ErrGetListData],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse{
		Status:  true,
		Code:    apicode.Success,
		Message: apicode.MapZH[apicode.Success],
		Data:    merchants,
	})
}

// merchantLogin 商户登录
// @Summary 商户登录
// @Description merchant login
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body model.MerchantLoginVO true "参数"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /merchants/login [post]
func merchantLogin(c *gin.Context) {

}
