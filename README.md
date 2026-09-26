# TuneTrend

วิดีโอยอดนิยมบน YouTube แยกตามประเทศ ทั้งเพลงและหมวดอื่นๆ อัปเดตอัตโนมัติทุก 3 ชั่วโมง ใช้งานได้ทั้งเว็บและแอปมือถือ

- เว็บ: https://tunetrend.pdouvch.com
- Android: [Google Play](https://play.google.com/store/apps/details?id=com.tunetrend.tunetrend_mobile)
- API: https://tunetrend-api.onrender.com

## สารบัญ

- [Features](#features)
- [สถาปัตยกรรม](#สถาปัตยกรรม)
- [Tech Stack](#tech-stack)
- [โครงสร้าง Repository](#โครงสร้าง-repository)
- [เริ่มต้นใช้งานบนเครื่อง](#เริ่มต้นใช้งานบนเครื่อง)
- [Environment Variables](#environment-variables)
- [API](#api)
- [Background Workers](#background-workers)
- [การทดสอบและคุณภาพโค้ด](#การทดสอบและคุณภาพโค้ด)
- [Deployment](#deployment)
- [การตั้งค่าขณะรัน (Runtime Configuration)](#การตั้งค่าขณะรัน-runtime-configuration)

## Features

### ฟีเจอร์หลัก (ทุก Platform)

- **ชาร์ตตามประเทศ**: ไทย (TH), เกาหลีใต้ (KR), ญี่ปุ่น (JP), สหรัฐอเมริกา (US), สหราชอาณาจักร (GB)
- **3 แท็บสำหรับหมวดเพลง**
  - **Trending**: เพลงยอดนิยม เรียงตามยอดวิว (สูงสุด 30 อันดับ)
  - **New Releases**: เพลงที่ปล่อยภายใน 7 วัน
  - **Music Videos**: เฉพาะคลิปที่ระบบจัดว่าเป็น Official MV
- **ตัวกรองหมวดหมู่**: 13 หมวดของ YouTube (Gaming, Entertainment, Sports, Comedy ฯลฯ) มีแท็บ Trending และ New Releases เปิด/ปิดแยกรายประเทศได้ และระบบซ่อนหมวดที่ดึงข้อมูลไม่สำเร็จให้อัตโนมัติ
- **ป้ายประเภทวิดีโอ**: MV, Lyric, Audio Track, Cover, Live Performance, General
- **เล่นวิดีโอในแอป** ผ่าน YouTube embedded player อย่างเป็นทางการ ยอดวิวและรายได้ของครีเอเตอร์นับตามปกติ
- **คิวเล่นต่อเนื่อง (Autoplay)**: ถามครั้งแรกเมื่อวิดีโอจบ แล้วเล่นตัวถัดไปหลังนับถอยหลัง 5 วินาที
- **แชร์วิดีโอ** เป็นลิงก์หน้าเล่นของ TuneTrend พร้อมลิงก์ดาวน์โหลดแอป
- **2 ภาษา**: ไทย / อังกฤษ (ตรวจจากภาษาเครื่อง และเปลี่ยนเองได้)
- **ธีมสว่าง / มืด** (ตามระบบ หรือเลือกเอง)
- **Onboarding tour** แนะนำการใช้งานครั้งแรก และเปิดดูซ้ำได้
- **ฟอร์มติดต่อผู้พัฒนา**: ส่งด้วยอีเมลหรือเบอร์โทรไทย มี rate limit และ honeypot กันสแปม
- **หน้า About** และลิงก์สนับสนุนผ่าน Buy Me a Coffee

### Android (แอป Flutter บน Google Play)

| หัวข้อ          | รายละเอียด                                                                                                                                                                        |
| --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| สถานะ           | เผยแพร่บน Google Play แล้ว (`com.tunetrend.tunetrend_mobile`)                                                                                                                     |
| เวอร์ชันระบบ    | Android 7.0 (API 24) ขึ้นไป, target API 36                                                                                                                                        |
| ประเทศเริ่มต้น  | เลือกจาก locale ของเครื่อง                                                                                                                                                        |
| ฟีเจอร์เฉพาะแอป | Push notification (Firebase Cloud Messaging, topic `general`), หน้าปิดปรับปรุงและหน้าบังคับอัปเดตผ่าน Remote Config, splash animation (Lottie), pull-to-refresh, skeleton loading |
| การเล่นวิดีโอ   | โทรศัพท์หมุนแนวนอนแล้วเข้าเต็มจออัตโนมัติ ส่วนแท็บเล็ตต้องกดปุ่มเต็มจอเอง                                                                                                         |
| Monitoring      | Crashlytics, Analytics (`country_changed`, `category_selected`, `video_played`), Performance Monitoring                                                                           |
| อัปเดต          | ผ่าน Google Play และรองรับ OTA patch ผ่าน Shorebird (เฉพาะ build ที่สร้างด้วย `shorebird release`)                                                                                |
| แคช             | แคชรูป thumbnail ในเครื่อง ส่วนข้อมูลชาร์ตดึงใหม่ทุกครั้ง                                                                                                                         |
| ข้อจำกัด        | ใช้ได้เฉพาะเมื่อออนไลน์ · ไม่จำประเทศที่เลือกไว้เมื่อเปิดแอปใหม่ · ไม่มี Discover carousel · แตะ notification แล้วยังไม่มี deep link                                              |

### iOS

| หัวข้อ                       | รายละเอียด                                                                                                                            |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| สถานะ                        | โค้ดแอป Flutter พร้อม build (`com.tunetrend.tunetrendMobile`) แต่**ยังไม่ได้เผยแพร่บน App Store**                                     |
| วิธีใช้ตอนนี้                | ผ่านเว็บใน Safari และติดตั้งเป็น PWA ได้ด้วย "Add to Home Screen" (เว็บแสดงวิธีทำให้อัตโนมัติ)                                        |
| ฟีเจอร์ของแอป (เมื่อเผยแพร่) | เหมือน Android ทุกข้อ รวมถึง push notification ผ่าน APNs                                                                              |
| ข้อจำกัด                     | หน้าบังคับอัปเดตยังลิงก์ไป Google Play · entitlement push ยังเป็น `development` · ต้องตั้ง APNs key ใน Firebase ก่อนถึงจะส่ง push ได้ |

### Web (Next.js บนมือถือ)

| หัวข้อ         | รายละเอียด                                                                                                                                                                                         |
| -------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| การเข้าใช้     | เบราว์เซอร์ทุกตัวบนมือถือ                                                                                                                                                                          |
| ประเทศเริ่มต้น | เลือกจาก `Accept-Language` ของเบราว์เซอร์ ทุกประเทศมี URL ของตัวเอง (เช่น `/th`, `/kr?tab=new&category=20`) จึงแชร์และ bookmark ได้                                                                |
| ติดตั้งแอป     | **Android**: modal และปุ่มชวนโหลดแอปจาก Google Play (manifest ประกาศ `prefer_related_applications` จึงไม่มีการติดตั้ง PWA บน Android) · **iOS**: ติดตั้งเป็น PWA ได้ · modal แสดงไม่เกินวันละครั้ง |
| Offline        | Service worker แคชเฉพาะ app shell จึงเปิดหน้าได้ตอนเน็ตหลุด แต่ไม่มีข้อมูลชาร์ตแบบ offline                                                                                                         |
| Layout         | รายการแบบ compact และเมนูแบบ drawer                                                                                                                                                                |
| ข้อจำกัด       | ไม่มี push notification · Discover carousel ไม่แสดงบนจอเล็ก · ข้อมูลอาจช้ากว่า backend สูงสุด 1 ชั่วโมงเพราะมีแคชฝั่ง server                                                                       |

### Desktop (Web บน Windows / macOS / Linux)

| หัวข้อ                | รายละเอียด                                                                                                                                                   |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| การเข้าใช้            | ผ่านเว็บในเบราว์เซอร์เท่านั้น **ไม่มีแอป desktop** (แอป Flutter รองรับเฉพาะ Android/iOS และ build บน macOS/Windows/Linux/Web จะ error ตอนเริ่ม Firebase)     |
| Layout หน้าเล่นวิดีโอ | player 75% + คิวด้านข้าง 25% ในหน้าจอเดียว ไม่ต้อง scroll                                                                                                    |
| ฟีเจอร์เฉพาะ desktop  | **Discover carousel**: วิดีโอฮิตข้ามหมวดและข้ามประเทศ เลื่อนอัตโนมัติทุก 3 วินาทีและหยุดเมื่อชี้เมาส์ · header มีปุ่ม About / ภาษา / ธีม / tour / โหลดแอปครบ |
| แชร์                  | ใช้ Web Share API ถ้าเบราว์เซอร์รองรับ ไม่งั้นคัดลอกลิงก์ลง clipboard                                                                                        |
| Monitoring            | Vercel Analytics และ Speed Insights                                                                                                                          |
| ข้อจำกัด              | ไม่มี push notification และไม่มีโหมด offline                                                                                                                 |

### ตารางเปรียบเทียบ

| ฟีเจอร์                                | Android |       iOS        | Web (มือถือ) |      Desktop (Web)      |
| -------------------------------------- | :-----: | :--------------: | :----------: | :---------------------: |
| ชาร์ต 5 ประเทศ / 3 แท็บ / ตัวกรองหมวด  |   ✅    |       ✅*        |      ✅      |           ✅            |
| เล่นวิดีโอ + Autoplay                  |   ✅    |       ✅*        |      ✅      |           ✅            |
| แชร์วิดีโอ                             |   ✅    |       ✅*        |      ✅      | ✅ (clipboard fallback) |
| ไทย / อังกฤษ, ธีมสว่าง/มืด, Onboarding |   ✅    |       ✅*        |      ✅      |           ✅            |
| ฟอร์มติดต่อ                            |   ✅    |       ✅*        |      ✅      |           ✅            |
| Discover carousel                      |   ❌    |        ❌        |      ❌      |           ✅            |
| URL แยกต่อหน้า (deep link / bookmark)  |   ❌    |        ❌        |      ✅      |           ✅            |
| Push notification                      |   ✅    |       ✅*        |      ❌      |           ❌            |
| หน้า Maintenance / Force update        |   ✅    |       ✅*        |      ❌      |           ❌            |
| Pull-to-refresh                        |   ✅    |       ✅*        |      ❌      |           ❌            |
| ติดตั้งแบบ PWA                         |   ❌    | ✅ (ผ่าน Safari) | iOS เท่านั้น |     ตามเบราว์เซอร์      |
| เปิดได้ตอน offline                     |   ❌    |        ❌        |  app shell   |        app shell        |
| OTA update (Shorebird)                 |   ✅    |       ✅*        |      –       |            –            |

\* ฟีเจอร์ของแอป iOS ที่ build ได้แล้วแต่ยังไม่ได้เผยแพร่บน App Store

## สถาปัตยกรรม

```
                 YouTube Data API v3
                         ▲  (worker ดึงตามรอบเวลา)
┌──────────────┐   REST  │                        SQL
│ Mobile app   │────────▶ Backend (Go / Fiber) ─────────▶ PostgreSQL
│ Flutter      │         │  - REST API (อ่านอย่างเดียว + contact)
└──────┬───────┘         │  - background workers
       │ Firebase        │  - rate limit (Upstash Redis)
       ▼                 │  - อีเมลแจ้งเตือน (Resend)
   Firebase              ▲
                         │ fetch ฝั่ง server (แคช 1 ชม.)
               Web (Next.js บน Vercel) ◀──── Browser / PWA
```

- Backend **ไม่ได้เรียก YouTube ตอนมี request เข้ามา** แต่อ่านจากฐานข้อมูลที่ worker ซิงก์ไว้ล่วงหน้า จึงตอบเร็วและประหยัด quota
- เว็บเรียก backend จากฝั่ง server เท่านั้น URL ของ backend จึงไม่ไปถึงเบราว์เซอร์
- แอปมือถือเรียก backend ตรง โดยกำหนด URL ตอน build ด้วย `--dart-define`
- Backend แบ่งชั้นแบบ Clean Architecture: `domain` → `usecase` → `repository` → `delivery/http`

## Tech Stack

| ส่วน           | เทคโนโลยี                                                                                                                                                                         |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Backend        | Go 1.26, Fiber v2, GORM, PostgreSQL 15, Swagger (swaggo), Air (hot reload)                                                                                                        |
| Web            | Next.js 16 (App Router), React 19, TypeScript, Tailwind CSS 4, react-youtube, driver.js, Upstash Ratelimit                                                                        |
| Mobile         | Flutter (Dart ^3.13), Dio, youtube_player_iframe, showcaseview, cached_network_image, Lottie, Firebase (Crashlytics, Analytics, Performance, Remote Config, Messaging), Shorebird |
| Infrastructure | Render (backend), Vercel (web), Supabase (Postgres), Upstash (Redis REST), Resend (อีเมล), Google Play                                                                            |
| CI/CD          | GitHub Actions, Codecov, Dependabot, Husky + commitlint + lint-staged + Prettier                                                                                                  |

## โครงสร้าง Repository

```
tunetrend/
├── apps/
│   ├── backend/     # Go API + background workers
│   ├── frontend/    # เว็บ Next.js
│   └── mobile/      # แอป Flutter (Android / iOS)
├── .github/         # CI workflow + Dependabot
├── .husky/          # git hooks
└── package.json     # tooling ระดับ root (husky, commitlint, prettier)
```

## เริ่มต้นใช้งานบนเครื่อง

### สิ่งที่ต้องมี

- Go 1.26+, Docker (สำหรับ PostgreSQL)
- Node.js 22+
- Flutter stable (CI ใช้ 3.47.2) และ Android Studio / Xcode
- YouTube Data API v3 key

### 0. Root tooling

```bash
npm install          # ติดตั้ง husky, commitlint, prettier, lint-staged
```

### 1. Backend

```bash
cd apps/backend
./setup.sh           # สร้าง .env, เปิด Postgres ใน Docker, generate Swagger
```

แก้ `apps/backend/.env`:

- ใส่ `YOUTUBE_API_KEY`
- ให้ `DB_USER` / `DB_PASSWORD` ตรงกับ `docker-compose.yml` (`root` / `secretpassword`)

```bash
go run ./cmd/api     # หรือ `air` สำหรับ hot reload
curl localhost:8080/health
open http://localhost:8080/docs/index.html   # ต้องตั้ง SWAGGER_ENABLED=true
```

รอบแรกที่ start ระบบจะสร้างตาราง ใส่ค่าเริ่มต้น และซิงก์ข้อมูลจาก YouTube ให้อัตโนมัติ

### 2. Web

```bash
cd apps/frontend
./setup.sh           # สร้าง .env + npm install
npm run dev          # http://localhost:3000
```

### 3. Mobile

```bash
cd apps/mobile
./scripts/start.sh                    # เปิด emulator Android และรันแอปกับ backend local (10.0.2.2:8080)
./scripts/start.sh --with-backend     # start Postgres + backend ให้ด้วย
./scripts/start.sh --api https://tunetrend-api.onrender.com
./scripts/start.sh --help
```

หรือรันด้วยคำสั่ง Flutter โดยตรง:

```bash
flutter pub get
flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080   # Android emulator
flutter run --dart-define=API_BASE_URL=http://localhost:8080  # iOS simulator
```

Release build ใช้ได้เฉพาะ HTTPS เพราะไม่ได้เปิด cleartext traffic

## Environment Variables

### Backend (`apps/backend/.env`)

| ตัวแปร                                                                                 | จำเป็น | ค่าเริ่มต้น | คำอธิบาย                                                           |
| -------------------------------------------------------------------------------------- | :----: | ----------- | ------------------------------------------------------------------ |
| `DATABASE_URL`                                                                         |  ✅*   | –           | Postgres URI (ถ้าตั้งไว้จะใช้ตัวนี้ก่อน)                           |
| `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `DB_SSLMODE`, `DB_TIMEZONE` |  ✅*   | –           | ใช้แทน `DATABASE_URL`                                              |
| `APP_PORT`                                                                             |        | `8080`      | พอร์ตที่ server listen                                             |
| `APP_ENV`                                                                              |        | –           | ตั้ง `prod` เพื่อให้ Swagger ใช้ `SWAGGER_HOST` และ https          |
| `SWAGGER_ENABLED`                                                                      |        | ปิด         | ตั้ง `true` เพื่อเปิด `/docs`                                      |
| `SWAGGER_HOST`                                                                         |        | –           | domain ที่แสดงใน Swagger ตอน prod                                  |
| `CORS_ALLOW_ORIGINS`                                                                   |        | `*`         | origin ที่อนุญาต คั่นด้วย comma                                    |
| `YOUTUBE_API_KEY`                                                                      |   ✅   | –           | YouTube Data API v3                                                |
| `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN`                                   |        | –           | ถ้าไม่ตั้ง API ฝั่งอ่านจะไม่จำกัด request และ `/contact` จะตอบ 503 |
| `RESEND_API_KEY`, `CONTACT_EMAIL_TO`                                                   |        | –           | อีเมลแจ้งเตือนเมื่อมีข้อความ contact                               |

\* ต้องตั้งอย่างใดอย่างหนึ่ง

### Web (`apps/frontend/.env`)

| ตัวแปร                                               | ฝั่ง   | ค่าเริ่มต้น             | คำอธิบาย                                        |
| ---------------------------------------------------- | ------ | ----------------------- | ----------------------------------------------- |
| `API_BASE_URL`                                       | server | `http://localhost:8080` | URL ของ backend (ห้ามใส่ prefix `NEXT_PUBLIC_`) |
| `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN` | server | –                       | rate limit ของเว็บ                              |
| `NEXT_PUBLIC_INSTALL_ANDROID_MODE`                   | client | `play-store`            | `play-store` / `pwa` / `disabled`               |
| `NEXT_PUBLIC_INSTALL_IOS_MODE`                       | client | `pwa`                   | `app-store` / `pwa` / `disabled`                |
| `NEXT_PUBLIC_PLAY_STORE_URL`                         | client | ลิงก์ Play Store ของแอป |                                                 |
| `NEXT_PUBLIC_APP_STORE_URL`                          | client | ว่าง                    | ต้องตั้งก่อนใช้โหมด `app-store`                 |

ค่า `NEXT_PUBLIC_*` ถูกฝังตอน build เปลี่ยนแล้วต้อง build หรือ deploy ใหม่

### Mobile (`--dart-define`)

| ตัวแปร         | ค่าเริ่มต้น            | คำอธิบาย        |
| -------------- | ---------------------- | --------------- |
| `API_BASE_URL` | `http://10.0.2.2:8080` | URL ของ backend |

## API

Response ทุกเส้นเป็นรูป `{ "success": true, "data": ... }` หรือ `{ "success": false, "error": "..." }`

| Method | Path          | Query / Body                                      | คำอธิบาย                                      |
| ------ | ------------- | ------------------------------------------------- | --------------------------------------------- |
| GET    | `/health`     | –                                                 | ตรวจสถานะ API และฐานข้อมูล                    |
| GET    | `/trends`     | `country`, `category`                             | ชาร์ตยอดนิยม (ไม่ใส่ `category` = หมวดเพลง)   |
| GET    | `/trends/new` | `country`, `category`                             | เผยแพร่ภายใน 7 วัน                            |
| GET    | `/trends/mv`  | `country`                                         | เฉพาะ Official MV                             |
| GET    | `/categories` | `country`                                         | หมวดหมู่ที่เปิดใช้งานในประเทศนั้น             |
| GET    | `/discover`   | –                                                 | วิดีโอฮิตไม่เกิน 2 รายการต่อหมวด รวมทุกประเทศ |
| POST   | `/contact`    | `name`, `message`, `contactEmail`, `contactPhone` | ส่งข้อความหาผู้พัฒนา                          |

- `country` มีค่าเริ่มต้นเป็น `TH` · ส่ง `category` ที่ระบบไม่รู้จักจะได้ `400`
- Rate limit: กลุ่ม GET 60 request/นาที/IP · `/contact` 5 ครั้ง/10 นาที/IP
- เอกสาร API ฉบับเต็ม: Swagger ที่ `/docs/index.html`

## Background Workers

Worker รันอยู่ใน process เดียวกับ API ทำงานทันทีเมื่อ start แล้ววนตามรอบที่ตั้งไว้ในตาราง `worker_settings`

| Job                          | รอบเริ่มต้น | หน้าที่                                               |
| ---------------------------- | ----------- | ----------------------------------------------------- |
| `youtube_sync`               | 3 ชม.       | ซิงก์ชาร์ตเพลงทุกประเทศ (มี advisory lock กันรันซ้อน) |
| `<category>_videos` (13 job) | 3 ชม.       | ซิงก์ชาร์ตของแต่ละหมวด                                |
| `video_category_sync`        | 24 ชม.      | ซิงก์รายชื่อหมวดและคำนวณสถานะเปิด/ปิด                 |
| category resume              | 7 วัน       | ลองซิงก์หมวดที่ระบบปิดอัตโนมัติอีกครั้ง               |
| api log cleanup              | 24 ชม.      | ลบ `api_logs` ที่เก่ากว่า 14 วัน                      |

หมวดที่ดึงข้อมูลล้มเหลวต่อเนื่องเกิน grace period (24 ชม.) จะถูกซ่อนอัตโนมัติ (`deactivated_reason = auto_fetch_failure`) และเปิดกลับเองเมื่อซิงก์สำเร็จ หมวดที่แอดมินปิดเองด้วยค่า reason อื่น ระบบจะไม่เปิดคืนให้

## การทดสอบและคุณภาพโค้ด

```bash
# Backend
cd apps/backend && go vet ./... && go test ./...

# Web
cd apps/frontend && npm run lint && npm run build

# Mobile
cd apps/mobile && flutter analyze && flutter test
```

- Commit message ต้องเป็นแบบ [Conventional Commits](https://www.conventionalcommits.org/) (ตรวจด้วย commitlint ทั้ง local และบน PR)
- pre-commit รัน ESLint (frontend), Prettier และ `gofmt` เฉพาะไฟล์ที่ stage
- CI (`.github/workflows/ci.yml`) รันเมื่อมี PR หรือ push เข้า `main`: commitlint, frontend lint + build, backend vet + build + test + coverage (Codecov)

## Deployment

| ส่วน     | แพลตฟอร์ม                 | วิธี deploy                                                                                                                                                                  |
| -------- | ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Database | Supabase (PostgreSQL)     | schema สร้างอัตโนมัติด้วย GORM AutoMigrate เมื่อ backend start                                                                                                               |
| Backend  | Render                    | deploy อัตโนมัติจาก `main` · root `apps/backend` · build `go build -o app ./cmd/api` · start `./app` · health check `/health`                                                |
| Web      | Vercel                    | deploy อัตโนมัติจาก `main` · root `apps/frontend` · มี preview deployment ทุก PR                                                                                             |
| Android  | Google Play               | GitHub Actions → **CI → Run workflow** (job `mobile-release`) ได้ `.aab` / `.apk` ที่ sign แล้ว พร้อม debug symbols หรือใช้ `shorebird release android` ถ้าต้องการ OTA patch |
| iOS      | App Store (ยังไม่เผยแพร่) | `flutter build ipa` หรือ `shorebird release ios`                                                                                                                             |

ข้อควรรู้:

- **Backend รันแค่ 1 instance** เพราะ worker หมวดหมู่ไม่มี lock กันรันซ้ำ
- บนแพลนที่ service หลับเมื่อไม่มีคนใช้ worker จะหยุดด้วย ข้อมูลจึงไม่อัปเดตจนกว่า service จะถูกปลุก
- Mobile release ต้องเพิ่มเลข build ใน `apps/mobile/pubspec.yaml` (`version: X.Y.Z+N`) ทุกครั้ง
- GitHub Secrets ที่ต้องมีสำหรับ `mobile-release`: `ANDROID_KEYSTORE_BASE64`, `ANDROID_KEYSTORE_PASSWORD`, `ANDROID_KEY_PASSWORD`, `ANDROID_KEY_ALIAS` (และ `CODECOV_TOKEN` สำหรับ coverage)
- Build แบบ `--obfuscate` ต้องอัปโหลด debug symbols ไป Crashlytics ถึงจะอ่าน stack trace ได้

## การตั้งค่าขณะรัน (Runtime Configuration)

| ต้องการ                          | ทำที่                                                                                    | ต้อง restart backend |
| -------------------------------- | ---------------------------------------------------------------------------------------- | :------------------: |
| ปิดแอปมือถือชั่วคราว             | Firebase Remote Config `is_maintenance = true`                                           |          –           |
| บังคับอัปเดตแอป                  | Firebase Remote Config `min_app_version = X.Y.Z`                                         |          –           |
| ส่ง push notification            | Firebase Messaging → topic `general`                                                     |          –           |
| ซ่อนหมวดในประเทศใดประเทศหนึ่ง    | `UPDATE video_categories SET is_active = false, deactivated_reason = 'manual' WHERE ...` |          ❌          |
| เปลี่ยนรอบซิงก์หรือรายชื่อประเทศ | ตาราง `worker_settings`                                                                  |          ✅          |
| เพิ่มหมวดหมู่ใหม่                | `INSERT INTO category_video_configs (category_id, table_name, label) ...`                |          ✅          |

การเพิ่มประเทศใหม่ต้องแก้ทั้ง `worker_settings.countries`, `apps/backend/internal/usecase/discover_usecase.go`, `apps/frontend/lib/countries.ts` และ `apps/mobile/lib/constants/countries.dart`
