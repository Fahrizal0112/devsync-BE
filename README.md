# DevSync Backend

DevSync adalah aplikasi backend untuk platform manajemen proyek kolaboratif yang memungkinkan tim developer untuk bekerja sama secara efektif dalam pengembangan software. Aplikasi ini menyediakan API RESTful dan komunikasi real-time melalui WebSocket.

## 🚀 Fitur Utama

### 🔐 Autentikasi & Otorisasi
- **GitHub OAuth Integration** - Login menggunakan akun GitHub
- **JWT Token Authentication** - Sistem autentikasi berbasis token
- **Token Refresh** - Pembaruan token otomatis

### 📁 Manajemen Proyek
- **CRUD Projects** - Buat, baca, update, dan hapus proyek
- **GitHub Repository Integration** - Integrasi dengan repository GitHub
- **Public/Private Projects** - Kontrol visibilitas proyek
- **Multi-user Collaboration** - Kolaborasi tim dalam satu proyek

### 📄 Manajemen File
- **File Upload ke Google Cloud Storage** - Penyimpanan file di cloud
- **Text File Management** - Manajemen file teks dengan konten di database
- **File Metadata** - Informasi lengkap file (ukuran, tipe, MIME type)
- **File Versioning** - Pelacakan perubahan file

### ✅ Manajemen Tugas
- **Task Management** - Buat dan kelola tugas proyek
- **Task Status Tracking** - Status: Todo, In Progress, Done
- **Task Assignment** - Assign tugas ke anggota tim
- **Priority System** - Sistem prioritas tugas
- **GitHub Issues Integration** - Integrasi dengan GitHub Issues

### 🏃‍♂️ Sprint Management
- **Sprint Planning** - Perencanaan sprint dengan tanggal mulai/selesai
- **Sprint Status** - Status: Active, Completed, Cancelled
- **Task Assignment to Sprints** - Assign tugas ke sprint tertentu

### 💬 Sistem Chat Real-time
- **Project Chat** - Chat per proyek
- **File-specific Chat** - Chat terkait file tertentu
- **Task-specific Chat** - Chat terkait tugas tertentu
- **WebSocket Integration** - Komunikasi real-time

### 📚 Dokumentasi
- **Project Documentation** - Dokumentasi proyek dengan Markdown
- **File-based Documentation** - Dokumentasi berbasis file

### 🚀 Deployment Management
- **Multi-environment Deployment** - Development, Staging, Production
- **Deployment Status Tracking** - Pelacakan status deployment
- **Version Management** - Manajemen versi deployment

### 🔄 Real-time Features
- **WebSocket Communication** - Komunikasi real-time
- **Live Updates** - Update langsung untuk perubahan file, tugas, dan chat
- **Multi-client Synchronization** - Sinkronisasi antar klien

## 🛠 Teknologi yang Digunakan

### Backend Framework
- **Go 1.24.1** - Bahasa pemrograman utama
- **Gin Framework** - Web framework untuk API REST
- **GORM** - ORM untuk database operations

### Database & Storage
- **PostgreSQL** - Database utama
- **Google Cloud Storage** - Penyimpanan file
- **Redis** - Caching dan session management

### Authentication & Security
- **JWT (JSON Web Tokens)** - Sistem autentikasi
- **GitHub OAuth2** - Integrasi login GitHub
- **CORS Middleware** - Cross-origin resource sharing

### Real-time Communication
- **Gorilla WebSocket** - WebSocket implementation
- **Custom WebSocket Hub** - Manajemen koneksi WebSocket

### Development & Deployment
- **Docker & Docker Compose** - Containerization
- **Environment Variables** - Konfigurasi aplikasi
- **Swagger Documentation** - API documentation

## 📋 Prasyarat

- Go 1.21 atau lebih baru
- PostgreSQL 15+
- Redis 7+
- Google Cloud Platform account (untuk GCS)
- GitHub OAuth App (untuk autentikasi)

## ⚙️ Instalasi & Setup

### 1. Clone Repository
```bash
git clone <repository-url>
cd DevSync-BE
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment Variables
Buat file `.env` di root directory:
```env
# Database
DATABASE_URL=postgres://user:password@localhost:5432/devsync?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key

# GitHub OAuth
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
REDIRECT_URL=http://localhost:3000/auth/callback

# Google Cloud Storage
GCP_PROJECT_ID=your-gcp-project-id
GCP_BUCKET_NAME=your-gcs-bucket-name
GCP_CREDENTIALS_PATH=path/to/service-account.json

# Server
PORT=8080
```

### 4. Setup Database
```bash
# Menggunakan Docker Compose
docker-compose up postgres redis -d

# Atau setup manual PostgreSQL dan Redis
```

### 5. Setup Google Cloud Storage
1. Buat project di Google Cloud Console
2. Enable Cloud Storage API
3. Buat service account dan download JSON credentials
4. Buat storage bucket
5. Set path ke credentials file di environment variable

### 6. Setup GitHub OAuth
1. Buka GitHub Settings > Developer settings > OAuth Apps
2. Buat New OAuth App
3. Set Authorization callback URL: `http://localhost:8080/auth/github/callback`
4. Copy Client ID dan Client Secret ke `.env`

