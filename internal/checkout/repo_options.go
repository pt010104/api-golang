package checkout

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateOption struct {
	ProductIDs []primitive.ObjectID
}
