package entity

import (
	"time"
)

type AuditLog struct {
	ID              string            `bson:"_id,omitempty"` // Biarkan Mongo generate ObjectID jika kosong
	Username        string            `bson:"username"`
	WarehouseID     string            `bson:"warehouse_id"`
	Role            string            `bson:"role"`
	Action          string            `bson:"action"`
	SKU             string            `bson:"sku"`
	ProductName     string            `bson:"product_name"`
	QuantityChanged int32             `bson:"quantity_changed"`
	FinalStock      int32             `bson:"final_stock"`
	Metadata        map[string]string `bson:"metadata"`
	CreatedAt       time.Time         `bson:"created_at"`
}
