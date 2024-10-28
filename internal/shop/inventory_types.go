package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Inventory
type CreateInventoryInput struct {
	ProductID       primitive.ObjectID
	StockLevel      uint
	ReorderLevel    *uint
	ReorderQuantity *uint
}

type CreateInventoryOutput struct {
	Inventory models.Inventory
}

type UpdateInventoryInput struct {
	ProductID       primitive.ObjectID
	StockLevel      *uint
	ReorderLevel    *uint
	ReorderQuantity *uint
}