package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// RunCommand executes shell commands and returns output
func RunCommand(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running %s: %v\nOutput: %s", cmd, err, out)
	}
	return string(out)
}

// DebugHandler checks file structure & logs env variables
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Debug request received")

	// Log Environment Variables (NOT exposed to browser)
	log.Println("🛠️ Environment Variables:")
	for _, e := range os.Environ() {
		log.Println(e)
	}

	// Check folders in /var/task (printed in browser)
	output := "📂 Folder structure in /var/task:\n"
	output += RunCommand("ls", "-lah", "/var/task") + "\n"
	output += "📂 Folder structure in /var/task/public:\n"
	output += RunCommand("ls", "-lah", "/var/task/public") + "\n"

	// Send folder structure to browser
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}
