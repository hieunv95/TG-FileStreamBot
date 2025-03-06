package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/zelenin/go-tdlib/client"
)

// Proxy file Telegram qua Vercel
func FileProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy file ID từ URL
	fileID := r.URL.Path[len("/file/"):]

	// Khởi tạo Telegram MTProto Client
	tdlibClient, err := client.NewClient(client.Config{
		APIID:   os.Getenv("API_ID"),
		APIHash: os.Getenv("API_HASH"),
	})
	if err != nil {
		log.Println("Lỗi khởi tạo Telegram client:", err)
		http.Error(w, "Lỗi server", http.StatusInternalServerError)
		return
	}

	// Lấy thông tin file
	file, err := tdlibClient.GetFile(&client.GetFileRequest{FileID: fileID})
	if err != nil {
		log.Println("Lỗi lấy file:", err)
		http.Error(w, "Không tìm thấy file", http.StatusNotFound)
		return
	}

	// Proxy file từ Telegram về client
	telegramURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv("BOT_TOKEN"), file.Local.Path)
	resp, err := http.Get(telegramURL)
	if err != nil {
		log.Println("Lỗi tải file từ Telegram:", err)
		http.Error(w, "Lỗi tải file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Gửi dữ liệu file về client
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Local.Path)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
