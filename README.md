# MRT Schedules API

Backend API service untuk informasi jadwal, tarif, dan fasilitas MRT Jakarta. API ini berfungsi sebagai proxy yang mengonsumsi data dari API resmi MRT Jakarta dan menyajikannya dalam format yang terstruktur dan mudah dikonsumsi.

## 📋 Overview

- **Purpose**: Menyediakan data terstruktur MRT Jakarta untuk integrasi dengan aplikasi pihak ketiga
- **Source Data**: [API MRT Jakarta](https://jakartamrt.co.id/id/val/stasiuns)
- **Architecture**: Clean Architecture (Handler → Usecase → Service)
- **Language**: Go 1.24.3 dengan Gin Framework

## 🚀 Fitur

### 📡 Available Endpoints

#### Stasiun
- `GET /v1/api/stations` - Daftar semua stasiun (dengan filter nama)
- `GET /v1/api/stations/{id}` - Jadwal keberangkatan stasiun
- `GET /v1/api/stations/{id}/details` - Detail lengkap stasiun (fasilitas, retail, transportasi)

#### Jadwal & Tarif
- `GET /v1/api/stations/{id}/next-train?destination=<LB|HI>` - 3 kereta berikutnya
- `GET /v1/api/stations/fare?from=<id>&to=<id>` - Tarif dan durasi perjalanan

## 🏗️ Arsitektur

### Struktur Project
```
mrt-schedules/
├── cmd/server/main.go           # Entry point aplikasi
├── internal/                    # Private application code
│   ├── config/config.go         # Konfigurasi aplikasi
│   └── api/
│       ├── handler/station.go   # HTTP handlers & routing
│       ├── service/station/     # Data fetching layer
│       └── usecase/station/     # Business logic layer
└── pkg/                        # Public/shared code
    ├── client/client.go        # HTTP client utility
    └── response/               # Standard API responses
```

### Alur Data
```
HTTP Request → Handler → Usecase → Service → External API
                    ↓           ↓         ↓
                Response ← Usecase ← Service ← JSON Response
```

## 🛠️ Teknologi

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- **Configuration**: [godotenv](https://github.com/joho/godotenv) - Environment variables
- **HTTP Client**: Standard library `net/http`
- **Time Parsing**: `time.ParseInLocation` untuk timezone-aware scheduling

## 📦 Installation & Setup

### Prerequisites
- Go 1.24.3 atau lebih tinggi
- Akses internet ke API MRT Jakarta

### Quick Start
```bash
# Clone repository
git clone https://github.com/IkrmMrbsy/mrt-schedules.git
cd mrt-schedules

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Run server
go run cmd/server/main.go
```

### Environment Variables
```env
SERVER_PORT=8080                     # Port server
HTTP_TIMEOUT=10                      # HTTP timeout (detik)
MRT_API_URL=https://jakartamrt.co.id/id/val/stasiuns  # Source API
```

## 📖 API Documentation

### Response Format
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### Examples

#### 1. Daftar Stasiun
```bash
curl "http://localhost:8080/v1/api/stations"
```
Response:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {"id": "21", "nama": "Stasiun Bundaran HI"},
    {"id": "22", "nama": "Stasiun Dukuh Atas"}
  ]
}
```

#### 2. Filter Stasiun
```bash
curl "http://localhost:8080/v1/api/stations?name=bundaran"
```

#### 3. Jadwal Stasiun
```bash
curl "http://localhost:8080/v1/api/stations/21"
```

#### 4. Kereta Berikutnya
```bash
curl "http://localhost:8080/v1/api/stations/21/next-train?destination=LB"
```

#### 5. Tarif Perjalanan
```bash
curl "http://localhost:8080/v1/api/stations/fare?from=21&to=1"
```

#### 6. Detail Stasiun
```bash
curl "http://localhost:8080/v1/api/stations/21/details"
```

## 🔄 Data Flow

### 1. Station Data
- **Source**: MRT API endpoint `/val/stasiuns`
- **Format**: JSON dengan nested objects (retails, fasilitas, antarmoda)
- **Process**: Parse → Filter → Group by category
- **Output**: Simplified station list dengan informasi dasar

### 2. Schedule Data
- **Source**: String waktu format `"HH:MM:SS,HH:MM:SS,..."`
- **Process**: Parse → Convert to `time.Time` → Filter upcoming trains
- **Output**: Array of upcoming departure times

### 3. Fare & Duration
- **Source**: Nested estimation objects
- **Process**: Matrix lookup → Find matching pairs
- **Output**: Fare amount + travel duration

## 🎯 Use Cases

### Mobile Apps
- Transit tracking apps
- Jakarta tourism guides
- Transportation aggregators

### Web Applications
- Dashboard monitoring
- Corporate travel management
- Public information systems

### Integration
- Third-party transportation platforms
- Smart city applications
- Logistics and delivery services

## 🛠️ Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o mrt-api cmd/server/main.go
```

### Docker Support
```dockerfile
FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o mrt-api cmd/server/main.go
EXPOSE 8080
CMD ["./mrt-api"]
```

## 📊 Performance

- **Response Time**: ~200-500ms (tergantung jaringan ke API MRT)
- **Rate Limiting**: Mengikuti policy API MRT Jakarta
- **Caching**: Tidak tersedia (direct proxy to real-time data)
- **Concurrent Support**: Standard Go HTTP server capabilities

## 🔧 Troubleshooting

### Common Issues
1. **Invalid Time Format**: API response mengandung format waktu yang tidak terduga
2. **Network Timeout**: Koneksi ke API MRT lambat atau down
3. **Station Not Found**: ID stasiun tidak valid dalam sistem MRT

### Error Responses
```json
{
  "code": 404,
  "message": "station not found",
  "data": null
}
```

## 📄 License

This project is provided as-is for educational and integration purposes. Data ownership remains with MRT Jakarta authorities.