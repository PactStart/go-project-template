package admin

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pquerna/otp/totp"
	"net/url"
	"orderin-server/internal/components"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/application"
	"orderin-server/pkg/common/auth"
	"orderin-server/pkg/common/cache"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/constant"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/customtypes"
	"orderin-server/pkg/common/email"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
	"orderin-server/pkg/common/weixin"
	"time"
)

type SysUser struct {
	api.Api
	AuthService     auth.AuthService
	VerifyCodeCache cache.VerifyCodeCache
}

// @Summary 分页查询用户
// @Description 分页查询用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserPageQueryReq false "用户筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysUser}}
// @Router /sys/user/page_query [post]
// @Security RequireLogin
func (e SysUser) PageQuery(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserPageQueryReq{}
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
	list := make([]models.SysUser, 0)
	var count int64

	err = s.PageQuery(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}

// @Summary 根据id获取用户信息
// @Description 根据id获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserGetReq true "用户id"
// @Success 200 {object} api.Response{data=models.SysUser}
// @Router /sys/user/get_by_id [post]
// @Security RequireLogin
func (e SysUser) GetById(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserGetReq{}
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
	user, err := s.GetById(req.ID)
	if err != nil {
		log.ZError(e.Context, "根据id获取用户失败", err)
		e.Error(err)
		return
	}
	e.OK(user)
}

// @Summary 添加用户
// @Description 添加用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserInsertReq true "用户信息"
// @Success 200 {object} api.Response{data=dto.SysUserAddResp}
// @Router /sys/user/add [post]
// @Security RequireLogin
func (e SysUser) Add(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserInsertReq{}
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
	user := models.SysUser{}
	utils.CopyStructFields(&user, req)
	user.Status = 1
	currentUserId := mcontext.GetOpUserID(context)
	user.CreatedBy = &currentUserId
	initPassword := utils.GenerateRandomString(8)
	user.PasswordSalt = utils.Md5(initPassword)

	err = s.Insert(&user)
	if err != nil {
		log.ZError(e.Context, "添加用户失败", err)
		e.Error(err)
		return
	}
	resp := dto.SysUserAddResp{ID: user.ID, InitPassword: initPassword}
	e.OK(resp)
}

// @Summary 修改用户状态
// @Description 修改用户状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserUpdateStatusReq true "用户id和状态"
// @Success 200 {object} api.Response
// @Router /sys/user/update_status [post]
// @Security RequireLogin
func (e SysUser) UpdateStatus(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserUpdateStatusReq{}
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
	updateUser := models.SysUser{}
	utils.CopyStructFields(&updateUser, req)
	currentUserId := mcontext.GetOpUserID(context)
	updateUser.UpdatedBy = &currentUserId

	err = s.UpdateById(req.ID, updateUser)
	if err != nil {
		log.ZError(e.Context, "更新状态失败", err)
		e.AddError(err)
		return
	}
	e.OK(nil)
}

// @Summary 账号密码登录
// @Description 账号密码登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserLoginByPasswordReq true "账号密码"
// @Success 200 {object} api.Response{data=dto.SysUserLoginResp}
// @Router /sys/user/login_by_password [post]
func (e SysUser) LoginByPassword(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserLoginByPasswordReq{}
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
	user, err := s.GetByUserName(req.Username)
	if err != nil {
		log.ZError(e.Context, "根据用户名查询用户失败", err)
		e.Error(errs.NewCodeError(errs.AccountOrPasswordError, "账号或密码错误"))
		return
	}
	password := utils.Md5(req.Password)
	if user.PasswordSalt != password {
		e.Error(errs.NewCodeError(errs.AccountOrPasswordError, "账号或密码错误"))
		return
	}
	resp, err := e.Login(context, s, user)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	e.OK(resp)
}

func (e SysUser) Login(context *gin.Context, userService services.SysUser, user *models.SysUser) (*dto.SysUserLoginResp, error) {
	now := customtypes.Time(time.Now())
	updateUser := models.SysUser{
		LastLoginIp:   mcontext.GetRemoteAddr(context),
		LastLoginTime: &now,
	}
	userService.UpdateById(user.ID, updateUser)
	token, err := e.AuthService.CreateToken(context, user.ID, user.SuperAdmin, constant.AdminPlatformID)
	loginResp := dto.SysUserLoginResp{
		Token:             token,
		ExpireTimeSeconds: config.Config.TokenPolicy.Expire * 24 * 60 * 60,
	}
	permissions, err := userService.GetPermissionsByUserId(user.ID)
	if err != nil {
		log.ZError(e.Context, "获取用户权限失败", err)
		return nil, err
	}
	application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).SetMyPermissions(user.ID, permissions)
	return &loginResp, nil
}

