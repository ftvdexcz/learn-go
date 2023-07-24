package handler

import (
	"github.com/gin-gonic/gin"
	"restaurant-management/model"
	"restaurant-management/model/entity"
	"restaurant-management/pkg/container"
	restaurantservice "restaurant-management/service/restaurant"
)

type Controller struct {
	restaurantService restaurantservice.IService `di:"inject"`
}

func New() *Controller {
	obj := &Controller{}
	container.Fill(obj)
	return obj
}

func (o Controller) CreateRestaurant(c *gin.Context) {
	var restaurantBody entity.RestaurantCreate

	if err := c.ShouldBind(&restaurantBody); err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	if err := o.restaurantService.CreateRestaurant(c.Request.Context(), &restaurantBody); err != nil {
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

func (o Controller) GetRestaurants(c *gin.Context) {
	var paging model.Paging

	if err := c.ShouldBind(&paging); err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	response, err := o.restaurantService.GetRestaurants(c.Request.Context(), &paging)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    response,
	})
}

func (o Controller) GetRestaurant(c *gin.Context) {
	id := c.Param("id")

	response, err := o.restaurantService.GetRestaurantById(c.Request.Context(), id)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    response,
	})
}

func (o Controller) UpdateRestaurant(c *gin.Context) {
	id := c.Param("id")

	var restaurantUpdate entity.RestaurantUpdate

	if err := c.ShouldBind(&restaurantUpdate); err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	updated, err := o.restaurantService.UpdateRestaurant(c.Request.Context(), id, &restaurantUpdate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    updated,
	})
	return
}

func (o Controller) DeleteRestaurant(c *gin.Context) {
	id := c.Param("id")

	if err := o.restaurantService.DeleteRestaurantById(c.Request.Context(), id); err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	c.JSON(204, nil)
	return
}
