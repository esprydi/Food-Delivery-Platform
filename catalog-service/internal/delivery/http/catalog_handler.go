package http

import (
	"net/http"

	"catalog-service/internal/usecase"

	"github.com/labstack/echo/v4"
)

type catalogHandler struct {
	usecase usecase.CatalogUsecase
}

func NewCatalogHandler(e *echo.Echo, usecase usecase.CatalogUsecase, jwtMiddleware echo.MiddlewareFunc) {
	handler := &catalogHandler{
		usecase: usecase,
	}

	v1 := e.Group("/api/v1")
	
	// Public routes
	v1.GET("/restaurants", handler.GetRestaurants)
	v1.GET("/restaurants/:id/menus", handler.GetRestaurantMenus)

	// Protected routes (Merchant only)
	merchantGroup := v1.Group("/merchant")
	merchantGroup.Use(jwtMiddleware)
	merchantGroup.Use(RoleMiddleware("MERCHANT"))
	merchantGroup.POST("/menus", handler.AddMenu)
	merchantGroup.PUT("/menus/:id", handler.UpdateMenu)
	merchantGroup.POST("/restaurants", handler.CreateRestaurant)
	merchantGroup.GET("/restaurants/me", handler.GetMyRestaurant)
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

func (h *catalogHandler) GetRestaurants(c echo.Context) error {
	restaurants, err := h.usecase.GetActiveRestaurants(c.Request().Context())
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return successResponse(c, http.StatusOK, "Active restaurants retrieved", restaurants)
}

func (h *catalogHandler) GetRestaurantMenus(c echo.Context) error {
	restaurantID := c.Param("id")
	menus, err := h.usecase.GetRestaurantMenus(c.Request().Context(), restaurantID)
	if err != nil {
		return errorResponse(c, http.StatusNotFound, err.Error())
	}
	return successResponse(c, http.StatusOK, "Menus retrieved", menus)
}

func (h *catalogHandler) AddMenu(c echo.Context) error {
	ownerID := c.Get("user_id").(string) // Extract from JWT via RoleMiddleware

	var req struct {
		RestaurantID string  `json:"restaurant_id"`
		Name         string  `json:"name"`
		Description  string  `json:"description"`
		Price        float64 `json:"price"`
	}

	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	menu, err := h.usecase.AddMenu(c.Request().Context(), ownerID, req.RestaurantID, req.Name, req.Description, req.Price)
	if err != nil {
		return errorResponse(c, http.StatusForbidden, err.Error()) // Forbidden or BadRequest
	}

	return successResponse(c, http.StatusCreated, "Menu item added successfully", menu)
}

func (h *catalogHandler) CreateRestaurant(c echo.Context) error {
	ownerID := c.Get("user_id").(string)

	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	restaurant, err := h.usecase.CreateRestaurant(c.Request().Context(), ownerID, req.Name, req.Address)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return successResponse(c, http.StatusCreated, "Restaurant created successfully", restaurant)
}

func (h *catalogHandler) GetMyRestaurant(c echo.Context) error {
	ownerID := c.Get("user_id").(string)

	restaurant, err := h.usecase.GetMyRestaurant(c.Request().Context(), ownerID)
	if err != nil {
		return errorResponse(c, http.StatusNotFound, "Restaurant not found")
	}

	return successResponse(c, http.StatusOK, "Restaurant retrieved", restaurant)
}

func (h *catalogHandler) UpdateMenu(c echo.Context) error {
	ownerID := c.Get("user_id").(string)
	menuID := c.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	menu, err := h.usecase.UpdateMenu(c.Request().Context(), ownerID, menuID, req.Name, req.Description, req.Price)
	if err != nil {
		return errorResponse(c, http.StatusForbidden, err.Error())
	}

	return successResponse(c, http.StatusOK, "Menu item updated successfully", menu)
}
