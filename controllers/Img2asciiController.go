package controllers

import (
	"aaimg2ascii/datamodels"
	"aaimg2ascii/services"
	"aaimg2ascii/validators"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Img2asciiController struct {
	Img2asciiService services.Img2asciiService
}

var img2asciiController *Img2asciiController
var img2asciiControllerOnce sync.Once

func GetImg2asciiController() *Img2asciiController {
	img2asciiControllerOnce.Do(func() {
		img2asciiController = &Img2asciiController{
			Img2asciiService: services.NewImg2asciiService(),
		}
	})
	return img2asciiController
}

//region 上傳圖片
func (c *Img2asciiController) Img2ascii(ctx *gin.Context) {

	//region 資料驗證
	var img2ascii *validators.Img2ascii
	ctx.ShouldBind(&img2ascii)
	if err := validators.GlobalValidator.Check(img2ascii); err != nil {
		ctx.JSON(http.StatusBadRequest, strings.Split(err.Error(), "|"))
		return
	}
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
	newImgresult, status := c.Img2asciiService.UploadImg(file, newFileName)

	if !status {
		ctx.JSON(http.StatusBadRequest, newImgresult)
	}

	img2asciiresult, status := c.Img2asciiService.Img2ascii(newFileName, img2ascii, userData)
	if status {
		ctx.JSON(http.StatusOK, img2asciiresult)
	} else {
		ctx.JSON(http.StatusBadRequest, img2asciiresult)
	}

}

//endregion

func (c *Img2asciiController) UpdateHot(ctx *gin.Context) {
	//region 資料驗證
	var updateHot *validators.UpdateHot
	ctx.ShouldBind(&updateHot)
	if err := validators.GlobalValidator.Check(updateHot); err != nil {
		ctx.JSON(http.StatusBadRequest, strings.Split(err.Error(), "|"))
		return
	}
	userData := ctx.MustGet("user").(*datamodels.User)
	//endregion

	result := c.Img2asciiService.UpdateHot(updateHot, userData)

	if result != "" {
		ctx.JSON(http.StatusBadRequest, result)
	} else {
		ctx.JSON(http.StatusOK, "已修改熱度")
	}
}

func (c *Img2asciiController) GetHotTop(ctx *gin.Context) {
	result := c.Img2asciiService.GetHotTop()
	ctx.JSON(http.StatusOK, result)
}

func (c *Img2asciiController) GetUserAsciiById(ctx *gin.Context) {
	//region 資料驗證
	var userAscii *validators.UserAscii
	ctx.ShouldBind(&userAscii)
	if err := validators.GlobalValidator.Check(userAscii); err != nil {
		ctx.JSON(http.StatusBadRequest, strings.Split(err.Error(), "|"))
		return
	}
	//endregion
	result := c.Img2asciiService.GetUserAsciiById(userAscii)
	ctx.JSON(http.StatusOK, result)
}

func (c *Img2asciiController) GetMyAscii(ctx *gin.Context) {
	//region 資料驗證
	userData := ctx.MustGet("user").(*datamodels.User)
	//endregion
	result := c.Img2asciiService.GetMyAscii(userData)
	ctx.JSON(http.StatusOK, result)
}
