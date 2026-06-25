package domain

import (
	"context"
)

type MenuItem struct {
	ID           string  `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RestaurantID string  `json:"restaurant_id" db:"restaurant_id"`
	Name         string  `json:"name" db:"name"`
	Description  string  `json:"description" db:"description"`
	Price        float64 `json:"price" db:"price"`
	IsAvailable  bool    `json:"is_available" db:"is_available" gorm:"default:true"`
}

type MenuRepository interface {
	Create(ctx context.Context, menu *MenuItem) error
	GetByRestaurantID(ctx context.Context, restaurantID string) ([]MenuItem, error)
	GetByID(ctx context.Context, id string) (*MenuItem, error)
	Update(ctx context.Context, menu *MenuItem) error
}
