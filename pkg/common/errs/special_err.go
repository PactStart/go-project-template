package errs

var handlers []func(err error) CodeError

func AddErrHandler(h func(err error) CodeError) {
	if h == nil {
		panic("nil handler")
	}
	handlers = append(handlers, h)
}

func AddReplace(target error, codeErr CodeError) {
	AddErrHandler(func(err error) CodeError {
		if err == target {
			return codeErr
		}
		return nil
	})
}

func ErrCode(err error) CodeError {
	if codeErr, ok := err.(CodeError); ok {
		return codeErr
	}
	for i := 0; i < len(handlers); i++ {
		if codeErr := handlers[i](err); codeErr != nil {
			return codeErr
		}
	}
	return nil
}
