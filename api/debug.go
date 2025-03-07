package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// listFiles lists all files (including hidden) in a directory
func listFiles(dir string) string {
	out, err := exec.Command("ls", "-lah", dir).Output()
	if err != nil {
		return fmt.Sprintf("Error listing %s: %v", dir, err)
	}
	return string(out)
}

// getAllEnv retrieves all OS environment variables
func getAllEnv() string {
	return strings.Join(os.Environ(), "\n")
}

// DebugHandler logs and returns system info
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	paths := []string{"/var/task/", "/var/runtime/", "/var/lang/", "/tmp/", "/"}
	logData := ""

	for _, path := range paths {
		logData += fmt.Sprintf("📂 Listing %s:\n%s\n\n", path, listFiles(path))
	}

	logData += fmt.Sprintf("🛠️ Environment Variables:\n%s\n", getAllEnv())

	// Log and respond
	log.Println(logData)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(logData))
}
