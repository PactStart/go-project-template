package dto

import (
	"orderin-server/internal/models"
	"orderin-server/pkg/common/dto"
)

type SysDictAddReq struct {
	Name        string `json:"name" vd:"@:len($)>0 && len($)<50"`
	Label       string `json:"label" vd:"@:len($)>0"`
	Description string `json:"description" vd:"@:len($)<100"`
}

type SysDictUpdateReq struct {
	ID          int64  `json:"id" vd:"@:$>0"`
	Description string `json:"description" vd:"@:len($)<100"`
}

type SysDictDeleteReq struct {
	ID int64 `json:"id" vd:"@:$>0"`
}

type SysDictPageQueryReq struct {
	dto.Pagination `search:"-"`
	Keyword        string `json:"keyword" search:"type:contains;column:name,label,description;table:sys_dicts"`
	DictOrder
}

type DictOrder struct {
	DictIdOrder string `json:"idOrder" search:"type:order;column:id;table:sys_dicts" `
}

type SysDictItemPageQueryReq struct {
	dto.Pagination `search:"-"`
	DictName       string `json:"dictName" vd:"@:len($)>0 && len($)<50" search:"type:exact;column:dict_name;table:sys_dict_items"`
	Keyword        string `json:"keyword" search:"type:contains;column:item_label,item_value,description;table:sys_dict_items"`
	DictItemOrder
}

type DictItemOrder struct {
	SortOrder string `json:"sortOrder" search:"type:order;column:sort;table:sys_dict_items" `
}

type SysDictItemAddReq struct {
	DictName    string `json:"dictName" vd:"@:len($)>0 && len($)<50"`
	ItemLabel   string `json:"itemLabel" vd:"@:len($)>0"`
	ItemValue   string `json:"itemValue" vd:"@:len($)>0"`
	Description string `json:"description" vd:"@:len($)<50"`
	IsDefault   bool   `json:"isDefault"`
	Sort        int    `json:"sort"`
}

type SysDictItemUpdateReq struct {
	ID          int64  `json:"id" vd:"@:$>0"`
	ItemLabel   string `json:"itemLabel" vd:"@:len($)>0"`
	ItemValue   string `json:"itemValue" vd:"@:len($)>0"`
	Description string `json:"description" vd:"@:len($)<50"`
	Status      int    `json:"status" vd:"@:in($,1,2)"`
	IsDefault   bool   `json:"isDefault"`
	Sort        int    `json:"sort"`
}

type SysDictGetItemsReq struct {
	Name string `json:"name" vd:"@:len($)>0 && len($)<50"`
}

type SysDictBatchGetItemsReq struct {
	Names []string `json:"names" vd:"@:len($)>0"`
}

type SysDictGetItemsResp struct {
	List []models.SysDictItem `json:"list"`
}
