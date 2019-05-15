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
	saveDir = flag.String("save-dir", "./upload", "folder for save images")
	httpClient *http.Client
	port       string

	// ATTENTION! if add new extension notice the validator and import library "image/*"
	allowedExtensionsFile = []string{"jpeg", "jpg", "png"}
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
	r.Run(":" + port)

}

func Init() error {
	imagick.Initialize()
	defer imagick.Terminate()

	flag.Parse()
	fmt.Printf("You selected \"%s\" directory for save images\n", *saveDir)
	if _, err := os.Stat(*saveDir); os.IsNotExist(err) {
		return errors.Wrap(err, "The selected directory does not exist")
	}

	// create http client for usage in app
	httpClient = &http.Client{Timeout: 60 * time.Second}
	port = os.Getenv("MY_SERVER_PORT")
	if port == "" {
		port = "58001"
	}

	return nil
}
