package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	/*
		createRestaurant := entity.Restaurant{
			Name: "Nha hang 1",
			Addr: "Ha Noi",
		}

		if err := db.Create(&createRestaurant).Error; err != nil {
			log.Fatalln(err)
		}
	*/

	/*
		var restaurant entity.Restaurant

		if err := db.Where("id = ?", 1).First(&restaurant).Error; err != nil {
			log.Fatalln(err)
		}

		log.Printf("Before update: %v", restaurant)

		restaurant.Name = "Nhà Hàng 1"
		//restaurant.Addr = "Hà Nội"

		db.Save(&restaurant)
		log.Printf("After update: %v", restaurant)
	*/

	/*
		if err := db.Where("id = ?", 1).Delete(&entity.Restaurant{}).Error; err != nil {
			log.Fatalln(err)
		}
	*/
}
