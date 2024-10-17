package shop

import "errors"

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrNoPermissionToDelete = errors.New("You dont have permission to delete this shop")
	ErrInvalidPhone         = errors.New("invalid phone")
	ErrShopExist            = errors.New("shop exist")
	ErrShopDoesNotExist     = errors.New("shop doesnot exist")
)
