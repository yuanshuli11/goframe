package datebase

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

const (
	USERNAME = "root"
	PASSWORD = "Shuli123666~"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "yourwar"
)

func init() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "" + defaultTableName
	}
	db, err := gorm.Open("mysql", dsn)
	//db, err := sql.Open("mysql", dsn)
	db.LogMode(true)
	db.SetC
	if err != nil {
		fmt.Printf("Open mysql failed,err:%v\n", err)
		return
	}
	//db.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	//db.SetMaxOpenConns(100)                  //设置最大连接数
	//db.SetMaxIdleConns(16)                   //设置闲置连接数
	//	defer DB.Close()

	DB = db
}
