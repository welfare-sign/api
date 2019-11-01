package server

import (
	"log"
	"net/http"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/wsgin"
	"welfare-sign/internal/service"

	_ "welfare-sign/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Svc .
var (
	svc *service.Service
)

// New new server.
func New(s *service.Service) (srv *http.Server) {
	svc = s
	router := wsgin.New()
	initRouter(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv = &http.Server{
		Addr:    viper.GetString(config.KeyHttpAddr),
		Handler: router,
	}
	return
}

func initRouter(e *gin.Engine) {
	e.GET("/health", ping)

	v1 := e.Group("/v1")

	// 商户
	merchants := v1.Group("/merchants")
	{
		merchants.POST("/login", wsgin.ProcessExec(&MerchantLoginRequest{}))
		merchants.POST("", wsgin.ProcessExec(&MerchantAddRequest{}))
		merchants.GET("", wsgin.ProcessExec(&MerchantListRequest{}))
		merchants.GET("/detail", wsgin.ProcessExec(&MerchantDetailRequest{}))
		merchants.GET("/writeoff", wsgin.ProcessExec(&WriteOffRequest{}))
		merchants.POST("/writeoff", wsgin.ProcessExec(&ExecWriteOffRequest{}))
	}

	// 后台用户
	users := v1.Group("/users")
	{
		users.POST("/login", wsgin.ProcessExec(&UserLoginRequest{}))
	}

	// 文件
	files := v1.Group("/files")
	{
		files.POST("/upload", uploadFile)
		files.GET("/download", downloadFile)
	}

	// 客户
	customers := v1.Group("/customers")
	{
		customers.GET("", wsgin.ProcessExec(&CustomerListRequest{}))
		customers.GET("/detail", wsgin.ProcessExec(&CustomerDetailRequest{}))
		customers.GET("/checkin_record", wsgin.ProcessExec(&CheckinRecordRequest{}))
	}

	// 验证码
	smscode := v1.Group("/verify_code")
	{
		smscode.GET("", wsgin.ProcessExec(&CodeRequest{}))
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

// BaseResponse .
type BaseResponse struct {
	Status  bool          `json:"status"`
	Code    wsgin.APICode `json:"code"`
	Message string        `json:"message"`
	Error   string        `json:"error,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
}
