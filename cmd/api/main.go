package main

import (
	"async-audit-service/internal/config"
	grpcHandler "async-audit-service/internal/handler/grpc"
	"async-audit-service/internal/middleware"
	"async-audit-service/internal/repository"
	"async-audit-service/internal/service"
	"async-audit-service/pb"
	"async-audit-service/pkg/database"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file (jika ada)
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, relying on system environment variables or defaults")
	}

	// 2. Load Config
	cfg := config.LoadConfig()

	// 3. Connect to MongoDB
	db, err := database.InitMongoDB(cfg.MongoURI, cfg.DbName)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	logger.Info("Connected to MongoDB")

	// 4. Init Layers (Dependency Injection)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo, logger)
	auditHandler := grpcHandler.NewAuditHandler(auditService, logger)

	// Init Middleware
	// Init Middleware (Pastikan ini ada)
	authInterceptor := middleware.NewAuthInterceptor(cfg)
	loggerInterceptor := middleware.NewLoggerInterceptor(logger)

	// 5. Setup gRPC Server
	lis, err := net.Listen("tcp", ":"+cfg.AppPort)
	if err != nil {
		logger.Error("Failed to listen", "port", cfg.AppPort, "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		// ChainUnaryInterceptor mengeksekusi urutan dari kiri ke kanan (atau atas ke bawah)
		grpc.ChainUnaryInterceptor(
			loggerInterceptor.Unary(), // Log dulu (Start Timer)
			authInterceptor.Unary(),   // Baru cek Auth
		),
	)
	pb.RegisterAuditServiceServer(grpcServer, auditHandler)

	// Enable reflection (agar bisa ditest pakai Postman / grpcurl)
	reflection.Register(grpcServer)

	// 6. Jalankan Server
	go func() {
		logger.Info("Starting gRPC server", "port", cfg.AppPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Failed to serve gRPC", "error", err)
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
	grpcServer.GracefulStop()
	logger.Info("Server stopped")
}
