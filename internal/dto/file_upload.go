package dto

type UploadFileResp struct {
	Url string `json:"url"`
}

type Base64ImageUploadReq struct {
	Base64Image string `json:"image"`
}

type BlobImageUploadReq struct {
	Blob     []byte `json:"blob"`
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
}

type ImageUrlUploadReq struct {
	Url string `json:"url"`
}
