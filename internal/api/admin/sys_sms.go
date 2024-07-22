package admin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/cache"
	"orderin-server/pkg/common/constant"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/sms"
	"orderin-server/pkg/common/utils"
	"strings"
	"time"
)

type SysSms struct {
	api.Api
	VerifyCodeCache cache.VerifyCodeCache
}

// @Summary 发送短信
// @Description 发送短信
// @Tags 短信管理
// @Accept json
// @Produce json
// @Param param body dto.SysSmsSendReq false "手机号和发送场景"
// @Success 200 {object} api.Response
// @Router /sys/sms/send [post]
func (e SysSms) Send(context *gin.Context) {
	s := services.SysSms{}
	req := dto.SysSmsSendReq{}
	err := e.MakeContext(context).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	smsLog := models.SysSmsLog{}
	smsLog.Phone = req.Phone
	switch req.Scene {
	case constant.LoginByPhone:
		userService := services.SysUser{}
		userService.Orm = e.Orm
		userService.Context = e.Context

		user, err := userService.GetByPhone(req.Phone)
		if err != nil {
			e.Error(errs.NewCodeError(errs.ServerInternalError, "发送短信验证码失败").WithDetail(err.Error()))
			return
		}
		if user == nil {
			e.Error(errs.NewCodeError(errs.AccountNotExistError, "手机号不存在"))
			return
		}
		err = e.sendValidateSms(req.Phone, req.Scene, &smsLog)
	case constant.BindPhone:
		userService := services.SysUser{}
		userService.Orm = e.Orm
		userService.Context = e.Context

		exists, err := userService.IsPhoneExist(req.Phone)
		if err != nil {
			log.ZError(e.Context, "校验手机号是否存在失败", err)
			e.Error(errs.ErrInternalServer)
			return
		}
		if exists {
			e.Error(errs.NewCodeError(errs.DuplicateKeyError, "该手机号已被其他用户绑定"))
			return
		}
		err = e.sendValidateSms(req.Phone, req.Scene, &smsLog)
	default:
		err = errs.NewCodeError(errs.ArgsError, "不支持的短信场景!")
	}
	if err != nil {
		e.Error(err)
		return
	}
	err = s.Insert(smsLog)
	if err != nil {
		log.ZError(e.Context, "保存短信发送记录失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 发送短信(登录后)
// @Description 发送短信(登录后)
// @Tags 短信管理
// @Accept json
// @Produce json
// @Param param body dto.SysSmsSendToMyselfReq false "发送场景"
// @Success 200 {object} api.Response
// @Router /sys/sms/send_to_myself [post]
// @Security RequireLogin
func (e SysSms) SendToMyself(context *gin.Context) {
	s := services.SysSms{}
	req := dto.SysSmsSendToMyselfReq{}
	err := e.MakeContext(context).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	smsLog := models.SysSmsLog{}
	switch req.Scene {
	case constant.CheckIdentity:
		currentUserId := mcontext.GetOpUserID(context)
		userService := services.SysUser{}
		userService.Context = context
		userService.Orm = e.Orm
		user, err := userService.GetById(currentUserId)
		if err != nil {
			log.ZError(e.Context, "", err)
			e.Error(err)
			return
		}
		smsLog.Phone = user.Phone
		err = e.sendValidateSms(user.Phone, req.Scene, &smsLog)
	default:
		err = errs.NewCodeError(errs.ArgsError, "unsupported sms scene!")
	}
	if err != nil {
		e.Error(err)
		return
	}
	err = s.Insert(smsLog)
	if err != nil {
		log.ZError(e.Context, "保存短信发送记录失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

func (e SysSms) sendValidateSms(phone string, scene string, smsLog *models.SysSmsLog) error {
	code := utils.GenerateDigitalString(6)
	variables := map[string]string{
		"code": code,
	}
	template := "您的验证码为：${code}，请勿泄露于他人"
	content := replaceTemplateVariable(template, variables)

	smsLog.Vendor = string(sms.SmsVendorInstance.GetVendorType())
	smsLog.TemplateCode = "SMS_296350564"
	smsLog.Content = content
	smsLog.Scene = scene

	//缓存验证码
	err := e.VerifyCodeCache.StoreCode(phone, scene, code, 10*time.Minute)
	if err != nil {
		return errs.NewCodeError(errs.ServerInternalError, "短信发送失败").Wrap()
	}
	//发送短信
	resp, err := sms.SmsVendorInstance.Send(phone, "小西西家政", "SMS_296350564", variables)

	json, _ := json.Marshal(resp)
	log.ZInfo(e.Context, "短信发送结果", "resp", string(json))
	if err != nil || resp == nil || resp.StatusCode == nil || *resp.StatusCode != 200 || resp.Code == nil || *resp.Code != "OK" {
		smsLog.Status = "1"
		smsLog.Remark = err.Error()
		return errs.NewCodeError(errs.ServerInternalError, "短信发送失败").WithDetail(err.Error())
	} else {
		smsLog.Status = "2"
		smsLog.Remark = "SUCCESS"
	}
	return nil
}

func replaceTemplateVariable(template string, variables map[string]string) string {
	for key, value := range variables {
		variable := "${" + key + "}"
		template = strings.ReplaceAll(template, variable, value)
	}
	return template
}

// @Summary 分页查询短信记录
// @Description 分页查询短信记录
// @Tags 短信管理
// @Accept json
// @Produce json
// @Param param body dto.SysSmsLogPageQueryReq false "短信记录筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysSmsLog}}
// @Router /sys/sms/page_query [post]
// @Security RequireLogin
func (e SysSms) PageQuery(context *gin.Context) {
	s := services.SysSms{}
	req := dto.SysSmsLogPageQueryReq{}
	err := e.MakeContext(context).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	list := make([]models.SysSmsLog, 0)
	var count int64

	err = s.PageQuery(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}
