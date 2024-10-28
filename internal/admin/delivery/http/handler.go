package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
	"github.com/pt010104/api-golang/pkg/util"
)

// VerifyShop godoc
// @Summary Verify shop
// @Schemes http https
// @Description verify shop by  IDs
// @Tags Shop
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param session-id header string true "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param ids body []shop_id
// @Success 200 {object} VerifyShopResp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/admin/verify-shop

func (h handler) VerifyShop(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processVerifyShopRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "admin.delivery.http.handler.Create.processCreateRequest: %v", err)
		response.Error(c, err)
		return
	}
	h.l.Debugf(ctx, "role:", sc.Role)
	s, err := h.uc.VerifyShop(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "admin.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}
	res := h.toResList(s)

	util.PrintJson(res)

	response.OK(c, res)
}