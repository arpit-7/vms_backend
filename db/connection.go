package db 

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_"github.com/go-sql-driver/mysql"
)

var DB *gorm.DB
var DBname = "vms_backend"

func Connect(){
	dsnRoot:="root:Chauhan7@@tcp(localhost:3306)/"
	sqlDB, err := sql.Open("mysql",dsnRoot)
	if err!= nil {
		log.Fatal("Failed to connect to MySQL", err)
	}

	_,err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + DBname)
	if err!= nil {
		log.Fatal("Failed to craete databse", err)
	}
	fmt.Println("Databse Created",DBname)

	sqlDB.Close()

	dsn := fmt.Sprintf("root:Chauhan7@@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",DBname)
	DB, err = gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err!= nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	fmt.Println("Connected to database")
}