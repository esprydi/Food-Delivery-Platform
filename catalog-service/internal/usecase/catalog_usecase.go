package usecase

import (
	"context"
	"errors"

	"catalog-service/internal/domain"
)

type CatalogUsecase interface {
	GetActiveRestaurants(ctx context.Context) ([]domain.Restaurant, error)
	GetRestaurantMenus(ctx context.Context, restaurantID string) ([]domain.MenuItem, error)
	AddMenu(ctx context.Context, ownerID, restaurantID, name, description string, price float64) (*domain.MenuItem, error)
	UpdateMenu(ctx context.Context, ownerID, menuID, name, description string, price float64) (*domain.MenuItem, error)
	DeleteMenu(ctx context.Context, ownerID, menuID string) error
	CreateRestaurant(ctx context.Context, ownerID, name, address string) (*domain.Restaurant, error)
	GetMyRestaurant(ctx context.Context, ownerID string) (*domain.Restaurant, error)
}

type catalogUsecase struct {
	restaurantRepo domain.RestaurantRepository
	menuRepo       domain.MenuRepository
}

func NewCatalogUsecase(rRepo domain.RestaurantRepository, mRepo domain.MenuRepository) CatalogUsecase {
	return &catalogUsecase{
		restaurantRepo: rRepo,
		menuRepo:       mRepo,
	}
}

func (u *catalogUsecase) GetActiveRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	return u.restaurantRepo.GetAllActive(ctx)
}

func (u *catalogUsecase) GetRestaurantMenus(ctx context.Context, restaurantID string) ([]domain.MenuItem, error) {
	// Optional: verify if restaurant exists
	_, err := u.restaurantRepo.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	return u.menuRepo.GetByRestaurantID(ctx, restaurantID)
}

func (u *catalogUsecase) AddMenu(ctx context.Context, ownerID, restaurantID, name, description string, price float64) (*domain.MenuItem, error) {
	restaurant, err := u.restaurantRepo.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	// Validate ownership
	if restaurant.OwnerID != ownerID {
		return nil, errors.New("unauthorized: you don't own this restaurant")
	}

	menuItem := &domain.MenuItem{
		RestaurantID: restaurantID,
		Name:         name,
		Description:  description,
		Price:        price,
		IsAvailable:  true,
	}

	err = u.menuRepo.Create(ctx, menuItem)
	if err != nil {
		return nil, err
	}

	return menuItem, nil
}

func (u *catalogUsecase) CreateRestaurant(ctx context.Context, ownerID, name, address string) (*domain.Restaurant, error) {
	restaurant := &domain.Restaurant{
		OwnerID: ownerID,
		Name:    name,
		Address: address,
		IsOpen:  true,
	}

	err := u.restaurantRepo.Create(ctx, restaurant)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (u *catalogUsecase) GetMyRestaurant(ctx context.Context, ownerID string) (*domain.Restaurant, error) {
	return u.restaurantRepo.GetByOwnerID(ctx, ownerID)
}

func (u *catalogUsecase) UpdateMenu(ctx context.Context, ownerID, menuID, name, description string, price float64) (*domain.MenuItem, error) {
	menu, err := u.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return nil, errors.New("menu not found")
	}

	restaurant, err := u.restaurantRepo.GetByID(ctx, menu.RestaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	if restaurant.OwnerID != ownerID {
		return nil, errors.New("unauthorized: you don't own this restaurant")
	}

	menu.Name = name
	menu.Description = description
	menu.Price = price

	err = u.menuRepo.Update(ctx, menu)
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (u *catalogUsecase) DeleteMenu(ctx context.Context, ownerID, menuID string) error {
	menu, err := u.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return errors.New("menu not found")
	}

	restaurant, err := u.restaurantRepo.GetByID(ctx, menu.RestaurantID)
	if err != nil {
		return errors.New("restaurant not found")
	}

	if restaurant.OwnerID != ownerID {
		return errors.New("unauthorized: you don't own this restaurant")
	}

	return u.menuRepo.Delete(ctx, menuID)
}
