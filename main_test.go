package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestImageDownload_DownloadImage(t *testing.T) {

	file, err := os.Open("./test_data/test1.jpg")
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal("Error read test image")
		return
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
		return
	}))
	defer ts.Close()

	dw := NewDownloader(http.DefaultClient)
	u, _ := url.Parse(ts.URL)
	b, err := dw.DownloadImage(u)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != string(data) {
		t.Fatal("Error download")
	}
}

func TestImageDownload_DownloadImage_FailTimeout(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		return
	}))
	defer ts.Close()

	dw := NewDownloader(&http.Client{Timeout: 1 * time.Second})
	u, _ := url.Parse(ts.URL)
	_, err := dw.DownloadImage(u)
	switch err := err.(type) {
	case net.Error:
		if err.Timeout() {
			return
		}
	case *url.Error:
		if err, ok := err.Err.(net.Error); ok && err.Timeout() {
			return
		}
	default:
		t.Fatal("Fail setting timeout for download image")
	}

}

func TestNewImagickResizer(t *testing.T) {

}
