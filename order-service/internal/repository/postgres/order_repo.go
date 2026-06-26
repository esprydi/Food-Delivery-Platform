package postgres

import (
	"context"

	"order-service/internal/domain"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	// GORM will automatically insert the items array because of the foreign key relation
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByCustomerID(ctx context.Context, customerID string) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("customer_id = ?", customerID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetByRestaurantID(ctx context.Context, restaurantID string) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("restaurant_id = ?", restaurantID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}
