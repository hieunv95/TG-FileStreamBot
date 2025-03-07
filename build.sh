#!/bin/sh
set -e  # Dừng nếu có lỗi

# Thiết lập biến môi trường để dùng CGO
export CGO_ENABLED=1
export CGO_CFLAGS="-I/home/vercel/tdlib/include"
export CGO_LDFLAGS="-L/home/vercel/tdlib/lib -ltdjson"

# Tạo thư mục output cho Vercel
mkdir -p /vercel/output

# Build webhook.go
echo "🔧 Đang biên dịch webhook.go..."
go build -o /vercel/output/webhook api/webhook.go
echo "✅ Build webhook.go thành công!"

# Build file.go
echo "🔧 Đang biên dịch file.go..."
go build -o /vercel/output/file api/file.go
echo "✅ Build file.go thành công!"
