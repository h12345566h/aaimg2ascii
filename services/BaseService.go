package services

import (
	"aaimg2ascii/models"
	"github.com/jinzhu/gorm"
	"io"
	"mime/multipart"
	"os"
)

type BaseService interface {
	SaveImg(file multipart.File, newFileName string) (string, bool)
}

func NewBaseService() BaseService {
	return &baseService{
		db: models.DB.Mysql,
	}
}

type baseService struct {
	db *gorm.DB
}

func (s *baseService) SaveImg(file multipart.File, newFileName string) (string, bool) {
	newfilepath := "public/img/" + newFileName
	newfile, err := os.Create(newfilepath)
	if err != nil {
		return "新建圖片失敗", false
	}
	defer newfile.Close()

	_, err = io.Copy(newfile, file)
	if err != nil {
		return "圖片儲存失敗", false
	}
	println(newFileName)
	return newFileName, true
}
