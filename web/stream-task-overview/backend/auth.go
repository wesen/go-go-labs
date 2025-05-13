package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/agentstation/egothic"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/rs/zerolog/log"
)

// AuthHandler manages authentication
type AuthHandler struct {
	adminUsername string
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler() *AuthHandler {
	adminUsername := os.Getenv("ADMIN_GITHUB_USERNAME")
	if adminUsername == "" {
		log.Warn().Msg("ADMIN_GITHUB_USERNAME environment variable not set; admin functionality will be disabled")
	} else {
		log.Info().Str("admin_username", adminUsername).Msg("Admin GitHub username configured")
	}

	return &AuthHandler{
		adminUsername: adminUsername,
	}
}

// InitAuth initializes authentication
func InitAuth(e *echo.Echo) {
	// Set up session store
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Warn().Msg("SESSION_SECRET environment variable not set; using insecure default")
		sessionSecret = "stream-task-overview-default-secret"
	}

	store := sessions.NewCookieStore([]byte(sessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	e.Use(session.Middleware(store))
	log.Debug().Msg("Session middleware configured")

	// GitHub OAuth configuration
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	callbackURL := os.Getenv("GITHUB_CALLBACK_URL")

	log.Debug().
		Str("client_id_set", fmt.Sprintf("%t", clientID != "")).
		Str("client_secret_set", fmt.Sprintf("%t", clientSecret != "")).
		Str("callback_url_raw", callbackURL).
		Msg("GitHub OAuth configuration values")

	if callbackURL == "" {
		callbackURL = "http://localhost:8080/auth/github/callback"
		log.Warn().Str("callback_url", callbackURL).Msg("GITHUB_CALLBACK_URL not set, using default")
	}

	if clientID != "" && clientSecret != "" {
		log.Info().
			Str("client_id", clientID[:4]+"...").
			Str("callback_url", callbackURL).
			Msg("Initializing GitHub OAuth provider")

		provider := github.New(clientID, clientSecret, callbackURL, "read:user", "user:email")
		goth.UseProviders(provider)

		log.Debug().
			Str("provider_name", provider.Name()).
			Str("auth_url", provider.BeginAuthURL(nil)).
			Msg("GitHub provider initialized")
	} else {
		if clientID == "" {
			log.Error().Msg("GITHUB_CLIENT_ID environment variable not set")
		}
		if clientSecret == "" {
			log.Error().Msg("GITHUB_CLIENT_SECRET environment variable not set")
		}
		log.Warn().Msg("GitHub OAuth credentials incomplete; authentication will fail")
	}
}

// GithubAuthBegin initiates GitHub OAuth flow
func (h *AuthHandler) GithubAuthBegin(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	log.Debug().
		Str("request_id", reqID).
		Str("path", c.Request().URL.Path).
		Str("method", c.Request().Method).
		Str("user_agent", c.Request().UserAgent()).
		Msg("Beginning GitHub OAuth flow")

	// Check if providers are configured
	providers := goth.GetProviders()
	if len(providers) == 0 {
		log.Error().Str("request_id", reqID).Msg("No OAuth providers configured")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "OAuth not configured",
			"details": "GitHub OAuth credentials are not set. Check GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET environment variables.",
		})
	}

	// Check if GitHub provider exists
	if _, exists := providers["github"]; !exists {
		log.Error().Str("request_id", reqID).Msg("GitHub provider not found in configured providers")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":               "GitHub OAuth not configured",
			"details":             "GitHub OAuth provider is missing.",
			"available_providers": fmt.Sprintf("%v", getProviderNames(providers)),
		})
	}

	log.Debug().
		Str("request_id", reqID).
		Str("provider", "github").
		Msg("Calling BeginAuthHandler")

	// Attempt to begin auth and catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Error().
				Str("request_id", reqID).
				Interface("panic", r).
				Msg("Panic in BeginAuthHandler")
		}
	}()

	err := egothic.BeginAuthHandler(c)
	if err != nil {
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg("Failed to begin GitHub OAuth flow")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Authentication initialization failed",
			"details": err.Error(),
		})
	}

	log.Debug().
		Str("request_id", reqID).
		Msg("GitHub OAuth flow initiated successfully")
	return nil
}

// Helper function to get provider names
func getProviderNames(providers map[string]goth.Provider) []string {
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}

