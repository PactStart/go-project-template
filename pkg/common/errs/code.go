package errs

const (
	// 通用错误码.
	ServerInternalError = 500 // 服务器内部错误
	ArgsError           = 400 // 输入参数错误
	UnauthorizedError   = 401 //未授权
	NoPermissionError   = 403 // 权限不足
	RecordNotFoundError = 404 // 记录不存在
	DuplicateKeyError   = 409 //重复的key

	// token错误码.
	TokenExpiredError     = 1501
	TokenInvalidError     = 1502
	TokenMalformedError   = 1503
	TokenNotValidYetError = 1504
	TokenUnknownError     = 1505
	TokenKickedError      = 1506
	TokenNotExistError    = 1507

	AccountOrPasswordError = 1001
	AccountNotExistError   = 1002
	SmsCodeError           = 1003
	EmailCodeError         = 1004
	MFACodeError           = 1005
	CheckIdentityError     = 1006
	UnSupportedOperation   = 1007
)
