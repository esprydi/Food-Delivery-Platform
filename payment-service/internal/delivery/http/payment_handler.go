package http

import (
	"net/http"

	"payment-service/internal/domain"

	"github.com/labstack/echo/v4"
)

type paymentHandler struct {
	usecase domain.PaymentUsecase
}

func NewPaymentHandler(e *echo.Echo, usecase domain.PaymentUsecase) {
	handler := &paymentHandler{
		usecase: usecase,
	}

	v1 := e.Group("/api/v1")
	
	// This endpoint is meant to be called by Midtrans servers (Webhook)
	// It doesn't use JWT since Midtrans uses signature keys (simplified here)
	// It doesn't use JWT since Midtrans uses signature keys (simplified here)
	v1.POST("/payments/midtrans/notification", handler.MidtransNotification)
	v1.GET("/payments/order/:order_id", handler.GetPaymentByOrderID)
}

func successResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func errorResponse(c echo.Context, statusCode int, errorMsg string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": false,
		"message": "Operation failed",
		"error":   errorMsg,
	})
}

func (h *paymentHandler) MidtransNotification(c echo.Context) error {
	var notificationPayload map[string]interface{}

	if err := c.Bind(&notificationPayload); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Safely extract fields
	orderID, ok1 := notificationPayload["order_id"].(string)
	transactionStatus, ok2 := notificationPayload["transaction_status"].(string)
	transactionID, ok3 := notificationPayload["transaction_id"].(string)

	if !ok1 || !ok2 || !ok3 {
		// Just a dummy payload for local testing if some fields are missing
		return errorResponse(c, http.StatusBadRequest, "Missing required Midtrans fields")
	}

	err := h.usecase.HandleMidtransNotification(c.Request().Context(), orderID, transactionStatus, transactionID)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return successResponse(c, http.StatusOK, "Notification processed", nil)
}

func (h *paymentHandler) GetPaymentByOrderID(c echo.Context) error {
	orderID := c.Param("order_id")
	payment, err := h.usecase.GetPaymentByOrderID(c.Request().Context(), orderID)
	if err != nil {
		return errorResponse(c, http.StatusNotFound, "Payment not found")
	}

	return successResponse(c, http.StatusOK, "Payment retrieved", payment)
}
