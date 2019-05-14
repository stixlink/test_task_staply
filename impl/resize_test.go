package impl

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"
)

func TestNewImagickResizer(t *testing.T) {

	file, err := os.Open("../test_data/t.jpeg")
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	defer file.Close()

	thumbPath := "../test_data/thumbnail.jpg"
	fileOut, err := os.OpenFile(thumbPath, os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err, "Error create thumbnail file")
		return
	}
	defer os.Remove(thumbPath)

	resizer := NewImagickResizer()
	err = resizer.Resize(file, fileOut, 100, 100)
	if err != nil {
		fileOut.Close()
		t.Fatal(err, "Error resize input image file")
		return
	}
	fileOut.Close()

	thumb, err := os.Open(thumbPath)
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	defer thumb.Close()

	im, _, err := image.DecodeConfig(thumb)
	if err != nil {
		t.Fatal(err, "Error get info about resize image")
		return
	}
	if im.Height != 100 || im.Width != 100 {
		t.Fatalf("Thumbnail save with wrong size")
	}
}

func TestNewImagickResizer_FailOpenMainImage(t *testing.T) {

	fakePath := "../test_data/t2_empty.jpeg"
	fileFake, err := os.OpenFile(fakePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	defer os.Remove(fakePath)
	defer fileFake.Close()

	thumbPath := "../test_data/t2_fail.jpg"
	fileOut, err := os.OpenFile(thumbPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err, "Error create thumbnail fileFake")
		return
	}
	defer os.Remove(thumbPath)
	defer fileOut.Close()

	resizer := NewImagickResizer()
	err = resizer.Resize(fileFake, fileOut, 100, 100)
	if err == nil {
		t.Fatal(err, "Error resize empty fileFake")
		return
	}
}
