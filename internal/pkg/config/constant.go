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

	KeySMSDomain   = "sms.domain"
	KeySMSRegion   = "sms.region"
	KeySMSSignName = "sms.sign_name"
	KeySMSAK       = "sms.ak"
	KeySMSAS       = "sms.as"
	KeySMSTemplate = "sms.template"
	KeySMSLength   = "sms.length"
	KeySMSExpire   = "sms.expire"

	KeyWxAppID     = "wx.appid"
	KeyWxAppSecret = "wx.appsecret"

	KeyQRCodeURL = "qrcode.url"
)
