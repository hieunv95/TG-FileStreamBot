package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Lấy API ID, Hash và Bot Token từ biến môi trường
var (
	apiID   = os.Getenv("API_ID")
	apiHash = os.Getenv("API_HASH")
	botToken = os.Getenv("BOT_TOKEN")
	vercelDomain = os.Getenv("HOST") // VD: "your-vercel-app.vercel.app"
)

// Cấu trúc JSON nhận từ Webhook
type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int    `json:"message_id"`
		Text      string `json:"text"`
		Chat      struct {
			ID   int64  `json:"id"`
			Type string `json:"type"`
		} `json:"chat"`
		Document *struct { // Thêm field Document để chứa thông tin file
			FileID       string `json:"file_id"`
			FileUniqueID string `json:"file_unique_id"`
			FileName     string `json:"file_name"`
			MimeType     string `json:"mime_type"`
			FileSize     int    `json:"file_size"`
		} `json:"document,omitempty"` // `omitempty` để bỏ qua nếu không có file
	} `json:"message"`
}

// Gửi tin nhắn chứa direct link
func sendMessage(chatID int64, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	body, _ := json.Marshal(message)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Lỗi gửi tin nhắn:", err)
	}
}

// Xử lý Webhook từ Telegram
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Println("Lỗi decode JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Kiểm tra nếu message rỗng
	if update.Message.Chat.ID == 0 {
		log.Println("Không có message")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Kiểm tra chỉ lắng nghe private chat
	if update.Message.Chat.Type != "private" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Kiểm tra có file hay không
	if update.Message.Document.FileID == "" {
		log.Println("Không có file trong tin nhắn")
		w.WriteHeader(http.StatusOK)
		return
	}

	chatID := update.Message.Chat.ID
	fileID := update.Message.Document.FileID

	// Tạo direct link với domain của Vercel
	directLink := fmt.Sprintf("https://%s/file/%s", vercelDomain, fileID)

	// Gửi direct link cho người dùng
	sendMessage(chatID, "Direct link: "+directLink)
	w.WriteHeader(http.StatusOK)
}

