package lib

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

//DB .
var DB *gorm.DB

// GetDB 获取数据库单例
func GetDB() *gorm.DB {
	if DB != nil {
		return DB
	}
	conf := Conf()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.GetConf("db", "user"),
		conf.GetConf("db", "pwd"),
		conf.GetConf("db", "host"),
		conf.GetConf("db", "port"),
		conf.GetConf("db", "dbname"))
	db, err := gorm.Open("mysql", connStr)
	DB = db
	if err != nil {
		fmt.Println("connect mysql failed ", err)
		os.Exit(1)
	}

	if os.Getenv("GIN_MODE") != "release" {
		db.LogMode(true)
	}

	return db
}
