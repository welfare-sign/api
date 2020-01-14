package config

// config constant
const (
	KeyMysqlDSN = "mysql.dsn"

	KeyHTTPAddr = "http.addr"

	KeyJWTSign   = "jwt.sign"
	KeyJWTIssuer = "jwt.issuer"
	KeyJWTExpire = "jwt.expire"

	KeyRedisDB   = "redis.db"
	KeyRedisHost = "redis.host"
	KeyRedisPWD  = "redis.pwd"

	KeyLogEnable  = "log.enable"
	KeyLogProject = "log.project"
	KeyLogPath    = "log.path"

	KeySMSDomain      = "sms.domain"
	KeySMSRegion      = "sms.region"
	KeySMSSignName    = "sms.sign_name"
	KeySMSAK          = "sms.ak"
	KeySMSAS          = "sms.as"
	KeySMSTemplate    = "sms.template"
	KeySMSLength      = "sms.length"
	KeySMSExpire      = "sms.expire"
	KeySMSEnable      = "sms.enable"
	KeySMSSpecialCode = "sms.special_code"

	KeyWxAppID     = "wx.appid"
	KeyWxAppSecret = "wx.appsecret"

	KeyQRCodeURL = "qrcode.url"

	KeyWXPayMchID     = "wx.pay_mch_id"
	KeyWXPayAPI       = "wx.pay_api_key"
	KeyWXPayNotifyURL = "wx.pay_notify_url"
	KeyWXPayAmount    = "wx.pay_amount"

	KeyTaskEnable                          = "task.enable"
	KeyTaskCheckinExpiredTime              = "task.checkin_expired_time"
	KeyTaskCheckinExpiredTimeStartInterval = "task.checkin_expired_time_start_interval"
)
