package middleware

import (
	"async-audit-service/internal/config"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtSecret string
}

func NewAuthInterceptor(cfg *config.Config) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSecret: cfg.JWTSecret,
	}
}

// Unary mengembalikan interceptor untuk validasi JWT pada setiap request
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. Ambil Metadata (Header)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		// 2. Ambil Header Authorization
		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		// Format: "Bearer <token>"
		tokenString := strings.Replace(values[0], "Bearer ", "", 1)

		// 3. Parse & Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(i.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// 4. Extract Data dari Claims (Username & WarehouseID)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Helper function untuk ambil string dengan aman
			getString := func(key string) string {
				if v, ok := claims[key].(string); ok {
					return v
				}
				return ""
			}

			// Tambahkan info user ke Log LoggerInterceptor
			AddLogFields(ctx, slog.String("jwt_username", getString("username")), slog.String("jwt_warehouse_id", getString("warehouse_id")))
		}

		// Jika valid, lanjutkan ke handler utama
		return handler(ctx, req)
	}
}
