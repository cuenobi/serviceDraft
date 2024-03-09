package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
)

func UrlEncoded(originalText string) string {
	encodedText := url.QueryEscape(originalText)
	return encodedText
}

type Data struct {
	Name  string
	Value string
}

func QRCodeGenerate(storePath, url string, d []*Data) (string, error) {
	dataCount := len(d)
	if dataCount > 0 {
		url = url + "?"
	}

	// string for convert to qrcode
	for i, data := range d {
		data.Value = UrlEncoded(data.Value)
		if i+1 == dataCount {
			url = url + fmt.Sprintf("%s=%s", data.Name, data.Value)
		} else {
			url = url + fmt.Sprintf("%s=%s&", data.Name, data.Value)
		}
		i++
	}

	// Validate dir is defined
	if _, err := os.Stat("." + storePath); os.IsNotExist(err) {
		// If undefined, create it
		if err := os.MkdirAll("."+storePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// path to save qrcode file
	now := time.Now()
	storePath = fmt.Sprintf(storePath+"%v-qr.png", now.Unix())

	// generate qrcode
	err := qrcode.WriteFile(url, qrcode.Medium, 256, "."+storePath)
	if err != nil {
		return "", err
	}

	url = viper.GetString(`url.qrcode`)
	return url + storePath, nil
}

func QRCodeWithOutUrl(uid string) (string, error) {
	storePath := viper.GetString(`file.store`) // path to save qrcode file
	// Validate dir is defined
	if _, err := os.Stat("." + storePath); os.IsNotExist(err) {
		// If undefined, create it
		if err := os.MkdirAll("."+storePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// path to save qrcode file
	now := time.Now()
	id := uuid.New()
	storePath = fmt.Sprintf(storePath+"%v-%-qr.png", now.Unix(), id)

	// generate qrcode
	err := qrcode.WriteFile(uid, qrcode.Medium, 256, "."+storePath)
	if err != nil {
		return "", err
	}

	url := viper.GetString(`url.qrcode`)
	return url + storePath, nil
}
