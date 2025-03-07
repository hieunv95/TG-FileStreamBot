package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// listFiles executes the "ls -lah" command for a given directory
func listFiles(dir string) string {
	out, err := exec.Command("ls", "-lah", dir).Output()
	if err != nil {
		return fmt.Sprintf("Error listing files in %s: %v", dir, err)
	}
	return string(out)
}

// getEnvVar safely gets an environment variable
func getEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "(Not Set)"
	}
	return value
}

// DebugHandler is the exported function for Vercel
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	// Get current working directory
	wd, _ := os.Getwd()

	// List files in important directories
	taskFiles := listFiles("/var/task/")
	tmpFiles := listFiles("/tmp/")
	tdlibFiles := listFiles("/var/task/tdlib/lib/")
	vercelFiles := listFiles("/vercel/")
	vercelCacheFiles := listFiles("/.vercel/")
	outputFiles := listFiles("/output/")

	// Get relevant environment variables
	envVars := map[string]string{
		"CGO_ENABLED":    getEnvVar("CGO_ENABLED"),
		"CGO_CFLAGS":     getEnvVar("CGO_CFLAGS"),
		"CGO_LDFLAGS":    getEnvVar("CGO_LDFLAGS"),
		"LD_LIBRARY_PATH": getEnvVar("LD_LIBRARY_PATH"),
	}

	// Format environment variables as string
	var envOutput []string
	for key, value := range envVars {
		envOutput = append(envOutput, fmt.Sprintf("%s=%s", key, value))
	}

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

📂 **Files in /vercel/**
%s

📂 **Files in /.vercel/**
%s

📂 **Files in /output/**
%s

🛠 **Environment Variables**
%s

✅ **If you see missing files or variables, check build scripts & environment settings!**
`, wd, taskFiles, tmpFiles, tdlibFiles, vercelFiles, vercelCacheFiles, outputFiles, strings.Join(envOutput, "\n"))

	// Print debug info in logs
	log.Println(response)

	// Send response to client
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
