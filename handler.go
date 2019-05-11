package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stixlink/test_task_go_staply/iface"
	"github.com/stixlink/test_task_go_staply/utility"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	allowedContentType = []string{"image/jpeg", "image/pjpeg", "image/png"}
)

func NewHandler(imageSaveDir string, client *http.Client) *UploadHandler {
	return &UploadHandler{
		saveDir:        imageSaveDir,
		ImageValidator: NewFileValidator(),
		ImageResizer:   NewImagickResizer(),
		Downloader:     NewDownloader(client),
		Namer:          NewNamer(),
	}
}

type UploadHandler struct {
	saveDir string
	iface.ImageValidator
	iface.ImageResizer
	iface.Downloader
	iface.Namer
}

func (h *UploadHandler) SaveFormData(c *gin.Context) {
	response := NewResponse(c)
	// default max data size  32 << 20 // 32 MB
	// for change settings set router.MaxMultipartMemory before createserver
	// Example router.MaxMultipartMemory = N << 20  // N MiB
	formFile, err := c.FormFile("image")
	if err != nil {
		response.JSON(400, "Error get image", "")
		return
	}

	tmpFile, err := formFile.Open()
	if err != nil {
		response.JSON(500, "", "Internal error for get image")
		return
	}

	ext := strings.Trim(filepath.Ext(formFile.Filename), ".")
	_, filesPath := h.CreateName(h.saveDir, ext, []string{"100x100"})
	mainFilePath := filesPath[0]
	resizedFilePath := filesPath[1]

	file, err := os.OpenFile(mainFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		response.JSON(500, "", "Internal error in time open image for write")
		return
	}
	defer file.Close()

	buff := &bytes.Buffer{}
	_, err = io.Copy(buff, tmpFile)
	if err != nil {
		file.Close()
		response.JSON(500, "", "Internal error")
		return
	}

	ok, err := h.Validate(buff.Bytes())
	if err != nil || !ok {
		file.Close()
		errs := utility.RemoveErrorFiles(mainFilePath)
		if len(errs) > 0 {
			fmt.Println(fmt.Sprintf("Errors: %s", utility.JoinError(errs, " & ")))
		}
		response.JSON(400, "", "Invalid file type or data")
		return
	}

	_, err = file.Write(buff.Bytes())
	if err != nil {
		response.JSON(500, "", "Error write to main image")
		return
	}
	file.Close()

	err = h.resize(mainFilePath, resizedFilePath, 100, 100)
	if err != nil {
		errs := utility.RemoveErrorFiles([]string{mainFilePath, resizedFilePath}...)
		if len(errs) > 0 {
			fmt.Println(fmt.Sprintf("Errors: %s", utility.JoinError(errs, " & ")))
		}
		response.JSON(500, "", err.Error())
		return
	}

	response.JSON(200, "", "")
	return
}

func (h *UploadHandler) SaveLink(c *gin.Context) {

	response := NewResponse(c)
	imageUrl := c.Query("image")
	if imageUrl == "" {
		c.Writer.WriteHeader(400)
		c.Writer.Write([]byte("Invalid request"))
		return
	}

	urli, err := url.Parse(imageUrl)
	if err != nil {
		response.JSON(400, "", fmt.Sprintf("Invalid value parameter \"image\". \"%s\" is not valid url", imageUrl))
		return
	}

	ext := strings.Trim(filepath.Ext(imageUrl), ".")
	_, filesPath := h.CreateName(h.saveDir, ext, []string{"100x100"})
	mainFilePath := filesPath[0]
	resizedFilePath := filesPath[1]

	data, err := h.DownloadImage(urli)
	if err != nil {
		response.JSON(400, "", fmt.Sprintf("Error download image url: \"%s\"", imageUrl))
		return
	}
	ok, err := h.Validate(data)
	if err != nil || !ok {
		response.JSON(400, "", "Invalid file type or data")
		return
	}

	file, err := os.OpenFile(mainFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		response.JSON(500, "", "Error save image")
		return
	}

	_, err = file.Write(data)
	if err != nil {
		errs := utility.RemoveErrorFiles(mainFilePath)
		if len(errs) > 0 {
			fmt.Println(fmt.Sprintf("Errors: %s", utility.JoinError(errs, " & ")))
		}
		response.JSON(500, "", "Error save image")
		return
	}
	file.Sync()
	file.Close()

	err = h.resize(mainFilePath, resizedFilePath, 100, 100)
	if err != nil {
		errs := utility.RemoveErrorFiles([]string{mainFilePath, resizedFilePath}...)
		if len(errs) > 0 {
			fmt.Println(fmt.Sprintf("Errors: %s", utility.JoinError(errs, " & ")))
		}
		response.JSON(500, "", err.Error())
		return
	}

	response.JSON(200, "", "")
	return

}

type DataImageJSON struct {
	Data string `json:"data" binding:"required"`
}

func (h *UploadHandler) SaveBase64Json(c *gin.Context) {
	response := NewResponse(c)
	data := &DataImageJSON{}
	err := c.BindJSON(data)
	if err != nil {
		response.JSON(400, "", fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}
	buff, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		response.JSON(400, "", "Fail decode base64 data")
		return
	}
	ok, err := h.Validate(buff)
	if err != nil || !ok {
		response.JSON(400, "", "Invalid file type or data")
		return
	}
	fileType := http.DetectContentType(buff)
	t := strings.Split(fileType, "/")
	if len(t) < 2 {
		response.JSON(400, "", "Fail get file extension")
		return
	}

	_, filesPath := h.CreateName(h.saveDir, t[1], []string{"100x100"})
	mainFilePath := filesPath[0]
	resizedFilePath := filesPath[1]

	file, err := os.OpenFile(mainFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		response.JSON(500, "", "Error save image")
		return
	}

	_, err = file.Write(buff)
	if err != nil {
		errs := utility.RemoveErrorFiles(mainFilePath)
		if len(errs) > 0 {
			fmt.Println(fmt.Sprintf("Errors: %s", utility.JoinError(errs, " & ")))
		}
		response.JSON(500, "", "Error save image")
		return
	}
	file.Sync()
	file.Close()

	err = h.resize(mainFilePath, resizedFilePath, 100, 100)
	if err != nil {
		errs := utility.RemoveErrorFiles([]string{mainFilePath, resizedFilePath}...)
		if len(errs) > 0 {
			fmt.Println(fmt.Sprintf("Errors: %s", utility.JoinError(errs, " & ")))
		}
		response.JSON(500, "", err.Error())
		return
	}

	response.JSON(200, "", "")
	return

}

func (h *UploadHandler) resize(pathMainImage, pathResizeImage string, weight, height uint) error {

	file, err := os.Open(pathMainImage)
	if err != nil {
		return errors.New("Internal error")
	}
	defer file.Close()
	// create and open file for resize image
	resizedFile, err := os.OpenFile(pathResizeImage, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("Internal error")
	}
	defer resizedFile.Close()

	// resize and write resize image
	err = h.Resize(file, resizedFile, weight, height)
	if err != nil {
		return errors.New("Internal error")
	}

	return nil
}
