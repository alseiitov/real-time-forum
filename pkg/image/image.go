package image

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func ParseFromRequest(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	mimeTypes := make(map[string]string)
	mimeTypes["jpg"] = "image/jpeg"
	mimeTypes["jpeg"] = "image/jpeg"
	mimeTypes["jpe"] = "image/jpeg"
	mimeTypes["png"] = "image/png"
	mimeTypes["gid"] = "image/gif"

	file, stat, err := r.FormFile("image")
	if err != nil {
		return nil, nil, err
	}

	if stat.Size > 20*1024*1024 {
		return nil, nil, errors.New("uploaded image size is too big! (Maximum 20 Mb)")
	}

	reg := regexp.MustCompile(`\.(jpg|jpeg|jpe|png|gif)$`)
	if !reg.MatchString(stat.Filename) {
		return nil, nil, errors.New("only jpg, jpeg, jpe, png, gif files are allowed")
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return nil, nil, errors.New("can't read your image, please try again")
	}

	ext := getExtension(stat.Filename)
	mime := http.DetectContentType(buff)
	if mimeTypes[ext] != mime {
		return nil, nil, errors.New("invalid or not allowed file extension")
	}

	return file, stat, nil
}

func SaveImage(file multipart.File, stat *multipart.FileHeader, path string) (string, error) {
	ext := getExtension(stat.Filename)

	name := fmt.Sprintf("%s.%s", uuid.NewV4().String(), ext)
	filePath := filepath.Join(path, name)

	file.Seek(0, 0)

	img, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(img, file)
	if err != nil {
		return "", err
	}

	return name, nil
}

func getExtension(fileName string) string {
	arr := strings.Split(fileName, ".")
	return arr[len(arr)-1]
}
