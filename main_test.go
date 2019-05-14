package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	main_mocks "github.com/stixlink/test_task_staply/iface/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestUploadHandler_SaveBase64Json2(t *testing.T) {
	var ignoreResponceBodyValue = "none"

	ctrl := gomock.NewController(t)
	validator := main_mocks.NewMockImageValidator(ctrl)
	resizer := main_mocks.NewMockImageResizer(ctrl)
	downloader := main_mocks.NewMockDownloader(ctrl)
	namer := main_mocks.NewMockNamer(ctrl)

	type TValidator struct {
		arg0       []byte
		returnArg0 bool
		returnArg1 error
	}
	type TResponse struct {
		status int
		body   string
	}
	type TRequest struct {
		data string
	}

	type TCase struct {
		validator TValidator
		response  TResponse
		request   TRequest
	}
	tables := []TCase{
		{
			validator: TValidator{
				arg0:       []byte("test"),
				returnArg0: true,
				returnArg1: nil,
			},
			response: TResponse{
				status: 400,
				body:   `{"data":"","error":"Fail get file extension"}`,
			},
			request: TRequest{
				data: "{\"data\":\"dGVzdA==\"}",
			},
		},
		{
			validator: TValidator{
				arg0:       []byte("test"),
				returnArg0: false,
				returnArg1: nil,
			},
			response: TResponse{
				status: 400,
				body:   `{"data":"","error":"Invalid file type or data"}`,
			},
			request: TRequest{
				data: "{\"data\":\"dGVzdA==\"}",
			},
		},
		{
			validator: TValidator{
				arg0:       []byte("test"),
				returnArg0: true,
				returnArg1: nil,
			},
			response: TResponse{
				status: 400,
				body:   `{"data":"","error":"Fail decode base64 data"}`,
			},
			request: TRequest{
				data: "{\"data\":\"111asd\"}",
			},
		},
		{
			validator: TValidator{
				arg0:       []byte("test"),
				returnArg0: true,
				returnArg1: nil,
			},
			response: TResponse{
				status: 400,
				body:   ignoreResponceBodyValue,
			},
			request: TRequest{
				data: "{}",
			},
		},
	}
	gin.SetMode(gin.TestMode)
	for _, c := range tables {

		reader := strings.NewReader(c.request.data)

		// prepare validator
		validator.EXPECT().Validate(c.validator.arg0).Return(c.validator.returnArg0, c.validator.returnArg1)

		h := UploadHandler{
			saveDir:        "./test_data/upload",
			ImageValidator: validator,
			ImageResizer:   resizer,
			Downloader:     downloader,
			Namer:          namer,
		}
		router := gin.New()
		router.POST("/json", h.SaveBase64Json)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/json", reader)
		router.ServeHTTP(w, req)
		assert.Equal(t, c.response.status, w.Code)
		if c.response.body != ignoreResponceBodyValue {
			assert.Equal(t, c.response.body, w.Body.String())
		}
	}
}

func TestUploadHandler_SaveLink(t *testing.T) {
	var ignoreResponceBodyValue = "none"

	ctrl := gomock.NewController(t)
	validator := main_mocks.NewMockImageValidator(ctrl)
	resizer := main_mocks.NewMockImageResizer(ctrl)
	downloader := main_mocks.NewMockDownloader(ctrl)
	namer := main_mocks.NewMockNamer(ctrl)

	basePath := "./test_data/upload"
	os.MkdirAll(basePath+"/", 0777)
	defer os.Remove(basePath + "/")

	ext := "jpg"
	prefix := []string{"100x100"}
	baseName := "base_name"
	pathsCheck := append([]string{}, fmt.Sprintf("%s/%s.%s", basePath, baseName, ext))
	for _, p := range prefix {
		pathsCheck = append(pathsCheck, fmt.Sprintf("%s/%s_%s.%s", basePath, p, baseName, ext))
	}

	type TDownloader struct {
		arg0       *url.URL
		returnArg0 []byte
		returnArg1 error
	}
	type TValidator struct {
		arg0       []byte
		returnArg0 bool
		returnArg1 error
	}
	type TResponse struct {
		status int
		body   string
	}
	type TRequest struct {
		data     string
		urlQuery string
	}
	type TCase struct {
		downloader TDownloader
		validator  TValidator
		response   TResponse
		request    TRequest
	}
	tables := []TCase{
		{
			downloader: TDownloader{
				arg0:       &url.URL{Host: "localhost:8080", Path: "/image/1.jpg", Scheme: "http"},
				returnArg0: []byte("test"),
				returnArg1: nil,
			},
			validator: TValidator{
				arg0:       []byte("test"),
				returnArg0: false,
				returnArg1: nil,
			},
			response: TResponse{
				status: 400,
				body:   `none`,
			},
			request: TRequest{
				data:     "",
				urlQuery: "image=http://localhost:8080/image/1.jpg",
			},
		},
		{
			downloader: TDownloader{
				arg0:       &url.URL{Host: "localhost:8080", Path: "/image/1.jpg", Scheme: "http"},
				returnArg0: []byte(""),
				returnArg1: errors.New("test err downloader"),
			},
			validator: TValidator{
				arg0:       []byte("test"),
				returnArg0: false,
				returnArg1: nil,
			},
			response: TResponse{
				status: 400,
				body:   `none`,
			},
			request: TRequest{
				data:     "",
				urlQuery: "image=http://localhost:8080/image/1.jpg",
			},
		},
		{
			downloader: TDownloader{
				arg0:       &url.URL{Host: "", Path: "image2.jpg", Scheme: ""},
				returnArg0: []byte("test"),
				returnArg1: errors.New("test err downloader"),
			},
			response: TResponse{
				status: 400,
				body:   `none`,
			},
			request: TRequest{
				data:     "",
				urlQuery: "image=image2.jpg",
			},
		},
		{
			response: TResponse{
				status: 400,
				body:   `none`,
			},
			request: TRequest{
				data:     "",
				urlQuery: "",
			},
		},
	}
	gin.SetMode(gin.TestMode)
	for _, c := range tables {

		reader := strings.NewReader(c.request.data)

		namer.EXPECT().CreateName(basePath, ext, prefix).Return(baseName, pathsCheck)
		downloader.EXPECT().DownloadImage(c.downloader.arg0).Return(c.downloader.returnArg0, c.downloader.returnArg1)
		validator.EXPECT().Validate(c.validator.arg0).Return(c.validator.returnArg0, c.validator.returnArg1)

		h := UploadHandler{
			saveDir:        basePath,
			ImageValidator: validator,
			ImageResizer:   resizer,
			Downloader:     downloader,
			Namer:          namer,
		}
		router := gin.New()
		router.GET("/link", h.SaveLink)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/link?"+c.request.urlQuery, reader)
		router.ServeHTTP(w, req)
		assert.Equal(t, c.response.status, w.Code)
		if c.response.body != ignoreResponceBodyValue {
			assert.Equal(t, c.response.body, w.Body.String())
		}
	}
}
