package grpc

import (
	"async-audit-service/internal/entity"
	"async-audit-service/internal/middleware"
	"async-audit-service/internal/service"
	"async-audit-service/pb" // Asumsi kode generate ada di sini
	"context"
	"log/slog"
)

type AuditHandler struct {
	pb.UnimplementedAuditServiceServer
	service service.AuditService
	logger  *slog.Logger
}

func NewAuditHandler(service service.AuditService, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuditHandler) LogActivity(ctx context.Context, req *pb.AuditRequest) (*pb.AuditResponse, error) {
	// Konversi Proto ke Entity Domain
	auditLog := entity.AuditLog{
		UserID:          req.UserId, // Sekarang sudah int64 sesuai proto terbaru
		Username:        req.Username,
		WarehouseID:     req.WarehouseId,
		Role:            req.Role,
		Action:          req.Action,
		SKU:             req.Sku,
		ProductName:     req.ProductName,
		QuantityChanged: req.QuantityChanged,
		FinalStock:      req.FinalStock,
		Metadata:        req.Metadata,
		CreatedAt:       req.Timestamp.AsTime(),
	}

	// Ambil Request ID dari Context (yg digenerate middleware)
	reqID := middleware.GetRequestID(ctx)

	h.logger.Info("Received LogActivity request", "request_id", reqID, "action", req.Action, "sku", req.Sku, "user", req.Username)

	err := h.service.RecordLog(ctx, auditLog)
	if err != nil {
		h.logger.Error("Database operation failed for LogActivity", "request_id", reqID, "error", err)
		return &pb.AuditResponse{
			Success: false,
		}, err
	}

	return &pb.AuditResponse{
		Success: true,
	}, nil
}

func (h *AuditHandler) GetRecentLogs(ctx context.Context, req *pb.GetRecentLogsRequest) (*pb.GetRecentLogsResponse, error) {
	reqID := middleware.GetRequestID(ctx)
	h.logger.Info("Processing GetRecentLogs", "request_id", reqID, "limit", req.Limit)

	// 1. Ambil data dari repository (hasilnya []entity.RecentLog)
	logs, err := h.service.GetRecentLogs(ctx, int(req.Limit))
	if err != nil {
		h.logger.Error("Failed to fetch recent logs", "request_id", reqID, "error", err)
		return nil, err
	}

	// 2. Mapping & Konversi
	pbLogs := make([]*pb.RecentLog, 0, len(logs))
	for _, l := range logs {
		pbLogs = append(pbLogs, l.ToProto())
	}

	return &pb.GetRecentLogsResponse{Logs: pbLogs}, nil
}
