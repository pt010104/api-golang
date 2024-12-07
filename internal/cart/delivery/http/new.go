package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc cart.UseCase
}

func New(l log.Logger, uc cart.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}