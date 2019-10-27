package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
)

// userLogin 后台用户登录
// @Summary 后台用户登录
// @Description user login
// @Tags 后台用户
// @Accept json
// @Produce json
// @Param args body model.UserVO true "参数"
// @Success 200 {object} http.baseResponse	"{"status":true}"
// @Router /users/login [post]
func userLogin(c *gin.Context) {
	var vo model.UserVO
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrInvalidParame,
			Message: apicode.MapZH[apicode.ErrInvalidParame],
			Error:   err.Error(),
		})
		return
	}
	token, err := svc.UserLogin(c, vo)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse{
			Status:  false,
			Code:    apicode.ErrCreateToken,
			Message: apicode.MapZH[apicode.ErrCreateToken],
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
