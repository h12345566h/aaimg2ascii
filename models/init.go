package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Database struct {
	Mysql *gorm.DB
}

var DB *Database

func (db *Database) Init() {
	DB = &Database{
		Mysql: initMysql(),
	}
}

func (db *Database) Close() {
	log.Error("DB Close")

	db.Mysql.Close()
}

func initMysql() *gorm.DB {

	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.addr"),
		viper.GetString("mysql.name"),
		true,
		"Asia%2FTaipei")

	db, err := gorm.Open("mysql", config)

	if err != nil {
		panic(err)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("DB connect failed")
	} else {
		log.Info("DB connect success")
		db.LogMode(viper.GetBool("gormlog"))
		db.DB().SetMaxOpenConns(4000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
		db.DB().SetMaxIdleConns(2)    // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。

		//region AutoMigrate
		db.AutoMigrate(&User{}, &AsciiArt{})
		//AddForeignKey
		db.Model(&AsciiArt{}).AddForeignKey("user_id", "users(user_id)", "RESTRICT", "RESTRICT")
		//endregion

	}

	return db
}

func keepAlive(db *gorm.DB) {
	for {
		db.DB().Ping()
		time.Sleep(60 * time.Second)
	}
}
