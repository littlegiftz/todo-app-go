package db

import (
	"fmt"
	"os"

	"github.com/littlegiftz/todo-app-go/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func Init() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("DB Connection Error")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Todo{})

}

func DBManager() *gorm.DB {
	return db
}
