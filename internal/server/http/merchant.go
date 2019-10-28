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
	var vo model.MerchantLoginVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   err.Error(),
		})
		return
	}
	token, err := svc.MerchantLogin(c, vo)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrLogin,
			Message: apicode.MapZH[apicode.ErrLogin],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse{
		Status:  true,
		Code:    apicode.Success,
		Message: apicode.MapZH[apicode.Success],
		Data:    token,
	})
}

// merchantDetail 获取商户详情信息
// @Summary 获取商户详情信息
// @Description get merchant detail
// @Tags 商户
// @Accept json
// @Produce json
// @Param access_token query string true "商户token"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /merchants/detail [get]
func merchantDetail(c *gin.Context) {
	token := c.Query("access_token")
	if token == "" {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
		})
		return
	}
	merchant, err := svc.GetMerchantDetailBySelfAccessToken(c, token)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrDetail,
			Message: apicode.MapZH[apicode.ErrDetail],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse{
		Status:  true,
		Code:    apicode.Success,
		Message: apicode.MapZH[apicode.Success],
		Data:    merchant,
	})
}

// writeOff 获取商家核销信息
// @Summary 获取商户核销信息
// @Description get merchant write off
// @Tags 商户
// @Accept json
// @Produce json
// @Param access_token query string true "商户token"
// @Param customer_id query int true "客户ID"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /merchants/writeoff [get]
func getWriteOff(c *gin.Context) {
	var vo model.MerchantWriteOffVO
	if err := c.ShouldBindQuery(&vo); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   err.Error(),
		})
		return
	}
	resp, err := svc.GetWriteOff(c, vo)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrWriteOff,
			Message: apicode.MapZH[apicode.ErrWriteOff],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse{
		Status:  true,
		Code:    apicode.Success,
		Message: apicode.MapZH[apicode.Success],
		Data:    resp,
	})
}

// writeOff 核销
// @Summary 商户核销
// @Description merchant exec write off
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body model.MerchantExecWriteOffVO true "参数"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /merchants/writeoff [post]
func writeOff(c *gin.Context) {
	var vo model.MerchantExecWriteOffVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   err.Error(),
		})
		return
	}
	resp, err := svc.ExecWriteOff(c, vo)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrExecWriteOff,
			Message: apicode.MapZH[apicode.ErrExecWriteOff],
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse{
		Status:  true,
		Code:    apicode.Success,
		Message: apicode.MapZH[apicode.Success],
		Data:    resp,
	})
}
