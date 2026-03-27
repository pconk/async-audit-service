package service_test

import (
	"async-audit-service/internal/entity"
	"async-audit-service/internal/service"
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuditRepository adalah mock untuk repository
type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) Insert(ctx context.Context, log entity.AuditLog) error {
	// Merekam pemanggilan argumen untuk dicek nanti
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAuditRepository) GetRecentLogs(ctx context.Context, limit int) ([]entity.RecentLog, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.RecentLog), args.Error(1)
}

func (m *MockAuditRepository) CreateIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestRecordLog(t *testing.T) {
	// 1. Setup
	mockRepo := new(MockAuditRepository)
	// Gunakan logger discard atau stdout untuk test
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	auditService := service.NewAuditService(mockRepo, logger)

	ctx := context.Background()
	dummyLog := entity.AuditLog{
		UserID:          "user-123",
		Username:        "test_user",
		Action:          "ADJUST_STOCK",
		SKU:             "ITEM-001",
		ProductName:     "Indomie Goreng",
		WarehouseID:     "1",
		Role:            "admin",
		QuantityChanged: -5,
		FinalStock:      10,
		CreatedAt:       time.Now(),
	}

	t.Run("Success Record Log", func(t *testing.T) {
		// Expectation: Repository Insert dipanggil 1 kali dengan argumen apa saja, return nil (sukses)
		mockRepo.On("Insert", ctx, dummyLog).Return(nil).Once()

		// Execute
		err := auditService.RecordLog(ctx, dummyLog)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Record Log (DB Error)", func(t *testing.T) {
		// Expectation: Repository return error
		mockRepo.On("Insert", ctx, dummyLog).Return(errors.New("db connection error")).Once()

		err := auditService.RecordLog(ctx, dummyLog)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRecentLogs(t *testing.T) {
	// 1. Setup
	mockRepo := new(MockAuditRepository)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	auditService := service.NewAuditService(mockRepo, logger)
	ctx := context.Background()

	t.Run("Success Get Recent Logs", func(t *testing.T) {
		limit := 2
		expectedLogs := []entity.RecentLog{
			{
				UserID:   "user-1",
				Username: "admin",
				Action:   "STOCK_IN",
				SKU:      "SKU-001",
			},
			{
				UserID:   "user-2",
				Username: "staff",
				Action:   "STOCK_OUT",
				SKU:      "SKU-002",
			},
		}

		mockRepo.On("GetRecentLogs", ctx, limit).Return(expectedLogs, nil).Once()

		result, err := auditService.GetRecentLogs(ctx, limit)

		assert.NoError(t, err)
		assert.Equal(t, len(expectedLogs), len(result))
		assert.Equal(t, "admin", result[0].Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Recent Logs (DB Error)", func(t *testing.T) {
		mockRepo.On("GetRecentLogs", ctx, 5).Return(nil, errors.New("query error")).Once()

		_, err := auditService.GetRecentLogs(ctx, 5)
		assert.Error(t, err)
	})
}
