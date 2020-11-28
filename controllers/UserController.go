package controllers

import (
	"aaimg2ascii/datamodels"
	"aaimg2ascii/services"
	"aaimg2ascii/validators"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"sync"
	"time"
)

type UserController struct {
	UserService services.UserService
}

var userController *UserController
var userControllerOnce sync.Once

func GetUserController() *UserController {
	userControllerOnce.Do(func() {
		userController = &UserController{
			UserService: services.NewUserService(),
		}
	})
	return userController
}

//region 註冊
func (c *UserController) Register(ctx *gin.Context) {
	var register validators.Register
	ctx.ShouldBind(&register)
	if err := validators.GlobalValidator.Check(register); err != nil {
		ctx.JSON(http.StatusBadRequest, strings.Split(err.Error(), "|"))
		return
	}
	//match, _ := regexp.MatchString("^([A-Za-z0-9]+[0-9]+)[A-Za-z0-9]*$", register.Password)
	//if !match {
	//	ctx.JSON(http.StatusBadRequest, []string{"密碼需包含英文數字"})
	//	return
	//}
	result, userId := c.UserService.Register(&register)
	if result == "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "aaimg2ascii/register",
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
			//"jti": "9527",
			"sub": userId,
		})
		secret := viper.GetString("auth.secret")
		tokenString, _ := token.SignedString([]byte(secret))
		ctx.JSON(http.StatusOK, tokenString)
	} else {
		ctx.JSON(http.StatusBadRequest, []string{result})
	}

}

//endregion

//region 登入
func (c *UserController) Login(ctx *gin.Context) {
	var login validators.Login
	ctx.ShouldBind(&login)
	if err := validators.GlobalValidator.Check(login); err != nil {
		ctx.JSON(http.StatusBadRequest, strings.Split(err.Error(), "|"))
		return
	}
	result, userId := c.UserService.Login(&login)
	if result == "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "aaimg2ascii/login",
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
			"sub": userId,
		})
		secret := viper.GetString("auth.secret")
		tokenString, _ := token.SignedString([]byte(secret))
		ctx.JSON(http.StatusOK, tokenString)
	} else {
		ctx.JSON(http.StatusBadRequest, []string{result})
	}

}

//endregion

//region 個人資料
func (c *UserController) GetUserData(ctx *gin.Context) {
	userData := ctx.MustGet("user").(*datamodels.User)
	ctx.JSON(http.StatusOK, userData)
}

//endregion

//region 上傳大頭貼
func (c *UserController) UpdatePhoto(ctx *gin.Context) {

	//region 資料驗證
	userData := ctx.MustGet("user").(*datamodels.User)
	//endregion

	//region 圖片驗證
	var maxSize int64 = 8 * 1024 * 1024 * 1024
	info, err := ctx.FormFile("UploadImg")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, []string{err.Error()})
		return
	}
	file, err := info.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, []string{err.Error()})
		return
	}
	if file == nil {
		ctx.JSON(http.StatusBadRequest, []string{"上傳失敗，請重新再試"})
		return
	}
	defer file.Close()
	fileType := info.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" {
		ctx.JSON(http.StatusBadRequest, []string{"請上傳jpeg、png 圖片檔"})
		return
	}
	if info.Size > maxSize {
		ctx.JSON(http.StatusBadRequest, []string{"檔案不可超過8MB"})
		return
	}
	//endregion

	newFileName := time.Now().Format("2006-01-02") + "_" + info.Filename

	//圖片儲存
	newImgresult, status := c.UserService.UpdatePhoto(file, newFileName, userData)

	if status {
		ctx.JSON(http.StatusOK, newImgresult)
	} else {
		ctx.JSON(http.StatusBadRequest, newImgresult)
	}
}

//region 測試
func (c *UserController) Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello")
}
