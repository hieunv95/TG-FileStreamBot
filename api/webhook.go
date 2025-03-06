package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const telegramAPI = "https://api.telegram.org/bot"

type TelegramUpdate struct {
	Message struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		ForwardFromChat struct {
			ID int64 `json:"id"`
		} `json:"forward_from_chat"`
		Document struct {
			FileID string `json:"file_id"`
		} `json:"document"`
	} `json:"message"`
}

func getFileDirectLink(fileID string) (string, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return "", fmt.Errorf("missing bot token")
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

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.OK {
		return "", fmt.Errorf("failed to get file path")
	}

	// Trả về link trực tiếp
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, result.Result.FilePath), nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var update TelegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Kiểm tra tin nhắn có chứa file không
	if update.Message.Document.FileID == "" || update.Message.ForwardFromChat.ID == 0 {
		http.Error(w, "No forwarded file found", http.StatusBadRequest)
		return
	}

	// Lấy link tải trực tiếp
	fileURL, err := getFileDirectLink(update.Message.Document.FileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Gửi link về Telegram
	botToken := os.Getenv("BOT_TOKEN")
	sendURL := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=%s", telegramAPI, botToken, update.Message.Chat.ID, fileURL)
	http.Get(sendURL)

	// Phản hồi Vercel
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"file_url": fileURL})
}
