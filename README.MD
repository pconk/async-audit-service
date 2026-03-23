[![Go CI](https://github.com/pconk/async-audit-service/actions/workflows/ci.yml/badge.svg)](https://github.com/pconk/async-audit-service/actions/workflows/ci.yml)

# Async Audit Service (gRPC + Go)

**Async Audit Service** adalah microservice yang dirancang untuk mencatat jejak audit (log aktivitas) secara asinkronus dari layanan utama. Dibangun menggunakan **Go**, **gRPC**, dan **MongoDB** dengan menerapkan prinsip **Clean Architecture**.

## 🚀 Fitur Utama
* **gRPC Framework:** Komunikasi antar-service yang cepat dan efisien menggunakan Protocol Buffers.
* **Asynchronous Logging:** Didesain untuk dipanggil via goroutine agar tidak menghambat performa aplikasi utama.
* **Clean Architecture:** Pemisahan tanggung jawab yang jelas antara Handler, Service, Repository, dan Entity.
* **Middleware Chain:** 
    * **Logger Interceptor:** Mencatat latensi dan detail setiap request.
    * **Auth Interceptor:** Validasi keamanan menggunakan JWT.
* **MongoDB Storage:** Penyimpanan log yang fleksibel dan terukur.
* **Graceful Shutdown:** Memastikan semua proses selesai dengan aman sebelum aplikasi mati.

---

## 📂 Struktur Proyek
```text
.
├── audit/              # Definisi file .proto
├── internal/           # Logika privat aplikasi
│   ├── cmd/api/        # Entry point (main.go)
│   ├── config/         # Manajemen konfigurasi & env
│   ├── entity/         # Struktur data / model database
│   ├── handler/grpc/   # Implementasi gRPC server
│   ├── middleware/     # Interceptors (Auth, Logger)
│   ├── repository/     # Layer akses ke MongoDB
│   └── service/        # Business logic
├── pb/                 # Generated code dari protoc
├── go.mod              # Go dependencies
└── README.md
```

---

## 🛠️ Persiapan & Instalasi

### 1. Prasyarat
* Go 1.25+
* MongoDB (Running via Docker atau Lokal)
* Protocol Buffer Compiler (`protoc`)

### 2. Instalasi Dependensi
```bash
go mod tidy
```

### 3. Generate Protobuf (Jika ada perubahan)
```bash
protoc --go_out=. --go-grpc_out=. audit/audit.proto
```

---

## ⚙️ Konfigurasi
Aplikasi menggunakan variabel lingkungan (environment variables) yang dikelola di `internal/config`. Pastikan variabel berikut tersedia:

| Variabel | Deskripsi | Default |
| :--- | :--- | :--- |
| `APP_PORT` | Port untuk gRPC Server | `50051` |
| `MONGO_URI` | URI koneksi MongoDB | `mongodb://localhost:27017` |
| `DB_NAME` | Nama database di Mongo | `audit_db` |
| `JWT_SECRET` | Secret key untuk validasi Auth | `-` |

---

## 🚦 Cara Menjalankan

### Menjalankan Server
```bash
go run internal/cmd/api/main.go
```

## 🧪 Cara Testing (Manual)

Kamu bisa menggunakan tool seperti **grpcurl** atau **Postman** (fitur gRPC) untuk ngetes endpoint.

### Menggunakan Script PowerShell (Windows)
Telah disediakan script untuk memudahkan testing di Windows.
1. Buka `test-grpcurl.ps1`.
2. Masukkan token JWT yang valid di variabel `$token`.
3. Jalankan:
   ```powershell
   .\test-grpcurl.ps1
   ```

---

## 📝 Catatan Arsitektur
* **Middleware Order:** Logger dieksekusi terlebih dahulu untuk mencatat semua trafik, diikuti oleh Auth untuk keamanan.
* **Async Caller:** Disarankan memanggil service ini menggunakan `go func()` dari sisi client (Warehouse API) agar tidak menambah latensi pada transaksi utama.

---
