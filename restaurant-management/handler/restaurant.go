package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"restaurant-management/entity"
	"restaurant-management/model"
	"strconv"
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
		var paging model.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(400, gin.H{
				"message": "bad request",
			})
			return
		}

		fmt.Println(paging.Page, paging.Limit)

		if paging.Page <= 0 {
			paging.Page = 0
		}

		if paging.Limit <= 0 {
			paging.Limit = 5
		}

		var restaurants []entity.Restaurant

		db.Offset((paging.Page - 1) * paging.Limit).
			Order("id desc").
			Limit(paging.Limit).
			Find(&restaurants)

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

func UpdateRestaurant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		var restaurantUpdate entity.RestaurantUpdate

		if err := c.ShouldBind(&restaurantUpdate); err != nil {
			c.JSON(400, gin.H{
				"message": "bad request",
			})
			return
		}

		var restaurant entity.Restaurant
		if err := db.Where("id = ?", id).First(&restaurant).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "bad request",
			})
			return
		}

		if err := db.Model(&restaurant).Updates(&restaurantUpdate).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "bad request 11",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "success",
			"data":    restaurant,
		})
		return
	}
}

func DeleteRestaurant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).Delete(&entity.Restaurant{}).Error; err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(204, nil)
		return
	}
}
