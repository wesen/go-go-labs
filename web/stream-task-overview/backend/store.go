package main

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// StreamStore manages stream data
type StreamStore struct {
	mutex sync.RWMutex
	db    *sqlx.DB
	sql   squirrel.StatementBuilderType
}

// NewStreamStore creates a new store with SQLite persistence
func NewStreamStore() *StreamStore {
	log.Debug().Msg("Creating new StreamStore")

	// Connect to SQLite database
	log.Debug().Str("db_path", "./stream.db").Msg("Connecting to SQLite database")
	db, err := sqlx.Connect("sqlite3", "./stream.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connection established")

	// Create SQL builder with SQLite placeholder
	sql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	// Create store
	store := &StreamStore{
		db:  db,
		sql: sql,
	}

	// Initialize database schema
	log.Debug().Msg("Initializing database schema")
	store.initSchema()

	// Create default data if none exists
	log.Debug().Msg("Initializing default data")
	store.initDefaultData()

	log.Info().Msg("StreamStore successfully initialized")
	return store
}

// initSchema creates database tables if they don't exist
func (s *StreamStore) initSchema() {
	// Create stream_info table
	log.Debug().Msg("Creating stream_info table if not exists")
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS stream_info (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		start_time DATETIME NOT NULL,
		language TEXT NOT NULL,
		github_repo TEXT NOT NULL,
		viewer_count INTEGER NOT NULL
	);
	`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create stream_info table")
	}

	// Create steps table
	log.Debug().Msg("Creating steps table if not exists")
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS steps (
		id TEXT PRIMARY KEY,
		completed TEXT NOT NULL,
		active TEXT NOT NULL,
		upcoming TEXT NOT NULL
	);
	`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create steps table")
	}

	// Create transcript_entries table
	log.Debug().Msg("Creating transcript_entries table if not exists")
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS transcript_entries (
		id TEXT PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		type TEXT NOT NULL,
		content TEXT NOT NULL,
		task_name TEXT,
		commit_hash TEXT,
		commit_url TEXT,
		time_range_start DATETIME,
		time_range_end DATETIME,
		title TEXT,
		speaker TEXT,
		duration INTEGER
	);
	`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create transcript_entries table")
	}

	// Create github_integration table
	log.Debug().Msg("Creating github_integration table if not exists")
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS github_integration (
		id INTEGER PRIMARY KEY,
		token TEXT NOT NULL,
		repo_owner TEXT NOT NULL,
		repo_name TEXT NOT NULL,
		current_branch TEXT NOT NULL,
		latest_commit_hash TEXT,
		latest_commit_message TEXT,
		latest_commit_author TEXT,
		latest_commit_date DATETIME,
		latest_commit_url TEXT
	);
	`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create github_integration table")
	}

	log.Info().Msg("Database schema setup complete")
}

