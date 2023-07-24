package restaurantservice

import (
	"context"
	"restaurant-management/model"
	"restaurant-management/model/entity"
	"restaurant-management/pkg/container"
	restaurantrepository "restaurant-management/repository/restaurant"
)

type service struct {
	restaurantStore restaurantrepository.IRepository `di:"inject"`
}

func New() *service {
	obj := &service{}
	container.Fill(obj)
	return obj
}

func (s service) CreateRestaurant(ctx context.Context, restaurantBody *entity.RestaurantCreate) error {
	if err := s.restaurantStore.Create(ctx, restaurantBody); err != nil {
		return err
	}

	return nil
}

func (s service) GetRestaurants(ctx context.Context, paging *model.Paging) ([]*entity.Restaurant, error) {
	if paging.Page <= 0 {
		paging.Page = 0
	}

	if paging.Limit <= 0 {
		paging.Limit = 5
	}

	response, err := s.restaurantStore.GetAll(ctx, paging)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s service) GetRestaurantById(ctx context.Context, id string) (*entity.Restaurant, error) {
	reponse, err := s.restaurantStore.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return reponse, nil
}

func (s service) UpdateRestaurant(ctx context.Context, id string, update *entity.RestaurantUpdate) (*entity.Restaurant, error) {
	updated, err := s.restaurantStore.Update(ctx, id, update)

	if err != nil {
		return nil, err
	}

	return updated, nil
}
func (s service) DeleteRestaurantById(ctx context.Context, id string) error {
	if err := s.restaurantStore.DeleteById(ctx, id); err != nil {
		return err
	}

	return nil
}
