package http

import (
	"log"
	"net/http"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/service"

	_ "welfare-sign/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	svc *service.Service
)

// New new server.
func New(s *service.Service) (srv *http.Server) {
	svc = s
	gin.DisableConsoleColor()
	router := gin.Default()
	initRouter(router)

	srv = &http.Server{
		Addr:    viper.GetString(config.KeyHttpAddr),
		Handler: router,
	}
	return
}

func initRouter(e *gin.Engine) {
	e.GET("/health", ping)
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := e.Group("/v1")
	{
		// 商户
		v1.POST("/merchants", addMerchant)
		v1.GET("/merchants", getMerchantList)
		v1.POST("/merchants/login", merchantLogin)
		v1.GET("/merchants/detail", merchantDetail)
		v1.GET("/merchants/writeoff", getWriteOff)
		v1.POST("/merchants/writeoff", writeOff)
		// 后台用户
		v1.POST("/users/login", userLogin)
		// 文件
		v1.POST("/files/upload", uploadFile)
		v1.GET("/files/download", downloadFile)
		// 客户
		v1.GET("/customers", getCustomerList)
		v1.GET("/customers/detail", customerDetail)
		// 验证码
		v1.GET("/verify_code", getVerifyCode)
	}
}

func ping(c *gin.Context) {
	if err := svc.Ping(c); err != nil {
		log.Printf("ping error(%v)\n", err.Error())
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	c.JSON(http.StatusOK, "ping success")
}

// baseResponse 基础响应
type baseResponse struct {
	Status  bool         `json:"status"`
	Code    apicode.Code `json:"code"`
	Message string       `json:"message"`
	Error   string       `json:"error"`
	Data    interface{}  `json:"data"`
}
