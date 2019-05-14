package iface

//go:generate mockgen -package main_mocks -destination mocks/iface.go github.com/stixlink/test_task_staply/iface ImageResizer,ImageValidator,Namer,Downloader

import (
	"net/url"
	"os"
)

type
(
	ErrorDownloadNotFound = error
	ErrorDownloadWrongResponseStatus = error
)

type ImageResizer interface {
	Resize(inputFile *os.File, outFile *os.File, weight, height uint) error
}

type ImageValidator interface {
	Validate(data []byte) (result bool, err error)
}

type Namer interface {
	CreateName(basePath string, ext string, prefix []string) (baseName string, paths []string)
}

type Downloader interface {
	DownloadImage(imageUrl *url.URL) ([]byte, error)
}
