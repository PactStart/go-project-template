package file_store

import (
	"context"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"orderin-server/pkg/common/log"
)

type HuaWeiOBS struct {
	Client     interface{}
	BucketName string
	Domain     string
}

func (e *HuaWeiOBS) Setup(endpoint, accessKeyID, accessKeySecret, bucketName string, domain string, options ...ClientOption) error {
	// 创建ObsClient结构体
	client, err := obs.New(accessKeyID, accessKeySecret, endpoint)
	if err != nil {
		log.ZError(context.Background(), "set up huawei obs client fail", err)
		return err
	}
	e.Client = client
	e.BucketName = bucketName
	e.Domain = domain
	return nil
}

// UpLoad 文件上传
// yourObjectName 文件路径名称，与objectKey是同一概念，表示断点续传上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg
func (e *HuaWeiOBS) UpLoad(yourObjectName string, localFile interface{}) (*string, error) {
	// 获取存储空间。
	input := &obs.PutFileInput{}
	input.Bucket = e.BucketName
	input.Key = yourObjectName
	input.SourceFile = localFile.(string)
	output, err := e.Client.(*obs.ObsClient).PutFile(input)

	if err == nil {
		log.ZError(context.Background(), "fail to upload file to huawei obs", err, "requestId", output.RequestId, "Etag", output.ETag, "StorageClass", output.StorageClass)
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			log.ZError(context.Background(), "fail to upload file to huawei obs", err, "Code", obsError.Code, "Message", obsError.Message)
		} else {
			log.ZError(context.Background(), "fail to upload file to huawei obs", err)
		}
	}
	url := e.Domain + "/" + yourObjectName
	return &url, nil
}

func (e *HuaWeiOBS) GetTempToken() (string, error) {
	return "", nil
}
