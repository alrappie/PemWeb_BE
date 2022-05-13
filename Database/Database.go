package Database

import (
	"PemWeb_BE/Car"
	"PemWeb_BE/User"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func Open() *gorm.DB {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
	}
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))),
		&gorm.Config{})
	if err != nil {
		println(err.Error())
	}
	if err := db.AutoMigrate(&User.User{}, &Car.Car{}); err != nil {
		log.Fatal(err.Error())
	}
	return db
}
