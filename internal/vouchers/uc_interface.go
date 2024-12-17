package vouchers

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateVoucher(ctx context.Context, sc models.Scope, input CreateVoucherInput) (models.Voucher, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Voucher, error)
	List(ctx context.Context, sc models.Scope, opt GetVoucherFilter) ([]models.Voucher, error)
}
