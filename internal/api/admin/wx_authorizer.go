package admin

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/officialAccount"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/base/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"orderin-server/internal/components"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/application"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
	"orderin-server/pkg/common/weixin"
	"time"
)

type WxAuthorizer struct {
	api.Api
}

// @Summary 分页查询授权公众号
// @Description 分页查询授权公众号
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Param param body dto.WxAuthorizerPageQueryReq false "授权公众号筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.WxAuthorizer}}
// @Router /wx/authorizer/page_query [post]
// @Security RequireLogin
func (e WxAuthorizer) PageQuery(context *gin.Context) {
	s := services.WxAuthorizer{}
	req := dto.WxAuthorizerPageQueryReq{}
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
	list := make([]models.WxAuthorizer, 0)
	var count int64

	err = s.PageQuery(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}

// @Summary 获取公众号授权连接
// @Description 获取授权链接，微信扫码访问该连接，会跳转到选择公众号授权给第三方平台页面
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=dto.GetAuthorizeUrlResp}
// @Router /wx/authorizer/get_authorize_url [post]
// @Security RequireLogin
func (e WxAuthorizer) GetAuthorizeUrl(context *gin.Context) {
	e.MakeContext(context)
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	rs, err := weixin.OpenPlatformApp.Base.CreatePreAuthorizationCode(context.Request.Context())
	if err != nil {
		e.Error(err)
		return
	}
	fmt2.Dump(rs)
	//拼接授权链接
	componentAppid := config.Config.WxOpenPlatform.AppID
	redirectUrl := context.Request.Host + "/api/v1/wx/authorizer/authorize/redirect"
	authorizeUrl := fmt.Sprintf("https://open.weixin.qq.com/wxaopen/safe/bindcomponent?action=bindcomponent&no_scan=1&component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=1#wechat_redirect",
		componentAppid, rs.PreAuthCode, redirectUrl)

	e.OK(dto.GetAuthorizeUrlResp{
		Url: authorizeUrl,
	})
}

// @Summary 授权回调
// @Description 获取auth_code
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Success 200
// @Router /wx/authorizer/authorize/redirect [get]
func (e WxAuthorizer) AuthorizeRedirect(context *gin.Context) {
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	authCode := context.DefaultQuery("auth_code", "")
	expiresIn := context.DefaultQuery("expires_in", "")

	log.ZInfo(context, "授权回调", "authCode", authCode, "expiresIn", expiresIn)
	res, err := weixin.OpenPlatformApp.Base.HandleAuthorize(context.Request.Context(), authCode)
	if err != nil {
		panic(err)
	}
	fmt2.Dump(res)
	s := services.WxAuthorizer{}
	err = e.MakeContext(context).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	duration := time.Duration(res.AuthorizationInfo.ExpiresIn) * time.Second

	authorizer := &models.WxAuthorizer{
		ComponentAppid:  config.Config.WxOpenPlatform.AppID,
		AuthorizerAppid: res.AuthorizationInfo.AuthorizerAppid,

		AuthorizationCode:          authCode,
		AuthorizationCodeExpiresAt: time.Now().Add(time.Duration(utils.StringToInt(expiresIn)) * time.Second),

		AuthorizerAccessToken: res.AuthorizationInfo.AuthorizerAccessToken,
		AccessTokenExpiresAt:  time.Now().Add(duration),

		AuthorizerRefreshToken: res.AuthorizationInfo.AuthorizerRefreshToken,
	}
	err = s.Save(authorizer)
	if err != nil {
		e.Error(err)
		return
	}
	syncAuthorizerInfo(context, s, authorizer)
	components.RegisterAuthorizer(context, *authorizer)
	context.HTML(http.StatusOK, "authorize_success.html", nil)
}

