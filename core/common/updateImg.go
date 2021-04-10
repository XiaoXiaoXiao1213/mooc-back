package common

import (
	"github.com/google/uuid"
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

func Upload(ctx iris.Context, key string) (string, error) {
	file, info, err := ctx.FormFile(key)
	if err != nil {
		log.Error(err)
		return "", err
	}

	dir := "/upload/" + time.Now().Format("20060102") + "/"
	err = os.MkdirAll("/var/www/"+dir, os.ModePerm)
	if err != nil {
		log.Error( err)
		return "", err
	}

	url := dir + uuid.New().String() + strings.ToLower(path.Ext(info.Filename))
	out, err := os.OpenFile("var/www/"+url, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error("ad", err)
		return "", err
	}

	defer file.Close()
	defer out.Close()
	io.Copy(out, file)

	return url, nil
}
