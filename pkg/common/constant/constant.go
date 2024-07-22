package constant

const (
	OpUserPlatform = "platform"
	OpUserID       = "opUserID"
	SuperAdmin     = "superAdmin"
	Token          = "token"
	RemoteAddr     = "remoteAddr"
	RequestId      = "X-Request-Id"

	// token.
	NormalToken  = 0
	InValidToken = 1
	KickedToken  = 2
	ExpiredToken = 3

	// flag parse.
	FlagPort           = "port"
	FlagPrometheusPort = "prometheus_port"
	FlagConf           = "config_folder_path"

	// scenes of sending verify code
	LoginByPhone  = "LoginByPhone"
	BindPhone     = "BindPhone"
	BindEmail     = "BindEmail"
	CheckIdentity = "CheckIdentity"

	//these operations need check identity
	ChangeEmail     = "change_email"
	ChangePhone     = "change_phone"
	BindMFADevice   = "bind_mfa_device"
	UnBindMFADevice = "unbind_mfa_device"
	ChangePassword  = "change_password"
	SwitchEnableMFA = "switch_enable_mfa"
	BindWechat      = "bind_wechat"
	UnbindWechat    = "unbind_wechat"
)

const LocalHost = "0.0.0.0"
