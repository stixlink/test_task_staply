package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v2/imagick"
	"net/http"
	"os"
	"time"
)

var (
	saveDir    *string
	httpClient *http.Client
)

func main() {

	err := Init()
	if err != nil {
		panic(errors.Wrap(err, "Error initialization"))
	}

	h := NewHandler(*saveDir, httpClient)

	r := gin.Default()
	r.POST("/form", h.SaveFormData)
	r.POST("/json", h.SaveBase64Json)
	r.GET("/link", h.SaveLink)
	r.Run(":58001")

}

func Init() error {

	imagick.Initialize()
	defer imagick.Terminate()

	saveDir = flag.String("save-dir", "./upload", "folder for save images")
	fmt.Printf("You selected \"%s\" directory for save images\n", *saveDir)
	if _, err := os.Stat(*saveDir); os.IsNotExist(err) {
		return errors.Wrap(err, "The selected directory does not exist")
	}

	// create http client for usage in app
	httpClient = &http.Client{Timeout: 60 * time.Second}

	return nil
}
