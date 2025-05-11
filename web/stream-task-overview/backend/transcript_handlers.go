package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// GetTranscript returns all transcript entries
func (h *StreamHandler) GetTranscript(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Handling GetTranscript request")
	
	entries, err := h.store.GetTranscriptEntries()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transcript entries")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve transcript entries"})
	}
	
	log.Debug().Int("count", len(entries)).Msg("Returning transcript entries")
	return c.JSON(http.StatusOK, entries)
}

// GetTranscriptByType returns transcript entries filtered by type
func (h *StreamHandler) GetTranscriptByType(c echo.Context) error {
	typeParam := c.Param("type")
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Str("type", typeParam).Msg("Handling GetTranscriptByType request")
	
	// Validate type
	var entryType TranscriptEntryType
	switch typeParam {
	case string(TranscriptTypeTaskStarted):
		entryType = TranscriptTypeTaskStarted
	case string(TranscriptTypeTaskCompleted):
		entryType = TranscriptTypeTaskCompleted
	case string(TranscriptTypeCommit):
		entryType = TranscriptTypeCommit
	case string(TranscriptTypeNote):
		entryType = TranscriptTypeNote
	case string(TranscriptTypeParagraph):
		entryType = TranscriptTypeParagraph
	case string(TranscriptTypeTranscript):
		entryType = TranscriptTypeTranscript
	default:
		log.Error().Str("type", typeParam).Msg("Invalid transcript entry type")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transcript entry type"})
	}
	
	entries, err := h.store.GetTranscriptEntriesByType(entryType)
	if err != nil {
		log.Error().Err(err).Str("type", typeParam).Msg("Failed to get transcript entries by type")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve transcript entries"})
	}
	
	log.Debug().Str("type", typeParam).Int("count", len(entries)).Msg("Returning transcript entries by type")
	return c.JSON(http.StatusOK, entries)
}

// AddTranscriptEntry adds a new transcript entry
func (h *StreamHandler) AddTranscriptEntry(c echo.Context) error {
	log.Debug().Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).Msg("Handling AddTranscriptEntry request")
	
	var entry TranscriptEntry
	if err := c.Bind(&entry); err != nil {
		log.Error().Err(err).Msg("Failed to parse AddTranscriptEntry request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}
	
	// Validate required fields
	if entry.Content == "" {
		log.Error().Msg("Empty content received")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Content cannot be empty"})
	}
	
	if entry.Type == "" {
		log.Error().Msg("Empty type received")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type cannot be empty"})
	}
	
	// Validate type
	switch entry.Type {
	case TranscriptTypeTaskStarted, 
	     TranscriptTypeTaskCompleted, 
	     TranscriptTypeCommit, 
	     TranscriptTypeNote, 
	     TranscriptTypeParagraph, 
	     TranscriptTypeTranscript:
		// Valid type
	default:
		log.Error().Str("type", string(entry.Type)).Msg("Invalid transcript entry type")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transcript entry type"})
	}
	
	// Add entry to database
	createdEntry, err := h.store.AddTranscriptEntry(entry)
	if err != nil {
		log.Error().Err(err).Interface("entry", entry).Msg("Failed to add transcript entry")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add transcript entry"})
	}
	
	log.Debug().Str("id", createdEntry.ID).Str("type", string(createdEntry.Type)).Msg("Transcript entry added successfully")
	return c.JSON(http.StatusCreated, createdEntry)
}