# Konfigurasi
$GrpcAddress = "localhost:50052"
$ProtoFile = "proto/audit.proto"
$requestID = "reqid-001"
$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzQ4NzIxNTYsImlzcyI6ImF1ZGl0LXNlcnZpY2UtdGVzdCIsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOiIxIiwidXNlcm5hbWUiOiJhZG1pbl9ndWRhbmciLCJ3YXJlaG91c2VfaWQiOiJXSC1KS1QtMDk5In0.J2M_tRQPCUi_AGJrWZJwTxNUxNh_HJD-pKLRONOtp"

# Cek apakah grpcurl terinstall
if (-not (Get-Command grpcurl -ErrorAction SilentlyContinue)) {
    Write-Error "Tool 'grpcurl' tidak ditemukan. Silakan install: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit
}

Write-Host "--- Audit Service gRPC Tester ---" -ForegroundColor Cyan
Write-Host "1. LogActivity (Simulasi Trigger dari Warehouse API)"
Write-Host "2. GetRecentLogs (Ambil Data untuk Gateway)"
Write-Host "------------------------------------------------"
$choice = Read-Host "Pilih menu (1/2)"

switch ($choice) {
    "1" {
        Write-Host "`n--- Menjalankan LogActivity ---" -ForegroundColor Yellow
        
        # Menyiapkan timestamp dalam format RFC3339 yang diterima google.protobuf.Timestamp
        $timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

        $payload = @{
            user_id          = 1
            username         = "admin_gudang"
            warehouse_id     = "WH-JKT-01"
            role             = "admin"
            action           = "STOCK_IN"
            sku              = "ELC-LAP-001"
            product_name     = "MacBook Pro M3"
            quantity_changed = 10
            final_stock      = 55
            timestamp        = $timestamp
            metadata         = @{
                source     = "manual-test-ps"
                request_id = [guid]::NewGuid().ToString()
            }
        } | ConvertTo-Json -Compress
        
        Write-Output $payload | grpcurl -plaintext -H "Authorization: Bearer $token" -H "x-request-id: $requestID" -proto $ProtoFile -d "@" $GrpcAddress audit.AuditService/LogActivity
    }

    "2" {
        Write-Host "`n--- Menjalankan GetRecentLogs ---" -ForegroundColor Yellow
        $limit = Read-Host "Masukkan jumlah limit (default 5)"
        if (-not $limit) { $limit = 5 }

        $payload = @{ limit = [int]$limit } | ConvertTo-Json -Compress

        $payload | grpcurl -plaintext -H "Authorization: Bearer $token" -H "x-request-id: $requestID" -proto $ProtoFile -d "@" $GrpcAddress audit.AuditService/GetRecentLogs
    }

    Default {
        Write-Host "Pilihan tidak valid." -ForegroundColor Red
    }
}

Write-Host "`nSelesai." -ForegroundColor Green