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

# Build ứng dụng
echo "🔧 Đang biên dịch ứng dụng..."
go build -o /vercel/output/server ./cmd/main.go
echo "✅ Build thành công!"
