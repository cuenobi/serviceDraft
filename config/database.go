package config

import (
	"fmt"
	"net/url"

	"github.com/cuenobi/serviceDraft/service/entity"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() *gorm.DB {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.username`)
	dbPass := viper.GetString(`database.password`)
	dbName := viper.GetString(`database.dbname`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Bangkok")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	Db.AutoMigrate(entity.User{})

	return Db
}
