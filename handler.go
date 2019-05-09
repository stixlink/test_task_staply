package main

import (
	"github.com/gin-gonic/gin"
)

var (
	allowedFileExtension = []string{"jpg", "jpeg", "png"}
	allowedContentType   = []string{"image/jpeg", "image/pjpeg", "image/png"}
)

func NewHandler(imageSaveDir string) *UploadHandler {
	return &UploadHandler{
		saveDir:      imageSaveDir,
		ImageResizer: NewImagickResizer(),
	}
}

type UploadHandler struct {
	saveDir string
	ImageResizer
}

func (h *UploadHandler) SaveFormData(c *gin.Context) {

	// default max data size  32 << 20 // 32 MB
	// for change settings set router.MaxMultipartMemory before createserver
	// Example router.MaxMultipartMemory = N << 20  // N MiB
	formFile, err := c.FormFile("image")
	if err != nil {
		c.Writer.WriteHeader(400)
		c.Writer.Write([]byte("Error get image"))
		return
	}

	validator := NewFormDataValidator(formFile)
	ok, err := validator.Validate()
	if err != nil || !ok {
		c.Writer.WriteHeader(400)
		c.Writer.Write([]byte("Invalid file type or data"))
		return
	}

	filePath := h.saveDir + formFile.Filename
	err = c.SaveUploadedFile(formFile, filePath)
	if err != nil || !ok {
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte("Error save image"))
		return
	}

	//
	//file, err := formFile.Open()
	//if err != nil {
	//	c.Writer.WriteHeader(400)
	//	c.Writer.Write([]byte("Error get file"))
	//	return
	//}
	//defer file.Close()
	//
	//reader := bytes.Buffer{}
	//reader.ReadFrom(file)
	//
	//c.Writer.Write(reader.Bytes())

}

func (h *UploadHandler) SaveBase64Json(c *gin.Context) {

}

func (h *UploadHandler) SaveLink(c *gin.Context) {

}
