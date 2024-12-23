package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildInventoryDetailQuery(ctx context.Context, id primitive.ObjectID) (bson.M, error) {
	filter := bson.M{}

	filter["_id"] = id

	return filter, nil
}

func (repo implRepo) buildInventoryQuery(ctx context.Context, sc models.Scope, IDs []primitive.ObjectID) (bson.M, error) {
	filter, err := mongo.BuildShopScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.buildGetInventoryQuery.BuildShopScopeQuery: %v", err)
		return bson.M{}, err
	}

	if len(IDs) > 0 {
		filter["_id"] = bson.M{"$in": IDs}
	}

	return filter, nil
}
