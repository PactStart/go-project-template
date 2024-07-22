package errs

var (
	ErrArgs                 = NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission         = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrInternalServer       = NewCodeError(ServerInternalError, "ServerInternalError")
	ErrRecordNotFound       = NewCodeError(RecordNotFoundError, "RecordNotFoundError")
	ErrDuplicateKey         = NewCodeError(DuplicateKeyError, "DuplicateKeyError")
	ErrMFACode              = NewCodeError(MFACodeError, "MFACodeError")
	ErrCheckIdentity        = NewCodeError(CheckIdentityError, "CheckIdentityError")
	ErrUnSupportedOperation = NewCodeError(UnSupportedOperation, "UnSupportedOperation")

	ErrTokenExpired     = NewCodeError(TokenExpiredError, "TokenExpiredError")         //token已过期
	ErrTokenInvalid     = NewCodeError(TokenInvalidError, "TokenInvalidError")         //token不合法
	ErrTokenMalformed   = NewCodeError(TokenMalformedError, "TokenMalformedError")     // 格式错误
	ErrTokenNotValidYet = NewCodeError(TokenNotValidYetError, "TokenNotValidYetError") // 还未生效
	ErrTokenUnknown     = NewCodeError(TokenUnknownError, "TokenUnknownError")         // 未知错误
	ErrTokenKicked      = NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist    = NewCodeError(TokenNotExistError, "TokenNotExistError") // 在redis中不存在

)
