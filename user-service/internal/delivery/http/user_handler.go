package http

import (
	"net/http"

	"user-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, usecase domain.UserUsecase) {
	handler := &userHandler{
		usecase: usecase,
	}

	v1 := e.Group("/api/v1")
	
	v1.POST("/auth/register", handler.Register)
	v1.POST("/auth/login", handler.Login)
	v1.GET("/users/me", handler.GetProfile) // TODO: Add JWT Middleware
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

func (h *userHandler) Register(c echo.Context) error {
	var req struct {
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Phone    string      `json:"phone"`
		Role     domain.Role `json:"role"`
	}

	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	user, err := h.usecase.Register(c.Request().Context(), req.Name, req.Email, req.Password, req.Phone, req.Role)
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, err.Error())
	}

	return successResponse(c, http.StatusCreated, "User registered successfully", user)
}

func (h *userHandler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	token, err := h.usecase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return errorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return successResponse(c, http.StatusOK, "Login successful", map[string]string{"token": token})
}

func (h *userHandler) GetProfile(c echo.Context) error {
	// Normally, we'd use middleware to extract this.
	// For simplicity in this scaffold without middleware implementation yet, we'll extract it manually or mock it.
	// Ideally, we'd get the user ID from the JWT token in context.
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	user, err := h.usecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return errorResponse(c, http.StatusNotFound, err.Error())
	}

	return successResponse(c, http.StatusOK, "Profile retrieved", user)
}