### 7. Jalankan Aplikasi
```bash
# Development
go run main.go

# Atau menggunakan Docker
docker-compose up
```

## 🌐 API Endpoints

### Authentication
- `GET /auth/github` - Redirect ke GitHub OAuth
- `GET /auth/github/callback` - GitHub OAuth callback
- `GET /auth/me` - Get current user info
- `POST /auth/refresh` - Refresh JWT token
- `POST /auth/dev-login` - Development login (tanpa GitHub)

### Users
- `GET /api/v1/users/search` - Search users by username or email
- `GET /api/v1/users` - Get all users (with pagination)

### Projects
- `GET /api/v1/projects` - Get all projects
- `POST /api/v1/projects` - Create new project
- `GET /api/v1/projects/:id` - Get project by ID
- `PUT /api/v1/projects/:id` - Update project
- `DELETE /api/v1/projects/:id` - Delete project

### Project Members
- `GET /api/v1/projects/:id/members` - Get project members
- `POST /api/v1/projects/:id/members` - Add member to project
- `DELETE /api/v1/projects/:id/members/:userId` - Remove member from project

### Files
- `GET /api/v1/projects/:id/files` - Get project files
- `POST /api/v1/projects/:id/files` - Create text file
- `POST /api/v1/projects/:id/upload` - Upload file to GCS
- `GET /api/v1/files/:id` - Get file by ID
- `PUT /api/v1/files/:id` - Update file
- `DELETE /api/v1/files/:id` - Delete file

### Tasks
- `GET /api/v1/projects/:id/tasks` - Get project tasks
- `POST /api/v1/projects/:id/tasks` - Create new task
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task

### Sprints
- `GET /api/v1/projects/:id/sprints` - Get project sprints
- `POST /api/v1/projects/:id/sprints` - Create new sprint

### Chat
- `GET /api/v1/projects/:id/messages` - Get project messages
- `POST /api/v1/projects/:id/messages` - Send message

### WebSocket
- `GET /ws` - WebSocket connection endpoint

## 🧪 Testing

### Postman Collection
Import koleksi Postman yang tersedia di `postman/DevSync-API.postman_collection.json` untuk testing API endpoints.

### Environment Variables untuk Testing
```json
{
  "base_url": "http://localhost:8080",
  "jwt_token": "",
  "user_id": "",
  "project_id": "",
  "file_id": "",
  "task_id": "",
  "sprint_id": "",
  "github_code": ""
}
```

### Testing Guide
Lihat dokumentasi lengkap di `docs/API_TESTING_GUIDE.md` untuk panduan testing yang detail.

## 🏗 Struktur Proyek

```
DevSync-BE/
├── cmd/                    # Application entrypoints
├── internal/
│   ├── api/
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Middleware functions
│   │   └── routes.go       # Route definitions
│   ├── auth/              # Authentication logic
│   ├── config/            # Configuration management
│   ├── database/          # Database connection & migrations
│   ├── models/            # Database models
│   ├── storage/           # File storage (GCS)
│   └── websocket/         # WebSocket hub & client management
├── docs/                  # Documentation
├── postman/              # Postman collections
├── docker-compose.yml    # Docker services
├── Dockerfile           # Application container
├── go.mod              # Go modules
└── main.go            # Application entry point
```

## 🔧 Konfigurasi

### Database Models
Aplikasi menggunakan auto-migration GORM untuk model berikut:
- `User` - Data pengguna
- `Project` - Data proyek
- `File` - Metadata file
- `Task` - Data tugas
- `Sprint` - Data sprint
- `Comment` - Komentar tugas
- `Documentation` - Dokumentasi proyek
- `ChatMessage` - Pesan chat
- `Deployment` - Data deployment

### WebSocket Events
- `file_updated` - File telah diperbarui
- `task_updated` - Tugas telah diperbarui
- `new_message` - Pesan chat baru
- `user_joined` - User bergabung ke proyek
- `user_left` - User meninggalkan proyek

## 🚀 Deployment

### Docker Deployment
```bash
# Build dan jalankan dengan Docker Compose
docker-compose up --build

# Atau build manual
docker build -t devsync-be .
docker run -p 8080:8080 devsync-be
```

### Production Considerations
- Set `GIN_MODE=release` untuk production
- Gunakan database PostgreSQL yang terpisah
- Setup load balancer untuk multiple instances
- Implement proper logging dan monitoring
- Setup backup untuk database dan file storage

## 🤝 Kontribusi

1. Fork repository
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📞 Support

Untuk pertanyaan atau dukungan, silakan buat issue di repository ini atau hubungi tim development.

---

**DevSync Backend** - Empowering collaborative software development 🚀