package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6" // Using gofakeit as it's more widely used than go-faker
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Define models similar to the main application

// TranscriptEntryType defines the type of transcript entry
type TranscriptEntryType string

const (
	TranscriptTypeTaskStarted   TranscriptEntryType = "task_started"
	TranscriptTypeTaskCompleted TranscriptEntryType = "task_completed"
	TranscriptTypeCommit        TranscriptEntryType = "commit"
	TranscriptTypeNote          TranscriptEntryType = "note"
	TranscriptTypeParagraph     TranscriptEntryType = "paragraph"
	TranscriptTypeTranscript    TranscriptEntryType = "transcript"
)

// Step represents a single task step
type Step struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

func main() {
	// Check command line args for database path
	dbPath := "../../stream.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	// Connect to the database
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Seed data
	seedStreamInfo(db)
	seedSteps(db)
	seedTranscriptEntries(db)
	seedGitHubInfo(db)

	fmt.Println("Database seeded successfully!")
}

func seedStreamInfo(db *sqlx.DB) {
	// Check if data already exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM stream_info")
	if err != nil {
		log.Fatalf("Failed to check stream_info data: %v", err)
	}

	// Only seed if table is empty
	if count == 0 {
		streamID := uuid.NewString()
		_, err = db.Exec(`
			INSERT INTO stream_info (
				id, title, description, start_time, language, github_repo, viewer_count
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`, 
			streamID,
			gofakeit.ProductName(),
			gofakeit.Sentence(10),
			time.Now().Add(-3 * time.Hour),
			"TypeScript/React",
			"https://github.com/" + gofakeit.Username() + "/" + gofakeit.AppName(),
			gofakeit.Number(10, 200),
		)
		if err != nil {
			log.Fatalf("Failed to seed stream_info: %v", err)
		}
		fmt.Println("Seeded stream_info table")
	} else {
		fmt.Println("stream_info table already has data, skipping")
	}
}

func seedSteps(db *sqlx.DB) {
	// Check if data already exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM steps")
	if err != nil {
		log.Fatalf("Failed to check steps data: %v", err)
	}

	// Only seed if table is empty
	if count == 0 {
		stepsID := uuid.NewString()

		// Create completed steps with IDs
		completedSteps := []Step{}
		for i := 0; i < 3; i++ {
			step := Step{
				ID:          uuid.NewString(),
				Description: fmt.Sprintf("Complete %s", gofakeit.JobTitle()),
				CreatedAt:   time.Now().Add(time.Duration(-24*(3-i)) * time.Hour),
			}
			completedSteps = append(completedSteps, step)
		}

		// Create active step
		activeStep := Step{
			ID:          uuid.NewString(),
			Description: fmt.Sprintf("Implement %s component", gofakeit.ProgrammingLanguage()),
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		}

		// Create upcoming steps
		upcomingSteps := []Step{}
		for i := 0; i < 5; i++ {
			step := Step{
				ID:          uuid.NewString(),
				Description: fmt.Sprintf("Build %s for %s", gofakeit.NounAbstract(), gofakeit.AppName()),
				CreatedAt:   time.Now(),
			}
			upcomingSteps = append(upcomingSteps, step)
		}

		// Marshal to JSON
		completedJSON, err := json.Marshal(completedSteps)
		if err != nil {
			log.Fatalf("Failed to marshal completed steps: %v", err)
		}

		activeJSON, err := json.Marshal(activeStep)
		if err != nil {
			log.Fatalf("Failed to marshal active step: %v", err)
		}

		upcomingJSON, err := json.Marshal(upcomingSteps)
		if err != nil {
			log.Fatalf("Failed to marshal upcoming steps: %v", err)
		}

		// Insert into database
		_, err = db.Exec(`
			INSERT INTO steps (id, completed, active, upcoming) VALUES (?, ?, ?, ?)
		`, 
			stepsID, 
			string(completedJSON), 
			string(activeJSON), 
			string(upcomingJSON),
		)
		if err != nil {
			log.Fatalf("Failed to seed steps: %v", err)
		}
		fmt.Println("Seeded steps table")
	} else {
		fmt.Println("steps table already has data, skipping")
	}
}

