package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
    "path/filepath"

	"github.com/zelenin/go-tdlib/client"
)

func initTDLib() (*client.Client, error) {
	// Chuyển đổi API_ID từ string → int
	apiID, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		log.Fatalf("❌ API_ID không hợp lệ: %v", err)
		return nil, err
	}

	// Lấy API_HASH từ biến môi trường
	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		log.Fatal("❌ Thiếu API_HASH")
		return nil, err
	}

	// Khởi tạo TDLib Client
	tdlibParameters := &client.SetTdlibParametersRequest{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  int32(apiID),
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
	}

	// botToken := "000000000:gsVCGG5YbikxYHC7bP5vRvmBqJ7Xz6vG6td"
    authorizer := client.BotAuthorizer(tdlibParameters, os.Getenv("BOT_TOKEN"))
	
	if err != nil {
		log.Fatalf("SetLogVerbosityLevel error: %s", err)
	}
	
    tdlibClient, err := client.NewClient(authorizer)
    if err != nil {
        log.Fatalf("NewClient error: %s", err)
    }

	versionOption, err := client.GetOption(&client.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		log.Fatalf("GetOption error: %s", err)
	}

	commitOption, err := client.GetOption(&client.GetOptionRequest{
		Name: "commit_hash",
	})
	if err != nil {
		log.Fatalf("GetOption error: %s", err)
	}

	log.Printf("TDLib version: %s (commit: %s)", versionOption.(*client.OptionValueString).Value, commitOption.(*client.OptionValueString).Value)

    me, err := tdlibClient.GetMe()
    if err != nil {
        log.Fatalf("GetMe error: %s", err)
    }

    log.Printf("Me: %s %s", me.FirstName, me.LastName)
	return tdlibClient, err
}

// Proxy file Telegram qua Vercel
func FileProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy file ID từ URL
	fileID := r.URL.Path[len("/file/"):]

	// Khởi tạo Telegram MTProto Client
	tdlibClient, err := initTDLib()
	if err != nil {
		log.Println("Lỗi khởi tạo Telegram client:", err)
		http.Error(w, "Lỗi server", http.StatusInternalServerError)
		return
	}

	// Chuyển đổi fileID từ string → int
	id, err := strconv.Atoi(fileID)
	if err != nil {
		log.Printf("❌ Lỗi chuyển đổi fileID: %v", err)
		return
	}

	// Lấy thông tin file
	file, err := tdlibClient.GetFile(&client.GetFileRequest{FileId: int32(id)})
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
