package api

import (
	"github.com/gin-gonic/gin"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/errs"
	"reflect"
)

type Response struct {
	RequestId string `json:"requestId"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Detail    string `json:"detail"`
	Data      any    `json:"data,omitempty"`
}

type PageData struct {
	Total     int         `json:"total"`
	PageIndex int         `json:"pageNum"`
	PageSize  int         `json:"pageSize"`
	List      interface{} `json:"list"`
}

type ResponseFormat interface {
	Format()
}

func isAllFieldsPrivate(v any) bool {
	typeOf := reflect.TypeOf(v)
	if typeOf == nil {
		return false
	}
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return false
	}
	num := typeOf.NumField()
	for i := 0; i < num; i++ {
		c := typeOf.Field(i).Name[0]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func Success(c *gin.Context, data any) *Response {
	if format, ok := data.(ResponseFormat); ok {
		format.Format()
	}
	if isAllFieldsPrivate(data) {
		return &Response{}
	}
	return &Response{
		RequestId: mcontext.GetRequestId(c),
		Code:      0,
		Msg:       "success",
		Data:      data,
	}
}

func ParseError(c *gin.Context, err error) *Response {
	if err == nil {
		return Success(c, nil)
	}
	unwrap := errs.Unwrap(err)
	if codeErr, ok := unwrap.(errs.CodeError); ok {
		resp := Response{Code: codeErr.Code(), Msg: codeErr.Msg(), Detail: codeErr.Detail()}
		return &resp
	}
	return &Response{RequestId: mcontext.GetRequestId(c), Code: errs.ServerInternalError, Msg: "服务器错误", Detail: err.Error()}
}
