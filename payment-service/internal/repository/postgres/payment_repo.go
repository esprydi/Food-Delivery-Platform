package postgres

import (
	"context"

	"payment-service/internal/domain"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdateStatusAndTransactionID(ctx context.Context, orderID string, status domain.PaymentStatus, transactionID string) error {
	return r.db.WithContext(ctx).Model(&domain.Payment{}).Where("order_id = ?", orderID).Updates(map[string]interface{}{
		"status":         status,
		"transaction_id": transactionID,
	}).Error
}
