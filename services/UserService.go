package services

import (
	"aaimg2ascii/datamodels"
	"aaimg2ascii/models"
	"aaimg2ascii/validators"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
)

type UserService interface {
	Register(register *validators.Register) (string, uint32)
	Login(login *validators.Login) (string, uint32)
	GetUserData(userId uint32) *datamodels.User
	UpdatePhoto(file multipart.File, newFileName string, userData *datamodels.User) (string, bool)
}

func NewUserService() UserService {
	return &userService{
		db:          models.DB.Mysql,
		baseService: NewBaseService(),
	}

}

type userService struct {
	db          *gorm.DB
	baseService BaseService
}

func (s *userService) Register(register *validators.Register) (string, uint32) {
	var user models.User
	if err := s.db.First(&user, "email = ?", register.Email).Error; err == nil {
		return "此電子郵件已存在", 0
	}
	if newPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost); err == nil {
		newUser := models.User{Account: register.Account, Password: string(newPassword),
			Email: register.Email}
		s.db.Create(&newUser)
		return "", newUser.UserId
	}
	return "密碼加密失敗", 0

}

func (s *userService) Login(login *validators.Login) (string, uint32) {
	var user models.User
	if err := s.db.First(&user, "account = ?", login.Account).Error; err != nil {
		return "查無此帳號，請先註冊", 0
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) == nil {
		return "", user.UserId
	}
	return "密碼錯誤", 0
}

func (s *userService) GetUserData(userId uint32) *datamodels.User {
	var user datamodels.User
	if err := s.db.Select("user_id, account, email, photo,created_at").First(&user, userId).Error; err != nil {
		println(user.UserId)
		return nil
	}
	return &user
}
func (s *userService) UpdatePhoto(file multipart.File, newFileName string, userData *datamodels.User) (string, bool) {
	saveResult, status := s.baseService.SaveImg(file, newFileName)
	if status == false {
		return saveResult, false
	}
	userData.Photo = saveResult

	if err := s.db.Save(&userData).Error; err != nil {
		return "更新頭貼失敗", false
	}
	return userData.Photo, true
}
