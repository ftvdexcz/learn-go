package main

import (
	"github.com/gin-gonic/gin"
	"restaurant-management/handler"
)

func main() {
	bootstrap()

	h := handler.New()

	r := gin.Default()
	v1 := r.Group("/v1/restaurants")
	{
		v1.GET("/", h.GetRestaurants)
		v1.GET("/:id", h.GetRestaurant)
		v1.POST("/", h.CreateRestaurant)
		v1.PATCH("/:id", h.UpdateRestaurant)
		v1.DELETE("/:id", h.DeleteRestaurant)
	}

	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
