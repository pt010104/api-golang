package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildCheckoutDetailQuery(ctx context.Context, sc models.Scope, checkoutID string) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "Checkout.Repo.buildCheckoutDetailQuery.BuildScopeQuery", err)
		return nil, err
	}

	filter = mongo.BuildQueryWithSoftDelete(filter)

	if checkoutID != "" {
		filter["_id"] = mongo.ObjectIDFromHexOrNil(checkoutID)
	}

	return filter, nil
}
