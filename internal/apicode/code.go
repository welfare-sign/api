package apicode

import "welfare-sign/internal/pkg/wsgin"

// api code define
const (
	ErrModelCreate       wsgin.APICode = "ERR_MODEL_CREATE"
	ErrCreateToken       wsgin.APICode = "ERR_CREATE_TOKEN"
	ErrUploadFile        wsgin.APICode = "ERR_UPLOAD_FILE"
	ErrDownloadFile      wsgin.APICode = "ERR_DOWNLOAD_FILE"
	ErrGetListData       wsgin.APICode = "ERR_GET_LIST_DATA"
	ErrMobile            wsgin.APICode = "ERR_MOBILE"
	ErrSendSMS           wsgin.APICode = "ERR_SEND_SMS"
	ErrLogin             wsgin.APICode = "ERR_LOGIN"
	ErrDetail            wsgin.APICode = "ERR_DETAIL"
	ErrWriteOff          wsgin.APICode = "ERR_WRITE_OFF"
	ErrExecWriteOff      wsgin.APICode = "ERR_EXEC_WRITE_OFF"
	ErrGetCheckinRecord  wsgin.APICode = "ERR_GET_CHECKIN_RECORD"
	ErrMobileExists      wsgin.APICode = "ERR_MOBILE_EXISTS"
	ErrGetNearMerchant   wsgin.APICode = "ERR_NEAR_MERCHANT"
	ErrExecCheckinRecord wsgin.APICode = "ERR_CHECKIN_RECORD"
	ErrHasCheckin        wsgin.APICode = "ERR_CHECKIN"
)

func init() {
	wsgin.APICodeMapZH[ErrModelCreate] = "创建失败"
	wsgin.APICodeMapZH[ErrCreateToken] = "token创建失败"
	wsgin.APICodeMapZH[ErrUploadFile] = "上传文件失败"
	wsgin.APICodeMapZH[ErrDownloadFile] = "下载文件失败"
	wsgin.APICodeMapZH[ErrGetListData] = "获取列表数据失败"
	wsgin.APICodeMapZH[ErrMobile] = "手机号不正确"
	wsgin.APICodeMapZH[ErrSendSMS] = "验证码发送失败"
	wsgin.APICodeMapZH[ErrLogin] = "登录失败"
	wsgin.APICodeMapZH[ErrDetail] = "获取详情失败"
	wsgin.APICodeMapZH[ErrWriteOff] = "获取核销信息失败"
	wsgin.APICodeMapZH[ErrExecWriteOff] = "核销失败"
	wsgin.APICodeMapZH[ErrGetCheckinRecord] = "获取签到记录失败"
	wsgin.APICodeMapZH[ErrMobileExists] = "手机号已存在"
	wsgin.APICodeMapZH[ErrGetNearMerchant] = "获取附近商家失败"
	wsgin.APICodeMapZH[ErrExecCheckinRecord] = "签到失败"
	wsgin.APICodeMapZH[ErrHasCheckin] = "今天已经签到了,请明天再来"
}
