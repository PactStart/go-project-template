package sms

// VendorType 供应商类型
type VendorType string

const (
	// AliYunDysms 阿里云
	AliYunDysms VendorType = "AliyunDysms"
)

type SendResp struct {
	StatusCode *int32
	Code       *string
	Msg        *string
	RequestId  *string
	BizId      *string
}

type SmsVendor interface {
	GetVendorType() VendorType
	// Setup
	Setup(endpoint, accessKeyID, accessKeySecret string) error
	// 发送
	Send(phone string, signName string, templateCode string, templateParams map[string]string) (*SendResp, error)
}
