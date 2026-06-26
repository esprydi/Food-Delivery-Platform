package http

import (
	"net/http"

	"order-service/internal/domain"

	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	usecase domain.OrderUsecase
}

func NewOrderHandler(e *echo.Echo, usecase domain.OrderUsecase, jwtMiddleware echo.MiddlewareFunc) {
	handler := &orderHandler{
		usecase: usecase,
	}

	v1 := e.Group("/api/v1")
	
	// Protected routes (Customer only)
	customerGroup := v1.Group("/orders")
	customerGroup.Use(jwtMiddleware)
	customerGroup.Use(RoleMiddleware("CUSTOMER"))
	
	customerGroup.POST("", handler.Checkout)
	customerGroup.GET("/customer", handler.GetCustomerOrders)

	// Protected routes (Merchant only)
	merchantGroup := v1.Group("/orders/merchant")
	merchantGroup.Use(jwtMiddleware)
	merchantGroup.Use(RoleMiddleware("MERCHANT"))
	
	merchantGroup.GET("/:restaurant_id", handler.GetMerchantOrders)
	merchantGroup.PUT("/:order_id/status", handler.UpdateOrderStatus)
}

func successResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
		"error":   nil,
	})
}

func errorResponse(c echo.Context, statusCode int, errorMsg string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": false,
		"message": "Operation failed",
		"data":    nil,
		"error":   errorMsg,
	})
}

type CheckoutRequest struct {
	RestaurantID    string `json:"restaurant_id"`
	DeliveryAddress string `json:"delivery_address"`
	Items           []struct {
		MenuItemID   string  `json:"menu_item_id"`
		MenuItemName string  `json:"menu_item_name"`
		Quantity     int     `json:"quantity"`
		UnitPrice    float64 `json:"unit_price"`
	} `json:"items"`
}

func (h *orderHandler) Checkout(c echo.Context) error {
	customerID := c.Get("user_id").(string) // Extract from JWT via RoleMiddleware
	customerEmail := c.Get("email").(string)

	var req CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	var items []domain.OrderItem
	for _, reqItem := range req.Items {
		items = append(items, domain.OrderItem{
			MenuItemID:   reqItem.MenuItemID,
			MenuItemName: reqItem.MenuItemName,
			Quantity:     reqItem.Quantity,
			UnitPrice:    reqItem.UnitPrice,
		})
	}

	order, err := h.usecase.Checkout(c.Request().Context(), customerID, customerEmail, req.RestaurantID, req.DeliveryAddress, items)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return successResponse(c, http.StatusCreated, "Order placed successfully", order)
}

func (h *orderHandler) GetCustomerOrders(c echo.Context) error {
	customerID := c.Get("user_id").(string)

	orders, err := h.usecase.GetCustomerOrders(c.Request().Context(), customerID)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return successResponse(c, http.StatusOK, "Customer orders retrieved", orders)
}

func (h *orderHandler) GetMerchantOrders(c echo.Context) error {
	restaurantID := c.Param("restaurant_id")

	orders, err := h.usecase.GetMerchantOrders(c.Request().Context(), restaurantID)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return successResponse(c, http.StatusOK, "Merchant orders retrieved", orders)
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (h *orderHandler) UpdateOrderStatus(c echo.Context) error {
	orderID := c.Param("order_id")
	var req UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Validate status conversion
	status := domain.OrderStatus(req.Status)
	
	err := h.usecase.UpdateOrderStatus(c.Request().Context(), orderID, status)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return successResponse(c, http.StatusOK, "Order status updated successfully", nil)
}
