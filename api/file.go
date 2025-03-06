package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Đường dẫn lưu TDLib
const tdlibPath = "/tmp/tdlib"

// Kiểm tra TDLib đã cài chưa
func isTDLibInstalled() bool {
	_, err := os.Stat(filepath.Join(tdlibPath, "lib", "libtdjson.so"))
	return err == nil
}

// Tải và cài đặt TDLib nếu chưa có
func installTDLib() error {
	if isTDLibInstalled() {
		log.Println("✅ TDLib đã được cài đặt")
		return nil
	}

	log.Println("⏳ Đang cài đặt TDLib...")

	// Tải thư viện TDLib từ GitHub Release hoặc nguồn khác
	cmd := exec.Command("sh", "-c", `
		mkdir -p /tmp/tdlib &&
		curl -L https://github.com/tdlib/td/releases/download/v1.8.0/tdlib-linux-x64.tar.gz | tar -xz -C /tmp/tdlib
	`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("❌ Lỗi khi cài đặt TDLib:", err)
		return err
	}

	log.Println("✅ TDLib đã được cài đặt thành công!")
	return nil
}

// Proxy file qua Vercel
func ProxyFileHandler(w http.ResponseWriter, r *http.Request) {
	// Cài đặt TDLib trước khi xử lý request
	if err := installTDLib(); err != nil {
		http.Error(w, "Lỗi cài đặt TDLib", http.StatusInternalServerError)
		return
	}

	// Lấy file ID từ URL
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		http.Error(w, "Thiếu file_id", http.StatusBadRequest)
		return
	}

	// Gọi MTProto API để lấy direct link (giả sử có hàm GetDirectLink)
	directLink, err := GetDirectLink(fileID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lỗi lấy direct link: %v", err), http.StatusInternalServerError)
		return
	}

	// Proxy dữ liệu từ Telegram
	resp, err := http.Get(directLink)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lỗi proxy file: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy dữ liệu sang response
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
