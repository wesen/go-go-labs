package main

import (
	"time"
)

// StreamInfo represents stream metadata
type StreamInfo struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	StartTime   time.Time `json:"startTime" db:"start_time"`
	Language    string    `json:"language" db:"language"`
	GithubRepo  string    `json:"githubRepo" db:"github_repo"`
	ViewerCount int       `json:"viewerCount" db:"viewer_count"`
}

// Step represents a single task step
type Step struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// StepInfo represents all task steps
type StepInfo struct {
	ID        string `json:"id" db:"id"` 
	Completed []Step `json:"completed"`
	Active    *Step  `json:"active"`
	Upcoming  []Step `json:"upcoming"`
}

// Stream represents the complete stream state
type Stream struct {
	Info  StreamInfo `json:"info"`
	Steps StepInfo   `json:"steps"`
}