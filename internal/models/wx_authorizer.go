package models

import "time"

type WxAuthorizer struct {
	SnowID
	ComponentAppid   string `gorm:"column:component_appid;type:varchar(100);comment:第三方平台 appid" json:"componentAppid"`
	AuthorizerAppid  string `gorm:"column:authorizer_appid;type:varchar(100);comment:授权方 appid" json:"authorizerAppid"`
	DefaultKfAccount string `gorm:"column:default_kf_account;type:varchar(100);comment:默认的客服账号" json:"defaultKfAccount"`

	NickName      string `gorm:"column:nick_name;type:varchar(100);comment:授权方昵称" json:"nickName"`
	HeadImg       string `gorm:"column:Head_img;type:varchar(255);comment:授权方头像" json:"headImg"`
	UserName      string `gorm:"column:user_name;type:varchar(100);comment:授权方原始ID" json:"userName"`
	QrcodeUrl     string `gorm:"column:qrcode_url;type:varchar(255);comment:授权方二维码图片的 URL" json:"qrcodeUrl"`
	ServiceType   string `gorm:"column:service_type;type:varchar(100);comment:服务类型" json:"serviceType"`
	VerifyType    string `gorm:"column:verify_type;type:varchar(100);comment:认证类型" json:"verifyType"`
	PrincipalName string `gorm:"column:principal_name;type:varchar(100);comment:主体名称" json:"principalName"`

	AuthorizationCode          string    `gorm:"column:authorization_code;type:varchar(200);comment:授权码" json:"-"`
	AuthorizationCodeExpiresAt time.Time `gorm:"column:authorization_code_expires_at;comment:令牌过期时间" json:"-"`

	AuthorizerAccessToken string    `gorm:"column:authorizer_access_token;type:varchar(200);comment:授权方令牌" json:"-"`
	AccessTokenExpiresAt  time.Time `gorm:"column:access_token_expires_at;comment:令牌过期时间" json:"-"`

	AuthorizerRefreshToken string `gorm:"column:authorizer_refresh_token;type:varchar(200);comment:刷新令牌" json:"-"`
	BaseModel
}

type WxAuthorizerMedia struct {
	SnowID
	AppId     string `gorm:"column:app_id;type:varchar(100);comment:公众号appId" json:"appId"`
	UserID    *int64 `gorm:"column:user_id;comment:用户id;" json:"userId"`
	Category  string `gorm:"column:category;type:varchar(20);comment:分类" json:"category"`
	MediaType string `gorm:"column:media_type;type:varchar(20);comment:媒体类型" json:"mediaType"`
	MediaID   string `gorm:"column:me;type:varchar(50);comment:媒体id" json:"mediaId"`
	BaseModel
}
