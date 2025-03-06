package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"EverythingSuckz/fsb/config"

	"EverythingSuckz/fsb/lib/bot"
	"EverythingSuckz/fsb/lib/cache"
	"EverythingSuckz/fsb/lib/routes"
	"EverythingSuckz/fsb/lib/types"
	"EverythingSuckz/fsb/lib/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"


	"EverythingSuckz/fsb/pkg/qrlogin"

	"github.com/spf13/cobra"
)

type Response struct {
	Message string `json:"message"`
}

const versionString = "3.0.0"

var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Run the bot with the given configuration.",
	DisableSuggestions: false,
	Run:                runApp,
}

var startTime time.Time = time.Now()

func runApp(cmd *cobra.Command, args []string) {
	utils.InitLogger()
	log := utils.Logger
	mainLogger := log.Named("Main")
	mainLogger.Info("Starting server")
	config.Load(log, cmd)
	router := getRouter(log)

	mainBot, err := bot.StartClient(log)
	if err != nil {
		log.Panic("Failed to start main bot", zap.Error(err))
	}
	cache.InitCache(log)
	workers, err := bot.StartWorkers(log)
	if err != nil {
		log.Panic("Failed to start workers", zap.Error(err))
		return
	}
	workers.AddDefaultClient(mainBot, mainBot.Self)
	bot.StartUserBot(log)
	mainLogger.Info("Server started", zap.Int("port", config.ValueOf.Port))
	mainLogger.Info("File Stream Bot", zap.String("version", versionString))
	mainLogger.Sugar().Infof("Server is running at %s", config.ValueOf.Host)
	err = router.Run(fmt.Sprintf(":%d", config.ValueOf.Port))
	if err != nil {
		mainLogger.Sugar().Fatalln(err)
	}
}

func getRouter(log *zap.Logger) *gin.Engine {
	if config.ValueOf.Dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(gin.ErrorLogger())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, types.RootResponse{
			Message: "Server is running.",
			Ok:      true,
			Uptime:  utils.TimeFormat(uint64(time.Since(startTime).Seconds())),
			Version: versionString,
		})
	})
	routes.Load(log, router)
	return router
}

var sessionCmd = &cobra.Command{
	Use:                "session",
	Short:              "Generate a string session.",
	DisableSuggestions: false,
	Run:                generateSession,
}

func init() {
	sessionCmd.Flags().StringP("login-type", "T", "qr", "The login type to use. Can be either 'qr' or 'phone'")
	sessionCmd.Flags().Int32P("api-id", "I", 0, "The API ID to use for the session (required).")
	sessionCmd.Flags().StringP("api-hash", "H", "", "The API hash to use for the session (required).")
	sessionCmd.MarkFlagRequired("api-id")
	sessionCmd.MarkFlagRequired("api-hash")
}

func generateSession(cmd *cobra.Command, args []string) {
	loginType, _ := cmd.Flags().GetString("login-type")
	apiId, _ := cmd.Flags().GetInt32("api-id")
	apiHash, _ := cmd.Flags().GetString("api-hash")
	if loginType == "qr" {
		qrlogin.GenerateQRSession(int(apiId), apiHash)
	} else if loginType == "phone" {
		generatePhoneSession()
	} else {
		fmt.Println("Invalid login type. Please use either 'qr' or 'phone'")
	}
}

func generatePhoneSession() {
	fmt.Println("Phone session is not implemented yet.")
}

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
		return
	}

	// Trả về JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Telegram File Stream Bot is running on Vercel!"})
}
