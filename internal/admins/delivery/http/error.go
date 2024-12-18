package http

import (
	"github.com/pt010104/api-golang/internal/admins"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errNoPermission = pkgErrors.NewHTTPError(130001, "you dont have permision to dothis")
	errWrongInput   = pkgErrors.NewHTTPError(130002, "category must have name and description")
	errWrongBody    = pkgErrors.NewHTTPError(130003, "invalid object ID")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case admins.ErrInvalidInput:
		return errWrongInput

	case admins.ErrNoPermission:
		return errNoPermission
	case admins.ErrWrongBody:
		return errWrongBody
	}

	return e
}
