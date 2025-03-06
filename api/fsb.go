package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/bot"

	"github.com/spf13/cobra"
)

type Response struct {
	Message string `json:"message"`
}

const versionString = "3.0.0"

var rootCmd = &cobra.Command{
	Use:               "fsb [command]",
	Short:             "Telegram File Stream Bot",
	Long:              "Telegram Bot to generate direct streamable links for telegram media.",
	Example:           "fsb run --port 8080",
	Version:           versionString,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	config.SetFlagsFromConfig(runCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.SetVersionTemplate(fmt.Sprintf(`Telegram File Stream Bot version %s`, versionString))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Failed to start bot: %v", err), http.StatusInternalServerError)
		return os.Exit(1)
	}

	// Trả về JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Telegram File Stream Bot is running on Vercel!"})
}
