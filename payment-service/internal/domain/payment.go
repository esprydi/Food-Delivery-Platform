package domain

import (
	"context"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID            string        `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID       string        `json:"order_id" db:"order_id"`
	CustomerEmail string        `json:"customer_email" db:"customer_email"`
	Amount        float64       `json:"amount" db:"amount"`
	Status        PaymentStatus `json:"status" db:"status"`
	SnapURL       string        `json:"snap_url" db:"snap_url"`
	TransactionID string        `json:"transaction_id" db:"transaction_id"` // From Midtrans
	CreatedAt     time.Time     `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByOrderID(ctx context.Context, orderID string) (*Payment, error)
	UpdateStatusAndTransactionID(ctx context.Context, orderID string, status PaymentStatus, transactionID string) error
}

type PaymentUsecase interface {
	ProcessOrderCreated(ctx context.Context, orderID string, customerID string, customerEmail string, amount float64) error
	HandleMidtransNotification(ctx context.Context, orderID, transactionStatus, transactionID string) error
	GetPaymentByOrderID(ctx context.Context, orderID string) (*Payment, error)
}

type PaymentEventPublisher interface {
	PublishPaymentSuccess(ctx context.Context, payment *Payment) error
}
