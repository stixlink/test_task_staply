package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v2/imagick"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ImageResizer interface {
	Resize(inputFile *os.File, outFile *os.File) error
}

func NewImagickResizer() *ImagickResizer {
	return &ImagickResizer{}
}

type ImagickResizer struct {
}

func (h *ImagickResizer) Resize(inputFile *os.File, outFile *os.File) error {

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageFile(inputFile)
	if err != nil {

		return errors.Wrap(err, "error read image")
	}

	err = mw.ResizeImage(100, 100, imagick.FILTER_UNDEFINED, 1)
	if err != nil {

		return errors.Wrap(err, "error resize image")
	}

	return mw.WriteImageFile(outFile)
}

type ImageValidator interface {
	Validate() (bool, error)
}

func NewFormDataValidator(fileHeader *multipart.FileHeader) *ImageFormDataValidator {
	return &ImageFormDataValidator{FileHeader: fileHeader}
}

type ImageFormDataValidator struct {
	FileHeader *multipart.FileHeader
}

func (h *ImageFormDataValidator) Validate() (result bool, err error) {

	ext := filepath.Ext(h.FileHeader.Filename)
	if ext == "" {
		err = errors.New(fmt.Sprintf("Unknown file extension \"%s\" ", ext))
		return
	}
	for _, t := range allowedFileExtension {
		if t == ext {
			result = true
		}
	}

	if !result {
		err = errors.New(fmt.Sprintf("Not allowed file extension \"%s\" ", ext))
		return
	}

	ct := h.FileHeader.Header.Get("Content-type")
	if ct == "" {
		err = errors.New(fmt.Sprintf("Unknown content type \"%s\" ", ct))

		return
	}
	for _, t := range allowedFileExtension {
		if t == ct {
			result = true
		}
	}
	if !result {
		err = errors.New(fmt.Sprintf("Not allowed content type \"%s\" ", ct))
	}

	return
}
