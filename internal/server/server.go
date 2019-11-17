package server

import (
	"log"
	"net/http"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/wsgin"
	"welfare-sign/internal/service"

	_ "welfare-sign/docs" // swagger docs

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
		Addr:    viper.GetString(config.KeyHTTPAddr),
		Handler: router,
	}
	return
}

func initRouter(e *gin.Engine) {
	e.GET("/health", ping)
	// 检查微信 服务器签名
	e.StaticFile("/favicon.ico", "./public/favicon.ico")
	e.StaticFile("/MP_verify_6IOxVtGiF56arfjR.txt", "./public/MP_verify_6IOxVtGiF56arfjR.txt")

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
		customers.GET("/checkin_record", wsgin.ProcessExec(&CheckinRecordRequest{}))      // 获取签到记录
		customers.POST("/checkin_record", wsgin.ProcessExec(&ExecCheckinRecordRequest{})) // 签到
		customers.GET("/qrcode", wsgin.ProcessExec(&QRCodeRequest{}))                     // 获取二维码
		customers.POST("/login", wsgin.ProcessExec(&CustomerLoginRequest{}))
		customers.GET("/near_merchant", wsgin.ProcessExec(&NearMerchantRequest{}))                      // 获取附近商家
		customers.GET("/issue_records", wsgin.ProcessExec(&IssueRecordRequest{}))                       // 查看我的福利
		customers.POST("/issue_records", wsgin.ProcessExec(&ExecIssueRecordRequest{}))                  // 领取福利
		customers.POST("/checkin_record/refresh", wsgin.ProcessExec(&RefreshCheckinRecordRequest{}))    // 用户重新签到
		customers.POST("/checkin_record/help", wsgin.ProcessExec(&HelpCheckinRequest{}))                // 帮助他人签到
		customers.GET("/issue_records/is_supplement", wsgin.ProcessExec(&IsSupplementCheckinRequest{})) // 是否是补签
	}

	// 验证码
	smscode := v1.Group("/verify_code")
	{
		smscode.GET("", wsgin.ProcessExec(&CodeRequest{}))
	}

	wx := v1.Group("/wx")
	{
		wx.GET("/config", wsgin.ProcessExec(&WXConfigRequest{}))
		wx.POST("/pay", wsgin.ProcessExec(&WXPayRequest{}))
		wx.POST("/pay/notify", wsgin.ProcessExec(&WxpayCallbackRequest{}))
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
