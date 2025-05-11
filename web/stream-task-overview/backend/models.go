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

// TranscriptEntryType defines the type of transcript entry
type TranscriptEntryType string

const (
	TranscriptTypeTaskStarted  TranscriptEntryType = "task_started"
	TranscriptTypeTaskCompleted TranscriptEntryType = "task_completed"
	TranscriptTypeCommit        TranscriptEntryType = "commit"
	TranscriptTypeNote          TranscriptEntryType = "note"
	TranscriptTypeParagraph     TranscriptEntryType = "paragraph"
	TranscriptTypeTranscript    TranscriptEntryType = "transcript"
)

// TimeRange represents a start and end time
type TimeRange struct {
	Start time.Time `json:"start" db:"time_range_start"`
	End   time.Time `json:"end" db:"time_range_end"`
}

// TranscriptEntry represents an entry in the stream transcript
type TranscriptEntry struct {
	ID        string             `json:"id" db:"id"`
	Timestamp time.Time          `json:"timestamp" db:"timestamp"`
	Type      TranscriptEntryType `json:"type" db:"type"`
	Content   string             `json:"content" db:"content"`
	
	// Optional fields based on type
	TaskName    *string    `json:"taskName,omitempty" db:"task_name"`
	CommitHash  *string    `json:"commitHash,omitempty" db:"commit_hash"`
	CommitURL   *string    `json:"commitUrl,omitempty" db:"commit_url"`
	TimeRange   *TimeRange `json:"timeRange,omitempty"` // Not directly stored in DB
	Title       *string    `json:"title,omitempty" db:"title"`
	Speaker     *string    `json:"speaker,omitempty" db:"speaker"`
	Duration    *int       `json:"duration,omitempty" db:"duration"`
}

// GitHubInfo represents GitHub repository information
type GitHubInfo struct {
	ID              int       `json:"id" db:"id"`
	Token           string    `json:"-" db:"token"` // Not exposed in JSON
	RepoOwner       string    `json:"repoOwner" db:"repo_owner"`
	RepoName        string    `json:"repoName" db:"repo_name"`
	CurrentBranch   string    `json:"currentBranch" db:"current_branch"`
	LatestCommitHash string    `json:"latestCommitHash,omitempty" db:"latest_commit_hash"`
	LatestCommitMsg  string    `json:"latestCommitMessage,omitempty" db:"latest_commit_message"`
	LatestCommitAuthor string   `json:"latestCommitAuthor,omitempty" db:"latest_commit_author"`
	LatestCommitDate time.Time `json:"latestCommitDate,omitempty" db:"latest_commit_date"`
	LatestCommitURL  string    `json:"latestCommitUrl,omitempty" db:"latest_commit_url"`
}

// CommitInfo represents a GitHub commit
type CommitInfo struct {
	Hash    string    `json:"hash"`
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
}