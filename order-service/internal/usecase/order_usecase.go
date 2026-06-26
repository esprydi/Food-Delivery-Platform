package usecase

import (
	"context"
	"errors"

	"order-service/internal/domain"
)

type orderUsecase struct {
	orderRepo      domain.OrderRepository
	eventPublisher domain.EventPublisher
}

func NewOrderUsecase(repo domain.OrderRepository, publisher domain.EventPublisher) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:      repo,
		eventPublisher: publisher,
	}
}

func (u *orderUsecase) Checkout(ctx context.Context, customerID, restaurantID, deliveryAddress string, items []domain.OrderItem) (*domain.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("cannot create order with empty items")
	}

	var totalAmount float64
	for _, item := range items {
		// IN A REAL PRODUCTION SCENARIO (Option 1):
		// We would make an HTTP Call to Catalog Service here to fetch the actual UnitPrice.
		// For now (Option 2 Prototyping), we trust the unit price passed in the request.
		totalAmount += float64(item.Quantity) * item.UnitPrice
	}

	order := &domain.Order{
		CustomerID:      customerID,
		RestaurantID:    restaurantID,
		Status:          domain.StatusPending,
		TotalAmount:     totalAmount,
		DeliveryAddress: deliveryAddress,
		Items:           items,
	}

	// 1. Save to Database
	err := u.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// 2. Publish Event to RabbitMQ (Asynchronous Workflow Kickoff)
	// Even if publishing fails, the order is safely stored as PENDING.
	// In an advanced setup, you'd use the Transactional Outbox Pattern here.
	_ = u.eventPublisher.PublishOrderCreated(ctx, order)

	return order, nil
}

func (u *orderUsecase) GetCustomerOrders(ctx context.Context, customerID string) ([]domain.Order, error) {
	return u.orderRepo.GetByCustomerID(ctx, customerID)
}

func (u *orderUsecase) GetMerchantOrders(ctx context.Context, restaurantID string) ([]domain.Order, error) {
	return u.orderRepo.GetByRestaurantID(ctx, restaurantID)
}

func (u *orderUsecase) UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	// Optional: add validation here if needed (e.g. only PAID can become PREPARING)
	return u.orderRepo.UpdateStatus(ctx, orderID, status)
}

func (u *orderUsecase) MarkOrderAsPaid(ctx context.Context, orderID string) error {
	return u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusPaid)
}

func (u *orderUsecase) MarkOrderAsFailed(ctx context.Context, orderID string) error {
	return u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusCanceled)
}
