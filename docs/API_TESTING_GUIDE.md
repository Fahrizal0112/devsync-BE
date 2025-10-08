# DevSync API Testing Guide

## Persiapan Testing

### 1. Setup Environment
Pastikan aplikasi DevSync backend sudah berjalan:
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

### 2. Import Postman Collection
1. Buka Postman
2. Click **Import**
3. Pilih file `postman/DevSync-API.postman_collection.json`
4. Collection akan muncul di sidebar Postman

### 3. Setup Environment Variables
Collection sudah include environment variables:
- `base_url`: http://localhost:8080
- `jwt_token`: (akan diisi otomatis setelah login)
- `project_id`: (akan diisi otomatis setelah create project)
- `file_id`: (akan diisi otomatis setelah create file)
- `task_id`: (akan diisi otomatis setelah create task)

## Urutan Testing

### Step 1: Authentication
Karena menggunakan GitHub OAuth, testing authentication agak tricky di Postman. Ada beberapa cara:

#### Cara 1: Manual OAuth Flow
1. Jalankan request **"GitHub OAuth Login"**
2. Copy URL response dan buka di browser
3. Login dengan GitHub
4. Setelah redirect, copy `code` parameter dari URL
5. Paste code tersebut ke request **"GitHub OAuth Callback"**
6. JWT token akan tersimpan otomatis di variable `jwt_token`

#### Cara 2: Generate JWT Manual (untuk testing)
Buat endpoint temporary untuk generate JWT:

```go
// Tambahkan di routes.go untuk testing
auth.POST("/test-login", func(c *gin.Context) {
    // Hardcode user untuk testing
    token, _ := auth.GenerateToken(1, "testuser", cfg.JWTSecret)
    c.JSON(200, gin.H{"token": token})
})
```

### Step 2: Test User Info
1. Jalankan **"Get Current User"**
2. Pastikan response menampilkan data user yang login

### Step 3: Project Management
1. **Create Project** - Buat project baru
2. **Get All Projects** - Lihat semua project
3. **Get Project by ID** - Detail project
4. **Update Project** - Update informasi project
5. **Delete Project** - Hapus project (opsional)

### Step 4: File Management
1. **Create File** - Buat file dalam project
2. **Get Project Files** - Lihat semua file
3. **Get File by ID** - Detail file
4. **Update File** - Edit konten file
5. **Delete File** - Hapus file (opsional)

### Step 5: Task Management
1. **Create Task** - Buat task baru
2. **Get Project Tasks** - Lihat semua task
3. **Update Task Status** - Ubah status task (todo → in_progress → done)
4. **Delete Task** - Hapus task (opsional)

### Step 6: Sprint Management
1. **Create Sprint** - Buat sprint baru
2. **Get Project Sprints** - Lihat semua sprint

### Step 7: Chat System
1. **Send Message** - Kirim pesan chat
2. **Get Project Messages** - Lihat semua pesan
3. **Get Messages by File** - Filter pesan berdasarkan file

## Testing WebSocket

Untuk testing WebSocket, gunakan tools seperti:

### 1. WebSocket King (Chrome Extension)
- URL: `ws://localhost:8080/ws`
- Send message:
```json
{
    "type": "file_updated",
    "project_id": 1,
    "user_id": 1,
    "data": {
        "id": 1,
        "content": "updated content"
    }
}
```

### 2. wscat (Command Line)
```bash
npm install -g wscat
wscat -c ws://localhost:8080/ws
```

## Expected Responses

### Success Response Format
```json
{
    "id": 1,
    "name": "Project Name",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

### Error Response Format
```json
{
    "error": "Error message description"
}
```

## Status Codes
- `200` - OK
- `201` - Created
- `204` - No Content (untuk DELETE)
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Testing Checklist

### Authentication ✅
- [ ] GitHub OAuth login redirect
- [ ] OAuth callback dengan valid code
- [ ] JWT token generation
- [ ] Get current user info
- [ ] Token refresh

### Projects ✅
- [ ] Create project
- [ ] Get all projects
- [ ] Get project by ID
- [ ] Update project
- [ ] Delete project

### Files ✅
- [ ] Create file in project
- [ ] Get all files in project
- [ ] Get file by ID
- [ ] Update file content
- [ ] Delete file

### Tasks ✅
- [ ] Create task
- [ ] Get all tasks in project
- [ ] Update task status
- [ ] Delete task

### Sprints ✅
- [ ] Create sprint
- [ ] Get all sprints in project

### Chat ✅
- [ ] Send message
- [ ] Get all messages
- [ ] Filter messages by file
- [ ] Filter messages by task

### WebSocket ✅
- [ ] Connect to WebSocket
- [ ] Receive real-time updates
- [ ] Send messages via WebSocket

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Pastikan JWT token valid
   - Check Authorization header format: `Bearer <token>`

2. **404 Not Found**
   - Pastikan endpoint URL benar
   - Check parameter ID valid

3. **500 Internal Server Error**
   - Check database connection
   - Check server logs untuk detail error

4. **CORS Issues**
   - Pastikan CORS middleware aktif
   - Check origin header

### Debug Tips
1. Enable Postman Console untuk melihat request/response detail
2. Check server logs di terminal
3. Gunakan Postman Tests untuk auto-validation
4. Set breakpoints di code untuk debugging

## Advanced Testing

### Load Testing
Gunakan Postman Runner untuk test multiple requests:
1. Select collection
2. Click "Run"
3. Set iterations dan delay
4. Monitor response times

### Automated Testing
Buat test scripts di Postman:
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has token", function () {
    const response = pm.response.json();
    pm.expect(response).to.have.property('token');
});
```