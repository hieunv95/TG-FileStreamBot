package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// listFiles executes the "ls -lah" command for a given directory
func listFiles(dir string) string {
	out, err := exec.Command("ls", "-lah", dir).Output()
	if err != nil {
		return fmt.Sprintf("Error listing files in %s: %v", dir, err)
	}
	return string(out)
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	// Get current working directory
	wd, _ := os.Getwd()

	// List files in /var/task (where Vercel deploys the app)
	taskFiles := listFiles("/var/task/")

	// List files in /tmp/ (temporary storage in Vercel)
	tmpFiles := listFiles("/tmp/")

	// List files in /var/task/tdlib/lib/ (if TDLib is being used)
	tdlibFiles := listFiles("/var/task/tdlib/lib/")

	// Prepare the response
	response := fmt.Sprintf(`
🖥️ **Vercel Serverless Debug Info**
-----------------------------------
📂 **Current Working Directory:** %s

📂 **Files in /var/task/**
%s

📂 **Files in /tmp/**
%s

📂 **Files in /var/task/tdlib/lib/**
%s

✅ **If you see missing files, check build scripts & environment settings!**
`, wd, taskFiles, tmpFiles, tdlibFiles)

	// Print debug info in logs
	log.Println(response)

	// Send response to client
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
