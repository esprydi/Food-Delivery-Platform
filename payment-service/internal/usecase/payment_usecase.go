package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"payment-service/internal/domain"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type paymentUsecase struct {
	paymentRepo    domain.PaymentRepository
	eventPublisher domain.PaymentEventPublisher
}

func NewPaymentUsecase(repo domain.PaymentRepository, publisher domain.PaymentEventPublisher) domain.PaymentUsecase {
	return &paymentUsecase{
		paymentRepo:    repo,
		eventPublisher: publisher,
	}
}

func (u *paymentUsecase) ProcessOrderCreated(ctx context.Context, orderID string, customerID string, customerEmail string, amount float64) error {
	// 1. Check if payment already exists
	_, err := u.paymentRepo.GetByOrderID(ctx, orderID)
	if err == nil {
		// Already processed
		return nil
	}

	// 2. Request Snap URL to Midtrans
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "Customer",
			Email: customerEmail,
		},
	}

	snapResp, snapErr := snap.CreateTransaction(req)
	var snapUrl string
	if snapErr != nil {
		slog.Error("Failed to create Midtrans Snap transaction", "error", snapErr.GetMessage())
		// We can still create a pending payment without URL and retry later, 
		// or use a dummy URL for pure local testing if Midtrans is down/not configured.
		snapUrl = fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/dummy_%s", orderID)
	} else {
		snapUrl = snapResp.RedirectURL
	}

	// 3. Save Payment to DB
	payment := &domain.Payment{
		OrderID:       orderID,
		CustomerEmail: customerEmail,
		Amount:        amount,
		Status:        domain.PaymentStatusPending,
		SnapURL:       snapUrl,
	}

	return u.paymentRepo.Create(ctx, payment)
}

func (u *paymentUsecase) HandleMidtransNotification(ctx context.Context, orderID, transactionStatus, transactionID string) error {
	// 1. Verify payment exists
	payment, err := u.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return errors.New("payment not found for order")
	}

	// 2. Determine new status
	newStatus := payment.Status
	
	switch transactionStatus {
	case "capture", "settlement":
		newStatus = domain.PaymentStatusSuccess
	case "deny", "cancel", "expire":
		newStatus = domain.PaymentStatusFailed
	case "pending":
		newStatus = domain.PaymentStatusPending
	}

	// 3. Update DB
	err = u.paymentRepo.UpdateStatusAndTransactionID(ctx, orderID, newStatus, transactionID)
	if err != nil {
		return err
	}

	// 4. Publish Event if Success
	if newStatus == domain.PaymentStatusSuccess && payment.Status != domain.PaymentStatusSuccess {
		// Publish event
		err = u.eventPublisher.PublishPaymentSuccess(ctx, payment)
		if err != nil {
			slog.Error("Failed to publish payment success", "error", err)
			return err
		}
	}

	return nil
}

func (u *paymentUsecase) GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return u.paymentRepo.GetByOrderID(ctx, orderID)
}
