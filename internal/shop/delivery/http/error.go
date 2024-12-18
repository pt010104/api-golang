package http

import (
	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/shop"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errWrongPaginationQuery        = pkgErrors.NewHTTPError(130001, "Wrong pagination query")
	errWrongQuery                  = pkgErrors.NewHTTPError(130002, "Wrong query")
	errWrongBody                   = pkgErrors.NewHTTPError(130003, "Wrong body")
	errWrongHeader                 = pkgErrors.NewHTTPError(130004, "Wrong header")
	errRequireField                = pkgErrors.NewHTTPError(130010, "Require field shop id or product id")
	ErrInvalidPhone                = pkgErrors.NewHTTPError(130005, "Invalid phone")
	errShopDoesNotExist            = pkgErrors.NewHTTPError(130006, "we cant find this shop")
	ErrNoPermissionToDelete        = pkgErrors.NewHTTPError(130005, "No permission to delete")
	ErrNonExistCategory            = pkgErrors.NewHTTPError(130005, "wrong category ID")
	ErrNoPermissionToDeleteProduct = pkgErrors.NewHTTPError(130005, "No permission to delete product")
)

func (h handler) mapErrors(e error) error {
	switch e {

	case shop.ErrInvalidPhone:
		return ErrInvalidPhone
	case shop.ErrNonExistCategory:
		return ErrNonExistCategory
	case shop.ErrRequireField:
		return errRequireField
	case shop.ErrShopDoesNotExist:
		return errShopDoesNotExist
	case shop.ErrNoPermissionToDeleteProduct:
		return ErrNoPermissionToDeleteProduct
	case admins.ErrNoPermission:
		return ErrNoPermissionToDeleteProduct
	}

	return e
}
