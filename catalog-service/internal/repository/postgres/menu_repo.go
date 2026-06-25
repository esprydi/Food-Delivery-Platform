package postgres

import (
	"context"

	"catalog-service/internal/domain"

	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) domain.MenuRepository {
	return &menuRepository{
		db: db,
	}
}

func (r *menuRepository) Create(ctx context.Context, menu *domain.MenuItem) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *menuRepository) GetByRestaurantID(ctx context.Context, restaurantID string) ([]domain.MenuItem, error) {
	var menus []domain.MenuItem
	err := r.db.WithContext(ctx).Where("restaurant_id = ? AND is_available = ?", restaurantID, true).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetByID(ctx context.Context, id string) (*domain.MenuItem, error) {
	var menu domain.MenuItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) Update(ctx context.Context, menu *domain.MenuItem) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *menuRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.MenuItem{}, "id = ?", id).Error
}