// initDefaultData inserts default data if tables are empty
func (s *StreamStore) initDefaultData() {
	// Check if stream_info has data
	log.Debug().Msg("Checking if default data needs to be created")
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM stream_info")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to check stream_info data")
	}

	if count == 0 {
		log.Info().Msg("No existing data found, creating default data")

		// Generate unique IDs
		streamID := uuid.NewString()
		stepsID := uuid.NewString()

		// Insert default stream info
		log.Debug().Msg("Creating default stream info")
		query := s.sql.Insert("stream_info").Columns(
			"id", "title", "description", "start_time",
			"language", "github_repo", "viewer_count",
		).Values(
			streamID,
			"Building a React Component Library",
			"Creating reusable UI components with TailwindCSS",
			time.Now(),
			"JavaScript/React",
			"https://github.com/yourusername/component-library",
			42,
		)

		sql, args, err := query.ToSql()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to build SQL for stream info")
		}

		_, err = s.db.Exec(sql, args...)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to insert default stream info")
		}
		log.Debug().Msg("Default stream info created successfully")

		// Create default steps with unique IDs
		log.Debug().Msg("Creating default steps")

		// Create completed steps with IDs
		completedSteps := []Step{
			{ID: uuid.NewString(), Description: "Project setup and initialization", CreatedAt: time.Now().Add(-2 * time.Hour)},
			{ID: uuid.NewString(), Description: "Design system planning", CreatedAt: time.Now().Add(-1 * time.Hour)},
		}
		completedJSON, err := json.Marshal(completedSteps)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to marshal completed steps")
		}

		// Create active step with ID
		activeStep := Step{
			ID:          uuid.NewString(),
			Description: "Setting up component architecture",
			CreatedAt:   time.Now(),
		}
		activeJSON, err := json.Marshal(activeStep)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to marshal active step")
		}

		// Create upcoming steps with IDs
		upcomingSteps := []Step{
			{ID: uuid.NewString(), Description: "Implement Button component", CreatedAt: time.Now()},
			{ID: uuid.NewString(), Description: "Create Card component", CreatedAt: time.Now()},
			{ID: uuid.NewString(), Description: "Build Form elements", CreatedAt: time.Now()},
			{ID: uuid.NewString(), Description: "Add dark mode toggle", CreatedAt: time.Now()},
		}
		upcomingJSON, err := json.Marshal(upcomingSteps)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to marshal upcoming steps")
		}

		query = s.sql.Insert("steps").Columns(
			"id", "completed", "active", "upcoming",
		).Values(
			stepsID,
			string(completedJSON),
			string(activeJSON),
			string(upcomingJSON),
		)

		sql, args, err = query.ToSql()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to build SQL for steps")
		}

		_, err = s.db.Exec(sql, args...)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to insert default steps")
		}
		log.Debug().Msg("Default steps created successfully")
		log.Info().Msg("Default data creation complete")
	} else {
		log.Info().Msg("Existing data found, skipping default data creation")
	}
}

// GetStreamInfo returns the current stream info
func (s *StreamStore) GetStreamInfo() StreamInfo {
	log.Debug().Msg("Getting stream info")
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	query := s.sql.Select(
		"id", "title", "description", "start_time",
		"language", "github_repo", "viewer_count",
	).From("stream_info").Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build SQL for GetStreamInfo")
		return StreamInfo{}
	}

	var info StreamInfo
	err = s.db.Get(&info, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get stream info from database")
		return StreamInfo{}
	}

	log.Debug().Interface("info", info).Msg("Retrieved stream info")
	return info
}

// UpdateStreamInfo updates stream info
func (s *StreamStore) UpdateStreamInfo(info StreamInfo) {
	log.Debug().Interface("info", info).Msg("Updating stream info")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := s.sql.Update("stream_info").SetMap(map[string]interface{}{
		"title":        info.Title,
		"description":  info.Description,
		"start_time":   info.StartTime,
		"language":     info.Language,
		"github_repo":  info.GithubRepo,
		"viewer_count": info.ViewerCount,
	}).Where(squirrel.Eq{"id": 1})

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build SQL for UpdateStreamInfo")
		return
	}

	_, err = s.db.Exec(sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update stream info in database")
		return
	}

	log.Info().Msg("Stream info updated successfully")
}

// GetSteps returns all steps
func (s *StreamStore) GetSteps() StepInfo {
	log.Debug().Msg("Getting steps")
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Get the first steps record (there should only be one)
	query := s.sql.Select("id", "completed", "active", "upcoming").From("steps").Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build SQL for GetSteps")
		return StepInfo{}
	}

	var row struct {
		ID        string `db:"id"`
		Completed string `db:"completed"`
		Active    string `db:"active"`
		Upcoming  string `db:"upcoming"`
	}

	err = s.db.Get(&row, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get steps from database")
		return StepInfo{}
	}

	// Parse JSON arrays
	var steps StepInfo
	steps.ID = row.ID

	// Parse active step
	var activeStep Step
	err = json.Unmarshal([]byte(row.Active), &activeStep)
	if err != nil {
		log.Error().Err(err).Str("json", row.Active).Msg("Failed to unmarshal active step")
		steps.Active = nil
	} else {
		steps.Active = &activeStep
	}

	// Parse completed steps
	err = json.Unmarshal([]byte(row.Completed), &steps.Completed)
	if err != nil {
		log.Error().Err(err).Str("json", row.Completed).Msg("Failed to unmarshal completed steps")
		steps.Completed = []Step{}
	}

	// Parse upcoming steps
	err = json.Unmarshal([]byte(row.Upcoming), &steps.Upcoming)
	if err != nil {
		log.Error().Err(err).Str("json", row.Upcoming).Msg("Failed to unmarshal upcoming steps")
		steps.Upcoming = []Step{}
	}

	log.Debug().Interface("steps", steps).Msg("Retrieved steps")
	return steps
}

