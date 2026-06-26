package domain

import (
	"context"
	"time"
)

type OrderStatus string

const (
	StatusPending        OrderStatus = "PENDING"
	StatusPaid           OrderStatus = "PAID"
	StatusPreparing      OrderStatus = "PREPARING"
	StatusReadyForPickup OrderStatus = "READY_FOR_PICKUP"
	StatusDelivering     OrderStatus = "DELIVERING"
	StatusCompleted      OrderStatus = "COMPLETED"
	StatusCanceled       OrderStatus = "CANCELED"
)

type Order struct {
	ID              string      `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID      string      `json:"customer_id" db:"customer_id"`
	RestaurantID    string      `json:"restaurant_id" db:"restaurant_id"`
	Status          OrderStatus `json:"status" db:"status"`
	TotalAmount     float64     `json:"total_amount" db:"total_amount"`
	DeliveryAddress string      `json:"delivery_address" db:"delivery_address"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID           string  `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID      string  `json:"order_id" db:"order_id"`
	MenuItemID   string  `json:"menu_item_id" db:"menu_item_id"`
	MenuItemName string  `json:"menu_item_name" db:"menu_item_name"`
	Quantity     int     `json:"quantity" db:"quantity"`
	UnitPrice    float64 `json:"unit_price" db:"unit_price"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]Order, error)
	GetByRestaurantID(ctx context.Context, restaurantID string) ([]Order, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
}

type OrderUsecase interface {
	Checkout(ctx context.Context, customerID, restaurantID, deliveryAddress string, items []OrderItem) (*Order, error)
	GetCustomerOrders(ctx context.Context, customerID string) ([]Order, error)
	GetMerchantOrders(ctx context.Context, restaurantID string) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status OrderStatus) error
	MarkOrderAsPaid(ctx context.Context, orderID string) error
	MarkOrderAsFailed(ctx context.Context, orderID string) error
}

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, order *Order) error
}
