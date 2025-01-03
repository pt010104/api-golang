package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/paginator"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary		Create a shop
// @Schemes		http https
// @Description	Create a shop
// @Tags			Shop
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string	false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string	true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string	true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string	true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			name						body		string	true	"Name"
// @Param			phone						body		string	true	"Phone"
// @Param			city						body		string	true	"City"
// @Param			street						body		string	true	"Street"
// @Param			district					body		string	false	"District"
//
// @Success		200							{object}	registerResponse
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops [POST]
func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Create: %v", err)
		response.Error(c, err)
		return
	}

	shop, err := h.uc.Create(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Create: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newCreateResponse(shop))
}

// @Summary		Get shop
// @Schemes		http https
// @Description	Get shop by Search, IDs,...
// @Tags			Shop
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			search						query		string		false	"Search"
// @Param			ids							query		[]string	false	"IDs"
// @Param			is_verified					query		bool		false	"Is Verified"
// @Param			page						query		int			false	"Page"	default(1)
// @Param			limit						query		int			false	"Limit"	default(10)
//
// @Success		200							{object}	getShopResp
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops [GET]
func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := h.processGetRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Get: %v", err)
		response.Error(c, err)
		return
	}

	var pagQuery paginator.PaginatorQuery
	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Get.ShouldBindQuery: %v", err)
		response.Error(c, errWrongPaginationQuery)
		return
	}

	pagQuery.Adjust()

	s, err := h.uc.Get(ctx, models.Scope{}, shop.GetShopInput{
		PagQuery:       pagQuery,
		GetShopsFilter: req.toInput(),
	})
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.Get: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newGetShopsResp(s))
}

// @Summary		Get shop detail
// @Schemes		http https
// @Description	Get shop detail by ID
// @Tags			Shop
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string	false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string	true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string	true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string	true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			id							path		string	true	"Shop ID"
//
// @Success		200							{object}	getDetailResp
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops/{id} [GET]
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := h.processDetailRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Detail: %v", err)
		response.Error(c, err)
		return
	}

	shop, err := h.uc.Detail(ctx, models.Scope{}, id)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Detail: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newDetailResponse(shop))
}

// @Summary		Delete a shop
// @Schemes		http https
// @Description	Delete a shop
// @Tags			Shop
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string	false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string	true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string	true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string	true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
//
// @Success		200							{object}	nil
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops [DELETE]
func (h handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processDeleteShopRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Delete: %v", err)
		response.Error(c, err)
		return
	}

	err = h.uc.Delete(ctx, sc)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Delete: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)

}

// @Summary		Update a shop
// @Schemes		http https
// @Description	Update a shop
// @Tags			Shop
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			request						body		updateShopRequest	true	"Update shop request"
//
// @Success		200							{object}	updateResp
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops [PATCH]
func (h handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processUpdateShopRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Update: %v", err)
		response.Error(c, err)
		return
	}

	updatedShop, err := h.uc.Update(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Update: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newUpdateShopResp(updatedShop))
}

// @Summary		Get shop ID by user ID
// @Schemes		http https
// @Description	Get shop ID by user ID
// @Tags			Shop
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			id							path		string		true	"User ID"
//
// @Success		200							{string}	string		"Shop ID"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops/get-shop-id-by-user-id/{id} [GET]
func (h handler) GetShopIDByUserID(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processGetShopIDByUserIDRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.GetShopIDByUserID: %v", err)
		response.Error(c, err)
		return
	}

	shopID, err := h.uc.GetIDByUserID(ctx, sc, req.ID)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.GetShopIDByUserID: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, shopID)
}

// @Summary		Get all products
// @Schemes		http https
// @Description	Get all products
// @Tags			Shop
// @Accept			json
// @Produce		json
// @Success		200	{object}	getAllProductsResp
// @Failure		400	{object}	response.Resp	"Bad Request"
// @Failure		500	{object}	response.Resp	"Internal Server Error"
// @Router			/api/v1/shops/products/all [GET]
func (h handler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processGetAllRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.GetAll: %v", err)
		response.Error(c, err)
		return
	}
	//print req

	var pagQuery paginator.PaginatorQuery
	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Get.ShouldBindQuery: %v", err)
		response.Error(c, errWrongPaginationQuery)
		return
	}

	pagQuery.Adjust()

	p, err := h.uc.GetAll(ctx, sc, shop.GetProductOption{
		PagQuery:         pagQuery,
		GetProductFilter: req.toInput()})
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.GetAll: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newGetAllProductsResp(p))
}

func (h handler) Report(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processReportRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Report: %v", err)
		response.Error(c, err)
		return
	}

	report, err := h.uc.Report(ctx, sc)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Report: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newReportResponse(report))
}
