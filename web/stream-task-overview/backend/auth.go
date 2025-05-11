package main

import (
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

	// GitHub OAuth configuration
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	callbackURL := os.Getenv("GITHUB_CALLBACK_URL")

	if callbackURL == "" {
		callbackURL = "http://localhost:8080/auth/github/callback"
		log.Warn().Str("callback_url", callbackURL).Msg("GITHUB_CALLBACK_URL not set, using default")
	}

	if clientID != "" && clientSecret != "" {
		log.Info().Msg("Initializing GitHub OAuth provider")
		goth.UseProviders(
			github.New(clientID, clientSecret, callbackURL, "read:user", "user:email"),
		)
	} else {
		log.Warn().Msg("GitHub OAuth credentials not provided; authentication will be disabled")
	}
}

// GithubAuthBegin initiates GitHub OAuth flow
func (h *AuthHandler) GithubAuthBegin(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Beginning GitHub OAuth flow")
	return egothic.BeginAuthHandler(c)
}

// GithubAuthCallback handles GitHub OAuth callback
func (h *AuthHandler) GithubAuthCallback(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Handling GitHub OAuth callback")

	user, err := egothic.CompleteUserAuth(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to complete user authentication")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Authentication failed"})
	}

	// Store user info in session
	sess, _ := session.Get("stream-session", c)
	sess.Values["user_id"] = user.UserID
	sess.Values["login"] = user.NickName
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["avatar"] = user.AvatarURL

	// Check if user is admin
	if h.adminUsername != "" && strings.EqualFold(user.NickName, h.adminUsername) {
		log.Info().Str("username", user.NickName).Msg("Admin user authenticated")
		sess.Values["is_admin"] = true
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		log.Error().Err(err).Msg("Failed to save session")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
	}

	// Redirect to frontend
	return c.Redirect(http.StatusFound, "/")
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Getting current user")

	sess, _ := session.Get("stream-session", c)

	if sess.Values["user_id"] == nil {
		log.Debug().Msg("No authenticated user found")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
	}

	user := map[string]interface{}{
		"authenticated": true,
		"id":           sess.Values["user_id"],
		"login":        sess.Values["login"],
		"name":         sess.Values["name"],
		"email":        sess.Values["email"],
		"avatar":       sess.Values["avatar"],
		"is_admin":     sess.Values["is_admin"] == true,
	}

	log.Debug().Interface("user", user).Msg("Returning current user")
	return c.JSON(http.StatusOK, user)
}

// Logout logs out the current user
func (h *AuthHandler) Logout(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Logging out user")

	sess, _ := session.Get("stream-session", c)
	sess.Options.MaxAge = -1 // Expire the cookie
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		log.Error().Err(err).Msg("Failed to save session during logout")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to logout"})
	}

	log.Info().Msg("User logged out successfully")
	return c.Redirect(http.StatusFound, "/")
}

// AdminOnly middleware checks if the user is an admin
func (h *AuthHandler) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("stream-session", c)
		if sess.Values["is_admin"] != true {
			log.Warn().Str("path", c.Path()).Str("method", c.Request().Method).Msg("Unauthorized admin access attempt")
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
		}
		return next(c)
	}
}

// RequireAuth middleware ensures a user is authenticated
func (h *AuthHandler) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("stream-session", c)
		if sess.Values["user_id"] == nil {
			log.Warn().Str("path", c.Path()).Str("method", c.Request().Method).Msg("Unauthenticated access attempt")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		}
		return next(c)
	}
}