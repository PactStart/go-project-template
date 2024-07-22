package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/api"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
)

type SysDict struct {
	api.Api
}

// @Summary 添加字典
// @Description 添加字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictAddReq false "字典信息"
// @Success 200 {object} api.Response
// @Router /sys/dict/add [post]
// @Security RequireLogin
func (e SysDict) Add(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictAddReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}

	model := models.SysDict{}
	utils.CopyStructFields(&model, req)
	currentUserId := mcontext.GetOpUserID(context)
	model.CreatedBy = &currentUserId

	err = s.Insert(model)
	if err != nil {
		log.ZError(e.Context, "添加字典失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)

}

// @Summary 修改字典
// @Description 根据id修改字典信息
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictUpdateReq false "字典信息"
// @Success 200 {object} api.Response
// @Router /sys/dict/update [post]
// @Security RequireLogin
func (e SysDict) Update(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictUpdateReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	model := models.SysDict{}
	utils.CopyStructFields(&model, req)
	currentUserId := mcontext.GetOpUserID(context)
	model.UpdatedBy = &currentUserId

	err = s.UpdateSelectiveById(model.ID, model)
	if err != nil {
		log.ZError(e.Context, "修改字典失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 删除字典
// @Description 根据id删除字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictDeleteReq false "要删除的字典ID"
// @Success 200 {object} api.Response
// @Router /sys/dict/delete [post]
// @Security RequireLogin
func (e SysDict) Delete(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictDeleteReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	s.Delete(req.ID)
	e.OK(nil)
}

// @Summary 分页查询字典
// @Description 分页查询字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictPageQueryReq false "字典筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysDict}}
// @Router /sys/dict/page_query [post]
// @Security RequireLogin
func (e SysDict) PageQuery(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictPageQueryReq{}
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
	list := make([]models.SysDict, 0)
	var count int64

	err = s.PageQuery(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}

// @Summary 添加字典项
// @Description 添加字典项
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictItemAddReq true "字典项"
// @Success 200 {object} api.Response
// @Router /sys/dict/item/add [post]
// @Security RequireLogin
func (e SysDict) AddItem(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictItemAddReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}

	model := models.SysDictItem{}
	utils.CopyStructFields(&model, req)
	currentUserId := mcontext.GetOpUserID(context)
	model.CreatedBy = &currentUserId

	err = s.AddItem(model)
	if err != nil {
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 修改字典项
// @Description 修改字典项
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictItemUpdateReq true "字典项"
// @Success 200 {object} api.Response
// @Router /sys/dict/item/update [post]
// @Security RequireLogin
func (e SysDict) UpdateItem(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictItemUpdateReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	dictItem := models.SysDictItem{}
	utils.CopyStructFields(&dictItem, req)
	currentUserId := mcontext.GetOpUserID(e.Context)
	dictItem.UpdatedBy = &currentUserId

	err = s.UpdateItem(dictItem)
	if err != nil {
		e.Error(err)
	}
	e.OK(nil)
}

// @Summary 分页查询字典
// @Description 分页查询字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictPageQueryReq false "字典筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysDict}}
// @Router /sys/dict/item/page_query [post]
// @Security RequireLogin
func (e SysDict) PageQueryItem(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictItemPageQueryReq{}
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
	list := make([]models.SysDictItem, 0)
	var count int64

	err = s.PageQueryItem(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}

// @Summary 根据字典名称获取所有项目
// @Description 根据字典名称获取所有项目
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictGetItemsReq true "字典名称"
// @Success 200 {object} api.Response{data=dto.SysDictGetItemsResp}
// @Router /sys/dict/item/get_by_name [post]
// @Security RequireLogin
func (e SysDict) GetItemsByName(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictGetItemsReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	list, err := s.GetItems(req.Name)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	e.OK(dto.SysDictGetItemsResp{List: list})
}

// @Summary 根据字典名称数组获取所有项目
// @Description 根据字典名称数组获取所有项目
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param param body dto.SysDictGetItemsReq true "字典名称数组"
// @Success 200 {object} api.Response{data=map[string][]models.SysDictItem}
// @Router /sys/dict/item/batch_get_by_names [post]
// @Security RequireLogin
func (e SysDict) BatchGetItemsByNames(context *gin.Context) {
	s := services.SysDict{}
	req := dto.SysDictBatchGetItemsReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	mapping, err := s.BatchGetItems(req.Names)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	e.OK(mapping)
}
