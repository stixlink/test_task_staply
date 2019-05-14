package impl

import (
	"github.com/stixlink/test_task_staply/iface"
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

	file, err := os.Open("../test_data/t.jpeg")
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

func TestImageDownload_DownloadImage_NotFound(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		return
	}))
	defer ts.Close()

	dw := NewDownloader(&http.Client{Timeout: 10 * time.Second})
	u, _ := url.Parse(ts.URL)
	_, err := dw.DownloadImage(u)
	switch err.(type) {
	case iface.ErrorDownloadNotFound:

	default:
		t.Fatal("Wrong error for is case")
	}

}

func TestImageDownload_DownloadImage_WrongStatus(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(444)
		return
	}))
	defer ts.Close()

	dw := NewDownloader(&http.Client{Timeout: 10 * time.Second})
	u, _ := url.Parse(ts.URL)
	_, err := dw.DownloadImage(u)
	switch err.(type) {
	case iface.ErrorDownloadWrongResponseStatus:

	default:
		t.Fatal("Wrong error for is case")
	}

}
