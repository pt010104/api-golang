package http

import "github.com/gin-gonic/gin"

func MapPublicRoutes(r *gin.RouterGroup, h Handler) {
	r.GET("", h.Get)
	r.GET("/:id", h.Detail)
}
