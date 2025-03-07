package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// listHiddenFiles lists all files, including hidden ones, in a directory
func listHiddenFiles(dir string) string {
	out, err := exec.Command("ls", "-lah", dir).Output()
	if err != nil {
		return fmt.Sprintf("Error listing files in %s: %v", dir, err)
	}
	return string(out)
}

// getAllEnv returns all environment variables as a formatted string
func getAllEnv() string {
	envs := os.Environ()
	return strings.Join(envs, "\n")
}

// DebugHandler is an exported function for Vercel
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	taskFiles := listHiddenFiles("/var/runtime/")
	envVars := getAllEnv()

	response := fmt.Sprintf(`
📂 **Hidden Files in /var/runtime/**
--------------------------------
%s

🛠️ **Environment Variables**
--------------------------------
%s
`, taskFiles, envVars)

	// Log output for debugging
	log.Println(response)

	// Send response to client
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
