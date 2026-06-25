package http

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return errorResponse(c, http.StatusUnauthorized, "missing or invalid token")
			}

			token := user.(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			
			role, ok := claims["role"].(string)
			if !ok || !strings.EqualFold(role, requiredRole) {
				return errorResponse(c, http.StatusForbidden, "you don't have permission to access this resource")
			}

			// Store user ID in context for easy access
			if sub, ok := claims["sub"].(string); ok {
				c.Set("user_id", sub)
			}

			return next(c)
		}
	}
}
