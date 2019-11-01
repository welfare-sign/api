package wsgin

import (
	"context"
)

// BaseResponse 基础返回对象
type BaseResponse struct {
	Status  bool    `json:"status"`          // 状态
	Code    APICode `json:"code"`            // 业务状态码
	Message string  `json:"message"`         // 提示消息
	Error   string  `json:"error,omitempty"` // Error信息
}

// NewResponse 根据业务状态码和err信息创建新的结构返回
func NewResponse(ctx context.Context, code APICode, err error) BaseResponse {
	if err == nil && (code == APICodeSuccess || code == APICodeDefault) {
		return NewSuccessResponse(ctx)
	}
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return BaseResponse{
		Status:  code == APICodeSuccess,
		Code:    code,
		Message: APICodeMapZH[code],
		Error:   errMsg,
	}
}

// NewSuccessResponse success response
func NewSuccessResponse(ctx context.Context) BaseResponse {
	return BaseResponse{
		Status:  true,
		Code:    APICodeSuccess,
		Message: APICodeMapZH[APICodeSuccess],
	}
}
