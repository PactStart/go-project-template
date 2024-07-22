package file_store

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"io"
	"orderin-server/pkg/common/log"
)

type ALiYunOSS struct {
	Client     interface{}
	BucketName string
	Domain     string
}

// Setup 装载
// endpoint
func (e *ALiYunOSS) Setup(endpoint, accessKeyID, accessKeySecret, bucketName string, domain string, options ...ClientOption) error {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.ZError(context.Background(), "set up aliyun oss client fail", err)
		return err
	}
	e.Client = client
	e.BucketName = bucketName
	e.Domain = domain
	return nil
}

// UpLoad 文件上传
func (e *ALiYunOSS) UpLoad(yourObjectName string, file interface{}) (*string, error) {
	// 获取存储空间。
	bucket, err := e.Client.(*oss.Client).Bucket(e.BucketName)
	if err != nil {
		log.ZError(context.Background(), "fail to get oss bucket obj", err)
		return nil, err
	}
	// 设置分片大小为100 KB，指定分片上传并发数为3，并开启断点续传上传。
	// 其中<yourObjectName>与objectKey是同一概念，表示断点续传上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// "LocalFile"为filePath，100*1024为partSize。
	if filePath, ok := file.(string); ok {
		err = bucket.UploadFile(yourObjectName, filePath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	} else if reader, ok := file.(io.Reader); ok {
		err = bucket.PutObject(yourObjectName, reader)
	} else {
		err = errors.New("unsupported data format")
		log.ZError(context.Background(), "fail to get oss bucket obj", err)
		return nil, err
	}
	if err != nil {
		log.ZError(context.Background(), "fail to upload file to aliyun oss", err)
		return nil, err
	}
	url := e.Domain + "/" + yourObjectName
	return &url, nil
}

func (e *ALiYunOSS) GetTempToken() (string, error) {
	return "", nil
}
