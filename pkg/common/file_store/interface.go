package file_store

// DriverType 驱动类型
type DriverType string

const (
	// HuaweiOBS 华为云OBS
	HuaweiOBS DriverType = "HuaweiOBS"
	// AliYunOSS 阿里云OSS
	AliYunOSS DriverType = "AliYunOSS"
	// QiNiuKodo 七牛云kodo
	QiNiuKodo DriverType = "QiNiuKodo"
)

type ClientOption map[string]interface{}

// FileStoreType OXS
type FileStoreType interface {
	// Setup 装载 endpoint sss
	Setup(endpoint, accessKeyID, accessKeySecret, bucketName string, domain string, options ...ClientOption) error
	// UpLoad 上传
	UpLoad(yourObjectName string, file interface{}) (*string, error)
	// GetTempToken 获取临时Token
	GetTempToken() (string, error)
}
