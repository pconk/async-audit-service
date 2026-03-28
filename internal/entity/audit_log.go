package entity

import (
	"async-audit-service/pb"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuditLog struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          int64              `bson:"user_id"`
	Username        string             `bson:"username"`
	WarehouseID     string             `bson:"warehouse_id"`
	Role            string             `bson:"role"`
	Action          string             `bson:"action"`
	SKU             string             `bson:"sku"`
	ProductName     string             `bson:"product_name"`
	QuantityChanged int32              `bson:"quantity_changed"`
	FinalStock      int32              `bson:"final_stock"`
	Metadata        map[string]string  `bson:"metadata"`
	CreatedAt       time.Time          `bson:"created_at"`
}

// RecentLog digunakan khusus untuk Query pengambilan data terbaru
type RecentLog struct {
	ID              primitive.ObjectID `bson:"_id"`
	UserID          int64              `bson:"user_id"`
	Username        string             `bson:"username"`
	WarehouseID     string             `bson:"warehouse_id"`
	Role            string             `bson:"role"`
	Action          string             `bson:"action"`
	SKU             string             `bson:"sku"`
	ProductName     string             `bson:"product_name"`
	QuantityChanged int32              `bson:"quantity_changed"`
	FinalStock      int32              `bson:"final_stock"`
	CreatedAt       time.Time          `bson:"created_at"`
}

func (l *RecentLog) ToProto() *pb.RecentLog {
	return &pb.RecentLog{
		Id:              l.ID.Hex(),
		UserId:          l.UserID,
		Username:        l.Username,
		WarehouseId:     l.WarehouseID,
		Role:            l.Role,
		Action:          l.Action,
		Sku:             l.SKU,
		ProductName:     l.ProductName,
		QuantityChanged: l.QuantityChanged,
		FinalStock:      l.FinalStock,
		CreatedAt:       timestamppb.New(l.CreatedAt),
	}
}
