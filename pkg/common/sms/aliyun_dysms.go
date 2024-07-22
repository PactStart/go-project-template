package sms

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AliyunDysms struct {
	Client *dysmsapi.Client
}

func (a *AliyunDysms) GetVendorType() VendorType {
	return AliYunDysms
}

func (a *AliyunDysms) Setup(endpoint, accessKeyId, accessKeySecret string) error {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(endpoint)
	//实例化一个客户端，从 &dysmsapi.Client 类生成对象 client 。
	client, err := dysmsapi.NewClient(config)
	if err == nil {
		a.Client = client
	}
	return err
}

func (a *AliyunDysms) Send(phone string, signName string, templateCode string, templateParams map[string]string) (*SendResp, error) {
	request := &dysmsapi.SendSmsRequest{}
	request.SetPhoneNumbers(phone)
	request.SetSignName(signName)
	request.SetTemplateCode(templateCode)
	if templateParams != nil && len(templateParams) > 0 {
		jsonString, err := json.Marshal(templateParams)
		if err != nil {
			return nil, err
		}
		request.SetTemplateParam(string(jsonString))
	}
	resp := SendResp{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		response, err := a.Client.SendSms(request)
		if err != nil {
			return err
		}
		resp.StatusCode = response.StatusCode
		if response.Body != nil {
			resp.Code = response.Body.Code
			resp.Msg = response.Body.Message
			resp.RequestId = response.Body.RequestId
			resp.BizId = response.Body.BizId
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		return nil, error
	}
	return &resp, nil
}
