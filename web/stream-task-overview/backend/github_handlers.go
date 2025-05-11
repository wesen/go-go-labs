package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// GitHubHandler handles GitHub-related requests
type GitHubHandler struct {
	store *StreamStore
}

// NewGitHubHandler creates a new GitHub handler
func NewGitHubHandler(store *StreamStore) *GitHubHandler {
	log.Debug().Msg("Creating new GitHubHandler")
	return &GitHubHandler{store: store}
}

// Helper function to check errors and log them
func logError(err error, msg string) {
	if err != nil {
		log.Error().Err(err).Msg(msg)
	}
}

// ConnectGitHub connects to a GitHub repository
func (h *GitHubHandler) ConnectGitHub(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Handling ConnectGitHub request")
	
	var data struct {
		Token     string `json:"token"`
		RepoOwner string `json:"repoOwner"`
		RepoName  string `json:"repoName"`
	}
	
	if err := c.Bind(&data); err != nil {
		log.Error().Err(err).Msg("Failed to parse ConnectGitHub request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}
	
	// Validate required fields
	if data.Token == "" {
		log.Error().Msg("Empty token received")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token cannot be empty"})
	}
	
	if data.RepoOwner == "" {
		log.Error().Msg("Empty repo owner received")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Repository owner cannot be empty"})
	}
	
	if data.RepoName == "" {
		log.Error().Msg("Empty repo name received")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Repository name cannot be empty"})
	}
	
	// Create GitHub info object
	info := GitHubInfo{
		Token:         data.Token,
		RepoOwner:     data.RepoOwner,
		RepoName:      data.RepoName,
		CurrentBranch: "main", // Default to main branch
	}
	
	// Connect to GitHub repository
	err := h.store.ConnectGitHub(info)
	if err != nil {
		log.Error().Err(err).Str("repo", data.RepoOwner+"/"+data.RepoName).Msg("Failed to connect to GitHub")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to GitHub"})
	}
	
	log.Info().Str("repo", data.RepoOwner+"/"+data.RepoName).Msg("Connected to GitHub successfully")
	return c.JSON(http.StatusOK, map[string]string{"message": "Connected to GitHub successfully"})
}

// GetGitHubInfo returns GitHub repository information
func (h *GitHubHandler) GetGitHubInfo(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Handling GetGitHubInfo request")
	
	// Get GitHub info from database
	info, err := h.store.GetGitHubInfo()
	if err != nil {
		// If not found, return empty response with isConnected=false
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Debug().Msg("No GitHub connection found")
			return c.JSON(http.StatusOK, map[string]interface{}{
				"isConnected": false,
			})
		}
		
		log.Error().Err(err).Msg("Failed to get GitHub info")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get GitHub info"})
	}
	
	// Return info with isConnected=true
	response := map[string]interface{}{
		"isConnected":   true,
		"repoUrl":       "https://github.com/" + info.RepoOwner + "/" + info.RepoName,
		"repoOwner":     info.RepoOwner,
		"repoName":      info.RepoName,
		"currentBranch": info.CurrentBranch,
	}
	
	// Add latest commit if available
	if info.LatestCommitHash != "" {
		response["latestCommit"] = map[string]interface{}{
			"hash":    info.LatestCommitHash,
			"message": info.LatestCommitMsg,
			"author":  info.LatestCommitAuthor,
			"date":    info.LatestCommitDate,
			"url":     info.LatestCommitURL,
		}
	}
	
	log.Debug().Interface("response", response).Msg("Returning GitHub info")
	return c.JSON(http.StatusOK, response)
}

// GetGitHubCommits returns recent commits
func (h *GitHubHandler) GetGitHubCommits(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Handling GetGitHubCommits request")
	
	// Get limit parameter, default to 10
	limit := 10
	limitParam := c.QueryParam("limit")
	if limitParam != "" {
		limitInt, err := strconv.Atoi(limitParam)
		if err == nil && limitInt > 0 {
			limit = limitInt
		}
	}
	
	// Get commits from database
	commits, err := h.store.GetCommits(limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get GitHub commits")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get GitHub commits"})
	}
	
	log.Debug().Int("count", len(commits)).Msg("Returning GitHub commits")
	return c.JSON(http.StatusOK, commits)
}

// HandleGitHubWebhook processes GitHub webhook events
func (h *GitHubHandler) HandleGitHubWebhook(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Handling GitHub webhook")
	
	// Parse webhook payload
	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to parse webhook payload")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}
	
	// Extract event type from header
	eventType := c.Request().Header.Get("X-GitHub-Event")
	log.Debug().Str("event_type", eventType).Msg("Received GitHub webhook event")
	
	// Process push event (commit)
	if eventType == "push" {
		// Extract commit information
		if commits, ok := payload["commits"].([]interface{}); ok && len(commits) > 0 {
			latestCommit := commits[0].(map[string]interface{})
			
			// Process the latest commit
			author := ""
			if authorInfo, ok := latestCommit["author"].(map[string]interface{}); ok {
				if name, ok := authorInfo["name"].(string); ok {
					author = name
				}
			}
			
			commitHash := ""
			if id, ok := latestCommit["id"].(string); ok {
				commitHash = id
			}
			
			commitMessage := ""
			if message, ok := latestCommit["message"].(string); ok {
				commitMessage = message
			}
			
			// Update commit information
			if commitHash != "" {
				log.Info().Str("hash", commitHash).Msg("Processing new commit from webhook")
				
				// Create a commit info object
				commit := CommitInfo{
					Hash:    commitHash,
					Message: commitMessage,
					Author:  author,
					Date:    time.Now(), // Use current time as webhook might not include timestamp
					URL:     "", // Would normally construct from repo URL and commit hash
				}
				
				// Update GitHub commit info
				_ = h.store.UpdateGitHubCommit(commit)
				
				// Create a transcript entry for this commit
				commitHashCopy := commitHash
				commitURLCopy := commit.URL
				transcriptEntry := TranscriptEntry{
					Type:       TranscriptTypeCommit,
					Content:    commitMessage,
					Timestamp:  time.Now(),
					CommitHash: &commitHashCopy,
					CommitURL:  &commitURLCopy,
				}
				
				_, _ = h.store.AddTranscriptEntry(transcriptEntry)
			}
		}
	}
	
	// Acknowledge webhook
	log.Info().Str("event_type", eventType).Msg("Processed GitHub webhook event")
	return c.JSON(http.StatusOK, map[string]string{"message": "Webhook received successfully"})
}