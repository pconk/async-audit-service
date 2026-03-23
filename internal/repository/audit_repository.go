package repository

import (
	"async-audit-service/internal/entity"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// AuditRepository interface mendefinisikan kontrak database
type AuditRepository interface {
	Insert(ctx context.Context, log entity.AuditLog) error
}

type auditRepository struct {
	db *mongo.Database
}

func NewAuditRepository(db *mongo.Database) AuditRepository {
	return &auditRepository{db: db}
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