// @Summary 手机号登录
// @Description 手机号登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserLoginByPhoneReq true "手机号和验证码"
// @Success 200 {object} api.Response{data=dto.SysUserLoginResp}
// @Router /sys/user/login_by_phone [post]
func (e SysUser) LoginByPhone(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserLoginByPhoneReq{}
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
	user, err := s.GetByPhone(req.Phone)
	if err != nil {
		log.ZError(e.Context, "根据手机号查询用户失败", err)
		e.Error(errs.NewCodeError(errs.ServerInternalError, "服务器错误").WithDetail(err.Error()))
		return
	}
	if user == nil {
		e.Error(errs.NewCodeError(errs.AccountNotExistError, "账号不存在"))
		return
	}
	pass, err := e.VerifyCodeCache.CheckCode(e.Context, req.Phone, constant.LoginByPhone, req.SmsCode)
	if err != nil {
		e.Error(err)
		return
	}
	if !pass {
		e.Error(errs.NewCodeError(errs.SmsCodeError, "验证码错误"))
		return
	}

	resp, err := e.Login(context, s, user)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	e.OK(resp)
}

// @Summary 退出登录
// @Description 退出登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Router /sys/user/logout [post]
// @Security RequireLogin
func (e SysUser) Logout(context *gin.Context) {
	currentUserId := mcontext.GetOpUserID(context)
	e.AuthService.DeleteToken(context, currentUserId, constant.AdminPlatformID)
	e.MakeContext(context)
	e.OK(nil)
}