// GithubAuthCallback handles GitHub OAuth callback
func (h *AuthHandler) GithubAuthCallback(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	log.Debug().
		Str("request_id", reqID).
		Str("path", c.Request().URL.Path).
		Str("query", c.Request().URL.RawQuery).
		Msg("Handling GitHub OAuth callback")

	// Check error query param
	if errorMsg := c.QueryParam("error"); errorMsg != "" {
		description := c.QueryParam("error_description")
		log.Error().
			Str("request_id", reqID).
			Str("error", errorMsg).
			Str("description", description).
			Msg("GitHub OAuth error returned")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "GitHub OAuth error: " + errorMsg,
			"description": description,
		})
	}

	// Check code query param
	code := c.QueryParam("code")
	if code == "" {
		log.Error().
			Str("request_id", reqID).
			Msg("No code parameter in callback")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Invalid callback request",
			"details": "Missing 'code' parameter",
		})
	}

	log.Debug().
		Str("request_id", reqID).
		Str("code_length", fmt.Sprintf("%d", len(code))).
		Msg("Found authorization code in callback")

	// Complete user auth
	user, err := egothic.CompleteUserAuth(c)
	if err != nil {
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg("Failed to complete user authentication")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Authentication failed",
			"details": err.Error(),
		})
	}

	log.Debug().
		Str("request_id", reqID).
		Str("user_id", user.UserID).
		Str("nickname", user.NickName).
		Str("email", user.Email).
		Msg("User authenticated successfully")

	// Store user info in session
	sess, err := session.Get("stream-session", c)
	if err != nil {
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg("Failed to get session")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Session error",
			"details": err.Error(),
		})
	}

	sess.Values["user_id"] = user.UserID
	sess.Values["login"] = user.NickName
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["avatar"] = user.AvatarURL

	// Check if user is admin
	isAdmin := false
	if h.adminUsername != "" && strings.EqualFold(user.NickName, h.adminUsername) {
		log.Info().
			Str("request_id", reqID).
			Str("username", user.NickName).
			Msg("Admin user authenticated")
		sess.Values["is_admin"] = true
		isAdmin = true
	}

	log.Debug().
		Str("request_id", reqID).
		Str("user_id", user.UserID).
		Str("login", user.NickName).
		Bool("is_admin", isAdmin).
		Msg("Saving user session")

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg("Failed to save session")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to save session",
			"details": err.Error(),
		})
	}

	log.Info().
		Str("request_id", reqID).
		Str("user", user.NickName).
		Msg("User authenticated and session saved")

	// Redirect to frontend
	log.Debug().
		Str("request_id", reqID).
		Msg("Redirecting to frontend")
	return c.Redirect(http.StatusFound, "/")
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	log.Debug().
		Str("request_id", reqID).
		Str("path", c.Request().URL.Path).
		Msg("Getting current user")

	sess, err := session.Get("stream-session", c)
	if err != nil {
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg("Failed to get session")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"error":         "Session error",
		})
	}

	if sess.Values["user_id"] == nil {
		log.Debug().
			Str("request_id", reqID).
			Msg("No authenticated user found")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
	}

	user := map[string]interface{}{
		"authenticated": true,
		"id":            sess.Values["user_id"],
		"login":         sess.Values["login"],
		"name":          sess.Values["name"],
		"email":         sess.Values["email"],
		"avatar":        sess.Values["avatar"],
		"is_admin":      sess.Values["is_admin"] == true,
	}

	log.Debug().
		Str("request_id", reqID).
		Interface("user", user).
		Msg("Returning current user")
	return c.JSON(http.StatusOK, user)
}

// Logout logs out the current user
func (h *AuthHandler) Logout(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	log.Debug().
		Str("request_id", reqID).
		Str("path", c.Request().URL.Path).
		Msg("Logging out user")

	sess, err := session.Get("stream-session", c)
	if err != nil {
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg("Failed to get session during logout")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Session error",
			"details": err.Error(),
		})
	}

	// Log user being logged out
	if login, ok := sess.Values["login"].(string); ok {
		log.Info().
			Str("request_id", reqID).
			Str("user", login).
			Msg("Logging out user")
	}

	sess.Options.MaxAge = -1 // Expire the cookie
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg("Failed to save session during logout")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to logout",
			"details": err.Error(),
		})
	}

	log.Info().
		Str("request_id", reqID).
		Msg("User logged out successfully")
	return c.Redirect(http.StatusFound, "/")
}

// AdminOnly middleware checks if the user is an admin
func (h *AuthHandler) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := c.Response().Header().Get(echo.HeaderXRequestID)
		sess, err := session.Get("stream-session", c)
		if err != nil {
			log.Error().
				Str("request_id", reqID).
				Err(err).
				Str("path", c.Path()).
				Msg("Failed to get session in AdminOnly middleware")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "Session error",
				"details": err.Error(),
			})
		}

		if sess.Values["is_admin"] != true {
			log.Warn().
				Str("request_id", reqID).
				Str("path", c.Path()).
				Str("method", c.Request().Method).
				Interface("user_id", sess.Values["user_id"]).
				Interface("login", sess.Values["login"]).
				Msg("Unauthorized admin access attempt")
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
		}

		log.Debug().
			Str("request_id", reqID).
			Str("path", c.Path()).
			Interface("login", sess.Values["login"]).
			Msg("Admin access granted")
		return next(c)
	}
}

// RequireAuth middleware ensures a user is authenticated
func (h *AuthHandler) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := c.Response().Header().Get(echo.HeaderXRequestID)
		sess, err := session.Get("stream-session", c)
		if err != nil {
			log.Error().
				Str("request_id", reqID).
				Err(err).
				Str("path", c.Path()).
				Msg("Failed to get session in RequireAuth middleware")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "Session error",
				"details": err.Error(),
			})
		}

		if sess.Values["user_id"] == nil {
			log.Warn().
				Str("request_id", reqID).
				Str("path", c.Path()).
				Str("method", c.Request().Method).
				Msg("Unauthenticated access attempt")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		}

		log.Debug().
			Str("request_id", reqID).
			Str("path", c.Path()).
			Interface("user_id", sess.Values["user_id"]).
			Interface("login", sess.Values["login"]).
			Msg("Authenticated access")
		return next(c)
	}
}
