package main

import (
	"aaimg2ascii/config"
	"aaimg2ascii/models"
	"aaimg2ascii/route"
	"aaimg2ascii/validators"
	"context"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	log.Info("run main")
	models.DB.Init()
	defer models.DB.Close()

	os.MkdirAll("./public", 0777)
	if viper.GetString("ginmode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	app := newApp()
	route.InitRouter(app)
	validators.InitValidator()

	if viper.GetBool("autotls") {
		//app.Run(iris.AutoTLS(":443", viper.GetString("domain"), viper.GetString("email"), func(h *iris.Supervisor) {
		//	h.RegisterOnServe(func(host host.TaskHost) {
		//		log.Info("server running!")
		//	})
		//}))
	} else {
		//app.Run(iris.Addr(":80", func(h *iris.Supervisor) {
		//	h.RegisterOnServe(func(host host.TaskHost) {
		//		log.Info("server running!")
		//	})
		//}))
		srv := &http.Server{
			Addr:    ":8000",
			Handler: app,
		}

		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")
	}
}
func newApp() *gin.Engine {
	app := gin.New()
	log.Info("new app finished")
	app.Use(CORSMiddleware())
	app.Use(static.Serve("/", static.LocalFile("./public", false)))
	return app
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
