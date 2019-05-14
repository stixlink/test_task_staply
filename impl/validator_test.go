package impl

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestNewImagickResizer_WrongTypeFile(t *testing.T) {

	wrongExtPath := "../test_data/test.txt"
	fileFake, err := os.Open(wrongExtPath)
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	defer fileFake.Close()

	validator := NewFileValidator([]string{"png", "jpeg", "jpg"})
	buff := &bytes.Buffer{}
	if _, err := io.Copy(buff, fileFake); err != nil {
		t.Fatal("Error copy file to bytes.Buffer")
	}

	ok, err := validator.Validate(buff.Bytes())
	if err == nil || ok {
		t.Fatal("Missed wrong file")
		return
	}
}

func TestNewImagickResizer_WrongTypeImage(t *testing.T) {

	wrongExtPath := "../test_data/giphy.gif"
	fileFake, err := os.Open(wrongExtPath)
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	defer fileFake.Close()

	validator := NewFileValidator([]string{"png", "jpeg", "jpg"})
	buff := &bytes.Buffer{}
	if _, err := io.Copy(buff, fileFake); err != nil {
		t.Fatal("Error copy file to bytes.Buffer")
	}

	ok, err := validator.Validate(buff.Bytes())
	if err == nil || ok {
		t.Fatal("Missed wrong file")
		return
	}
}

func TestNewImagickResizer_WrongTypeImage2(t *testing.T) {

	wrongExtPath := "../test_data/t.jpeg"
	fileFake, err := os.Open(wrongExtPath)
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	defer fileFake.Close()

	validator := NewFileValidator([]string{"png"})
	buff := &bytes.Buffer{}
	if _, err := io.Copy(buff, fileFake); err != nil {
		t.Fatal("Error copy file to bytes.Buffer")
	}

	ok, err := validator.Validate(buff.Bytes())
	if err == nil || ok {
		t.Fatal("Missed wrong file")
		return
	}
}

func TestNewImagickResizer_Success(t *testing.T) {

	wrongExtPath := "../test_data/t.jpeg"
	fileFake, err := os.Open(wrongExtPath)
	if err != nil {
		t.Fatal("Error open test image")
		return
	}
	defer fileFake.Close()

	validator := NewFileValidator([]string{"jpeg"})
	buff := &bytes.Buffer{}
	if _, err := io.Copy(buff, fileFake); err != nil {
		t.Fatal("Error copy file to bytes.Buffer")
	}

	ok, err := validator.Validate(buff.Bytes())
	if err != nil || !ok {
		t.Fatal("Not allowed suitable file")
		return
	}
}
