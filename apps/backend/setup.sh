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
if ! command -v swag &> /dev/null; then
    echo "   swag CLI not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
    export PATH="$(go env GOPATH)/bin:$PATH"
fi

swag init -g cmd/api/main.go --parseDependency --parseInternal
echo "✅ Swagger docs generated successfully."

echo ""
echo "🚀 Setup complete! Backend is ready."
echo "👉 Run 'go run cmd/api/main.go' to start the API."