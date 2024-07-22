package api

import (
	"errors"
	"fmt"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
)

type Api struct {
	Context *gin.Context
	Orm     *gorm.DB
	Errors  error
}

func (e *Api) AddError(err error) {
	if e.Errors == nil {
		e.Errors = err
	} else if err != nil {
		log.ZError(e.Context, "", err)
		e.Errors = fmt.Errorf("%v; %w", e.Errors, err)
	}
}

// MakeContext 设置http上下文
func (e *Api) MakeContext(c *gin.Context) *Api {
	e.Context = c
	return e
}

// Bind 参数校验
func (e *Api) Bind(d interface{}, bindings ...binding.Binding) *Api {
	var err error
	if len(bindings) == 0 {
		bindings = constructor.GetBindingForGin(d)
	}
	for i := range bindings {
		if bindings[i] == nil {
			err = e.Context.ShouldBindUri(d)
		} else {
			err = e.Context.ShouldBindWith(d, bindings[i])
		}
		if err != nil && err.Error() == "EOF" {
			log.ZWarn(e.Context, "request body is not present anymore. ", err)
			err = nil
			continue
		}
		if err != nil {
			e.AddError(err)
			break
		}
	}
	if err1 := vd.Validate(d); err1 != nil {
		if vdError, ok := err1.(*vd.Error); ok {
			e.AddError(errs.NewCodeError(errs.ArgsError, fmt.Sprintf("validate fail: %s", vdError.FailPath)))
		} else {
			e.AddError(errs.NewCodeError(errs.ArgsError, err1.Error()))
		}
	}
	return e
}

// GetOrm 获取Orm DB
func (e Api) GetOrm() *gorm.DB {
	return e.Orm
}

// MakeOrm 设置Orm DB
func (e *Api) MakeOrm() *Api {
	idb, exist := e.Context.Get("db")
	if !exist {
		err := errors.New("数据库连接获取失败")
		log.ZError(e.Context, "cannot get db from gin.context", err)
		e.AddError(err)
	}
	switch idb.(type) {
	case *gorm.DB:
		e.Orm = idb.(*gorm.DB)
	default:
		err := errors.New("数据库连接获取失败")
		log.ZError(e.Context, "the type of db from gin.context isn't *gorm.DB", err)
		e.AddError(err)
	}
	return e
}

func (e *Api) MakeService(c *service.Service) *Api {
	c.Orm = e.Orm
	c.Context = e.Context
	return e
}

// Error 通常错误数据处理
func (e Api) Error(err error) {
	GinError(e.Context, err)
}

// OK 通常成功数据处理
func (e Api) OK(data interface{}) {
	GinSuccess(e.Context, data)
}

// PageOK 分页数据处理
func (e Api) PageOK(result interface{}, count int, pageIndex int, pageSize int) {
	var data PageData
	data.List = result
	data.Total = count
	data.PageIndex = pageIndex
	data.PageSize = pageSize
	e.OK(data)
}
