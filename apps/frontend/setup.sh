#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

if [ -f .env ]; then
  echo "⚠️  .env already exists, skipping copy"
else
  cp .env.example .env
  echo "✅ Created .env from .env.example"
fi

echo "📦 Installing npm dependencies..."
npm install
echo "✅ Dependencies installed."

echo ""
echo "🚀 Setup complete! Frontend is ready."
echo "👉 Run 'npm run dev' to start the app (make sure the backend from apps/backend is running too)."
