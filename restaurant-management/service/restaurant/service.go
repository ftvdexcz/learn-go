package restaurantservice

import (
	"context"
	"restaurant-management/model"
	"restaurant-management/model/entity"
)

type IService interface {
	CreateRestaurant(ctx context.Context, restaurantBody *entity.RestaurantCreate) error
	GetRestaurants(ctx context.Context, paging *model.Paging) ([]*entity.Restaurant, error)
	GetRestaurantById(ctx context.Context, id string) (*entity.Restaurant, error)
	UpdateRestaurant(ctx context.Context, id string, update *entity.RestaurantUpdate) (*entity.Restaurant, error)
	DeleteRestaurantById(ctx context.Context, id string) error
}
