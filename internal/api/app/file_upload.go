package app

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"mime/multipart"
	"orderin-server/internal/dto"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/file_store"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
	"path/filepath"
	"strings"
)

type FileUpload struct {
	api.Api
}

const maxFileSize = 3 * 1024 * 1024 // 3MB

// @Summary 上传图片
// @Description 上传图片
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} api.Response{data=dto.UploadFileResp}
// @Router /file/upload_image [post]
func (e FileUpload) UploadImage(c *gin.Context) {
	e.MakeContext(c)
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		e.Error(errs.NewCodeError(errs.ArgsError, "No file uploaded"))
		return
	}
	// 验证文件类型是否为图片
	if !isImageFile(file) {
		e.Error(errs.NewCodeError(errs.ArgsError, "Only image files are allowed"))
		return
	}
	// 检查文件大小
	if file.Size > maxFileSize {
		e.Error(errs.NewCodeError(errs.ArgsError, "File size exceeds the limit（3M）"))
		return
	}
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		e.Error(errs.NewCodeError(errs.ServerInternalError, "Failed to open file"))
		return
	}
	defer src.Close()
	// 生成目标文件名
	destination := generateDestination(file.Filename)
	url, err := file_store.FileStore.UpLoad(destination, src)
	if err != nil {
		e.Error(err)
		return
	}
	if url == nil {
		e.Error(errs.NewCodeError(errs.ServerInternalError, "Failed to Get uploaded file url"))
		return
	}
	e.OK(dto.UploadFileResp{
		Url: *url,
	})
}

// @Summary 上传Base64格式图片
// @Description 上传Base64格式图片
// @Tags 文件上传
// @Accept json
// @Produce json
// @Param param body dto.Base64ImageUploadReq false "base64格式图片"
// @Success 200 {object} api.Response{data=dto.UploadFileResp}
// @Router /file/upload_base64_image [post]
func (e FileUpload) UploadBase64Image(c *gin.Context) {
	req := dto.Base64ImageUploadReq{}
	err := e.MakeContext(c).Bind(&req, binding.JSON).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	// 解码 Base64 图片数据
	imageData, err := base64.StdEncoding.DecodeString(req.Base64Image)
	if err != nil {
		log.ZError(c, "Failed to decode Base64 image", err, "image", req.Base64Image)
		e.Error(errs.NewCodeError(errs.ArgsError, "Failed to decode Base64 image"))
		return
	}
	src := strings.NewReader(string(imageData))
	// 生成目标文件名
	destination := generateDestination("base64.png")
	url, err := file_store.FileStore.UpLoad(destination, src)
	if err != nil {
		e.Error(err)
		return
	}
	if url == nil {
		e.Error(errs.NewCodeError(errs.ServerInternalError, "Failed to Get uploaded file url"))
		return
	}
	e.OK(dto.UploadFileResp{
		Url: *url,
	})
}

// 验证文件类型是否为图片
func isImageFile(file *multipart.FileHeader) bool {
	extension := filepath.Ext(file.Filename)
	switch extension {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	default:
		return false
	}
}

// 生成目标文件名
func generateDestination(filename string) string {
	// 可根据需要自定义生成目标文件名的逻辑
	// 这里简单地以当前时间戳作为文件名
	return fmt.Sprintf("xxxjz-app/%s/%d%s", config.Config.Env.Profiles, utils.GenID(), filepath.Ext(filename))
}
