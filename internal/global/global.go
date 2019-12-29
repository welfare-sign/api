package global

// 全局常量
const (
	InactiveStatus = "U" // 未激活
	ActiveStatus   = "A" // 已激活，可以是状态
	DeleteStatus   = "X" // 记录被删除、无效
	Readed         = "Y" // 是否读取了补签消息
	UnRead         = "N" // 没有阅读补签消息
	ScopeOpen      = "O" // 参与方式：无限制
	ScopeReceive   = "R" // 参与方式：必须在福利签领取过签到奖励
)
