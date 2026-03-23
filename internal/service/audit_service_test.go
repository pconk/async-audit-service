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

func TestRecordLog(t *testing.T) {
	// 1. Setup
	mockRepo := new(MockAuditRepository)
	// Gunakan logger discard atau stdout untuk test
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	auditService := service.NewAuditService(mockRepo, logger)

	ctx := context.Background()
	dummyLog := entity.AuditLog{
		Username:  "test_user",
		Action:    "TEST_ACTION",
		SKU:       "ITEM-001",
		CreatedAt: time.Now(),
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
