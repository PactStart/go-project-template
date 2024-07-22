package email

import (
	"orderin-server/pkg/common/config"
)

type EmailConfig struct {
	From     string
	Account  string
	Password string
	SmtpHost string
	SmtpPort int
}

var EmailVendorInstance EmailVendor

func InitEmailVendor() {
	e := EmailConfig{
		From:     config.Config.Mail.From,
		Account:  config.Config.Mail.Account,
		Password: config.Config.Mail.Password,
		SmtpHost: config.Config.Mail.SmtpHost,
		SmtpPort: config.Config.Mail.SmtpPort,
	}
	EmailVendorInstance = e.Setup(VendorType(config.Config.Mail.DriverType))
}

// Setup 配置文件存储driver
func (e *EmailConfig) Setup(vendorType VendorType) EmailVendor {
	var result EmailVendor
	switch vendorType {
	case QQ:
		result = new(QQEmail)
		err := result.Setup(e)
		if err != nil {
			panic(err)
		}
		return result
	case Gmail:
		result = new(QQEmail)
		err := result.Setup(e)
		if err != nil {
			panic(err)
		}
		return result
	}
	return nil
}
