package domain

import (
	"context"
	"time"
)

type Role string

const (
	RoleCustomer Role = "CUSTOMER"
	RoleMerchant Role = "MERCHANT"
	RoleDriver   Role = "DRIVER"
	RoleAdmin    Role = "ADMIN"
)

type User struct {
	ID           string    `json:"id" db:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email" gorm:"uniqueIndex"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never return password hash in JSON
	Role         Role      `json:"role" db:"role"`
	Phone        string    `json:"phone" db:"phone"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, name, email, password, phone string, role Role) (*User, error)
	Login(ctx context.Context, email, password string) (string, error) // Returns JWT Token
	GetProfile(ctx context.Context, id string) (*User, error)
}
