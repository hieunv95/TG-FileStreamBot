package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const telegramAPI = "https://api.telegram.org/bot"

type TelegramUpdate struct {
	Message struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		ForwardFrom struct {
			ID int64 `json:"id"`
		} `json:"forward_from"`
		Document struct {
			FileID string `json:"file_id"`
		} `json:"document"`
	} `json:"message"`
}

func getFileDirectLink(fileID string) (string, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return "", fmt.Errorf("missing BOT_TOKEN")
	}

	// Lấy đường dẫn file từ Telegram API
	url := fmt.Sprintf("%s%s/getFile?file_id=%s", telegramAPI, botToken, fileID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}

	body, _ := io.ReadAll(resp.Body)
	log.Println("Telegram getFile response:", string(body)) // Debug API response

	if err := json.Unmarshal(body, &result); err != nil || !result.OK {
		return "", fmt.Errorf("failed to get file path: %s", string(body))
	}

	// Trả về link trực tiếp
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, result.Result.FilePath), nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Đọc toàn bộ request body để debug
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Println("Received request body:", string(body)) // Debug Request Body

	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		log.Println("Invalid request JSON:", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	// ✅ Kiểm tra file forward đúng format không
	if update.Message.Document.FileID == "" || update.Message.ForwardFrom.ID == 0 {
		log.Println("No forwarded file found:", update)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Lấy link tải trực tiếp
	fileURL, err := getFileDirectLink(update.Message.Document.FileID)
	if err != nil {
		log.Println("Error getting file URL:", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Gửi link về Telegram
	botToken := os.Getenv("BOT_TOKEN")
	sendURL := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=%s", telegramAPI, botToken, update.Message.Chat.ID, fileURL)
	resp, err := http.Get(sendURL)
	if err != nil {
		log.Println("Failed to send message:", err)
		w.WriteHeader(http.StatusOK)
		return
	}
	defer resp.Body.Close()

	// Debug Telegram sendMessage response
	body, _ = io.ReadAll(resp.Body)
	log.Println("Telegram sendMessage response:", string(body))

	// Phản hồi Vercel (luôn HTTP 200)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"file_url": fileURL})
}
