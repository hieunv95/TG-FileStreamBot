#!/bin/sh
set -e  # Dừng nếu có lỗi

# Cài đặt Go nếu chưa có
if ! command -v go &> /dev/null
then
    echo "⏳ Đang cài đặt Go..."
    curl -LO https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo "✅ Go đã được cài đặt!"
fi

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
