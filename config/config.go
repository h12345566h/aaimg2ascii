package config

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
	"time"
)

type Config struct{}

func Init() error {
	c := Config{}

	if err := c.initConfig(); err != nil {
		return err
	}
	c.initLog()
	c.watchConfig()
	return nil
}

func (c *Config) initConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("conf")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) initLog() {
	baseLogPath := path.Join(viper.GetString("log.path"), viper.GetString("log.filename"))
	writer, err := rotatelogs.New(
		baseLogPath+"-%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(baseLogPath),     // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(120*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(3*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{DisableColors: true})
	log.AddHook(lfHook)
}

func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("config file changed")
	})
}