func seedTranscriptEntries(db *sqlx.DB) {
	// Check if data already exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM transcript_entries")
	if err != nil {
		log.Fatalf("Failed to check transcript_entries data: %v", err)
	}

	// Only seed if table is empty
	if count == 0 {
		// Create 20 random transcript entries
		for i := 0; i < 20; i++ {
			id := uuid.NewString()
			timestamp := time.Now().Add(-time.Duration(gofakeit.Number(1, 180)) * time.Minute)
			
			// Randomize the entry type
			var entryType TranscriptEntryType
			switch gofakeit.Number(1, 6) {
			case 1:
				entryType = TranscriptTypeTaskStarted
			case 2:
				entryType = TranscriptTypeTaskCompleted
			case 3:
				entryType = TranscriptTypeCommit
			case 4:
				entryType = TranscriptTypeNote
			case 5:
				entryType = TranscriptTypeParagraph
			case 6:
				entryType = TranscriptTypeTranscript
			}
			
			// Generate the content based on type
			content := ""
			switch entryType {
			case TranscriptTypeTaskStarted:
				content = fmt.Sprintf("Started working on %s", gofakeit.JobTitle())
			case TranscriptTypeTaskCompleted:
				content = fmt.Sprintf("Completed %s implementation", gofakeit.ProgrammingLanguage())
			case TranscriptTypeCommit:
				content = fmt.Sprintf("%s: %s", gofakeit.HackerVerb(), gofakeit.HackerNoun())
			case TranscriptTypeNote:
				content = gofakeit.Sentence(10)
			case TranscriptTypeParagraph:
				content = gofakeit.Paragraph(3, 5, 10, " ")
			case TranscriptTypeTranscript:
				content = gofakeit.Paragraph(1, 2, 10, " ")
			}
			
			// Prepare the query
			query := `
				INSERT INTO transcript_entries (
					id, timestamp, type, content
				`
			values := []interface{}{id, timestamp, entryType, content}
			
			// Add optional fields based on type
			if entryType == TranscriptTypeTaskStarted || entryType == TranscriptTypeTaskCompleted {
				taskName := gofakeit.JobTitle()
				query += `, task_name`
				values = append(values, taskName)
			}
			
			if entryType == TranscriptTypeCommit {
				commitHash := gofakeit.UUID()
				commitURL := fmt.Sprintf("https://github.com/%s/%s/commit/%s", 
					gofakeit.Username(), gofakeit.AppName(), commitHash)
				query += `, commit_hash, commit_url`
				values = append(values, commitHash, commitURL)
			}
			
			if entryType == TranscriptTypeParagraph || entryType == TranscriptTypeTranscript {
				title := gofakeit.ProductName()
				timeRangeStart := timestamp.Add(-time.Duration(gofakeit.Number(1, 10)) * time.Minute)
				timeRangeEnd := timestamp
				query += `, title, time_range_start, time_range_end`
				values = append(values, title, timeRangeStart, timeRangeEnd)
			}
			
			if entryType == TranscriptTypeTranscript {
				speaker := gofakeit.Name()
				duration := gofakeit.Number(30, 300)
				query += `, speaker, duration`
				values = append(values, speaker, duration)
			}
			
			// Finish the query
			query += `) VALUES (?` + strings.Repeat(", ?", len(values)-1) + `)`
			
			// Execute the query
			_, err = db.Exec(query, values...)
			if err != nil {
				log.Fatalf("Failed to seed transcript_entry %d: %v", i, err)
			}
		}
		fmt.Println("Seeded transcript_entries table with 20 entries")
	} else {
		fmt.Println("transcript_entries table already has data, skipping")
	}
}

func seedGitHubInfo(db *sqlx.DB) {
	// Check if data already exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM github_integration")
	if err != nil {
		log.Fatalf("Failed to check github_integration data: %v", err)
	}

	// Only seed if table is empty
	if count == 0 {
		// Generate GitHub info
		owner := gofakeit.Username()
		repo := gofakeit.AppName()
		branch := "main"
		commitHash := gofakeit.UUID()
		commitMsg := fmt.Sprintf("%s: %s %s", 
			gofakeit.HackerVerb(), 
			gofakeit.HackerAdjective(), 
			gofakeit.HackerNoun())
		commitAuthor := gofakeit.Name()
		commitDate := time.Now().Add(-time.Duration(gofakeit.Number(1, 24)) * time.Hour)
		commitURL := fmt.Sprintf("https://github.com/%s/%s/commit/%s", owner, repo, commitHash)
		
		// Insert into database
		_, err = db.Exec(`
			INSERT INTO github_integration (
				token, repo_owner, repo_name, current_branch, 
				latest_commit_hash, latest_commit_message, latest_commit_author, 
				latest_commit_date, latest_commit_url
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, 
			"fake_github_token", owner, repo, branch,
			commitHash, commitMsg, commitAuthor, commitDate, commitURL,
		)
		if err != nil {
			log.Fatalf("Failed to seed github_integration: %v", err)
		}
		fmt.Println("Seeded github_integration table")
	} else {
		fmt.Println("github_integration table already has data, skipping")
	}
}