package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
)

// getCustomerList 获取客户列表
// @Summary 获取客户列表
// @Description get customer list
// @Tags 客户
// @Accept json
// @Produce json
// @Param name query string false "用户名"
// @Param mobile query string false "联系电话"
// @Param page_no query int false "页码" default(1)
// @Param page_size query int false "页数" default(10)
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /customers [get]
func getCustomerList(c *gin.Context) {
	var vo model.CustomerListVO
	c.ShouldBindQuery(&vo)

	customers, err := svc.GetCustomerList(c, vo)
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
		Data:    customers,
	})
}

// customerDetail 获取客户详情
// @Summary 获取客户详情
// @Description get customer detail
// @Tags 客户
// @Accept json
// @Produce json
// @Param customer_id query string true "客户ID"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /customers/detail [get]
func customerDetail(c *gin.Context) {
	customerId := c.Query("customer_id")
	if customerId == "" {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
		})
		return
	}
	_customerId, err := strconv.ParseUint(customerId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
		})
		return
	}
	customer, err := svc.GetCustomerDetail(c, _customerId)
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
		Data:    customer,
	})
}
