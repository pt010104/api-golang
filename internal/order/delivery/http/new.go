package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	CreateCheckout(c *gin.Context)
	CreateOrder(c *gin.Context)
	ListOrder(c *gin.Context)
	ListOrderShop(c *gin.Context)
	UpdateOrder(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc order.UseCase
}

func New(l log.Logger, uc order.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
