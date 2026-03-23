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

	h.logger.Info("Received audit log", "request_id", reqID, "action", req.Action, "sku", req.Sku)

	err := h.service.RecordLog(ctx, auditLog)
	if err != nil {
		h.logger.Error("Failed to record log", "error", err)
		return &pb.AuditResponse{
			Success: false,
		}, err
	}

	return &pb.AuditResponse{
		Success: true,
	}, nil
}