// @Summary 更新密码（需要旧密码）
// @Description 更新密码（需要旧密码）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserUpdatePasswordReq true "手机号和验证码"
// @Success 200 {object} api.Response
// @Router /sys/user/update_password [post]
// @Security RequireLogin
func (e SysUser) UpdatePassword(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserUpdatePasswordReq{}
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
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	if err != nil {
		log.ZError(e.Context, "修改密码失败，无此id用户", err)
		e.Error(err)
		return
	}
	password := utils.Md5(req.OldPassword)
	if user.PasswordSalt != password {
		e.Error(errs.NewCodeError(errs.AccountOrPasswordError, "旧密码错误"))
		return
	}
	now := customtypes.Time(time.Now())
	updateUser := models.SysUser{
		PasswordSalt:      utils.Md5(req.NewPassword),
		PasswordResetTime: &now,
	}
	updateUser.UpdatedBy = &currentUserId

	err = s.UpdateById(user.ID, updateUser)
	if err != nil {
		log.ZError(e.Context, "修改密码失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 更新密码（无需旧密码，但是需要先进行身份验证）
// @Description 更新密码（无需旧密码，但是需要先进行身份验证）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserChangePasswordReq true "手机号和验证码"
// @Success 200 {object} api.Response
// @Router /sys/user/change_password [post]
// @Security RequireLogin
func (e SysUser) ChangePassword(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserChangePasswordReq{}
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
	currentUserId := mcontext.GetOpUserID(context)
	checked, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).IsCheckedIdentity(currentUserId, constant.ChangePassword)
	if err != nil {
		log.ZError(e.Context, "获取身份验证缓存失败", err)
		e.Error(err)
		return
	}
	if !checked {
		e.Error(errs.ErrCheckIdentity)
		return
	}

	now := customtypes.Time(time.Now())
	updateUser := models.SysUser{
		PasswordSalt:      utils.Md5(req.Password),
		PasswordResetTime: &now,
	}
	updateUser.UpdatedBy = &currentUserId

	err = s.UpdateById(currentUserId, updateUser)
	if err != nil {
		log.ZError(e.Context, "修改密码失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 重置密码
// @Description 重置密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserResetPasswordReq true "用户ID"
// @Success 200 {object} api.Response
// @Router /sys/user/reset_password [post]
// @Security RequireLogin
func (e SysUser) ResetPassword(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserResetPasswordReq{}
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
	superAdmin := mcontext.GetSuperAdmin(context)
	if !superAdmin {
		e.Error(errs.NewCodeError(errs.NoPermissionError, "只有超级管理员才能重置密码哦"))
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(req.ID)
	if err != nil {
		log.ZError(e.Context, "重置密码失败，无此id用户", err)
		e.Error(err)
		return
	}
	newPassword := utils.GenerateRandomString(8)
	now := customtypes.Time(time.Now())
	updateUser := models.SysUser{
		PasswordSalt:      utils.Md5(newPassword),
		PasswordResetTime: &now,
	}
	updateUser.UpdatedBy = &currentUserId

	err = s.UpdateById(user.ID, updateUser)
	if err != nil {
		log.ZError(e.Context, "重置密码失败", err)
		e.Error(err)
		return
	}
	resp := dto.SysUserResetPasswordResp{NewPassword: newPassword}
	e.OK(resp)
}

// @Summary 完善基本信息
// @Description 完善基本信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserPerfectInfoReq true "用户基本信息"
// @Success 200 {object} api.Response
// @Router /sys/user/perfect_info [post]
// @Security RequireLogin
func (e SysUser) PerfectInfo(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserPerfectInfoReq{}
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
	currentUserId := mcontext.GetOpUserID(context)
	updateUser := models.SysUser{}
	utils.CopyStructFields(&updateUser, req)
	updateUser.UpdatedBy = &currentUserId

	err = s.UpdateById(currentUserId, updateUser)
	if err != nil {
		log.ZError(e.Context, "完善信息失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 生成多因素认证密钥
// @Description 生成多因素认证密钥
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=dto.SysUserGenerateMFAKeyResp}
// @Router /sys/user/generate_mfa_device [post]
// @Security RequireLogin
func (e SysUser) GenerateMFADevice(context *gin.Context) {
	s := services.SysUser{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	if err != nil {
		log.ZError(e.Context, "获取当前用户信息失败", err)
		e.Error(err)
		return
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "xiaoxixijz",
		AccountName: fmt.Sprintf("%s@xiaoxixijz.com", user.Username),
	})
	if err != nil {
		log.ZError(e.Context, "生成topt key失败", err)
		e.Error(err)
		return
	}
	resp := dto.SysUserGenerateMFAKeyResp{
		Url:    key.URL(),
		Secret: key.Secret(),
	}
	e.OK(resp)
}

// @Summary 绑定多因素认证设备
// @Description 绑定多因素认证设备
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserBindMFAReq true "mfa信息"
// @Success 200 {object} api.Response{data=dto.SysUserBindMFAKeyResp}
// @Router /sys/user/bind_mfa_device [post]
// @Security RequireLogin
func (e SysUser) BindMFADevice(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserBindMFAReq{}
	err := e.MakeContext(context).
		Bind(&req, binding.JSON).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}

	currentUserId := mcontext.GetOpUserID(context)
	checked, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).IsCheckedIdentity(currentUserId, constant.BindMFADevice)
	if err != nil {
		log.ZError(e.Context, "获取身份验证缓存失败", err)
		e.Error(err)
		return
	}
	if !checked {
		e.Error(errs.ErrCheckIdentity)
		return
	}

	// 生成 TOTP 密码
	totpCode, err := totp.GenerateCode(req.Secret, time.Now())
	if err != nil {
		log.ZError(e.Context, "生成TOTP密码失败", err)
		e.Error(err)
		return
	}
	//校验TOTP密码
	if totpCode != req.Code {
		log.ZError(e.Context, "TOTP密码校验失败", err)
		e.Error(errs.ErrMFACode)
		return
	}
	recoverCode := utils.GenerateDigitalString(6)
	updateUser := models.SysUser{
		MFAKey:      req.Secret,
		EnableMFA:   req.EnableMFA,
		RecoverCode: recoverCode,
	}
	updateUser.UpdatedBy = &currentUserId

	err = s.UpdateById(currentUserId, updateUser)
	if err != nil {
		log.ZError(e.Context, "更新用户多因素认证密钥失败", err)
		e.Error(err)
		return
	}
	resp := dto.SysUserBindMFAKeyResp{
		RecoverCode: recoverCode,
	}
	e.OK(resp)
}

// @Summary 解绑多因素认证设备
// @Description 解绑多因素认证设备
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Router /sys/user/unbind_mfa_device [post]
// @Security RequireLogin
func (e SysUser) UnBindMFADevice(context *gin.Context) {
	s := services.SysUser{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	checked, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).IsCheckedIdentity(currentUserId, constant.UnBindMFADevice)
	if err != nil {
		log.ZError(e.Context, "获取身份验证缓存失败", err)
		e.Error(err)
		return
	}
	if !checked {
		e.Error(errs.ErrCheckIdentity)
		return
	}

	updateUser := models.SysUser{
		MFAKey:      "",
		EnableMFA:   false,
		RecoverCode: "",
	}
	updateUser.UpdatedBy = &currentUserId

	err = s.UpdateColumnsById(currentUserId, updateUser, "mfa_key", "enable_mfa", "recover_code")
	if err != nil {
		log.ZError(e.Context, "更新用户多因素认证密钥失败", err)
		e.Error(err)
		return
	}

	e.OK(nil)
}

// @Summary 获取个人信息
// @Description 获取个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=dto.SysUserPersonalInfoResp}
// @Router /sys/user/get_personal_info [post]
// @Security RequireLogin
func (e SysUser) GetPersonalInfo(context *gin.Context) {
	s := services.SysUser{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	user.IsBindMfaDevice = len(user.MFAKey) > 0
	if err != nil {
		log.ZError(e.Context, "获取用户失败", err)
		e.Error(err)
		return
	}
	roles, err := s.GetRolesByUserId(currentUserId)
	if err != nil {
		log.ZError(e.Context, "获取用户角色失败", err)
		e.Error(err)
		return
	}
	permissions, err := s.GetPermissionsByUserId(currentUserId)
	if err != nil {
		log.ZError(e.Context, "获取用户权限失败", err)
		e.Error(err)
		return
	}

	resp := dto.SysUserPersonalInfoResp{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
	}
	e.OK(resp)
}

// @Summary 获取绑定微信的授权链接
// @Description 获取绑定微信的授权链接
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=dto.SysUserGetOAuth2UrlResp}
// @Router /sys/user/get_oauth2_url [post]
// @Security RequireLogin
func (e SysUser) GetOAuth2Url(context *gin.Context) {
	e.MakeContext(context)

	appId := config.Config.WxOfficialAccount.AppID
	componentAppid := config.Config.WxOpenPlatform.AppID
	redirectURL := "https://" + context.Request.Host + "/api/v1/sys/user/bind_wechat"
	encodedURL := url.QueryEscape(redirectURL)

	currentUserId := mcontext.GetOpUserID(context)
	scope := "snsapi_base"
	oauth2Url := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&component_appid=%s#wechat_redirect",
		appId, encodedURL, scope, utils.Int64ToString(currentUserId), componentAppid)
	e.OK(dto.SysUserGetOAuth2UrlResp{
		Url: oauth2Url,
	})
}

// @Summary 获取微信登录链接
// @Description 获取微信登录链接
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=dto.SysUserGetWxLoginUrlResp}
// @Router /sys/user/get_wechat_login_url [post]
// @Security RequireLogin
func (e SysUser) GetWxLoginUrl(context *gin.Context) {
	e.MakeContext(context)

	appId := config.Config.WxOfficialAccount.AppID
	componentAppid := config.Config.WxOpenPlatform.AppID
	redirectURL := "https://" + context.Request.Host + "/api/v1/sys/user/login_by_wechat"
	encodedURL := url.QueryEscape(redirectURL)

	//随机生成一个id
	state := utils.Int64ToString(utils.GenID())
	scope := "snsapi_base"
	oauth2Url := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&component_appid=%s#wechat_redirect",
		appId, encodedURL, scope, state, componentAppid)
	e.OK(dto.SysUserGetWxLoginUrlResp{
		Url:   oauth2Url,
		State: state,
	})
}

// @Summary 绑定微信
// @Description 微信授权重定向
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=dto.SysUserGetOAuth2UrlResp}
// @Router /sys/user/bind_wechat [get]
func (e SysUser) BindWechat(context *gin.Context) {
	code := context.DefaultQuery("code", "")
	state := context.DefaultQuery("state", "")
	if len(code) == 0 {
		e.Error(errs.NewCodeError(errs.ArgsError, "code cannot be blank"))
		return
	}
	if len(state) == 0 {
		e.Error(errs.NewCodeError(errs.ArgsError, "state cannot be blank"))
		return
	}
	authorizerService := services.WxAuthorizer{}
	userService := services.SysUser{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&authorizerService.Service).
		MakeService(&userService.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	appId := config.Config.WxOfficialAccount.AppID
	authorizer, err := authorizerService.GetByAppID(appId)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	if authorizer == nil {
		e.Error(errs.NewCodeError(errs.ArgsError, fmt.Sprintf("appId %s isn't authorize yet", appId)))
		return
	}
	officialAccount, err := weixin.OpenPlatformApp.OfficialAccount(authorizer.AuthorizerAppid, authorizer.AuthorizerRefreshToken, nil)
	if err != nil {
		e.Error(err)
		return
	}
	officialAccount.OAuth.SetScopes([]string{"snsapi_base"})
	user, err := officialAccount.OAuth.UserFromCode(code)
	fmt2.Dump(user)

	userId := utils.StringToInt64(state)
	sysUser := models.SysUser{}
	sysUser.OpenId = user.GetOpenID()
	userService.UpdateColumnsById(userId, sysUser, "open_id")

	e.OK(nil)
}

// @Summary 微信登录
// @Description 微信授权重定向
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=dto.SysUserGetOAuth2UrlResp}
// @Router /sys/user/login_by_wechat [get]
func (e SysUser) LoginByWechat(context *gin.Context) {
	code := context.DefaultQuery("code", "")
	state := context.DefaultQuery("state", "")
	if len(code) == 0 {
		e.Error(errs.NewCodeError(errs.ArgsError, "code cannot be blank"))
		return
	}
	if len(state) == 0 {
		e.Error(errs.NewCodeError(errs.ArgsError, "state cannot be blank"))
		return
	}
	authorizerService := services.WxAuthorizer{}
	userService := services.SysUser{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&authorizerService.Service).
		MakeService(&userService.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	appId := config.Config.WxOfficialAccount.AppID
	authorizer, err := authorizerService.GetByAppID(appId)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	if authorizer == nil {
		e.Error(errs.NewCodeError(errs.ArgsError, fmt.Sprintf("appId %s isn't authorize yet", appId)))
		return
	}
	officialAccount, err := weixin.OpenPlatformApp.OfficialAccount(authorizer.AuthorizerAppid, authorizer.AuthorizerRefreshToken, nil)
	if err != nil {
		e.Error(err)
		return
	}
	officialAccount.OAuth.SetScopes([]string{"snsapi_base"})
	user, err := officialAccount.OAuth.UserFromCode(code)
	fmt2.Dump(user)

	openID := user.GetOpenID()
	sysUser, err := userService.GetByOpenID(openID)
	if err != nil {
		e.Error(err)
		return
	}
	if sysUser == nil {
		e.Error(errs.NewCodeError(errs.AccountNotExistError, "Cannot find account bind with current wechat"))
	}
	application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).SetWechatLoginResult(state, sysUser.ID)
	e.OK(nil)
}

// @Summary 获取微信登录结果
// @Description 获取微信登录结果
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserGetWxLoginResultReq true "登录id"
// @Success 200 {object} api.Response{data=dto.SysUserGetOAuth2UrlResp}
// @Router /sys/user/get_wechat_login_result [get]
func (e SysUser) GetWechatLoginResult(context *gin.Context) {
	req := dto.SysUserGetWxLoginResultReq{}
	s := services.SysUser{}
	err := e.MakeContext(context).
		Bind(&req, binding.JSON).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	userId, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).GetWechatLoginResult(req.State)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	if len(userId) == 0 {
		//未成功扫码
		e.OK(nil)
		return
	}
	user, err := s.GetById(utils.StringToInt64(userId))
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	resp, err := e.Login(context, s, user)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	e.OK(resp)
}

// @Summary 解绑微信
// @Description 解绑微信
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Router /sys/user/unbind_wechat [get]
func (e SysUser) UnBindWechat(context *gin.Context) {
	userService := services.SysUser{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&userService.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}

	userId := mcontext.GetOpUserID(context)
	sysUser := models.SysUser{}
	sysUser.OpenId = ""
	userService.UpdateColumnsById(userId, sysUser, "open_id")

	e.OK(nil)
}

// @Summary 发送验证邮件
// @Description 发送验证邮件
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserSendVerifyEmailReq true "邮箱"
// @Success 200 {object} api.Response
// @Router /sys/user/send_verify_email [post]
// @Security RequireLogin
func (e SysUser) SendVerifyEmail(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserSendVerifyEmailReq{}
	err := e.MakeContext(context).
		Bind(&req, binding.JSON).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	switch req.Scene {
	case constant.BindEmail:
		exist, err := s.IsEmailExist(req.Email)
		if err != nil {
			log.ZError(e.Context, "", err)
			e.Error(err)
			return
		}
		if exist {
			e.Error(errs.NewCodeError(errs.DuplicateKeyError, "邮箱已被其他用户绑定"))
			return
		}
		code := utils.GenerateDigitalString(6)
		expireTime := int64(10)
		variables := map[string]string{
			"code":       code,
			"expireTime": utils.Int64ToString(expireTime),
		}
		if err != nil {
			log.ZError(e.Context, "", err)
			e.Error(err)
			return
		}
		//发送邮件
		request := email.EmailRequest{
			To:               []string{req.Email},
			Subject:          "验证码",
			TemplateFileName: "verify_email.html",
			Params:           variables,
		}
		err = email.EmailVendorInstance.Send(request)
		if err != nil {
			log.ZError(e.Context, "", err)
			e.Error(err)
			return
		}
		//缓存验证码
		err = e.VerifyCodeCache.StoreCode(req.Email, constant.BindEmail, code, time.Duration(expireTime)*time.Minute)
	case constant.CheckIdentity:
		currentUserId := mcontext.GetOpUserID(context)
		user, err := s.GetById(currentUserId)
		if err != nil {
			log.ZError(e.Context, "", err)
			e.Error(err)
			return
		}
		code := utils.GenerateDigitalString(6)
		expireTime := int64(10)
		variables := map[string]string{
			"code":       code,
			"expireTime": utils.Int64ToString(expireTime),
		}
		if err != nil {
			log.ZError(e.Context, "", err)
			e.Error(err)
			return
		}
		//发送邮件
		request := email.EmailRequest{
			To:               []string{req.Email},
			Subject:          "验证码",
			TemplateFileName: "verify_email.html",
			Params:           variables,
		}
		err = email.EmailVendorInstance.Send(request)
		if err != nil {
			log.ZError(e.Context, "", err)
			e.Error(err)
			return
		}
		//缓存验证码
		err = e.VerifyCodeCache.StoreCode(user.Email, constant.CheckIdentity, code, time.Duration(expireTime)*time.Minute)
	default:
		e.Error(errs.NewCodeError(errs.ArgsError, "不支持的邮件场景"))
		return
	}
	e.OK(nil)
}

// @Summary 绑定邮箱
// @Description 绑定邮箱
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserBindEmailReq true "邮箱和邮箱验证码"
// @Success 200 {object} api.Response
// @Router /sys/user/bind_email [post]
// @Security RequireLogin
func (e SysUser) BindEmail(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserBindEmailReq{}
	err := e.MakeContext(context).
		Bind(&req, binding.JSON).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}

	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	if len(user.Phone) > 0 {
		checked, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).IsCheckedIdentity(currentUserId, constant.ChangeEmail)
		if err != nil {
			log.ZError(e.Context, "获取身份验证缓存失败", err)
			e.Error(err)
			return
		}
		if !checked {
			e.Error(errs.ErrCheckIdentity)
			return
		}
	}

	pass, err := e.VerifyCodeCache.CheckCode(context, req.Email, constant.BindEmail, req.Code)
	if err != nil {
		e.Error(err)
		return
	}
	if !pass {
		e.Error(errs.NewCodeError(errs.EmailCodeError, "验证码错误"))
		return
	}
	updateUser := models.SysUser{}
	updateUser.Email = req.Email
	s.UpdateById(currentUserId, updateUser)

	e.OK(nil)
}

