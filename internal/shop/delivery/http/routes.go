package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	MapShopRouters(r, h, mw)
	MapProductRouters(r.Group("/products"), h, mw)
}

func MapShopRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {

	r.Use(mw.Auth())
	r.GET("/:id", h.Detail)
	r.GET("", h.Get)
	r.POST("", h.Create)

	r.DELETE("", h.Delete)
	r.PATCH("", h.Update)
	r.GET("get-shop-id-by-user-id/:id", h.GetShopIDByUserID)

	r.Use(mw.AuthShop())
	r.GET("/report", h.Report)
}

func MapProductRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {

	r.Use(mw.Auth(), mw.AuthShop())
	r.DELETE("", h.DeleteProduct)
	r.POST("", h.CreateProduct)
	r.PUT("", h.UpdateProduct)

}
