# ==========================================
# Test Script gRPC Audit Service (PowerShell)
# ==========================================
# Requirement:
# 1. Service Audit jalan (go run internal/cmd/api/main.go)
# 2. Tool 'grpcurl' terinstall (go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest)

# --- 1. SETUP TOKEN ---
# Karena service ini diproteksi Auth Middleware, kita butuh JWT Token valid.
# Cara cepat untuk testing manual:
# 1. Buka https://jwt.io/
# 2. Di bagian VERIFY SIGNATURE, masukkan secret key dari .env: "rahasia-super-aman"
# 3. Copy string "Encoded" (token) dari box kiri dan paste di bawah ini:

$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzQzMzI0MDMsImlzcyI6ImF1ZGl0LXNlcnZpY2UtdGVzdCIsInJvbGUiOiJhZG1pbiIsInVzZXJuYW1lIjoiYWRtaW5fZ3VkYW5nIiwid2FyZWhvdXNlX2lkIjoiV0gtSktULTA5OSJ9.168sV5Lgvh3UOdUYhoeu6C_LEgNVgPR1-SNpcXiAcU8"

# --- 2. PREPARE PAYLOAD ---
$data = '{\"username\": \"admin_gudang\", \"warehouse_id\": \"WH-JKT-099\", \"role\": \"admin\", \"action\": \"STOCK_OUT\", \"sku\": \"MAC-001\", \"product_name\": \"Macbook Pro M2 14-inch\", \"quantity_changed\": 1, \"final_stock\": 6, \"timestamp\": \"2026-03-23T13:21:00Z\", \"metadata\": {\"request_id\": \"test-uuid-8888\", \"source\": \"powershell-script\"}}'

# --- 3a. DEBUG ---
Write-Host "Perintah yang akan dijalankan:"
Write-Host "grpcurl -plaintext -H 'Authorization: Bearer $token' -d '$data' localhost:50051 audit.AuditService/LogActivity"

# --- 3b. EXECUTE ---
Write-Host "Mengirim request ke localhost:50051..."
grpcurl -plaintext -H "Authorization: Bearer $token" -d $data localhost:50051 audit.AuditService/LogActivity