// @Summary 绑定手机
// @Description 绑定手机
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserBindPhoneReq true "手机号和短信验证码"
// @Success 200 {object} api.Response
// @Router /sys/user/bind_phone [post]
// @Security RequireLogin
func (e SysUser) BindPhone(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserBindPhoneReq{}
	err := e.MakeContext(context).
		Bind(&req, binding.JSON).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	if len(user.Phone) > 0 {
		checked, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).IsCheckedIdentity(currentUserId, constant.ChangePhone)
		if err != nil {
			log.ZError(e.Context, "获取身份验证缓存失败", err)
			e.Error(err)
			return
		}
		if !checked {
			e.Error(errs.ErrCheckIdentity)
			return
		}
	}

	pass, err := e.VerifyCodeCache.CheckCode(context, req.Phone, constant.BindPhone, req.Code)
	if err != nil {
		e.Error(err)
		return
	}
	if !pass {
		e.Error(errs.NewCodeError(errs.SmsCodeError, "验证码错误"))
		return
	}
	updateUser := models.SysUser{}
	updateUser.Phone = req.Phone
	s.UpdateById(currentUserId, updateUser)

	e.OK(nil)
}

// @Summary 切换是否开启MFA
// @Description 切换是否开启MFA
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Router /sys/user/switch_enable_mfa [post]
// @Security RequireLogin
func (e SysUser) SwitchEnableMFA(context *gin.Context) {
	s := services.SysUser{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	checked, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).IsCheckedIdentity(currentUserId, constant.SwitchEnableMFA)
	if err != nil {
		log.ZError(e.Context, "获取身份验证缓存失败", err)
		e.Error(err)
		return
	}
	if !checked {
		e.Error(errs.ErrCheckIdentity)
		return
	}
	updateUser := models.SysUser{
		EnableMFA: !user.EnableMFA,
	}
	updateUser.UpdatedBy = &currentUserId
	err = s.UpdateColumnsById(currentUserId, updateUser, "enable_mfa")
	if err != nil {
		log.ZError(e.Context, "更新数据库失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 当前操作是否需要进行身份验证
// @Description 当前操作是否需要进行身份验证
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserShouldCheckIdentityReq true "操作"
// @Success 200 {object} api.Response{data=dto.SysUserShouldCheckIdentityResp}
// @Router /sys/user/should_check_identity [post]
// @Security RequireLogin
func (e SysUser) ShouldCheckIdentity(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserShouldCheckIdentityReq{}
	err := e.MakeContext(context).
		Bind(&req, binding.JSON).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	if err != nil {
		log.ZError(e.Context, "获取当前用户信息失败", err)
		e.Error(err)
		return
	}
	resp := dto.SysUserShouldCheckIdentityResp{}
	checked, err := application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).IsCheckedIdentity(currentUserId, req.Operation)
	if err != nil {
		log.ZError(e.Context, "获取身份验证缓存失败", err)
		e.Error(err)
		return
	}
	if checked {
		resp.ShouldCheck = false
		e.OK(resp)
		return
	}

	switch req.Operation {
	case constant.ChangeEmail:
		shouldCheck := len(user.Email) > 0
		resp.ShouldCheck = shouldCheck
		if shouldCheck {
			resp.DefaultCheckType = "email"
			resp.CheckTypes = []string{"email"}

			if len(user.Phone) > 0 {
				resp.CheckTypes = append(resp.CheckTypes, "phone")
			}
			if user.EnableMFA {
				resp.CheckTypes = append(resp.CheckTypes, "mfa")
			}
		}
	case constant.ChangePhone:
		shouldCheck := len(user.Phone) > 0
		resp.ShouldCheck = shouldCheck
		if shouldCheck {
			resp.DefaultCheckType = "phone"
			resp.CheckTypes = []string{"phone"}

			if len(user.Email) > 0 {
				resp.CheckTypes = append(resp.CheckTypes, "email")
			}
			if user.EnableMFA {
				resp.CheckTypes = append(resp.CheckTypes, "mfa")
			}
		}
	case constant.BindMFADevice, constant.UnBindMFADevice, constant.BindWechat, constant.UnbindWechat:
		resp.ShouldCheck = true
		resp.DefaultCheckType = "phone"
		resp.CheckTypes = []string{}
		if len(user.Phone) > 0 {
			resp.CheckTypes = append(resp.CheckTypes, "phone")
		}
		if len(user.Email) > 0 {
			resp.CheckTypes = append(resp.CheckTypes, "email")
		}
		if len(user.MFAKey) > 0 {
			resp.DefaultCheckType = "mfa"
			resp.CheckTypes = append(resp.CheckTypes, "mfa")
		}
	case constant.ChangePassword:
		resp.ShouldCheck = true
		resp.DefaultCheckType = "password"
		resp.CheckTypes = []string{"password"}
	case constant.SwitchEnableMFA:
		resp.ShouldCheck = true
		resp.DefaultCheckType = "mfa"
		resp.CheckTypes = []string{}
		if len(user.MFAKey) > 0 {
			resp.CheckTypes = append(resp.CheckTypes, "mfa")
		}
		if len(user.Phone) > 0 {
			resp.CheckTypes = append(resp.CheckTypes, "phone")
		}
		if len(user.Email) > 0 {
			resp.CheckTypes = append(resp.CheckTypes, "email")
		}
	}
	e.OK(resp)
}

// @Summary 身份验证
// @Description 身份验证
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param param body dto.SysUserCheckIdentityReq true "验证方式和验证码"
// @Success 200 {object} api.Response
// @Router /sys/user/check_identity [post]
// @Security RequireLogin
func (e SysUser) CheckIdentity(context *gin.Context) {
	s := services.SysUser{}
	req := dto.SysUserCheckIdentityReq{}
	err := e.MakeContext(context).
		Bind(&req, binding.JSON).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	currentUserId := mcontext.GetOpUserID(context)
	user, err := s.GetById(currentUserId)
	if err != nil {
		log.ZError(e.Context, "获取当前用户信息失败", err)
		e.Error(err)
		return
	}
	switch req.CheckType {
	case "phone":
		pass, err := e.VerifyCodeCache.CheckCode(e.Context, user.Phone, constant.CheckIdentity, req.Code)
		if err != nil {
			e.Error(err)
			return
		}
		if !pass {
			e.Error(errs.NewCodeError(errs.SmsCodeError, "验证码错误"))
			return
		}
	case "email":
		pass, err := e.VerifyCodeCache.CheckCode(e.Context, user.Email, constant.CheckIdentity, req.Code)
		if err != nil {
			e.Error(err)
			return
		}
		if !pass {
			e.Error(errs.NewCodeError(errs.SmsCodeError, "验证码错误"))
			return
		}
	case "mfa":
		// 生成 TOTP 密码
		totpCode, err := totp.GenerateCode(user.MFAKey, time.Now())
		if err != nil {
			log.ZError(e.Context, "生成TOTP密码失败", err)
			e.Error(err)
			return
		}
		//校验TOTP密码
		if totpCode != req.Code {
			log.ZError(e.Context, "TOTP密码校验失败", err)
			e.Error(errs.ErrMFACode)
			return
		}
	case "password":
		password := utils.Md5(req.Code)
		if user.PasswordSalt != password {
			e.Error(errs.NewCodeError(errs.AccountOrPasswordError, "密码错误"))
			return
		}
	}
	err = application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).SetSkipCheckIdentity(currentUserId, req.Operation)
	if err != nil {
		e.AddError(err)
		return
	}
	e.OK(nil)
}

func (e SysUser) updateCurrentUserPerms(id int64) {

}
