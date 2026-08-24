# 🎵 TuneTrend Backend API

ระบบ Backend สำหรับแอปพลิเคชัน TuneTrend ทำหน้าที่ดึงข้อมูลเพลงที่กำลังมาแรง (Trending Music) ในแต่ละประเทศจาก YouTube นำมาจัดเก็บและให้บริการผ่าน API ด้วยความเร็วสูง

## 🛠 Tech Stack (เทคโนโลยีที่ใช้)

- **ภาษา:** Go (Golang 1.25+)
- **Web Framework:** Fiber v2 (เน้นประสิทธิภาพและความเร็ว)
- **Database:** PostgreSQL (จัดการผ่าน GORM)
- **External API:** YouTube Data API v3
- **Architecture:** Clean Architecture (โครงสร้างที่รองรับการเติบโตของโปรเจกต์)

---

## ✨ ฟีเจอร์เด่นระดับ Production

1. **🤖 Smart Video Categorization (AI-Lite Tokenizer):** ระบบวิเคราะห์ชื่อคลิปแบบ Custom Tokenization เพื่อแยกประเภท MV, Lyric, Live, Cover ได้อย่างแม่นยำ ป้องกันการจับคำผิดพลาด (เช่น `discover` หรือ `delivery`)
2. **⚡ Automated Sync Engine & Quota Optimized:** ระบบ Background Worker ดึงข้อมูล 50 อันดับแรกทุกๆ 3 ชั่วโมง ใช้โควตา YouTube API อย่างคุ้มค่าที่สุด (1 Unit / ประเทศ / รอบ)
3. **🗄️ Robust Database Layer:** ใช้ GORM `ON CONFLICT` เพื่อทำ Upsert ข้อมูลยอดวิวให้สดใหม่เสมอ และมีการ `CAST(view_count AS BIGINT)` เพื่อให้จัดอันดับยอดวิวได้อย่างถูกต้องแม่นยำ 100%
4. **🛡️ Security:** ส่ง API Key ผ่าน `X-goog-api-key` Header แทนการใส่ใน URL ป้องกันคีย์หลุดรั่วไหลลงไปใน Server Access Logs

---

## 📂 โครงสร้างโปรเจกต์ (Clean Architecture)

ระบบถูกออกแบบโดยยึดหลัก Clean Architecture เพื่อแยก Business Logic ออกจาก Database และ Framework ทำให้เขียน Unit Test ได้ง่ายและปรับสเกลในอนาคตได้สะดวก

```text
tunetrend-backend/
├── cmd/api/                  # จุดเริ่มต้นโปรเจกต์ (รวบรวม Dependency Injection)
├── docs/                     # ไฟล์ Swagger Documentation (สร้างอัตโนมัติ)
├── internal/                 # โค้ดหลักของระบบ (Private)
│   ├── core/database/        # การเชื่อมต่อ PostgreSQL
│   ├── domain/               # หัวใจของระบบ: Entity และ Interfaces
│   ├── repository/           # เลเยอร์จัดการข้อมูล (Postgres, YouTube API)
│   ├── usecase/              # เลเยอร์กฎทางธุรกิจ (Business Logic)
│   ├── delivery/http/        # เลเยอร์จัดการ HTTP Request/Response (Fiber)
│   └── worker/               # ระบบ Background Job (ซิงก์ข้อมูล YouTube)
├── .air.toml                 # ตั้งค่า Live Reload
├── .env.example              # ตัวอย่างไฟล์ตั้งค่า Environment
├── setup.sh                  # สคริปต์สำหรับติดตั้งระบบครั้งแรก
└── docker-compose.yml        # ไฟล์สำหรับรัน Database ใน Local

```

---

## 🚀 วิธีการรันโปรเจกต์แบบด่วน (Quick Start)

ทีมงานสามารถรันสคริปต์เดียวเพื่อตั้งค่า Environment, เปิด Database และ Generate Swagger ได้ทันที:

```bash
./setup.sh

```

---

## 🔄 ระบบ Live Reload & Swagger (สำหรับนักพัฒนา)

เพื่อให้การเขียนโค้ดสะดวกที่สุด ระบบนี้รองรับ **Air** (Live Reload) เมื่อมีการแก้ไขไฟล์ เซิร์ฟเวอร์จะทำการอัปเดต Swagger และ Restart ตัวเองอัตโนมัติ