// updateSteps updates all steps in database
func (s *StreamStore) updateSteps(steps StepInfo) error {
	log.Debug().Interface("steps", steps).Msg("Updating steps in database")

	// Marshal arrays to JSON
	completed, err := json.Marshal(steps.Completed)
	if err != nil {
		log.Error().Err(err).Interface("completed", steps.Completed).Msg("Failed to marshal completed steps")
		return errors.Wrap(err, "marshal completed steps")
	}

	// Marshal active step to JSON
	var activeJSON []byte
	if steps.Active != nil {
		activeJSON, err = json.Marshal(steps.Active)
		if err != nil {
			log.Error().Err(err).Interface("active", steps.Active).Msg("Failed to marshal active step")
			return errors.Wrap(err, "marshal active step")
		}
	} else {
		// Empty object if no active step
		activeJSON = []byte("{}")
	}

	upcoming, err := json.Marshal(steps.Upcoming)
	if err != nil {
		log.Error().Err(err).Interface("upcoming", steps.Upcoming).Msg("Failed to marshal upcoming steps")
		return errors.Wrap(err, "marshal upcoming steps")
	}

	// Update database
	query := s.sql.Update("steps").SetMap(map[string]interface{}{
		"completed": string(completed),
		"active":    string(activeJSON),
		"upcoming":  string(upcoming),
	}).Where(squirrel.Eq{"id": steps.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build SQL for updateSteps")
		return errors.Wrap(err, "build SQL")
	}

	_, err = s.db.Exec(sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute SQL for updateSteps")
		return errors.Wrap(err, "execute SQL")
	}

	log.Info().Msg("Steps updated successfully")
	return nil
}

// SetActiveStep sets a new active step
func (s *StreamStore) SetActiveStep(stepText string) {
	log.Debug().Str("step", stepText).Msg("Setting active step")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get current steps
	steps := s.GetSteps()

	// Add current active to completed if it exists
	if steps.Active != nil {
		log.Debug().Str("previous_active", steps.Active.Description).Msg("Moving previous active step to completed")
		steps.Completed = append(steps.Completed, *steps.Active)
	}

	// Create new active step with unique ID
	newStep := &Step{
		ID:          uuid.NewString(),
		Description: stepText,
		CreatedAt:   time.Now(),
	}

	// Set new active step
	steps.Active = newStep

	// Update database
	err := s.updateSteps(steps)
	if err != nil {
		log.Error().Err(err).Str("step", stepText).Msg("Failed to set active step")
		return
	}

	log.Info().Str("step", stepText).Msg("Active step set successfully")
}

// AddUpcomingStep adds a new upcoming step
func (s *StreamStore) AddUpcomingStep(stepText string) {
	log.Debug().Str("step", stepText).Msg("Adding upcoming step")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get current steps
	steps := s.GetSteps()

	// Create new step with unique ID
	newStep := Step{
		ID:          uuid.NewString(),
		Description: stepText,
		CreatedAt:   time.Now(),
	}

	// Add new upcoming step
	steps.Upcoming = append(steps.Upcoming, newStep)

	// Update database
	err := s.updateSteps(steps)
	if err != nil {
		log.Error().Err(err).Str("step", stepText).Msg("Failed to add upcoming step")
		return
	}

	log.Info().Str("step", stepText).Msg("Upcoming step added successfully")
}

// CompleteActiveStep completes the current active step
func (s *StreamStore) CompleteActiveStep() {
	log.Debug().Msg("Completing active step")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get current steps
	steps := s.GetSteps()

	// Process only if there's an active step
	if steps.Active != nil {
		// Add to completed
		log.Debug().Str("active", steps.Active.Description).Msg("Moving active step to completed")
		steps.Completed = append(steps.Completed, *steps.Active)

		// Set next step as active if available
		if len(steps.Upcoming) > 0 {
			log.Debug().Str("next_step", steps.Upcoming[0].Description).Msg("Setting next step as active")
			steps.Active = &steps.Upcoming[0]
			steps.Upcoming = steps.Upcoming[1:]
		} else {
			log.Debug().Msg("No upcoming steps, setting active to nil")
			steps.Active = nil
		}

		// Update database
		err := s.updateSteps(steps)
		if err != nil {
			log.Error().Err(err).Msg("Failed to complete active step")
			return
		}

		log.Info().Msg("Active step completed successfully")
	} else {
		log.Warn().Msg("No active step to complete")
	}
}

// ReactivateStep moves a step from completed/upcoming to active
func (s *StreamStore) ReactivateStep(stepID string, source string) {
	log.Debug().Str("stepID", stepID).Str("source", source).Msg("Reactivating step")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get current steps
	steps := s.GetSteps()

	// Add current active to completed if it exists
	if steps.Active != nil {
		log.Debug().Str("previous_active", steps.Active.Description).Msg("Moving previous active step to completed")
		steps.Completed = append(steps.Completed, *steps.Active)
	}

	// Find and set step as active
	var foundStep *Step

	// Remove from source list
	if source == "upcoming" {
		for i, s := range steps.Upcoming {
			if s.ID == stepID {
				log.Debug().Int("index", i).Str("description", s.Description).Msg("Removing step from upcoming list")
				foundStep = &s
				steps.Upcoming = append(steps.Upcoming[:i], steps.Upcoming[i+1:]...)
				break
			}
		}
	} else if source == "completed" {
		for i, s := range steps.Completed {
			if s.ID == stepID {
				log.Debug().Int("index", i).Str("description", s.Description).Msg("Removing step from completed list")
				foundStep = &s
				steps.Completed = append(steps.Completed[:i], steps.Completed[i+1:]...)
				break
			}
		}
	}

	// Set step as active if found
	if foundStep != nil {
		steps.Active = foundStep
		log.Debug().Str("description", foundStep.Description).Msg("Setting as active step")
	} else {
		log.Warn().Str("stepID", stepID).Msg("Step not found in source list")
		return
	}

	// Update database
	err := s.updateSteps(steps)
	if err != nil {
		log.Error().Err(err).Str("stepID", stepID).Str("source", source).Msg("Failed to reactivate step")
		return
	}

	log.Info().Str("stepID", stepID).Str("description", steps.Active.Description).Str("source", source).Msg("Step reactivated successfully")
}

// GetTranscriptEntries returns all transcript entries
func (s *StreamStore) GetTranscriptEntries() ([]TranscriptEntry, error) {
	log.Debug().Msg("Getting transcript entries")
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Use explicit column selection instead of * to avoid mapping issues with nested structs
	query := s.sql.Select("id", "timestamp", "type", "content", "task_name",
		"commit_hash", "commit_url", "time_range_start", "time_range_end",
		"title", "speaker", "duration").From("transcript_entries").OrderBy("timestamp DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build SQL for GetTranscriptEntries")
		return nil, errors.Wrap(err, "build SQL")
	}

	// Query raw data first to avoid struct mapping issues
	rows, err := s.db.Queryx(sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transcript entries from database")
		return nil, errors.Wrap(err, "select from database")
	}
	defer rows.Close()

	// Manually construct entries to handle TimeRange properly
	var entries []TranscriptEntry
	for rows.Next() {
		var entry TranscriptEntry

		// Create a map to scan into
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan transcript entry row")
			continue
		}

		// Set fields from the map
		if v, ok := result["id"]; ok {
			entry.ID = string(v.([]byte))
		}
		if v, ok := result["timestamp"]; ok {
			ts, err := time.Parse(time.RFC3339, string(v.([]byte)))
			if err == nil {
				entry.Timestamp = ts
			}
		}
		if v, ok := result["type"]; ok {
			entry.Type = TranscriptEntryType(string(v.([]byte)))
		}
		if v, ok := result["content"]; ok {
			entry.Content = string(v.([]byte))
		}

		// Handle optional fields
		if v, ok := result["task_name"]; ok && v != nil {
			str := string(v.([]byte))
			entry.TaskName = &str
		}
		if v, ok := result["commit_hash"]; ok && v != nil {
			str := string(v.([]byte))
			entry.CommitHash = &str
		}
		if v, ok := result["commit_url"]; ok && v != nil {
			str := string(v.([]byte))
			entry.CommitURL = &str
		}
		if v, ok := result["title"]; ok && v != nil {
			str := string(v.([]byte))
			entry.Title = &str
		}
		if v, ok := result["speaker"]; ok && v != nil {
			str := string(v.([]byte))
			entry.Speaker = &str
		}
		if v, ok := result["duration"]; ok && v != nil {
			if dur, err := strconv.Atoi(string(v.([]byte))); err == nil {
				entry.Duration = &dur
			}
		}

		// Handle TimeRange
		var start, end time.Time
		if v, ok := result["time_range_start"]; ok && v != nil {
			ts, err := time.Parse(time.RFC3339, string(v.([]byte)))
			if err == nil {
				start = ts
			}
		}
		if v, ok := result["time_range_end"]; ok && v != nil {
			ts, err := time.Parse(time.RFC3339, string(v.([]byte)))
			if err == nil {
				end = ts
			}
		}

		if !start.IsZero() && !end.IsZero() {
			entry.TimeRange = &TimeRange{
				Start: start,
				End:   end,
			}
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating transcript entry rows")
		return nil, errors.Wrap(err, "iterate rows")
	}

	log.Debug().Int("count", len(entries)).Msg("Retrieved transcript entries")
	return entries, nil
}

