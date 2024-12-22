package checkout

import "errors"

var (
	ErrCartNotFound          = errors.New("cart not found")
	ErrProductNotFoundInCart = errors.New("product not found in cart")
	ErrProductNotFound       = errors.New("product not found")
	ErrProductNotEnoughStock = errors.New("product not enough stock")
)