# 🧪 WebSocket Real-time Chat Testing Guide

Panduan lengkap untuk testing WebSocket real-time chat di DevSync.

## 📋 Prerequisites

1. **Server DevSync harus berjalan**
   ```bash
   go run main.go
   # atau
   docker-compose up
   ```

2. **Pastikan port 8080 terbuka**

---

## 🎯 Method 1: Testing dengan HTML File (RECOMMENDED)

### Langkah-langkah:

1. **Buka file `test-websocket.html`** di browser (Chrome/Firefox/Safari)
   ```bash
   # Di macOS
   open test-websocket.html
   
   # Atau double-click file tersebut
   ```

2. **Konfigurasi koneksi:**
   - WebSocket URL: `ws://localhost:8080/ws` (sudah default)
   - User ID: `1` (atau user ID yang valid)
   - Project ID: `1` (atau project ID yang valid)

3. **Klik tombol "Connect"**
   - Status akan berubah menjadi "Connected" dengan dot hijau
   - Akan muncul pesan sistem: "✅ Connected to WebSocket server!"

4. **Test Real-time Chat:**

   **Buka 2 tab/window browser berbeda** dengan file yang sama:
   
   - **Tab 1 (User 1, Project 1)**
     - User ID: 1
     - Project ID: 1
     - Connect
   
   - **Tab 2 (User 2, Project 1)**  
     - User ID: 2
     - Project ID: 1
     - Connect

5. **Kirim pesan dari Tab 1:**
   - Ketik: "Hello from User 1!"
   - Tekan Enter atau klik "Send"
   - **Pesan harus langsung muncul di Tab 2!** ✨

6. **Balas dari Tab 2:**
   - Ketik: "Hi User 1, I received your message!"
   - **Pesan harus langsung muncul di Tab 1!** ✨

### ✅ Yang Harus Terjadi:

- ✅ Pesan yang dikirim dari satu tab **langsung muncul** di tab lain
- ✅ Tidak perlu refresh halaman
- ✅ Setiap pesan menampilkan User ID dan timestamp
- ✅ JSON raw data ditampilkan untuk debugging

### ❌ Testing Isolation per Project:

1. **Buka Tab 3 dengan Project ID berbeda:**
   - User ID: 3
   - **Project ID: 2** (berbeda!)
   - Connect

2. **Kirim pesan dari Tab 3:**
   - Ketik: "Message from Project 2"
   - Kirim

3. **Cek Tab 1 & 2:**
   - **Pesan TIDAK boleh muncul** di Tab 1 & 2
   - Karena mereka di Project 1, Tab 3 di Project 2
   - ✅ Ini membuktikan isolation per-project bekerja!

---

## 🎯 Method 2: Testing dengan Terminal (websocat)

### Install websocat (macOS):

```bash
brew install websocat
```

### Test Koneksi:

```bash
# Terminal 1 - Connect dan listen
websocat ws://localhost:8080/ws
```

### Kirim Pesan Manual:

```bash
# Terminal 2 - Send message
echo '{"type":"chat_message","project_id":1,"user_id":1,"data":{"message":"Hello from terminal!","user_id":1,"project_id":1}}' | websocat ws://localhost:8080/ws
```

### Atau gunakan script yang sudah dibuat:

```bash
# Beri permission dulu
chmod +x test_websocket.sh

# Jalankan
./test_websocket.sh
```

---

## 🎯 Method 3: Testing dengan Browser Console

### 1. Buka browser dan console (F12)

### 2. Paste dan jalankan code ini:

```javascript
// Connect ke WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
    console.log('✅ Connected!');
    
    // Kirim test message
    const message = {
        type: "chat_message",
        project_id: 1,
        user_id: 1,
        data: {
            message: "Hello from console!",
            user_id: 1,
            project_id: 1
        }
    };
    
    ws.send(JSON.stringify(message));
};

ws.onmessage = (event) => {
    console.log('📨 Received:', event.data);
    const data = JSON.parse(event.data);
    console.log('Parsed:', data);
};

ws.onerror = (error) => {
    console.error('❌ Error:', error);
};

ws.onclose = () => {
    console.log('🔌 Disconnected');
};
```

### 3. Buka console di tab lain dan lakukan hal yang sama

### 4. Kirim pesan dari salah satu tab:

```javascript
ws.send(JSON.stringify({
    type: "chat_message",
    project_id: 1,
    user_id: 2,
    data: {
        message: "Another message!",
        user_id: 2,
        project_id: 1
    }
}));
```

---

## 🎯 Method 4: Testing dengan Postman/Thunder Client

### Setup:

1. Buat New WebSocket Request
2. URL: `ws://localhost:8080/ws`
3. Connect

### Send Message:

```json
{
    "type": "chat_message",
    "project_id": 1,
    "user_id": 1,
    "data": {
        "message": "Hello from Postman!",
        "user_id": 1,
        "project_id": 1
    }
}
```

---

## 🔍 Testing Checklist

### ✅ Basic Functionality:

- [ ] Server berjalan tanpa error
- [ ] WebSocket connection berhasil (status 101 Switching Protocols)
- [ ] Client dapat connect ke `/ws` endpoint
- [ ] Client menerima confirmation setelah connect

