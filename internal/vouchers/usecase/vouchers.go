package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUsecase) CreateVoucher(ctx context.Context, sc models.Scope, input vouchers.CreateVoucherInput) (models.Voucher, error) {
	role := sc.Role

	if err := uc.validateCreateVoucher(ctx, input); err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.validateCreateVoucher: %v", err)
		return models.Voucher{}, err
	}

	opt := vouchers.CreateVoucherOption{
		Data: vouchers.VoucherData{
			Name:                   input.Name,
			Code:                   input.Code,
			Description:            input.Description,
			ValidFrom:              input.ValidFrom,
			ValidTo:                input.ValidTo,
			DiscountType:           input.DiscountType,
			DiscountAmount:         input.DiscountAmount,
			MaxDiscountAmount:      input.MaxDiscountAmount,
			UsageLimit:             input.UsageLimit,
			MinimumOrderAmount:     input.MinimumOrderAmount,
			ApplicableProductIDs:   input.ApplicableProductIDs,
			ApplicableCategorieIDs: input.ApplicableCategorieIDs,
			CreatedBy:              sc.UserID,
		},
	}

	opt.Data.ShopIDs = []string{sc.ShopID}
	opt.Data.Scope = models.ScopeShop

	if role == models.RoleAdmin {
		if len(input.ShopIDs) > 0 {
			isVerified := true
			s, err := uc.shopUc.ListShop(ctx, models.Scope{},
				shop.GetShopsFilter{
					IDs:        input.ShopIDs,
					IsVerified: &isVerified,
				},
			)
			if err != nil {
				uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.ListShop: %v", vouchers.ErrShopNotFound)
				return models.Voucher{}, vouchers.ErrShopNotFound
			}

			if len(s) != len(input.ShopIDs) {
				uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.ListShop: %v", vouchers.ErrShopNotFound)
				return models.Voucher{}, vouchers.ErrShopNotFound
			}

			opt.Data.ShopIDs = input.ShopIDs
		} else {
			opt.Data.Scope = models.ScopeAll
		}
	}

	v, err := uc.repo.CreateVoucher(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.CreateVoucher: %v", err)
		return models.Voucher{}, err
	}

	return v, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, input vouchers.DetailVoucherInput) (models.Voucher, error) {
	opt := vouchers.DetailVoucherOption{
		ID:   input.ID,
		Code: input.Code,
	}

	v, err := uc.repo.DetailVoucher(ctx, sc, opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Voucher{}, vouchers.ErrVoucherNotFound
		}
		uc.l.Errorf(ctx, "vouchers.usecase.Detail.DetailVoucher: %v", err)
		return models.Voucher{}, err

	}
	return v, nil
}

func (uc implUsecase) List(ctx context.Context, sc models.Scope, opt vouchers.GetVoucherFilter) ([]models.Voucher, error) {

	vouchers1, err := uc.repo.ListVoucher(ctx, sc, opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Voucher{}, vouchers.ErrVoucherNotFound
		}

		uc.l.Errorf(ctx, "vouchers.usecase.List.ListVoucher: %v", err)
		return []models.Voucher{}, err
	}
	return vouchers1, nil

}

func (uc implUsecase) ApplyVoucher(ctx context.Context, sc models.Scope, input vouchers.ApplyVoucherInput) (models.Voucher, error) {
	voucher, err := uc.repo.DetailVoucher(ctx, sc, vouchers.DetailVoucherOption{
		ID:   input.ID,
		Code: input.Code,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Voucher{}, vouchers.ErrVoucherNotFound
		}
		uc.l.Errorf(ctx, "vouchers.usecase.ApplyVoucher.DetailVoucher: %v", err)
		return models.Voucher{}, err
	}

	if voucher.ValidFrom.After(util.Now()) || voucher.ValidTo.Before(util.Now()) {
		return models.Voucher{}, vouchers.ErrVoucherExpired
	}
	if voucher.UsageLimit > 0 && uint(voucher.UsedCount) >= voucher.UsageLimit {
		return models.Voucher{}, vouchers.ErrVoucherExpired
	}
	if voucher.MinimumOrderAmount > 0 {
		if input.OrderAmount < voucher.MinimumOrderAmount {
			return models.Voucher{}, vouchers.ErrVoucherMinimumOrderAmount
		}
	}

	voucher.UsedCount++
	nv, err := uc.repo.UpdateVoucher(ctx, sc, vouchers.UpdateVoucherOption{
		Model: voucher,
	})
	if err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.ApplyVoucher.UpdateVoucher: %v", err)
		return models.Voucher{}, err
	}

	return nv, nil
}
