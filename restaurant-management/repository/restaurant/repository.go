package restaurantrepository

import (
	"context"
	"restaurant-management/model"
	"restaurant-management/model/entity"
)

type IRepository interface {
	Create(ctx context.Context, obj *entity.RestaurantCreate) error
	GetById(ctx context.Context, id string) (*entity.Restaurant, error)
	GetAll(ctx context.Context, paging *model.Paging) ([]*entity.Restaurant, error)
	Update(ctx context.Context, id string, obj *entity.RestaurantUpdate) (*entity.Restaurant, error)
	DeleteById(ctx context.Context, id string) error
}
