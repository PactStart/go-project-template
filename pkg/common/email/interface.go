package email

// VendorType 供应商类型
type VendorType string

const (
	QQ    VendorType = "QQEmail"
	Gmail VendorType = "GMail"
)

type EmailRequest struct {
	To               []string
	Subject          string
	Content          string
	TemplateFileName string
	Params           map[string]string
}

type EmailVendor interface {
	// Setup
	Setup(config *EmailConfig) error
	// 发送
	Send(request EmailRequest) error
}
