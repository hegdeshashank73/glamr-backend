package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hegdeshashank73/glamr-backend/handlers"
	"github.com/hegdeshashank73/glamr-backend/middlewares"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var ginEngine *gin.Engine
var db *sql.DB
var handler *handlers.Handler
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the HTTP Server on Port 1729",
	Run: func(cmd *cobra.Command, args []string) {
		// Open Telemetry
		// cleanup := common.InitTracer()
		// defer cleanup(context.Background())
		setupGin()
		setupRoutes()
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func setupGin() {
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	ginEngine = gin.Default()
	ginEngine.Use(otelgin.Middleware(viper.GetString("OTEL_SERVICE_NAME")))
	ginEngine.Use(middlewares.HandlePanic)
	ginEngine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})
	handler = handlers.NewHandler(db)
}

func setupRoutes() {
	authGroup := ginEngine.Group("/")
	authGroup.Use(middlewares.Auth())

	privateGroup := ginEngine.Group("/")
	privateGroup.Use(middlewares.Private(true))

	authNoAuthGroup := ginEngine.Group("/")
	authNoAuthGroup.Use(middlewares.AuthNoAuth())

	authOrPrivateGroup := ginEngine.Group("/")
	authOrPrivateGroup.Use(middlewares.AuthOrPrivate())
	// Health
	ginEngine.GET("/health", handlers.HealthHandler)
	ginEngine.GET("/", handlers.HealthHandler)

	// Metadata
	ginEngine.GET("/test", handlers.TestHandler)

	// Auth: Magiclinks
	// ginEngine.POST("/auth/magiclink", handlers.CreateMagiclinkHandler)
	// ginEngine.POST("/auth/magiclink/verify", handlers.VerifyMagiclinkHandler)

	authNoAuthGroup.POST("/assets/upload", handlers.AssetUploadHandler)
	authNoAuthGroup.GET("/search/options", handlers.GetSearchOptionsHandler)
	// authGroup.GET("/history/options/:search_id", handlers.GetSearchHistoryOptionsHandler)
	// authGroup.GET("/history", handlers.GetSearchHistoryHandler)

	privateGroup.POST("/templates/email", handlers.CreateEmailTemplateHandler)
}

func startServer() {
	logrus.Info("Starting the server on :1729")
	err := ginEngine.Run(":1729")
	fmt.Println(err)
}