### 1. การติดตั้งเครื่องมือ (ทำแค่ครั้งแรก)

ติดตั้งเครื่องมือสำหรับ Live Reload และ Swagger Documentation:

```bash
go install [github.com/air-verse/air@latest](https://github.com/air-verse/air@latest)
go install [github.com/swaggo/swag/cmd/swag@latest](https://github.com/swaggo/swag/cmd/swag@latest)

```

_(หากเรียกคำสั่งไม่ได้ ให้ตั้งค่า Path: `export PATH=$(go env GOPATH)/bin:$PATH`)_

### 2. การรันเซิร์ฟเวอร์ระหว่างพัฒนา

ใช้คำสั่งนี้แทน `go run` เพื่อเปิดเซิร์ฟเวอร์แบบ Live Reload:

```bash
air

```

_หลังจากเซิร์ฟเวอร์ทำงาน สามารถดูเอกสาร API ได้ที่: `http://localhost:8080/swagger/_`

---

## 🧪 การทดสอบระบบ (Testing)

ระบบมีการเขียน Unit Test แบบ Mocking (ไม่ต้องต่อ Database จริง) เพื่อทดสอบ Business Logic สามารถรันเทสต์ทั้งหมดได้ด้วยคำสั่ง:

```bash
go test -v ./...

```

---

## 🔌 API Specification (สำหรับ Frontend / Mobile)

ทีมพัฒนาสามารถใช้หน้าเว็บ **Swagger UI** (`http://localhost:8080/swagger/`) หรือใช้ [Bruno](https://www.usebruno.com/) ในการยิงทดสอบ API ได้ โดยกำหนด Environment `base_url` เป็น `http://localhost:8080`

### 1. Health Check

ตรวจสอบสถานะการทำงานของเซิร์ฟเวอร์

- **Method:** `GET`
- **URL:** `/health`
- **Response (200 OK):**

```json
"OK - API & Worker are running!"
```

### 2. Manual Sync (สั่งดึงข้อมูลทันที)

สั่งให้เซิร์ฟเวอร์ดึงข้อมูลจาก YouTube API และอัปเดตลง Database ทันที (ไม่ต้องรอรอบ Worker)

- **Method:** `POST`
- **URL:** `/sync`
- **Query Parameter:** `country` (ตัวย่อประเทศ เช่น TH / ค่าเริ่มต้น: TH)

### 3. Get Trending Songs (ดึงข้อมูลเพลงฮิตทั้งหมด)

ดึงข้อมูลเพลงที่กำลังมาแรงตามประเทศเรียงตามยอดวิว

- **Method:** `GET`
- **URL:** `/trends`
- **Query Parameter:** `country` (ตัวย่อประเทศ เช่น TH / ค่าเริ่มต้น: TH)

### 4. Get New Releases (เพลงฮิตมาใหม่)

คัดกรองเฉพาะเทรนด์เพลงที่เพิ่งเผยแพร่ไม่เกิน 7 วัน

- **Method:** `GET`
- **URL:** `/trends/new`
- **Query Parameter:** `country` (ตัวย่อประเทศ เช่น TH / ค่าเริ่มต้น: TH)

### 5. Get MVs (เฉพาะ Official MV)

คัดกรองเฉพาะเทรนด์เพลงที่ระบบวิเคราะห์ว่าเป็น Music Video

- **Method:** `GET`
- **URL:** `/trends/mv`
- **Query Parameter:** `country` (ตัวย่อประเทศ เช่น TH / ค่าเริ่มต้น: TH)

**ตัวอย่าง Response (สำหรับ API เส้น `/trends`, `/trends/new`, `/trends/mv`):**

```json
{
  "success": true,
  "data": [
    {
      "id": "dQw4w9WgXcQ",
      "title": "Rick Astley - Never Gonna Give You Up",
      "channelTitle": "Rick Astley",
      "thumbnailUrl": "[https://i.ytimg.com/vi/dQw4w9WgXcQ/hqdefault.jpg](https://i.ytimg.com/vi/dQw4w9WgXcQ/hqdefault.jpg)",
      "viewCount": "1500000000",
      "countryCode": "US",
      "categoryId": "10",
      "publishedAt": "2026-08-16T10:00:00Z",
      "videoType": "MV"
    }
  ]
}
```

---

_Developed by Puttipong Doungvichai_
