package repository

import (
	"async-audit-service/internal/entity"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditRepository interface mendefinisikan kontrak database
type AuditRepository interface {
	Insert(ctx context.Context, log entity.AuditLog) error
	GetRecentLogs(ctx context.Context, limit int) ([]entity.RecentLog, error)
	CreateIndexes(ctx context.Context) error
}

type auditRepository struct {
	db *mongo.Database
}

func NewAuditRepository(db *mongo.Database) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) CreateIndexes(ctx context.Context) error {
	collection := r.db.Collection("audit_logs")

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}}, // -1 = Descending
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (r *auditRepository) Insert(ctx context.Context, log entity.AuditLog) error {
	collection := r.db.Collection("audit_logs")

	// Pastikan timeout agar tidak hanging
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func (r *auditRepository) GetRecentLogs(ctx context.Context, limit int) ([]entity.RecentLog, error) {
	collection := r.db.Collection("audit_logs")

	// 1. Setup Timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 2. Tentukan Filter (kosong karena ambil semua) dan Opsi (Sort & Limit)
	filter := bson.M{}
	findOptions := options.Find()

	// Order By created_at DESC (-1 untuk descending, 1 untuk ascending)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Batasi jumlah baris
	findOptions.SetLimit(int64(limit))

	// 3. Eksekusi Find
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 4. Decode hasil ke slice entity
	var logs []entity.RecentLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
