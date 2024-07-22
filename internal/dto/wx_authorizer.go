package dto

import "orderin-server/pkg/common/dto"

type WxAuthorizerPageQueryReq struct {
	dto.Pagination `search:"-"`
	Keyword        string `json:"keyword" search:"type:contains;column:nick_name,principal_name;table:sys_users"`
}

type GetAuthorizeUrlResp struct {
	Url string `json:"url"`
}

type WxAuthorizerSyncInfoReq struct {
	ID int64 `json:"id,string" vd:"@:$>0"`
}

type KfAccountCreateReq struct {
	AppID    string `json:"appId" vd:"@:len($)>0"`
	Account  string `json:"account" vd:"@:len($)>0"`
	Nickname string `json:"nickname" vd:"@:len($)>0"`
}

type KfAccountGetAllReq struct {
	AppID string `json:"appId" vd:"@:len($)>0"`
}

type KfAccountDeleteReq struct {
	AppID   string `json:"appId" vd:"@:len($)>0"`
	Account string `json:"account" vd:"@:len($)>0"`
}

type KfAccountSetDefaultReq struct {
	AppID   string `json:"appId" vd:"@:len($)>0"`
	Account string `json:"account" vd:"@:len($)>0"`
}

type MenuClearReq struct {
	AppID string `json:"appId" vd:"@:len($)>0"`
}
