package models

import "orderin-server/pkg/common/customtypes"

type SysUser struct {
	SnowID
	Username          string            `gorm:"column:username;type:varchar(20);comment:账户名;unique;not null" json:"username"`
	PasswordSalt      string            `gorm:"column:password_salt;type:varchar(64);comment:密码加盐;not null" json:"-"`
	PasswordResetTime *customtypes.Time `gorm:"column:password_reset_time;type:datetime;comment:密码重置时间" json:"passwordResetTime"`
	Avatar            string            `gorm:"column:avatar;type:varchar(255);default:'';comment:头像;" json:"avatar"`
	Nickname          string            `gorm:"column:nickname;type:varchar(50);comment:昵称" json:"nickname"`
	RealName          string            `gorm:"column:real_name;type:varchar(50);comment:真实姓名" json:"realName"`
	Phone             string            `gorm:"column:phone;type:varchar(20);comment:手机号" json:"phone"`
	Email             string            `gorm:"column:email;type:varchar(100);comment:邮箱" json:"email"`
	OpenId            string            `gorm:"column:open_id;type:varchar(50);comment:微信openid" json:"openId"`
	Wechat            string            `gorm:"column:wechat;type:varchar(50);comment:微信号" json:"wechat"`
	WechatNickname    string            `gorm:"column:wechat_nickname;type:varchar(50);comment:微信昵称" json:"wechatNickname"`
	WechatQrCode      string            `gorm:"column:wechat_qr_code;type:varchar(255);comment:微信二维码" json:"wechatQrCode"`
	EnableMFA         bool              `gorm:"column:enable_mfa;type:tinyint(1);default:0;comment:是否开启二步验证" json:"enableMFA"`
	MFAKey            string            `gorm:"column:mfa_key;type:varchar(100);default:'';comment:多因素认证密钥" json:"-"`
	RecoverCode       string            `gorm:"column:recover_code;type:varchar(10);comment:恢复代码，用于重置二步验证" json:"-"`
	Status            int               `gorm:"column:status;type:int(8);default:1;comment:状态：1正常，2禁用" json:"status"`
	LastLoginTime     *customtypes.Time `gorm:"column:last_login_time;type:datetime;comment:最近登录时间" json:"lastLoginTime"`
	LastLoginIp       string            `gorm:"column:last_login_ip;type:varchar(50);comment:最近登录ip" json:"lastLoginIp"`
	SuperAdmin        bool              `gorm:"column:super_admin;type:tinyint(1);default:0;comment:是否超级管理员" json:"superAdmin"`
	BaseModel

	IsBindMfaDevice bool `gorm:"-" json:"isBindMfaDevice"`
}
