package commons

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
	连接数据库
**/

//具体的数据库连接信息
type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

//获取数据库信息
var databaseSetting = &Database{}
var db *gorm.DB

//初始化函数，结构体加载时调用，将数据库连接信息赋值，并初始化一个数据库连接
func init() {
	//加载数据库文件
	filePtr, err := os.Open("./config/db.json") //config的文件目录
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
		return
	}
	//关闭文件
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	//读取配置文件中的信息
	err = decoder.Decode(databaseSetting)
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
	}
	//初始化一个数据库连接
	db = newConnection()
}

//获取数据库连接，私有方法
func newConnection() *gorm.DB {
	var dbUri string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		databaseSetting.User,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.Port,
		databaseSetting.Name)
	// 获取数据库连接
	conn, err := gorm.Open(mysql.Open(dbUri), &gorm.Config{})
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
	}
	//设置数据库连接池信息
	setup(conn)
	return conn
}

//设置数据库连接池
func setup(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
	}
	sqlDB.SetMaxIdleConns(10)                   //最大空闲连接数
	sqlDB.SetMaxOpenConns(30)                   //最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 300) //设置连接空闲超时
	//db.LogMode(true)
}

//获取DB对象，当前连接未断开，直接返回，否则返回新连接
func GetDB() *gorm.DB {
	sqlDB, err := db.DB()
	if err != nil {
		db = newConnection()
		return db
	}
	e := sqlDB.Ping()
	if e != nil {
		db = newConnection()
	}
	return db
}
