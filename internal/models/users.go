package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"user_name"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	IsVerified bool               `bson:"is_verified"`
	Role       int                `bson:"role"`
	MediaID    primitive.ObjectID `bson:"media_id"`
}
