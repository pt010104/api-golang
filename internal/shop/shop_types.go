package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

// Shop
type CreateShop struct {
	Name     string
	City     string
	Street   string
	District string
	Phone    string
}

type GetShopsFilter struct {
	IDs        []string
	Search     string
	IsVerified *bool
}
type DeleteShopInput struct {
	ID string
}

type GetShopInput struct {
	GetShopsFilter
	PagQuery paginator.PaginatorQuery
}
type UpdateInput struct {
	ShopID     string
	ShopIDs    []string
	Name       string
	City       string
	Street     string
	District   string
	Phone      string
	IsVerified bool
}
type GetShopOutput struct {
	Shops []models.Shop
	Pag   paginator.Paginator
}
