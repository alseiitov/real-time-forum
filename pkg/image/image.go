package image

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func SaveAndGetName(base64string, path string) (string, error) {
	if base64string == "" {
		return "", nil
	}

	imageBytes, err := bytesFromBase64(base64string)
	if err != nil {
		return "", err
	}

	err = validate(imageBytes)
	if err != nil {
		return "", err
	}

	newImageName := uuid.NewV4().String() + getExtension(imageBytes)

	err = save(imageBytes, path, newImageName)
	if err != nil {
		return "", err
	}

	return newImageName, nil
}

func ReadImage(path, name string) (string, error) {
	if name == "" {
		return "", nil
	}

	data, err := ioutil.ReadFile(filepath.Join(path, name))
	if err != nil {
		return "", err
	}

	mimeType := getMimeType(data)
	base64string := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64string), nil
}

func bytesFromBase64(base64string string) ([]byte, error) {
	arr := strings.Split(base64string, ",")
	if len(arr) != 2 {
		return nil, errors.New("invalid base64 string")
	}

	return base64.StdEncoding.DecodeString(arr[1])
}

func validate(data []byte) error {
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

func getExtension(data []byte) string {
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

func save(data []byte, path, name string) error {
	err := checkDirExistance(path)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(path, name))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func checkDirExistance(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
