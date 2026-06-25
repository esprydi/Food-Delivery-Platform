package domain

import (
	"context"
)

type Restaurant struct {
	ID      string `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID string `json:"owner_id" db:"owner_id"` // Matches User ID
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	IsOpen  bool   `json:"is_open" db:"is_open" gorm:"default:true"`
}

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *Restaurant) error
	GetAllActive(ctx context.Context) ([]Restaurant, error)
	GetByID(ctx context.Context, id string) (*Restaurant, error)
	GetByOwnerID(ctx context.Context, ownerID string) (*Restaurant, error)
}
