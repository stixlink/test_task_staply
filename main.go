package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v2/imagick"
	"net/http"
	"os"
)

var saveDir *string

func main() {
	err := Init()
	if err != nil {
		panic(errors.Wrap(err, "Error initialization"))
	}
	h := NewHandler(*saveDir)

	r := gin.Default()
	r.POST("/form", h.SaveFormData)
	r.POST("/base64", h.SaveBase64Json)
	r.POST("/link", h.SaveLink)
	http.ListenAndServe(":58001", r)
}

func Init() error {

	imagick.Initialize()
	defer imagick.Terminate()

	saveDir = flag.String("save-dir", "./", "folder for save images")
	fmt.Printf("You selected \"%s\" directory for save images", *saveDir)
	if _, err := os.Stat(*saveDir); os.IsNotExist(err) {
		return errors.Wrap(err, "The selected directory does not exist")
	}

	return nil
}
