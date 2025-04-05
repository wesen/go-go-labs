package handlers

import (
	"net/http"
	"strings"

	"github.com/go-go-golems/go-go-labs/cmd/apps/friday-talks/internal/auth"
	"github.com/go-go-golems/go-go-labs/cmd/apps/friday-talks/internal/models"
	"github.com/go-go-golems/go-go-labs/cmd/apps/friday-talks/internal/templates"
)

// AuthHandler handles authentication related routes
type AuthHandler struct {
	userRepo models.UserRepository
	auth     *auth.Auth
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userRepo models.UserRepository, auth *auth.Auth) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		auth:     auth,
	}
}

// HandleLogin handles the login page and form submission
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// If user is already logged in, redirect to home
	if user := auth.UserFromContext(r.Context()); user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"

	// Function to render the form (potentially partial)
	renderForm := func(errorMsg string) {
		if isHTMX {
			w.Header().Set("Content-Type", "text/html")
			templates._loginForm(errorMsg).Render(r.Context(), w)
		} else {
			templates.Login(errorMsg).Render(r.Context(), w)
		}
	}

	// Process login form submission
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate input
		if email == "" || password == "" {
			renderForm("Email and password are required")
			return
		}

		// Authenticate user
		user, err := h.auth.Authenticate(r.Context(), email, password)
		if err != nil {
			renderForm("Invalid email or password")
			return
		}

		// Generate token and set cookie
		token, err := h.auth.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		h.auth.SetTokenCookie(w, token)

		if isHTMX {
			// Redirect to home page via HTMX
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
		} else {
			// Standard redirect
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	// Render initial login page (full page)
	renderForm("")
}

// HandleRegister handles the registration page and form submission
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// If user is already logged in, redirect to home
	if user := auth.UserFromContext(r.Context()); user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"

	// Function to render the form (potentially partial)
	renderForm := func(errorMsg string) {
		if isHTMX {
			w.Header().Set("Content-Type", "text/html")
			templates._registerForm(errorMsg).Render(r.Context(), w)
		} else {
			templates.Register(errorMsg).Render(r.Context(), w)
		}
	}

	// Process registration form submission
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Validate input
		if name == "" || email == "" || password == "" {
			renderForm("All fields are required")
			return
		}
		if password != confirmPassword {
			renderForm("Passwords do not match")
			return
		}
		if len(password) < 8 {
			renderForm("Password must be at least 8 characters")
			return
		}

		// Check if email is already in use
		_, err := h.userRepo.FindByEmail(r.Context(), email)
		if err == nil {
			renderForm("Email is already in use")
			return
		}

		// Create new user
		hashedPassword, err := models.HashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		newUser := &models.User{
			Name:         name,
			Email:        email,
			PasswordHash: hashedPassword,
		}

		if err := h.userRepo.Create(r.Context(), newUser); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Generate token and set cookie
		token, err := h.auth.GenerateToken(newUser.ID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		h.auth.SetTokenCookie(w, token)

		if isHTMX {
			// Redirect to home page via HTMX
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
		} else {
			// Standard redirect
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	// Render initial registration page (full page)
	renderForm("")
}

// HandleLogout handles user logout
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear authentication cookie
	h.auth.ClearTokenCookie(w)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleProfile handles the user profile page and updates
func (h *AuthHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"

	// Function to render the form (potentially partial)
	renderForm := func(successMsg, errorMsg string) {
		// Re-fetch user in case details changed
		updatedUser, err := h.userRepo.FindByID(r.Context(), user.ID)
		if err != nil {
			// Handle error fetching updated user (should ideally not happen)
			updatedUser = user // Fallback to potentially stale user data
		}
		if isHTMX {
			w.Header().Set("Content-Type", "text/html")
			templates._profileForm(updatedUser, successMsg, errorMsg).Render(r.Context(), w)
		} else {
			templates.Profile(updatedUser, successMsg, errorMsg).Render(r.Context(), w)
		}
	}

	// Process profile update form submission
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		// Check if email is already in use by another user
		if email != user.Email {
			existingUser, err := h.userRepo.FindByEmail(r.Context(), email)
			if err == nil && existingUser.ID != user.ID {
				renderForm("", "Email is already in use by another account")
				return
			}
		}

		// Update user information
		user.Name = name
		user.Email = email

		// Update password if provided
		if currentPassword != "" && newPassword != "" {
			if !models.CheckPassword(currentPassword, user.PasswordHash) {
				renderForm("", "Current password is incorrect")
				return
			}

			if newPassword != confirmPassword {
				renderForm("", "New passwords do not match")
				return
			}

			if len(newPassword) < 8 {
				renderForm("", "New password must be at least 8 characters")
				return
			}

			hashedPassword, err := models.HashPassword(newPassword)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}

			user.PasswordHash = hashedPassword
		}

		// Save user updates
		if err := h.userRepo.Update(r.Context(), user); err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		// Render the form partial again with a success message
		renderForm("Profile updated successfully", "")
		return
	}

	// Render initial profile page (full page)
	renderForm("", "")
}

// ValidateEmail checks if the email is in a valid format
func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
