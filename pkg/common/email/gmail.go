package email

type GMail struct {
	Config *EmailConfig
}

func (e *GMail) Setup(config *EmailConfig) error {
	e.Config = config
	return nil
}

func (Q *GMail) Send(request EmailRequest) error {
	return nil
}
