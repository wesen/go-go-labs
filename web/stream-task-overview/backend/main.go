package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Command-line flags
var (
	debugFlag = flag.Bool("debug", false, "Enable debug logging")
)

// Initialize logger
func initLogger() {
	// Pretty console logging for development
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	multi := zerolog.MultiLevelWriter(consoleWriter)

	// Set global logger
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()

	// Set log level based on flags and environment
	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled via command line flag")
	} else if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled via DEBUG environment variable")
	} else if os.Getenv("ENV") == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled via development environment")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msg("Logger initialized")
}

// Custom echo logger adapter for zerolog
type ZerologAdapter struct {
	log zerolog.Logger
}

func (zl ZerologAdapter) Write(p []byte) (n int, err error) {
	zl.log.Info().Msg(string(p))
	return len(p), nil
}

func (zl ZerologAdapter) Output() io.Writer {
	return zl
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Initialize logger
	initLogger()

	// Echo setup
	e := echo.New()
	e.Logger.SetOutput(ZerologAdapter{log: log.With().Str("component", "echo").Logger()})
	e.HideBanner = true

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Custom logger middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			log.Debug().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("id", reqID).
				Str("user_agent", req.UserAgent()).
				Str("remote_addr", req.RemoteAddr).
				Msg("Request received")

			err := next(c)
			if err != nil {
				log.Error().
					Err(err).
					Str("id", reqID).
					Str("path", req.URL.Path).
					Msg("Request error")
				c.Error(err)
			}

			latency := time.Since(start)
			log.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Dur("latency", latency).
				Str("id", reqID).
				Msg("Response")

			return err
		}
	})

	// Initialize GitHub OAuth
	log.Info().Msg("Initializing authentication")
	InitAuth(e)
	auth := NewAuthHandler()

	// Setup store and handlers
	log.Info().Msg("Initializing data store")
	store := NewStreamStore()

	log.Info().Msg("Creating request handlers")
	h := NewStreamHandler(store)

	// Setup GitHub poller
	log.Info().Msg("Setting up GitHub poller")
	poller := NewGithubPoller(store)
	poller.Start()
	defer poller.Stop()

	// Auth routes
	log.Info().Msg("Setting up authentication routes")
	e.GET("/auth/github", auth.GithubAuthBegin)
	e.GET("/auth/github/callback", auth.GithubAuthCallback)
	e.GET("/auth/user", auth.GetCurrentUser)
	e.GET("/auth/logout", auth.Logout)

	// API routes
	log.Info().Msg("Setting up API routes")
	api := e.Group("/api")

	// Create GitHub handler
	log.Info().Msg("Creating GitHub handler")
	githubHandler := NewGitHubHandler(store)

	// Public routes
	api.GET("/stream", h.GetStreamInfo)
	api.GET("/stream/steps", h.GetSteps)
	api.GET("/stream/transcript", h.GetTranscript)
	api.GET("/stream/transcript/types/:type", h.GetTranscriptByType)
	api.GET("/github/info", githubHandler.GetGitHubInfo)
	api.GET("/github/commits", githubHandler.GetGitHubCommits)

	// Protected routes - require admin authentication
	adminGroup := api.Group("", auth.RequireAuth, auth.AdminOnly)
	adminGroup.PUT("/stream", h.UpdateStreamInfo)
	adminGroup.PUT("/stream/steps/active", h.SetActiveStep)
	adminGroup.POST("/stream/steps/upcoming", h.AddUpcomingStep)
	adminGroup.POST("/stream/steps/complete", h.CompleteActiveStep)
	adminGroup.PUT("/stream/steps/reactivate", h.ReactivateStep)
	adminGroup.POST("/stream/transcript", h.AddTranscriptEntry)
	adminGroup.POST("/github/connect", githubHandler.ConnectGitHub)
	adminGroup.POST("/github/webhook", githubHandler.HandleGitHubWebhook)

	// Serve frontend files
	e.Static("/", "../ui/dist")

	// Setup graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		s := <-sig
		log.Info().Str("signal", s.String()).Msg("Shutting down server...")
		poller.Stop()
		e.Close()
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info().Str("address", ":"+port).Msg("Starting server")
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failed")
	}
}
