package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error al cargar el .env")
	}

	var DNS = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	var error error
	DB, error = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if error != nil {
		log.Fatal(error)
	} else {
		fmt.Println("DB Connected")
	}

}
