package dto

import (
	"orderin-server/internal/models"
	"orderin-server/pkg/common/dto"
)

type SysUserPageQueryReq struct {
	dto.Pagination `search:"-"`
	Username       string `json:"username" search:"type:exact;column:username;table:sys_users"`
	Nickname       string `json:"nickname" search:"type:contains;column:nickname;table:sys_users"`
	RealName       string `json:"realName" search:"type:exact;column:real_name;table:sys_users"`
	Phone          string `json:"phone" search:"type:exact;column:phone;table:sys_users"`

	Keyword string `json:"keyword" search:"type:contains;column:username,nickname,real_name,phone;table:sys_users"`

	Status        int `json:"status" search:"type:exact;column:status;table:sys_users"`
	OwnRole       `search:"type:inner;on:user_id:id;table:sys_users;join:sys_role_users"`
	ExcludeRoleId int64 `json:"excludeRoleId,string" search:"-"`

	UserOrder
}

type UserOrder struct {
	UserIdOrder string `json:"idOrder" search:"type:order;column:id;table:sys_users" `
}

type OwnRole struct {
	RoleId int64 `json:"roleId,string" search:"type:exact;column:role_id;table:sys_role_users"`
}

type SysUserGetReq struct {
	ID int64 `json:"id,string" vd:"@:$>0""`
}

type SysUserInsertReq struct {
	Username string `json:"username" vd:"@:len($)>4 && len($)<20"`
	RealName string `json:"realName" vd:"@:len($)==0 || len($)<20"`
	Nickname string `json:"nickname" vd:"@:len($) == 0 || len($)<20"`
	Phone    string `json:"phone" vd:"@:len($)==0 || (len($)>0 && phone($))"`
	Email    string `json:"email" vd:"@:len($)==0 || (len($)>0 && email($))"`
}

type SysUserUpdatePasswordReq struct {
	OldPassword string `json:"oldPassword" vd:"@:len($)>1"`
	NewPassword string `json:"newPassword" vd:"@:regexp('^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$')"`
}

type SysUserChangePasswordReq struct {
	Password string `json:"newPassword" vd:"@:regexp('^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$')"`
}

type SysUserResetPasswordReq struct {
	ID int64 `json:"id,string" vd:"@:$>0""`
}

type SysUserPerfectInfoReq struct {
	Nickname       string `json:"nickname" vd:"@:len($)<20"`
	Avatar         string `json:"avatar"`
	Wechat         string `json:"wechat"`
	WechatNickname string `json:"wechatNickname"`
	WechatQrCode   string `json:"wechatQrCode"`
}

type SysUserBindMFAReq struct {
	Secret    string `json:"secret" vd:"@:len($)>0"`
	Code      string `json:"code" vd:"@:regexp('^[0-9]{6}$')"`
	EnableMFA bool   `json:"enableMFA"`
}

type SysUserShouldCheckIdentityReq struct {
	Operation string `json:"operation" vd:"@:len($)>0"`
}

type SysUserCheckIdentityReq struct {
	Operation string `json:"operation" vd:"@:len($)>0"`
	CheckType string `json:"checkType" vd:"@:len($)>0"`
	Code      string `json:"code" vd:"@:len($)>0"`
}

type SysUserLoginByPasswordReq struct {
	Username string `json:"username" vd:"@:len($)>4 && len($)<20"`
	Password string `json:"password" vd:"@:regexp('^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$')"`
}

type SysUserLoginByPhoneReq struct {
	Phone   string `json:"phone" vd:"@:phone($)"`
	SmsCode string `json:"smsCode" vd:"@:regexp('^[0-9]{6}$')"`
}

type SysUserUpdateStatusReq struct {
	ID     int64 `json:"id,string" vd:"@:$>0"`
	Status int   `json:"status" vd:"@:$>0"`
}

type SysUserSendVerifyEmailReq struct {
	Email string `json:"email" vd:"@:len($)>0"`
	Scene string `json:"scene" vd:"@:len($)>0"`
}

type SysUserBindEmailReq struct {
	Email string `json:"email" vd:"@:len($)==0 || (len($)>0 && email($))"`
	Code  string `json:"code" vd:"@:len($)>0"`
}

type SysUserBindPhoneReq struct {
	Phone string `json:"phone" vd:"@:len($)==0 || (len($)>0 && phone($))"`
	Code  string `json:"code" vd:"@:len($)>0"`
}

type SysUserGetWxLoginResultReq struct {
	State string `json:"state" vd:"@:len($)>0"`
}

////////////////////响应////////////////////

type SysUserAddResp struct {
	ID           int64  `json:"id,string"`
	InitPassword string `json:"initPassword"`
}

type SysUserLoginResp struct {
	Token             string `json:"token"`
	ExpireTimeSeconds int64  `json:"expireTimeSeconds"`
}

type SysUserPersonalInfoResp struct {
	User        *models.SysUser `json:"user"`
	Roles       []string        `json:"roles"`
	Permissions []string        `json:"permissions"`
}

type SysUserOAuth2UrlResp struct {
	Url string `json:"url"`
}

type SysUserResetPasswordResp struct {
	NewPassword string `json:"newPassword"`
}

type SysUserGenerateMFAKeyResp struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
}

type SysUserBindMFAKeyResp struct {
	RecoverCode string `json:"recoverCode"`
}

type SysUserGetOAuth2UrlResp struct {
	Url string `json:"url"`
}

type SysUserGetWxLoginUrlResp struct {
	Url   string `json:"url"`
	State string `json:"state"`
}

type SysUserShouldCheckIdentityResp struct {
	ShouldCheck      bool     `json:"shouldCheck"`
	DefaultCheckType string   `json:"defaultCheckType"`
	CheckTypes       []string `json:"checkTypes"`
}