### ✅ Real-time Messaging:

- [ ] Pesan yang dikirim dari Client A **langsung muncul** di Client B
- [ ] Tidak ada delay signifikan (< 100ms)
- [ ] Format JSON pesan sesuai dengan Message struct
- [ ] User ID dan Project ID ter-track dengan benar

### ✅ Multi-client Support:

- [ ] Multiple clients bisa connect bersamaan
- [ ] Broadcast ke semua clients di project yang sama
- [ ] Setiap client menerima copy pesan yang sama

### ✅ Project Isolation:

- [ ] Client di Project 1 **tidak menerima** pesan dari Project 2
- [ ] Client di Project 2 **tidak menerima** pesan dari Project 1
- [ ] Filtering berdasarkan `project_id` bekerja

### ✅ Connection Management:

- [ ] Client bisa disconnect dengan graceful
- [ ] Server log menampilkan "Client connected"
- [ ] Server log menampilkan "Client disconnected"
- [ ] Tidak ada memory leak saat banyak connect/disconnect

### ✅ Error Handling:

- [ ] Server tidak crash saat client disconnect tiba-tiba
- [ ] Invalid JSON tidak crash server
- [ ] Missing fields di-handle dengan baik

---

## 📊 Expected Server Logs

Saat testing berhasil, Anda harus melihat log seperti ini:

```
[GIN] 2025/10/15 - 14:30:15 | 101 | WebSocket Upgrade
2025/10/15 14:30:15 Client connected: User 1, Project 1
2025/10/15 14:30:22 Client connected: User 2, Project 1
2025/10/15 14:30:30 Client disconnected: User 1, Project 1
```

---

## 🐛 Troubleshooting

### Problem: "WebSocket connection failed"

**Solution:**
```bash
# Check if server running
curl http://localhost:8080/

# Check if port is open
lsof -i :8080
```

### Problem: "CORS error"

**Solution:**
- CORS sudah di-configure untuk allow all origins di `routes.go`
- Pastikan server di-restart setelah perubahan

### Problem: "Pesan tidak muncul di client lain"

**Solution:**
```bash
# Check logs untuk melihat project_id
# Pastikan kedua client punya project_id yang SAMA
# Lihat hub.go line 65-76 untuk logic filtering
```

### Problem: "Connection closed immediately"

**Solution:**
```go
// Check hub.go - HandleWebSocket function
// Pastikan TODO diimplementasi untuk extract user & project dari JWT/query
```

---

## 🎓 Understanding the Flow

```
┌─────────────┐                    ┌─────────────┐
│  Browser 1  │                    │  Browser 2  │
│  (User 1)   │                    │  (User 2)   │
└──────┬──────┘                    └──────┬──────┘
       │                                   │
       │ ws://localhost:8080/ws           │
       │                                   │
       ▼                                   ▼
┌────────────────────────────────────────────────┐
│              WebSocket Hub                     │
│  ┌──────────────────────────────────────────┐ │
│  │  Clients Map (Project 1)                 │ │
│  │  - Client 1 (User 1, send channel)       │ │
│  │  - Client 2 (User 2, send channel)       │ │
│  └──────────────────────────────────────────┘ │
└────────────────────────────────────────────────┘
                       │
                       ▼
        User 1 sends: "Hello!"
                       │
                       ▼
          ┌───────────────────────┐
          │   Broadcast Channel   │
          └───────────────────────┘
                       │
          ┌────────────┴────────────┐
          ▼                         ▼
    Client 1.send            Client 2.send
          │                         │
          ▼                         ▼
    ✅ Received              ✅ Received
    (own echo)              (real-time!)
```

---

## 🎉 Success Criteria

Test dianggap **BERHASIL** jika:

1. ✅ Multiple clients bisa connect bersamaan
2. ✅ Pesan dari satu client **langsung muncul** di client lain (<100ms)
3. ✅ Project isolation bekerja (Project 1 tidak terima pesan Project 2)
4. ✅ Connect/disconnect tidak crash server
5. ✅ Server logs menunjukkan aktivitas yang benar

---

## 📝 Next Steps

Setelah testing berhasil, implementasi berikut perlu dilakukan:

1. **Extract User ID dari JWT Token** (saat ini hardcoded ke 1)
2. **Extract Project ID dari query parameter** (misal: `/ws?project_id=5`)
3. **Implement proper authentication** untuk WebSocket
4. **Add message persistence** - simpan ke database via ChatHandler
5. **Add typing indicators** - "User X is typing..."
6. **Add read receipts** - "Message read by 3 users"
7. **Add file/task-specific chat** - filter by file_id/task_id

---

## 🔗 Related Files

- `internal/websocket/hub.go` - WebSocket server logic
- `internal/api/handlers/chat.go` - Chat message handlers
- `internal/api/routes.go` - Route configuration
- `test-websocket.html` - HTML testing tool
- `test_websocket.sh` - Terminal testing script

---

**Happy Testing! 🚀**

Jika ada masalah, check server logs dan pastikan semua prerequisites terpenuhi.
