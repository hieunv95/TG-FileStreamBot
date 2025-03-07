#!/bin/sh
set -e  # Dừng nếu có lỗi

# Tạo thư mục chứa TDLib
mkdir -p /tmp/tdlib

# Kiểm tra nếu TDLib chưa tồn tại thì tải về
if [ ! -f "/tmp/tdlib/lib/libtdjson.so" ]; then
    echo "⏳ Đang tải TDLib..."
    curl -L https://github.com/tdlib/td/releases/download/v1.8.0/tdlib-linux-x64.tar.gz | tar -xz -C /tmp/tdlib
    echo "✅ TDLib đã được cài đặt!"
fi

# Thiết lập biến môi trường để dùng CGO
export CGO_ENABLED=1
export CGO_CFLAGS="-I/tmp/tdlib/include"
export CGO_LDFLAGS="-L/tmp/tdlib/lib -ltdjson"

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
