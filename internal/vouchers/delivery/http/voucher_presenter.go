package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
)

type createVoucherReq struct {
	Name                   string   `json:"name" binding:"required"`
	Code                   string   `json:"code" binding:"required"`
	ValidFrom              string   `json:"valid_from" binding:"required"`
	ValidTo                string   `json:"valid_to" binding:"required"`
	DiscountType           string   `json:"discount_type" binding:"required"`
	DiscountAmount         float64  `json:"discount_amount" binding:"required"`
	Description            string   `json:"description"`
	UsageLimit             uint     `json:"usage_limit"`
	ApplicableProductIDs   []string `json:"applicable_product_ids"`
	ApplicableCategorieIDs []string `json:"applicable_category_ids"`
	MinimumOrderAmount     float64  `json:"minimum_order_amount"`
	MaxDiscountAmount      float64  `json:"max_discount_amount"`
}

func (req createVoucherReq) validate() error {
	if _, err := util.StrToDateTime(req.ValidFrom); err != nil {
		return errWrongBody
	}
	if _, err := util.StrToDateTime(req.ValidTo); err != nil {
		return errWrongBody
	}

	if req.DiscountType != models.DiscountTypePercent && req.DiscountType != models.DiscountTypeFixed {
		return errWrongBody
	}

	if len(req.ApplicableProductIDs) > 0 {
		for _, id := range req.ApplicableProductIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	if len(req.ApplicableCategorieIDs) > 0 {
		for _, id := range req.ApplicableCategorieIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	return nil
}

func (h handler) newDetailVoucherResp(v models.Voucher) detailVoucherResp {
	return detailVoucherResp{
		ID:                     v.ID.Hex(),
		Name:                   v.Name,
		Code:                   v.Code,
		ValidFrom:              v.ValidFrom.Format(util.DateTimeFormat),
		ValidTo:                v.ValidTo.Format(util.DateTimeFormat),
		DiscountType:           v.DiscountType,
		DiscountAmount:         v.DiscountAmount,
		Description:            v.Description,
		UsageLimit:             v.UsageLimit,
		ApplicableProductIDs:   mongo.HexFromObjectIDsOrNil(v.ApplicableProductIDs),
		ApplicableCategorieIDs: mongo.HexFromObjectIDsOrNil(v.ApplicableCategorieIDs),
		MinimumOrderAmount:     v.MinimumOrderAmount,
		MaxDiscountAmount:      v.MaxDiscountAmount,
		ShopID:                 v.ShopID.Hex(),
		CreatedAt:              v.CreatedAt.Format(util.DateTimeFormat),
	}
}

func (req createVoucherReq) toInput() vouchers.CreateVoucherInput {
	validFrom, _ := util.StrToDateTime(req.ValidFrom)
	validTo, _ := util.StrToDateTime(req.ValidTo)

	return vouchers.CreateVoucherInput{
		Name:                   req.Name,
		Description:            req.Description,
		Code:                   req.Code,
		ValidFrom:              validFrom,
		ValidTo:                validTo,
		UsageLimit:             req.UsageLimit,
		ApplicableProductIDs:   req.ApplicableProductIDs,
		ApplicableCategorieIDs: req.ApplicableCategorieIDs,
		MinimumOrderAmount:     req.MinimumOrderAmount,
		DiscountType:           req.DiscountType,
		DiscountAmount:         req.DiscountAmount,
		MaxDiscountAmount:      req.MaxDiscountAmount,
	}
}

type detailVoucherResp struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Code                   string   `json:"code"`
	ValidFrom              string   `json:"valid_from"`
	ValidTo                string   `json:"valid_to"`
	DiscountType           string   `json:"discount_type"`
	DiscountAmount         float64  `json:"discount_amount"`
	Description            string   `json:"description,omitempty"`
	UsageLimit             uint     `json:"usage_limit"`
	ApplicableProductIDs   []string `json:"applicable_product_ids,omitempty"`
	ApplicableCategorieIDs []string `json:"applicable_category_ids,omitempty"`
	MinimumOrderAmount     float64  `json:"minimum_order_amount"`
	MaxDiscountAmount      float64  `json:"max_discount_amount"`
	CreatedAt              string   `json:"created_at"`
	ShopID                 string   `json:"shop_id"`
}
type DetailVoucherReq struct {
	ID   string `uri:"id"`
	Code string `uri:"code"`
}

func (req DetailVoucherReq) validate() error {
	if req.ID != "" {
		if !mongo.IsObjectID(req.ID) {
			return errWrongBody
		}
	}
	if req.ID == "" && req.Code == "" {
		return errWrongBody
	}

	return nil
}

func (req DetailVoucherReq) toInput() vouchers.DetailVoucherInput {
	return vouchers.DetailVoucherInput{
		ID:   req.ID,
		Code: req.Code,
	}
}

func (h handler) newDetailResponse(v models.Voucher) detailVoucherResp {
	return detailVoucherResp{
		ID:                     v.ID.Hex(),
		Name:                   v.Name,
		Code:                   v.Code,
		ValidFrom:              v.ValidFrom.Format(util.DateTimeFormat),
		ValidTo:                v.ValidTo.Format(util.DateTimeFormat),
		DiscountType:           v.DiscountType,
		DiscountAmount:         v.DiscountAmount,
		Description:            v.Description,
		UsageLimit:             v.UsageLimit,
		ApplicableProductIDs:   mongo.HexFromObjectIDsOrNil(v.ApplicableProductIDs),
		ApplicableCategorieIDs: mongo.HexFromObjectIDsOrNil(v.ApplicableCategorieIDs),
		MinimumOrderAmount:     v.MinimumOrderAmount,
		MaxDiscountAmount:      v.MaxDiscountAmount,

		CreatedAt: v.CreatedAt.Format(util.DateTimeFormat),
	}
}

type ListVoucherReq struct {
	IDs                    []string `form:"ids"`
	Codes                  []string `form:"codes"`
	ApplicableProductIDs   []string `form:"applicable_product_ids"`
	ApplicableCategorieIDs []string `form:"applicable_category_ids"`
	ShopID                 string   `form:"shop_id"`
	ValidFrom              string   `form:"valid_from"`
	ValidTo                string   `form:"valid_to"`
	Scope                  int      `form:"scope"`
}

type listVoucherResp struct {
	List []detailVoucherResp `json:"list"`
}

func (r ListVoucherReq) validate() error {
	if r.ValidFrom != "" {
		if _, err := util.StrToDateTime(r.ValidFrom); err != nil {
			return errWrongBody
		}
	}
	if r.ValidTo != "" {
		if _, err := util.StrToDateTime(r.ValidTo); err != nil {
			return errWrongBody
		}
	}

	if len(r.IDs) > 0 {
		for _, id := range r.IDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	if len(r.Codes) > 0 {
		for _, code := range r.Codes {
			if !mongo.IsObjectID(code) {
				return errWrongBody
			}
		}
	}

	if len(r.ApplicableProductIDs) > 0 {
		for _, id := range r.ApplicableProductIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	if len(r.ApplicableCategorieIDs) > 0 {
		for _, id := range r.ApplicableCategorieIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	return nil
}
func (r ListVoucherReq) toInput() vouchers.GetVoucherFilter {

	validTo, _ := util.StrToDateTime(r.ValidTo)
	validFrom, _ := util.StrToDateTime(r.ValidFrom)
	return vouchers.GetVoucherFilter{
		IDs:                    r.IDs,
		Codes:                  r.Codes,
		ApplicableProductIDs:   r.ApplicableProductIDs,
		ApplicableCategorieIDs: r.ApplicableCategorieIDs,
		ValidFrom:              &validFrom,
		ValidTo:                &validTo,
		ShopID:                 r.ShopID,
	}
}
func (h handler) newListResponse(vouchers []models.Voucher) listVoucherResp {
	var res listVoucherResp
	for _, v := range vouchers {
		res.List = append(res.List, h.newDetailResponse(v))
	}
	return res
}

type applyVoucherReq struct {
	ID          string  `json:"id"`
	Code        string  `json:"code""`
	OrderAmount float64 `json:"order_amount" binding:"required"`
}

func (req applyVoucherReq) validate() error {
	if req.ID != "" {
		if !mongo.IsObjectID(req.ID) {
			return errWrongBody
		}
	}
	if req.ID == "" && req.Code == "" {
		return errWrongBody
	}
	if req.OrderAmount < 0 {
		return errWrongBody
	}

	return nil
}

func (req applyVoucherReq) toInput() vouchers.ApplyVoucherInput {
	return vouchers.ApplyVoucherInput{
		ID:          req.ID,
		Code:        req.Code,
		OrderAmount: req.OrderAmount,
	}
}

type applyVoucherResp struct {
	Voucher        detailVoucherResp `json:"voucher"`
	DiscountAmount float64           `json:"discount_amount"`
	OrderAmount    float64           `json:"order_amount"`
}

func (h handler) newApplyVoucherResponse(v models.Voucher, orderAmount float64, discountAmount float64) applyVoucherResp {
	return applyVoucherResp{
		Voucher:        h.newDetailResponse(v),
		OrderAmount:    orderAmount,
		DiscountAmount: discountAmount,
	}
}
