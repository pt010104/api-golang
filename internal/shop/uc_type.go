package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type CreateInput struct {
	Name     string
	City     string
	Street   string
	District string
	Phone    string
}

type GetShopsFilter struct {
	IDs    []string
	Search string
}

type GetInput struct {
	GetShopsFilter
	PagQuery paginator.PaginatorQuery
}

type GetOutput struct {
	Shops []models.Shop
	Pag   paginator.Paginator
}