// GetTranscriptEntriesByType returns transcript entries filtered by type
func (s *StreamStore) GetTranscriptEntriesByType(entryType TranscriptEntryType) ([]TranscriptEntry, error) {
	log.Debug().Str("type", string(entryType)).Msg("Getting transcript entries by type")
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Use explicit column selection instead of * to avoid mapping issues with nested structs
	query := s.sql.Select("id", "timestamp", "type", "content", "task_name",
		"commit_hash", "commit_url", "time_range_start", "time_range_end",
		"title", "speaker", "duration").From("transcript_entries").Where(squirrel.Eq{"type": entryType}).OrderBy("timestamp DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Str("type", string(entryType)).Msg("Failed to build SQL for GetTranscriptEntriesByType")
		return nil, errors.Wrap(err, "build SQL")
	}

	// Query raw data first to avoid struct mapping issues
	rows, err := s.db.Queryx(sql, args...)
	if err != nil {
		log.Error().Err(err).Str("type", string(entryType)).Msg("Failed to get transcript entries by type from database")
		return nil, errors.Wrap(err, "select from database")
	}
	defer rows.Close()

	// Manually construct entries to handle TimeRange properly
	var entries []TranscriptEntry
	for rows.Next() {
		var entry TranscriptEntry

		// Create a map to scan into
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan transcript entry row")
			continue
		}

		// Set fields from the map
		if v, ok := result["id"]; ok {
			entry.ID = string(v.([]byte))
		}
		if v, ok := result["timestamp"]; ok {
			ts, err := time.Parse(time.RFC3339, string(v.([]byte)))
			if err == nil {
				entry.Timestamp = ts
			}
		}
		if v, ok := result["type"]; ok {
			entry.Type = TranscriptEntryType(string(v.([]byte)))
		}
		if v, ok := result["content"]; ok {
			entry.Content = string(v.([]byte))
		}

		// Handle optional fields
		if v, ok := result["task_name"]; ok && v != nil {
			str := string(v.([]byte))
			entry.TaskName = &str
		}
		if v, ok := result["commit_hash"]; ok && v != nil {
			str := string(v.([]byte))
			entry.CommitHash = &str
		}
		if v, ok := result["commit_url"]; ok && v != nil {
			str := string(v.([]byte))
			entry.CommitURL = &str
		}
		if v, ok := result["title"]; ok && v != nil {
			str := string(v.([]byte))
			entry.Title = &str
		}
		if v, ok := result["speaker"]; ok && v != nil {
			str := string(v.([]byte))
			entry.Speaker = &str
		}
		if v, ok := result["duration"]; ok && v != nil {
			if dur, err := strconv.Atoi(string(v.([]byte))); err == nil {
				entry.Duration = &dur
			}
		}

		// Handle TimeRange
		var start, end time.Time
		if v, ok := result["time_range_start"]; ok && v != nil {
			ts, err := time.Parse(time.RFC3339, string(v.([]byte)))
			if err == nil {
				start = ts
			}
		}
		if v, ok := result["time_range_end"]; ok && v != nil {
			ts, err := time.Parse(time.RFC3339, string(v.([]byte)))
			if err == nil {
				end = ts
			}
		}

		if !start.IsZero() && !end.IsZero() {
			entry.TimeRange = &TimeRange{
				Start: start,
				End:   end,
			}
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating transcript entry rows")
		return nil, errors.Wrap(err, "iterate rows")
	}

	log.Debug().Str("type", string(entryType)).Int("count", len(entries)).Msg("Retrieved transcript entries by type")
	return entries, nil
}

