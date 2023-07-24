package restaurantstore

import (
	"context"
	"gorm.io/gorm"
	"restaurant-management/model"
	"restaurant-management/model/entity"
)

type RestaurantStore struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RestaurantStore {
	return &RestaurantStore{
		db: db,
	}
}

func (r RestaurantStore) Create(ctx context.Context, obj *entity.RestaurantCreate) error {
	if err := r.db.Create(obj).Error; err != nil {
		return err
	}

	return nil
}

func (r RestaurantStore) GetById(ctx context.Context, id string) (*entity.Restaurant, error) {
	var restaurant entity.Restaurant

	if err := r.db.Where("id = ?", id).First(&restaurant).Error; err != nil {
		return nil, err
	}

	return &restaurant, nil
}

func (r RestaurantStore) GetAll(ctx context.Context, paging *model.Paging) ([]*entity.Restaurant, error) {
	var restaurants []*entity.Restaurant

	if err := r.db.Offset((paging.Page - 1) * paging.Limit).
		Order("id desc").
		Limit(paging.Limit).
		Find(&restaurants).Error; err != nil {
		return nil, err
	}

	return restaurants, nil
}

func (r RestaurantStore) Update(ctx context.Context, id string, obj *entity.RestaurantUpdate) (*entity.Restaurant, error) {
	var restaurant entity.Restaurant

	if err := r.db.Where("id = ?", id).First(&restaurant).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&restaurant).Updates(obj).Error; err != nil {
		return nil, err
	}

	return &restaurant, nil
}

func (r RestaurantStore) DeleteById(ctx context.Context, id string) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.Restaurant{}).Error; err != nil {
		return err
	}

	return nil
}
