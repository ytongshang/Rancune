package util

import (
	"os/exec"

	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

func AndroidImageCapture(filePath string) (image image.Image, err error) {
	err = exec.Command("adb", "shell", "rm -f /sdcard/screenshot*.png").Run()
	if err != nil {
		return
	}
	err = exec.Command("adb", "shell", "screencap", "-p", "/sdcard/screenshot.png").Run()
	if err != nil {
		return
	}
	err = exec.Command("adb", "shell", "mv /sdcard/screenshot*.png /sdcard/screenshot.png").Run()
	if err != nil {
		return
	}
	err = exec.Command("adb", "pull", "/sdcard/screenshot.png", filePath).Run()
	if err != nil {
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func ImageCut(src image.Image, left, top, right, bottom int) (image.Image, error) {
	var subImg image.Image

	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(left, top, right, bottom)).(*image.YCbCr)
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(left, top, right, bottom)).(*image.RGBA)
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(left, top, right, bottom)).(*image.NRGBA)
	} else {
		return subImg, fmt.Errorf("%f", "图片裁剪失败")
	}

	return subImg, nil
}

func SavePNGFile(filename string, pic image.Image) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, pic)
}

func ImageToBase64(path string) (string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(f), nil
}
