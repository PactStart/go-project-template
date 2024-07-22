package file_store

import (
	"fmt"
	"orderin-server/pkg/common/config"
)

type OXS struct {
	// Endpoint 访问域名
	Endpoint string
	// AccessKeyID AK
	AccessKeyID string
	// AccessKeySecret AKS
	AccessKeySecret string
	// BucketName 桶名称
	BucketName string
	// 域名
	Domain string
}

var FileStore FileStoreType

func InitFileStore() {
	e := OXS{
		Endpoint:        config.Config.FileStore.Endpoint,
		AccessKeyID:     config.Config.FileStore.AccessKeyID,
		AccessKeySecret: config.Config.FileStore.AccessKeySecret,
		BucketName:      config.Config.FileStore.BucketName,
		Domain:          config.Config.FileStore.Domain,
	}
	FileStore = e.Setup(DriverType(config.Config.FileStore.DriverType))
}

// Setup 配置文件存储driver
func (e *OXS) Setup(driver DriverType, options ...ClientOption) FileStoreType {
	fileStoreType := driver
	var fileStore FileStoreType
	switch fileStoreType {
	case AliYunOSS:
		fileStore = new(ALiYunOSS)
		err := fileStore.Setup(e.Endpoint, e.AccessKeyID, e.AccessKeySecret, e.BucketName, e.Domain)
		if err != nil {
			fmt.Println(err)
		}
		return fileStore
	case HuaweiOBS:
		fileStore = new(HuaWeiOBS)
		err := fileStore.Setup(e.Endpoint, e.AccessKeyID, e.AccessKeySecret, e.BucketName, e.Domain)
		if err != nil {
			fmt.Println(err)
		}
		return fileStore
	case QiNiuKodo:
		fileStore = new(QiNiuKODO)
		err := fileStore.Setup(e.Endpoint, e.AccessKeyID, e.AccessKeySecret, e.BucketName, e.Domain)
		if err != nil {
			fmt.Println(err)
		}
		return fileStore
	}
	return nil
}
