package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"restaurant-management/handler"
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

	r := gin.Default()
	v1 := r.Group("/v1/restaurants")
	{
		v1.GET("/", handler.GetRestaurants(db))
		v1.GET("/:id", handler.GetRestaurant(db))
		v1.POST("/", handler.CreateRestaurant(db))
	}

	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
