package wsgin

// APICode api code
type APICode string

// APICodeMapZH api code map zh
var APICodeMapZH map[APICode]string

func init() {
	APICodeMapZH = make(map[APICode]string)
	APICodeMapZH[APICodeDefault] = "未处理的错误"
	APICodeMapZH[APICodeSuccess] = "成功"
	APICodeMapZH[APICodeServerError] = "服务器处理异常"
	APICodeMapZH[APICodeInvalidParame] = "参数验证未通过"
	APICodeMapZH[APICodeNoPermission] = "无权访问"
	APICodeMapZH[APICodeDBError] = "数据库处理异常"
}

// api code define
const (
	APICodeDefault       APICode = ""
	APICodeSuccess       APICode = "SUCCESS"
	APICodeServerError   APICode = "SERVER_ERROR"
	APICodeInvalidParame APICode = "INVALID_PARAME"
	APICodeNoPermission  APICode = "NO_PERMISSION"
	APICodeDBError       APICode = "DB_ERROR"
)
