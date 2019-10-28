package apicode

// Code api code type
type Code string

// MapZH .
var MapZH map[Code]string

// api code define
const (
	Success          Code = "SUCCESS"
	ErrInvalidParame Code = "ERR_INVALID_PARAME"
	ErrModelCreate   Code = "ERR_MODEL_CREATE"
	ErrCreateToken   Code = "ERR_CREATE_TOKEN"
	ErrUploadFile    Code = "ERR_UPLOAD_FILE"
	ErrDownloadFile  Code = "ERR_DOWNLOAD_FILE"
	ErrGetListData   Code = "ERR_GET_LIST_DATA"
	ErrMobile        Code = "ERR_MOBILE"
	ErrSendSMS       Code = "ERR_SEND_SMS"
	ErrLogin         Code = "ERR_LOGIN"
	ErrDetail        Code = "ERR_DETAIL"
	ErrWriteOff      Code = "ERR_WRITE_OFF"
	ErrExecWriteOff  Code = "ERR_EXEC_WRITE_OFF"
)

func init() {
	MapZH = make(map[Code]string)
	MapZH[Success] = "请求成功"
	MapZH[ErrInvalidParame] = "请求参数不正确"
	MapZH[ErrModelCreate] = "创建失败"
	MapZH[ErrCreateToken] = "token创建失败"
	MapZH[ErrUploadFile] = "上传文件失败"
	MapZH[ErrDownloadFile] = "下载文件失败"
	MapZH[ErrGetListData] = "获取列表数据失败"
	MapZH[ErrMobile] = "手机号不正确"
	MapZH[ErrSendSMS] = "验证码发送失败"
	MapZH[ErrLogin] = "登录失败"
	MapZH[ErrDetail] = "获取详情失败"
	MapZH[ErrWriteOff] = "获取核销信息失败"
	MapZH[ErrExecWriteOff] = "核心失败"
}
