package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"restaurant-management/entity"
)

func CreateRestaurant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var restaurantBody entity.Restaurant

		if err := c.ShouldBind(&restaurantBody); err != nil {
			c.JSON(400, gin.H{
				"message": "bad request",
			})
			return
		}

		fmt.Println(restaurantBody)
		if err := db.Create(&restaurantBody).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "bad request",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "success",
			"data":    restaurantBody,
		})
	}
}

func GetRestaurants(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var restaurants []entity.Restaurant

		if err := db.Find(&restaurants).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "bad request",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "success",
			"data":    restaurants,
		})
	}
}

func GetRestaurant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var restaurant entity.Restaurant

		id := c.Param("id")

		if err := db.Where("id = ?", id).First(&restaurant).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "bad request",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "success",
			"data":    restaurant,
		})
	}
}
