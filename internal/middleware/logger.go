package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// key untuk context
type ctxKey string

const (
	LogFieldsKey ctxKey = "log_fields"
	RequestIDKey ctxKey = "request_id"
)

// AddLogFields memungkinkan middleware lain menambahkan field ke log utama
func AddLogFields(ctx context.Context, fields ...slog.Attr) {
	if v, ok := ctx.Value(LogFieldsKey).(*[]slog.Attr); ok {
		*v = append(*v, fields...)
	}
}

// CreateTraceGroup menggabungkan request_id dan field tambahan ke dalam satu grup log trace
func CreateTraceGroup(requestID string, extraFields []slog.Attr) slog.Attr {
	traceAttrs := make([]any, 0, len(extraFields)+1)
	traceAttrs = append(traceAttrs, slog.String("request_id", requestID))
	for _, attr := range extraFields {
		traceAttrs = append(traceAttrs, attr)
	}
	return slog.Group("trace", traceAttrs...)
}

// DurationToMs mengonversi time.Duration menjadi float64 milidetik
func DurationToMs(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / 1e6
}

// GetRequestID mengambil request ID dari context (untuk dipakai di handler/service)
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return v
	}
	return ""
}

type LoggerInterceptor struct {
	logger *slog.Logger
}

func NewLoggerInterceptor(logger *slog.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{logger: logger}
}

func (i *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. Cek apakah ada Request ID di metadata (dikirim oleh Gateway/Warehouse)
		md, ok := metadata.FromIncomingContext(ctx)
		var reqID string
		if ok {
			if ids := md.Get("x-request-id"); len(ids) > 0 {
				reqID = ids[0]
			}
		}

		// 2. Jika tidak ada di metadata, baru generate sendiri
		if reqID == "" {
			reqID = generateRequestID()
		}

		// Simpan ID ke context
		ctx = context.WithValue(ctx, RequestIDKey, reqID)
		start := time.Now()

		// 3. Siapkan wadah untuk log tambahan (username, role, dll)
		var extraFields []slog.Attr
		ctx = context.WithValue(ctx, LogFieldsKey, &extraFields)

		// 4. Kirim balik Request ID ke Client lewat Header Response gRPC
		header := metadata.Pairs("x-request-id", reqID)
		grpc.SendHeader(ctx, header)

		// Panggil handler berikutnya (bisa berupa AuthInterceptor atau Handler utama)
		resp, err := handler(ctx, req)

		// Tentukan status, level, dan message
		level := slog.LevelInfo
		statusStr := "OK"
		msg := "gRPC Request Success"

		if err != nil {
			st, _ := status.FromError(err)
			statusStr = st.Code().String()
			level = slog.LevelError
			msg = "gRPC Request Failed"
			// Masukkan error ke dalam extra fields agar masuk ke group trace
			extraFields = append(extraFields, slog.String("error", err.Error()))
		}

		// Log dengan struktur Grouping (trace & grpc)
		i.logger.LogAttrs(ctx, level, msg,
			CreateTraceGroup(reqID, extraFields),
			slog.Group("grpc",
				slog.String("method", info.FullMethod),
				slog.String("status", statusStr),
				slog.Float64("duration_ms", DurationToMs(time.Since(start))),
			),
		)

		return resp, err
	}
}

func generateRequestID() string {
	b := make([]byte, 8) // 16 karakter hex
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
