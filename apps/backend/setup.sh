#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

if [ -f .env ]; then
  echo "⚠️  .env already exists, skipping copy"
else
  cp .env.example .env
  echo "✅ Created .env from .env.example"
fi

echo "🐳 Starting docker-compose services..."
docker-compose up -d

echo "📚 Generating Swagger documentation..."
# เช็กว่าเครื่องนี้ติดตั้งคำสั่ง swag หรือยัง ถ้ายังให้ติดตั้งอัตโนมัติ
if ! command -v swag &> /dev/null; then
    echo "   swag CLI not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
    # นำ Go bin path เข้าสู่ระบบชั่วคราวเพื่อให้เรียกใช้คำสั่ง swag ได้ทันที
    export PATH="$(go env GOPATH)/bin:$PATH"
fi

# สั่งสร้างไฟล์ Swagger Docs
swag init -g cmd/api/main.go --parseDependency --parseInternal
echo "✅ Swagger docs generated successfully."

echo ""
echo "🚀 Setup complete! Backend is ready."
echo "👉 Run 'go run cmd/api/main.go' to start the API."