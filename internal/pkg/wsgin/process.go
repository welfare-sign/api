package wsgin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

var processExtract func(c *gin.Context, req Process) (ctx context.Context, status bool)

var response func(c *gin.Context, data interface{})

func init() {
	processExtract = ProcessExtract
	response = Response
}

// BindFunc 对象绑定方法 c.Bind等
type BindFunc func(obj interface{}) error

// Process 请求过程处理接口
type Process interface {
	New() Process                            // 创建自身空副本
	Extract(c *gin.Context) (APICode, error) // 请求参数提取方法
	Exec(ctx context.Context) interface{}    // 请求处理方法
}

// ProcessExec http请求流程处理 c.JSON返回
func ProcessExec(p Process) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := p.New()
		ctx, status := processExtract(c, req)
		if !status {
			return
		}
		data := req.Exec(ctx)
		response(c, data)
	}
}

// Response .
func Response(c *gin.Context, data interface{}) {
	if data != nil {
		c.JSON(http.StatusOK, data)
	}
}

// ProcessExecOutString http请求流程处理 c.String返回 如果出现错误 c.JSON返回
func ProcessExecOutString(p Process) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := p.New()
		ctx, status := processExtract(c, req)
		if !status {
			return
		}
		// 提取成功 处理请求
		response(c, req.Exec(ctx).(string))
	}
}

// ProcessExtract 参数提取流程 失败后直接c.JSON返回 成功后需调用业务处理
func ProcessExtract(c *gin.Context, req Process) (ctx context.Context, status bool) {
	ctx = ExtractContext(c)
	code, err := req.Extract(c) //提取参数
	if err != nil {             //提取错误直接返回
		if code == APICodeDefault {
			code = APICodeInvalidParame
		}
		status = false
		response(c, NewResponse(ctx, code, err))
		return
	}
	status = true
	return
}

// ExtractContext 提取Context
func ExtractContext(c *gin.Context) context.Context {
	return c
}
