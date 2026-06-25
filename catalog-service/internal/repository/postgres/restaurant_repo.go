package postgres

import (
	"context"

	"catalog-service/internal/domain"

	"gorm.io/gorm"
)

type restaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) domain.RestaurantRepository {
	return &restaurantRepository{
		db: db,
	}
}

func (r *restaurantRepository) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	return r.db.WithContext(ctx).Create(restaurant).Error
}

func (r *restaurantRepository) GetAllActive(ctx context.Context) ([]domain.Restaurant, error) {
	var restaurants []domain.Restaurant
	err := r.db.WithContext(ctx).Where("is_open = ?", true).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) GetByID(ctx context.Context, id string) (*domain.Restaurant, error) {
	var restaurant domain.Restaurant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&restaurant).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *restaurantRepository) GetByOwnerID(ctx context.Context, ownerID string) (*domain.Restaurant, error) {
	var restaurant domain.Restaurant
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).First(&restaurant).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}
