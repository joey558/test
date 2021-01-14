package common

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/tuotoo/qrcode"
)

func writePng(filename string, img image.Image) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(file, img)
	// err = jpeg.Encode(file, img, &jpeg.Options{100})      //图像质量值为100，是最好的图像显示
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	log.Println(file.Name())
}

/**
* 生成二维码
 */
func CreateQrCode(base64 string) (image.Image, error) {
	code, err := qr.Encode(base64, qr.L, qr.Unicode)
	// code, err := code39.Encode(base64)
	if err != nil {
		return nil, err
	}

	if base64 != code.Content() {
		return nil, err
	}

	code, err = barcode.Scale(code, 300, 300)
	if err != nil {
		return nil, err
	}

	return code, nil
}

/**
* 读取二维码图片内容
 */
func ReadQrCode(img_file string) (string, error) {
	fi, err := os.Open(img_file)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	qrmatrix, err := qrcode.Decode(fi)
	if err != nil {
		return "", err
	}

	return qrmatrix.Content, nil
}
