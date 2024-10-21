package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, input shop.CreateInput) (models.Shop, error) {
	if input.City == "" || input.Name == "" || input.Street == "" {
		return models.Shop{}, shop.ErrInvalidInput
	}

	_, err := uc.repo.Detail(ctx, sc, "")
	if err == nil {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", shop.ErrShopExist)
		return models.Shop{}, shop.ErrShopExist
	}

	if !util.IsValidPhone(input.Phone) {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", shop.ErrInvalidPhone)
		return models.Shop{}, shop.ErrInvalidPhone
	}

	opt := shop.CreateOption{
		Name:     input.Name,
		Alias:    util.BuildAlias(input.Name),
		City:     input.City,
		Street:   input.Street,
		District: input.District,
		Phone:    input.Phone,
	}

	sh, err := uc.repo.Create(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", err)
		return models.Shop{}, err
	}

	return sh, nil
}

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, input shop.GetInput) (shop.GetOutput, error) {
	opt := shop.GetOption{
		GetShopsFilter: input.GetShopsFilter,
		PagQuery:       input.PagQuery,
	}
	s, pag, err := uc.repo.Get(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Get: %v", err)
		return shop.GetOutput{}, err
	}

	return shop.GetOutput{
		Shops: s,
		Pag:   pag,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (models.Shop, error) {
	s, err := uc.repo.Detail(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Detail: %v", err)
		return models.Shop{}, err
	}

	return s, nil
}
func (uc implUsecase) Delete(ctx context.Context, sc models.Scope) error {

	err := uc.repo.Delete(ctx, sc)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Delete.Repodele", err)
		return shop.ErrShopDoesNotExist
	}

	return nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, input shop.UpdateInput) (models.Shop, error) {

	shop, err := uc.repo.Update(ctx, sc, shop.UpdateOption{
		Name:     input.Name,
		Alias:    input.Alias,
		City:     input.City,
		District: input.District,
		Street:   input.Street,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.update.repo.update:", err)
		return models.Shop{}, err
	}
	return shop, nil
}
