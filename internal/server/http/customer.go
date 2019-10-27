package http

import (
	"net/http"

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
