package image

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func BytesFromBase64(base64string string) ([]byte, error) {
	arr := strings.Split(base64string, ",")
	if len(arr) != 2 {
		return nil, errors.New("invalid base64 string")
	}

	return base64.StdEncoding.DecodeString(arr[1])
}

func Validate(data []byte) error {
	if len(data) > 20*1024*1024 {
		return errors.New("uploaded image size is too big! (Maximum 20 Mb)")
	}

	mimeType := getMimeType(data)

	regex := regexp.MustCompile(`^image/(jpeg|png|gif)$`)
	if !regex.Match([]byte(mimeType)) {
		return errors.New("only jpeg, png and gif images can be uploaded")
	}

	return nil
}

func GetExtension(data []byte) string {
	var ext string
	mimeType := getMimeType(data)

	switch strings.Split(mimeType, "/")[1] {
	case "png":
		ext = ".png"
	case "jpeg":
		ext = ".jpg"
	case "gif":
		ext = ".gif"
	}
	return ext
}

func getMimeType(data []byte) string {
	return http.DetectContentType(data)
}

func Save(data []byte, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return ioutil.WriteFile(name, data, 0644)
}

func ReadImage(name string) (string, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	mimeType := getMimeType(data)
	base64string := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64string), nil
}
