package sms

import (
	"orderin-server/pkg/common/config"
)

type SmsVendorConfig struct {
	// Endpoint 访问域名
	Endpoint string
	// AccessKeyID AK
	AccessKeyID string
	// AccessKeySecret AKS
	AccessKeySecret string
}

var SmsVendorInstance SmsVendor

func InitSmsVendor() {
	e := SmsVendorConfig{
		Endpoint:        config.Config.Sms.Endpoint,
		AccessKeyID:     config.Config.Sms.AccessKeyID,
		AccessKeySecret: config.Config.Sms.AccessKeySecret}
	SmsVendorInstance = e.Setup(VendorType(config.Config.Sms.DriverType))
}

// Setup 配置文件存储driver
func (e *SmsVendorConfig) Setup(vendorType VendorType) SmsVendor {
	var smsVendor SmsVendor
	switch vendorType {
	case AliYunDysms:
		smsVendor = new(AliyunDysms)
		err := smsVendor.Setup(e.Endpoint, e.AccessKeyID, e.AccessKeySecret)
		if err != nil {
			panic(err)
		}
		return smsVendor
	}
	return nil
}