// AddTranscriptEntry adds a new transcript entry
func (s *StreamStore) AddTranscriptEntry(entry TranscriptEntry) (TranscriptEntry, error) {
	log.Debug().Interface("entry", entry).Msg("Adding transcript entry")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate ID if not provided
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}

	// Ensure timestamp is set
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Prepare columns and values map
	columns := []string{"id", "timestamp", "type", "content"}
	values := []interface{}{entry.ID, entry.Timestamp, entry.Type, entry.Content}

	// Add optional fields if present
	if entry.TaskName != nil {
		columns = append(columns, "task_name")
		values = append(values, *entry.TaskName)
	}

	if entry.CommitHash != nil {
		columns = append(columns, "commit_hash")
		values = append(values, *entry.CommitHash)
	}

	if entry.CommitURL != nil {
		columns = append(columns, "commit_url")
		values = append(values, *entry.CommitURL)
	}

	if entry.TimeRange != nil {
		columns = append(columns, "time_range_start", "time_range_end")
		values = append(values, entry.TimeRange.Start, entry.TimeRange.End)
	}

	if entry.Title != nil {
		columns = append(columns, "title")
		values = append(values, *entry.Title)
	}

	if entry.Speaker != nil {
		columns = append(columns, "speaker")
		values = append(values, *entry.Speaker)
	}

	if entry.Duration != nil {
		columns = append(columns, "duration")
		values = append(values, *entry.Duration)
	}

	// Build and execute query
	query := s.sql.Insert("transcript_entries").Columns(columns...).Values(values...)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Interface("entry", entry).Msg("Failed to build SQL for AddTranscriptEntry")
		return TranscriptEntry{}, errors.Wrap(err, "build SQL")
	}

	_, err = s.db.Exec(sql, args...)
	if err != nil {
		log.Error().Err(err).Interface("entry", entry).Msg("Failed to insert transcript entry")
		return TranscriptEntry{}, errors.Wrap(err, "insert into database")
	}

	log.Info().Str("id", entry.ID).Str("type", string(entry.Type)).Msg("Transcript entry added successfully")
	return entry, nil
}

