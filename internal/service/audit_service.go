package service

import (
	"async-audit-service/internal/entity"
	"async-audit-service/internal/repository"
	"context"
	"log/slog"
)

type AuditService interface {
	RecordLog(ctx context.Context, log entity.AuditLog) error
}

type auditService struct {
	repo   repository.AuditRepository
	logger *slog.Logger
}

func NewAuditService(repo repository.AuditRepository, logger *slog.Logger) AuditService {
	return &auditService{
		repo:   repo,
		logger: logger,
	}
}

func (s *auditService) RecordLog(ctx context.Context, log entity.AuditLog) error {
	// Di sini bisa ditaruh logic tambahan sebelum save ke DB
	// Misal: Validasi data kosong, masking data sensitif, dll.
	return s.repo.Insert(ctx, log)
}
