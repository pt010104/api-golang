package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	MapShopRouters(r, h, mw)
}

func MapShopRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())

	r.POST("", h.Create)
	r.GET("", h.Get)
	r.DELETE("", h.Delete)
	r.PATCH("", h.Update)
	r.POST("/create-product", h.CreateProduct)
}