// GetGitHubInfo returns the GitHub repository information
func (s *StreamStore) GetGitHubInfo() (GitHubInfo, error) {
	log.Debug().Msg("Getting GitHub info")
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	query := s.sql.Select("*").From("github_integration").Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build SQL for GetGitHubInfo")
		return GitHubInfo{}, errors.Wrap(err, "build SQL")
	}

	var info GitHubInfo
	err = s.db.Get(&info, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get GitHub info from database")
		return GitHubInfo{}, errors.Wrap(err, "select from database")
	}

	log.Debug().Interface("info", info).Msg("Retrieved GitHub info")
	return info, nil
}

// ConnectGitHub connects to a GitHub repository
func (s *StreamStore) ConnectGitHub(info GitHubInfo) error {
	log.Debug().Interface("info", info).Msg("Connecting to GitHub")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if a connection already exists
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM github_integration")
	if err != nil {
		log.Error().Err(err).Msg("Failed to check GitHub connection")
		return errors.Wrap(err, "check connection")
	}

	if count > 0 {
		// Update existing connection
		log.Debug().Msg("Updating existing GitHub connection")
		query := s.sql.Update("github_integration").SetMap(map[string]interface{}{
			"token":          info.Token,
			"repo_owner":     info.RepoOwner,
			"repo_name":      info.RepoName,
			"current_branch": info.CurrentBranch,
		})

		sql, args, err := query.ToSql()
		if err != nil {
			log.Error().Err(err).Msg("Failed to build SQL for updating GitHub connection")
			return errors.Wrap(err, "build SQL")
		}

		_, err = s.db.Exec(sql, args...)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update GitHub connection")
			return errors.Wrap(err, "update connection")
		}
	} else {
		// Create new connection
		log.Debug().Msg("Creating new GitHub connection")
		query := s.sql.Insert("github_integration").Columns(
			"token", "repo_owner", "repo_name", "current_branch",
		).Values(
			info.Token, info.RepoOwner, info.RepoName, info.CurrentBranch,
		)

		sql, args, err := query.ToSql()
		if err != nil {
			log.Error().Err(err).Msg("Failed to build SQL for creating GitHub connection")
			return errors.Wrap(err, "build SQL")
		}

		_, err = s.db.Exec(sql, args...)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create GitHub connection")
			return errors.Wrap(err, "create connection")
		}
	}

	log.Info().Str("repo", info.RepoOwner+"/"+info.RepoName).Msg("GitHub connection saved successfully")
	return nil
}

