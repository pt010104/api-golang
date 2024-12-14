package http

import (
	"fmt"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/paginator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createProductReq struct {
	Name            string   `json:"name" binding:"required"`
	Price           float32  `json:"price" binding:"required"`
	StockLevel      uint     `json:"stock_level" binding:"required"`
	ReorderLevel    *uint    `json:"reorder_level" binding:"required" `
	ReorderQuantity *uint    `json:"reorder_quantity" binding:"required"`
	CategoryIDs     []string `json:"category_ids" binding:"required"`
	MediaIDs        []string `json:"media_ids"`
}

func (r createProductReq) toInput() shop.CreateProductInput {
	return shop.CreateProductInput{
		Name:  r.Name,
		Price: r.Price,

		StockLevel:      r.StockLevel,
		ReorderLevel:    r.ReorderLevel,
		ReorderQuantity: r.ReorderQuantity,
		CategoryID:      r.CategoryIDs,
		MediaIDs:        r.MediaIDs,
	}
}
func (r createProductReq) validate() error {

	if r.Name == "" {
		fmt.Errorf("wrong name")
		return errWrongBody
	}

	if r.Price <= 0 {
		fmt.Errorf("wrong price")
		return errWrongBody
	}

	if r.StockLevel == 0 {
		fmt.Errorf("wrong stock level")
		return errWrongBody
	}

	if r.ReorderLevel == nil || *r.ReorderLevel == 0 {
		fmt.Errorf("wrong reorder level")
		return errWrongBody
	}

	if r.ReorderQuantity == nil || *r.ReorderQuantity == 0 {
		fmt.Errorf("wrong reorder quantity")
		return errWrongBody
	}

	for _, id := range r.CategoryIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			fmt.Errorf("wrong ids")
			return errWrongBody
		}
	}
	if len(r.MediaIDs) > 0 {
		for _, id := range r.MediaIDs {
			if _, err := primitive.ObjectIDFromHex(id); err != nil {
				fmt.Errorf("wrong ids")
				return errWrongBody
			}
		}
	}

	return nil
}

type detailProductReq struct {
	ID string `uri:"id" binding:"required"`
}
type avatar_obj struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}
type detailProductResp struct {
	ID            string       `json:"id" binding:"required"`
	Name          string       `json:"name" binding:"required"`
	CategoryName  []string     `json:"category_name" binding:"required"`
	CategoryID    []string     `json:"category_id" binding:"required"`
	ShopName      string       `json:"shop_name" binding:"required"`
	ShopID        string       `json:"shop_id" binding:"required"`
	InventoryName string       `json:"inventory_name" binding:"required"`
	Price         float32      `json:"price" binding:"required"`
	Avatar        []avatar_obj `json:"avatar,omitempty"`
}

func (h handler) newDetailProductResponse(p shop.DetailProductOutput) detailProductResp {
	categoryIDs := make([]string, len(p.Category))
	for i, category := range p.Category {
		categoryIDs[i] = category.ID.Hex()
	}

	var images []avatar_obj
	for _, media := range p.Medias {
		images = append(images, avatar_obj{
			MediaID: media.ID.Hex(),
			URL:     media.URL,
		})
	}
	return detailProductResp{
		ID:            p.ID,
		Name:          p.Name,
		CategoryName:  p.CategoryName,
		CategoryID:    categoryIDs,
		ShopName:      p.Shop.Name,
		ShopID:        p.Shop.ID.Hex(),
		InventoryName: p.Inventory.ID.Hex(),
		Price:         p.Price,
		Avatar:        images,
	}

}

type deleteProductRequest struct {
	IDs []string `json:"ids"`
}

type listProductRequest struct {
	IDs    []string `form:"ids"`
	Search string   `form:"search"`
	ShopID string   `form:"shop_id"`
}
type getProductRequest struct {
	IDs     []string `form:"ids"`
	Search  string   `form:"search"`
	ShopID  string   `form:"shop_id"`
	CateIDs []string `form:"category_ids"`
}

func (r listProductRequest) validate() error {
	if len(r.IDs) > 0 {
		for _, id := range r.IDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	return nil
}
func (r getProductRequest) validate() error {
	if len(r.IDs) > 0 {
		for _, id := range r.IDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}
	if len(r.CateIDs) > 0 {
		for _, id := range r.CateIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}
	if r.ShopID != "" && !mongo.IsObjectID(r.ShopID) {

		return errWrongBody
	}

	return nil
}

func (r listProductRequest) toInput() shop.GetProductFilter {
	return shop.GetProductFilter{
		IDs:    r.IDs,
		Search: r.Search,
		ShopID: r.ShopID,
	}
}
func (r getProductRequest) toInput() shop.GetProductFilter {
	return shop.GetProductFilter{
		IDs:     r.IDs,
		Search:  r.Search,
		ShopID:  r.ShopID,
		CateIDs: r.CateIDs,
	}
}

type listProductMetaResponse struct {
	paginator.PaginatorResponse
}
type getProductResp struct {
	meta       listMetaResponse
	Items      []listProductItem `json:"products"`
	ShopObject shopObject        `json:"shop_object"`
}
type categoryObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type AvatarProduct_obj struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}
type listProductItem struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	ShopID          string           `json:"shop_id"`
	InventoryID     string           `json:"inventory_id"`
	Price           float32          `json:"price"`
	CategoryObjects []categoryObject `json:"category_objects"`
	Avatar          []avatar_obj     `json:"avatar,omitempty"`
}

type shopObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type listProductResp struct {
	Products   []listProductItem `json:"products"`
	ShopObject shopObject        `json:"shop_object"`
}

func (h handler) getProductResp(output shop.GetProductOutput) getProductResp {
	var list []listProductItem

	for _, s := range output.Products {
		var categories []categoryObject
		for _, category := range s.Cate {
			categories = append(categories, categoryObject{
				ID:   category.ID.Hex(),
				Name: category.Name,
			})
		}

		var image []avatar_obj
		for _, media := range s.Images {
			image = append(image, avatar_obj{
				MediaID: media.ID.Hex(),
				URL:     media.URL,
			})
		}
		item := listProductItem{
			ID:              s.P.ID.Hex(),
			Name:            s.P.Name,
			InventoryID:     s.Inven,
			Price:           s.P.Price,
			CategoryObjects: categories,
			ShopID:          s.P.ShopID.Hex(),
			Avatar:          image,
		}
		list = append(list, item)
	}

	shopObject := shopObject{
		ID:   output.Shop.ID.Hex(),
		Name: output.Shop.Name,
	}

	return getProductResp{
		meta: listMetaResponse{
			PaginatorResponse: output.Pag.ToResponse(),
		},
		Items:      list,
		ShopObject: shopObject,
	}
}

type UpdateProductReq struct {
	ID              string   `json:"id" binding:"required"`
	Name            string   `json:"name"`
	Price           float32  `json:"price"`
	StockLevel      uint     `json:"stock_level"`
	ReorderLevel    uint     `json:"reorder_level"`
	ReorderQuantity uint     `json:"reorder_quantity"`
	CategoryIDs     []string `json:"category_ids"`
	MediaIDs        []string `json:"media_ids"`
}

func (r UpdateProductReq) toInput() shop.UpdateProductOption {
	return shop.UpdateProductOption{
		ID:              mongo.ObjectIDFromHexOrNil(r.ID),
		Name:            r.Name,
		Price:           r.Price,
		StockLevel:      r.StockLevel,
		ReorderLevel:    r.ReorderLevel,
		ReorderQuantity: r.ReorderQuantity,
		CategoryID:      mongo.ObjectIDsFromHexOrNil(r.CategoryIDs),
		MediaIDs:        mongo.ObjectIDsFromHexOrNil(r.MediaIDs),
	}
}
func (r UpdateProductReq) validate() error {
	if r.ID == "" || !mongo.IsObjectID(r.ID) {
		return errWrongBody
	}
	if len(r.CategoryIDs) > 0 {
		for _, id := range r.CategoryIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}

	}
	if len(r.MediaIDs) > 0 {
		for _, id := range r.MediaIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody

			}
		}

	}
	return nil
}

type updateProductResp struct {
	ID            string       `json:"id" binding:"required"`
	Name          string       `json:"name,omitempty"`          // Empty strings won't show
	CategoryName  []string     `json:"category_name,omitempty"` // Empty slices won't show
	CategoryID    []string     `json:"category_id,omitempty"`
	ShopName      string       `json:"shop_name,omitempty"`
	ShopID        string       `json:"shop_id,omitempty"`
	InventoryName string       `json:"inventory_name,omitempty"`
	Price         float32      `json:"price,omitempty"` // Zero values won't show
	Avatar        []avatar_obj `json:"avatar,omitempty"`
}

func (h handler) newUpdateProductResponse(p models.Product) updateProductResp {

	categoryIDs := make([]string, len(p.CategoryID))
	for i, category := range p.CategoryID {
		categoryIDs[i] = category.Hex()
	}

	var images []avatar_obj
	for _, media := range p.MediaIDs {
		images = append(images, avatar_obj{
			MediaID: media.Hex(),
			URL:     media.Hex(),
		})
	}
	return updateProductResp{
		ID:            p.ID.Hex(),
		Name:          p.Name,
		CategoryName:  []string{},
		CategoryID:    categoryIDs,
		ShopName:      p.ShopID.Hex(),
		ShopID:        p.ShopID.Hex(),
		InventoryName: p.InventoryID.Hex(),
		Price:         p.Price,
		Avatar:        images,
	}

}
