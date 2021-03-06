package impl

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stixlink/test_task_staply/iface"
	"gopkg.in/gographics/imagick.v2/imagick"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

func NewDownloader(client *http.Client) *ImageDownload {
	return &ImageDownload{
		client: client,
	}
}

type ImageDownload struct {
	client *http.Client
}

// DownloadImage  downloades image by url and return byte slice or error if than happened fail
func (d *ImageDownload) DownloadImage(imageUrl *url.URL) ([]byte, error) {

	resp, err := d.client.Get(imageUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:

	case http.StatusNotFound:
		return nil, iface.ErrorDownloadNotFound(errors.New(fmt.Sprintf("not found image for url: \"%s\"", imageUrl.String())))
	default:
		return nil, iface.ErrorDownloadWrongResponseStatus(fmt.Errorf("wrong response code: %d", resp.StatusCode))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// NewResponse creating response struct for write response data
func NewResponse(c *gin.Context) *Response {
	return &Response{c: c}
}

type Response struct {
	Data  interface{} `json:"data, omitempty"`
	Error string      `json:"error"`
	c     *gin.Context
}

// JSON write response data in json format
func (r *Response) JSON(code int, data interface{}, errMessage string) {
	r.Error = errMessage
	r.Data = data
	r.c.JSON(code, r)
}

func NewNamer() *NameCreator {
	return &NameCreator{
	}
}

type NameCreator struct {
}

// CreateName creates name and names with all passed prefixes.
// returns base name and slice with all paths
func (n *NameCreator) CreateName(basePath string, ext string, prefix []string) (baseName string, paths []string) {
	baseName = fmt.Sprintf("%v_%v", time.Now().UnixNano(), rand.Int63n(100))
	basePath = path.Clean(basePath)
	paths = append(paths, fmt.Sprintf("%s/%s.%s", basePath, baseName, ext))

	for _, v := range prefix {
		paths = append(paths, fmt.Sprintf("%s/%s_%s.%s", basePath, v, baseName, ext))
	}

	return
}

func NewFileValidator(allowedExtension []string) *ImageDataValidator {
	return &ImageDataValidator{AllowedExtension: allowedExtension}
}

type ImageDataValidator struct {
	AllowedExtension []string
}

// Validate check content type in data
// ATTENTION! if add new extension notice the validator and import library "image/*"
// and see usage function https://golang.org/pkg/image/#DecodeConfig
func (h *ImageDataValidator) Validate(data []byte) (result bool, err error) {

	reader := bytes.NewReader(data)
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		return
	}

	for _, t := range h.AllowedExtension {
		if t == format {
			result = true
		}
	}

	if !result {
		err = errors.New(fmt.Sprintf("Not allowed content type \"%s\" ", format))
	}

	return
}

func NewImagickResizer() *ImagickResizer {
	return &ImagickResizer{}
}

type ImagickResizer struct {
}

// Resize input image file and save to output file
func (h *ImagickResizer) Resize(inputFile *os.File, outputFile *os.File, weight, height uint) error {

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageFile(inputFile)
	if err != nil {

		return errors.Wrap(err, "error read image")
	}
	// TODO: use resize more smart
	err = mw.ResizeImage(weight, height, imagick.FILTER_UNDEFINED, 1)
	if err != nil {

		return errors.Wrap(err, "error resize image")
	}

	return mw.WriteImageFile(outputFile)
}
