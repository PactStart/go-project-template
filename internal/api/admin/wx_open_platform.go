package admin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	wxfmt "github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/officialAccount"
	openplatform "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/server/callbacks"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/application"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/weixin"
)

type WxOpenPlatformServer struct {
	api.Api
}

func (e WxOpenPlatformServer) Callback(context *gin.Context) {

	requestXML, _ := io.ReadAll(context.Request.Body)
	context.Request.Body = io.NopCloser(bytes.NewBuffer(requestXML))
	println(string(requestXML))

	query, _ := json.Marshal(context.Request.URL.Query())
	log.ZInfo(context, "receive a callback", "query", string(query), "body", string(requestXML))

	var err error

	rs, err := weixin.OpenPlatformApp.Server.Notify(context.Request, func(event *openplatform.Callback, decrypted []byte, infoType string) (result interface{}) {

		result = kernel.SUCCESS_EMPTY_RESPONSE

		switch infoType {
		case openplatform.EVENT_COMPONENT_VERIFY_TICKET:
			msg := &openplatform.EventVerifyTicket{}
			err = xml.Unmarshal(decrypted, msg)
			//fmt.Dump(event)
			if err != nil {
				return err
			}
			// set ticket in redis
			err = weixin.OpenPlatformApp.VerifyTicket.SetTicket(msg.ComponentVerifyTicket)
			if err != nil {
				return err
			}

			wxfmt.Dump(msg)
		case openplatform.EVENT_AUTHORIZED:

		}
		return result
	})

	if err != nil {
		panic(err)
	}

	err = rs.Write(context.Writer)
	if err != nil {
		panic(err)
	}

}

func (e WxOpenPlatformServer) CallbackWithApp(context *gin.Context) {
	requestXML, _ := io.ReadAll(context.Request.Body)
	context.Request.Body = io.NopCloser(bytes.NewBuffer(requestXML))
	println(string(requestXML))

	query, _ := json.Marshal(context.Request.URL.Query())
	appId := context.Param("appID")
	log.ZInfo(context, "receive a app callback", "appId", appId, "query", string(query), "body", string(requestXML))

	account := application.AppContext.GetComponent(appId)
	if account == nil {
		log.ZError(context, "未找到对应的授权公众号", errors.New("wechat authorizer not found"), "appId", appId)
		return
	}
	rs, err := account.(*officialAccount.Application).Server.Notify(context.Request, func(event contract.EventInterface) interface{} {
		wxfmt.Dump("event", event)
		return kernel.SUCCESS_EMPTY_RESPONSE
	})
	if err != nil {
		log.ZError(context, "处理微信消息推送发生错误", err)
		return
	}
	err = helper.HttpResponseSend(rs, context.Writer)
	if err != nil {
		log.ZError(context, "处理微信消息推送发生错误", err)
	}
}