// UpdateGitHubCommit updates the latest commit information
func (s *StreamStore) UpdateGitHubCommit(commit CommitInfo) error {
	log.Debug().Interface("commit", commit).Msg("Updating GitHub commit info")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if a connection exists
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM github_integration")
	if err != nil {
		log.Error().Err(err).Msg("Failed to check GitHub connection")
		return errors.Wrap(err, "check connection")
	}

	if count == 0 {
		log.Error().Msg("No GitHub connection exists")
		return errors.New("no GitHub connection")
	}

	// Update commit information
	query := s.sql.Update("github_integration").SetMap(map[string]interface{}{
		"latest_commit_hash":    commit.Hash,
		"latest_commit_message": commit.Message,
		"latest_commit_author":  commit.Author,
		"latest_commit_date":    commit.Date,
		"latest_commit_url":     commit.URL,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build SQL for updating GitHub commit")
		return errors.Wrap(err, "build SQL")
	}

	_, err = s.db.Exec(sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update GitHub commit info")
		return errors.Wrap(err, "update commit info")
	}

	log.Info().Str("hash", commit.Hash).Msg("GitHub commit info updated successfully")
	return nil
}

// GetCommits returns the latest commits (mock implementation for now)
func (s *StreamStore) GetCommits(limit int) ([]CommitInfo, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// First try to get GitHub info for repo details
	githubInfo, err := s.GetGitHubInfo()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get GitHub info")
		return nil, errors.Wrap(err, "get GitHub info")
	}

	// This would normally fetch from GitHub API
	// For now, just return the latest commit we have stored
	if githubInfo.LatestCommitHash != "" {
		log.Debug().Str("hash", githubInfo.LatestCommitHash).Msg("Returning stored commit")
		return []CommitInfo{{
			Hash:    githubInfo.LatestCommitHash,
			Message: githubInfo.LatestCommitMsg,
			Author:  githubInfo.LatestCommitAuthor,
			Date:    githubInfo.LatestCommitDate,
			URL:     githubInfo.LatestCommitURL,
		}}, nil
	}

	// If no commits found, return empty array
	log.Debug().Msg("No commits found")
	return []CommitInfo{}, nil
}
