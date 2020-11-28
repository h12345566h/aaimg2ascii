package route

import (
	"aaimg2ascii/controllers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func InitRouter(app *gin.Engine) {
	apiRoute := app.Group("/api")
	{
		apiRoute.GET("/hello", controllers.GetUserController().Hello)
		apiRoute.POST("/register", controllers.GetUserController().Register)
		apiRoute.POST("/login", controllers.GetUserController().Login)
		apiRoute.GET("/getHotTop", controllers.GetImg2asciiController().GetHotTop)

		jwtRoute := apiRoute.Group("/").Use(JwtMiddleware)
		{
			jwtRoute.GET("/getUserData", controllers.GetUserController().GetUserData)
			jwtRoute.POST("/updatePhoto", controllers.GetUserController().UpdatePhoto)
			jwtRoute.POST("/img2ascii", controllers.GetImg2asciiController().Img2ascii)
			jwtRoute.POST("/updateHot", controllers.GetImg2asciiController().UpdateHot)
			jwtRoute.GET("/getMyAscii", controllers.GetImg2asciiController().GetMyAscii)
			jwtRoute.POST("/getUserAsciiById", controllers.GetImg2asciiController().GetUserAsciiById)

		}
		//jwt2Route := apiRoute.Group("/").Use(Jwt2Middleware)
		//{
		//}

	}
}

//jwt驗證未通過則返回
func JwtMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, [1]string{"未登入"})
		ctx.Abort()
		return
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		ctx.JSON(http.StatusBadRequest, [1]string{"未登入"})
		ctx.Abort()
		return
	}

	token := authHeaderParts[1]

	parsedToken, err := new(jwt.Parser).Parse(token, ValidationKeyGetter)
	// Check if there was an error in parsing...
	if err != nil {
		ctx.JSON(http.StatusBadRequest, [1]string{"解析權杖失敗"})
		ctx.Abort()
		return
	}

	if !parsedToken.Valid {
		ctx.JSON(http.StatusBadRequest, [1]string{"權杖不合法"})
		ctx.Abort()
		return
	}
	userId := cast.ToUint32(parsedToken.Claims.(jwt.MapClaims)["sub"])
	user := controllers.GetUserController().UserService.GetUserData(userId)
	if user == nil {
		ctx.JSON(http.StatusBadRequest, [1]string{"請重新登入！"})
		ctx.Abort()
		return
	}
	//if user.Stop {
	//	ctx.JSON(http.StatusBadRequest, [1]string{"已被停權！"})
	//	return
	//}
	ctx.Set("user", user)
}

//jwt驗證未通過不返回
func Jwt2Middleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	defer ctx.Next()
	if authHeader == "" {
		return
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return
	}

	token := authHeaderParts[1]

	parsedToken, err := new(jwt.Parser).Parse(token, ValidationKeyGetter)
	if err != nil {
		return
	}

	if !parsedToken.Valid {
		return
	}
	userId := cast.ToUint32(parsedToken.Claims.(jwt.MapClaims)["sub"])
	user := controllers.GetUserController().UserService.GetUserData(userId)
	if user != nil {
		ctx.Set("user", user)
	}
}
func ValidationKeyGetter(token *jwt.Token) (interface{}, error) {
	secret := viper.GetString("auth.secret")
	return []byte(secret), nil
}
