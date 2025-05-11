package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// GithubPoller fetches GitHub data periodically
type GithubPoller struct {
	store      *StreamStore
	client     *http.Client
	ticker     *time.Ticker
	stopCh     chan struct{}
	running    bool
	mutex      sync.Mutex
	viewerData map[string]int
}

// NewGithubPoller creates a new GitHub data poller
func NewGithubPoller(store *StreamStore) *GithubPoller {
	return &GithubPoller{
		store:      store,
		client:     &http.Client{Timeout: 10 * time.Second},
		stopCh:     make(chan struct{}),
		viewerData: make(map[string]int),
	}
}

// Start begins the polling process
func (p *GithubPoller) Start() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.running {
		log.Warn().Msg("GitHub poller is already running")
		return
	}

	// Check if GitHub token is available
	ghToken := os.Getenv("GITHUB_API_TOKEN")
	if ghToken == "" {
		log.Warn().Msg("GITHUB_API_TOKEN environment variable not set; GitHub polling disabled")
		return
	}

	// Default to polling every 2 minutes
	intervalStr := os.Getenv("GITHUB_POLLING_INTERVAL")
	interval := 2 * time.Minute
	if intervalStr != "" {
		if parsedInterval, err := time.ParseDuration(intervalStr); err == nil {
			interval = parsedInterval
		}
	}

	log.Info().Dur("interval", interval).Msg("Starting GitHub polling")
	p.ticker = time.NewTicker(interval)
	p.running = true

	// Do an initial poll immediately
	go func() {
		p.pollData()

		for {
			select {
			case <-p.ticker.C:
				p.pollData()
			case <-p.stopCh:
				p.ticker.Stop()
				log.Info().Msg("GitHub polling stopped")
				return
			}
		}
	}()
}

// Stop ends the polling process
func (p *GithubPoller) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.running {
		return
	}

	close(p.stopCh)
	p.running = false
}

// pollData fetches data from GitHub API
func (p *GithubPoller) pollData() {
	// Get current stream info to get the repo
	streamInfo := p.store.GetStreamInfo()
	if streamInfo.ID == "" {
		log.Debug().Msg("No stream info available, skipping GitHub poll")
		return
	}

	// Extract owner/repo from the GitHub URL
	repoPath := extractRepoPath(streamInfo.GithubRepo)
	if repoPath == "" {
		log.Warn().Str("url", streamInfo.GithubRepo).Msg("Could not parse GitHub repository URL")
		return
	}

	// Fetch viewer count (traffic API requires push access to repo)
	viewCount := p.fetchViewerCount(repoPath)

	// Update stream info with viewer count if available
	if viewCount > 0 {
		streamInfo.ViewerCount = viewCount
		p.store.UpdateStreamInfo(streamInfo)
		log.Info().Int("viewer_count", viewCount).Str("repo", repoPath).Msg("Updated viewer count from GitHub")
	}
}

// fetchViewerCount gets the view count from GitHub
func (p *GithubPoller) fetchViewerCount(repoPath string) int {
	token := os.Getenv("GITHUB_API_TOKEN")
	if token == "" {
		return 0
	}

	// GitHub traffic API endpoint
	url := fmt.Sprintf("https://api.github.com/repos/%s/traffic/views", repoPath)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request for GitHub API")
		return 0
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch data from GitHub API")
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Str("repo", repoPath).Msg("GitHub API returned non-200 status")
		return 0
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read GitHub API response")
		return 0
	}

	var data struct {
		Count int `json:"count"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Failed to parse GitHub API response")
		return 0
	}

	return data.Count
}

// extractRepoPath extracts the owner/repo part from a GitHub URL
func extractRepoPath(url string) string {
	// Handle formats like:
	// https://github.com/owner/repo
	// http://github.com/owner/repo
	// github.com/owner/repo
	
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "www.")
	
	if !strings.HasPrefix(url, "github.com/") {
		return ""
	}
	
	parts := strings.Split(strings.TrimPrefix(url, "github.com/"), "/")
	if len(parts) < 2 {
		return ""
	}
	
	// Get owner/repo path
	return parts[0] + "/" + parts[1]
}