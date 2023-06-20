package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"restaurant-management/pkg/container"
	restaurantrepository "restaurant-management/repository/restaurant"
	restaurantservice "restaurant-management/service/restaurant"
	restaurantstore "restaurant-management/store/restaurant"
)

func bootstrap() {
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

	restaurantStore := restaurantstore.New(db)
	container.Register(
		func() restaurantrepository.IRepository {
			return restaurantStore
		},
	)

	container.Register(
		func() restaurantservice.IService {
			return restaurantservice.New()
		},
	)

}