// @Summary 同步授权方信息
// @Description 同步授权方信息
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Success 200
// @Router /wx/authorizer/sync_info [post]
func (e WxAuthorizer) SyncInfo(context *gin.Context) {
	req := dto.WxAuthorizerSyncInfoReq{}
	s := services.WxAuthorizer{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	wxAuthorizer, err := s.GetById(req.ID)
	if err != nil {
		e.Error(err)
		return
	}
	if wxAuthorizer == nil {
		e.Error(errs.NewCodeError(errs.RecordNotFoundError, "authorizer not found"))
		return
	}
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	authorizer, err := syncAuthorizerInfo(context, s, wxAuthorizer)
	if err != nil {
		e.Error(err)
		return
	}
	e.OK(authorizer)
}

// @Summary 创建客服账号
// @Description 为公众号创建客服账号
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Param param body dto.KfAccountCreateReq false "客服账号和昵称"
// @Success 200 {object} api.Response
// @Security RequireLogin
// @Router /wx/authorizer/create_kf_account [post]
func (e WxAuthorizer) CreateKfAccount(context *gin.Context) {
	s := services.WxAuthorizer{}
	req := dto.KfAccountCreateReq{}
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
	authorizer, err := s.GetByAppID(req.AppID)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	if authorizer == nil {
		log.ZError(e.Context, "authorizer not found", errs.NewCodeError(errs.RecordNotFoundError, "authorizer not found"), "appId", req.AppID)
		e.Error(err)
		return
	}
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	account := application.AppContext.GetComponent(req.AppID)

	data, err := account.(*officialAccount.Application).CustomerService.Create(context, req.Account, req.Nickname)
	if err != nil {
		log.ZError(e.Context, "create customer service account fail", err, "account", req.Account, "nickname", req.Nickname)
		e.Error(err)
		return
	}
	e.OK(data)
}

// @Summary 获取公众号的所有客服账号
// @Description 获取公众号的所有客服账号
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Param param body dto.KfAccountGetAllReq false "公众号appId"
// @Success 200 {object} api.Response
// @Security RequireLogin
// @Router /wx/authorizer/get_all_kf_accounts [post]
func (e WxAuthorizer) GetAllKfAccounts(context *gin.Context) {
	s := services.WxAuthorizer{}
	req := dto.KfAccountGetAllReq{}
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
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	account := application.AppContext.GetComponent(req.AppID)

	data, err := account.(*officialAccount.Application).CustomerService.List(context)
	if err != nil {
		log.ZError(e.Context, "get all customer service account fail", err, "appId", req.AppID)
		e.Error(err)
		return
	}
	e.OK(data)
}

// @Summary 删除客服账号
// @Description 删除客服账号
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Param param body dto.KfAccountDeleteReq false "公众号appId"
// @Success 200 {object} api.Response
// @Security RequireLogin
// @Router /wx/authorizer/delete_kf_account [post]
func (e WxAuthorizer) DeleteKfAccount(context *gin.Context) {
	s := services.WxAuthorizer{}
	req := dto.KfAccountDeleteReq{}
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
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	account := application.AppContext.GetComponent(req.AppID)

	data, err := account.(*officialAccount.Application).CustomerService.Delete(context, req.Account)
	if err != nil {
		log.ZError(e.Context, "delete customer service account fail", err, "appId", req.AppID, "account", req.Account)
		e.Error(err)
		return
	}
	e.OK(data)
}

// @Summary 设置默认的客服账号（使用该账号给用户发送消息）
// @Description 设置默认的客服账号（使用该账号给用户发送消息）
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Param param body dto.KfAccountSetDefaultReq false "公众号appId"
// @Success 200 {object} api.Response
// @Security RequireLogin
// @Router /wx/authorizer/set_default_kf_account [post]
func (e WxAuthorizer) SetDefaultKfAccount(context *gin.Context) {
	s := services.WxAuthorizer{}
	req := dto.KfAccountSetDefaultReq{}
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
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	err = s.UpdateDefaultKfAccountByAppId(req.Account, req.AppID)
	if err != nil {
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 清空菜单
// @Description 清空菜单
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Param param body dto.MenuClearReq false "公众号appId"
// @Success 200 {object} api.Response
// @Security RequireLogin
// @Router /wx/authorizer/delete_menu [post]
func (e WxAuthorizer) DeleteMenu(context *gin.Context) {
	s := services.WxAuthorizer{}
	req := dto.MenuClearReq{}
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
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	account := application.AppContext.GetComponent(req.AppID)

	data, err := account.(*officialAccount.Application).Menu.Delete(context)
	if err != nil {
		log.ZError(e.Context, "clear menu fail", err, "appId", req.AppID)
		e.Error(err)
		return
	}
	e.OK(data)
}

// @Summary 获取菜单
// @Description 获取菜单
// @Tags 授权公众号管理
// @Accept json
// @Produce json
// @Param param body dto.MenuClearReq false "公众号appId"
// @Success 200 {object} api.Response
// @Security RequireLogin
// @Router /wx/authorizer/get_menu [post]
func (e WxAuthorizer) GetMenu(context *gin.Context) {
	s := services.WxAuthorizer{}
	req := dto.MenuClearReq{}
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
	if config.Config.Env.Profiles == "dev" {
		e.Error(errs.NewCodeError(errs.UnSupportedOperation, "当前环境不支持该操作"))
		return
	}
	account := application.AppContext.GetComponent(req.AppID)

	data, err := account.(*officialAccount.Application).Menu.Get(context)
	if err != nil {
		log.ZError(e.Context, "get menu fail", err, "appId", req.AppID)
		e.Error(err)
		return
	}
	e.OK(data)
}

func syncAuthorizerInfo(context *gin.Context, wxAuthorizerService services.WxAuthorizer, wxAuthorizer *models.WxAuthorizer) (*response.ResponseGetAuthorizer, error) {
	authorizer, err := weixin.OpenPlatformApp.Base.GetAuthorizer(context, wxAuthorizer.AuthorizerAppid)
	if err != nil {
		return nil, err
	}
	fmt2.Dump(authorizer)
	wxAuthorizer.NickName = authorizer.AuthorizerInfo.NickName
	wxAuthorizer.HeadImg = authorizer.AuthorizerInfo.HeadImg
	wxAuthorizer.UserName = authorizer.AuthorizerInfo.UserName
	wxAuthorizer.QrcodeUrl = authorizer.AuthorizerInfo.QrcodeURL
	if authorizer.AuthorizerInfo.ServiceTypeInfo != nil {
		wxAuthorizer.ServiceType = utils.IntToString(authorizer.AuthorizerInfo.ServiceTypeInfo.ID)
	}
	if authorizer.AuthorizerInfo.VerifyTypeInfo != nil {
		wxAuthorizer.VerifyType = utils.IntToString(authorizer.AuthorizerInfo.VerifyTypeInfo.ID)
	}
	wxAuthorizer.PrincipalName = authorizer.AuthorizerInfo.PrincipalName
	return authorizer, wxAuthorizerService.UpdateById(wxAuthorizer.ID, wxAuthorizer)
}
