package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-go-golems/go-go-labs/cmd/apps/friday-talks/internal/auth"
	"github.com/go-go-golems/go-go-labs/cmd/apps/friday-talks/internal/models"
	"github.com/go-go-golems/go-go-labs/internal/helpers"
)

// HandleVoteOnTalk handles voting on a talk
func (h *TalkHandler) HandleVoteOnTalk(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get talk ID from URL parameters
	talkIDStr := chi.URLParam(r, "id")
	talkID, err := strconv.Atoi(talkIDStr)
	if err != nil {
		http.Error(w, "Invalid talk ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Keep redirect for unauthenticated
		return
	}

	// Get talk from repository
	talk, err := h.talkRepo.FindByID(r.Context(), talkID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Reusable function to render the vote form partial
	renderVoteForm := func(errorMsg, successMsg string) {
		// Re-check if user has voted (state might have changed)
		_, err := h.voteRepo.FindByIDs(r.Context(), user.ID, talk.ID)
		voted := err == nil
		w.Header().Set("Content-Type", "text/html")
		templates._talkVoteForm(user, talk, voted, errorMsg, successMsg).Render(r.Context(), w)
	}

	// Check if talk is still in proposed state
	if talk.Status != models.TalkStatusProposed {
		renderVoteForm("Can only vote on proposed talks", "")
		return
	}

	// Check if user is not the speaker (can't vote on own talk)
	if talk.SpeakerID == user.ID {
		renderVoteForm("You cannot vote on your own talk", "")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get vote data
	interestLevelStr := r.FormValue("interest_level")
	interestLevel, err := strconv.Atoi(interestLevelStr)
	if err != nil || interestLevel < 1 || interestLevel > 5 {
		renderVoteForm("Invalid interest level", "")
		return
	}

	// Create availability map from form data
	availability := make(map[string]bool)
	for _, date := range talk.PreferredDates {
		value := r.FormValue("availability_" + date)
		availability[date] = (value == "true")
	}

	// Check if we're updating an existing vote
	existingVote, err := h.voteRepo.FindByIDs(r.Context(), user.ID, talkID)
	if err == nil {
		existingVote.InterestLevel = interestLevel
		existingVote.Availability = availability
		if err := h.voteRepo.Update(r.Context(), existingVote); err != nil {
			http.Error(w, "Failed to update vote", http.StatusInternalServerError)
			return
		}
	} else {
		vote := &models.Vote{
			UserID:        user.ID,
			TalkID:        talkID,
			InterestLevel: interestLevel,
			Availability:  availability,
		}
		if err := h.voteRepo.Create(r.Context(), vote); err != nil {
			http.Error(w, "Failed to create vote", http.StatusInternalServerError)
			return
		}
	}

	// Render the vote form partial with success message
	renderVoteForm("", "Vote submitted successfully")
}

// HandleManageAttendance handles managing attendance for a talk
func (h *TalkHandler) HandleManageAttendance(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get talk ID from URL parameters
	talkIDStr := chi.URLParam(r, "id")
	talkID, err := strconv.Atoi(talkIDStr)
	if err != nil {
		http.Error(w, "Invalid talk ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Keep redirect for unauthenticated
		return
	}

	// Get talk from repository
	talk, err := h.talkRepo.FindByID(r.Context(), talkID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Reusable function to render the attendance partial
	renderAttendance := func(errorMsg, successMsg string) {
		// Re-fetch attendance status
		attendance, _ := h.attendanceRepo.FindByIDs(r.Context(), talkID, user.ID)
		w.Header().Set("Content-Type", "text/html")
		templates._talkAttendance(user, talk, attendance, errorMsg, successMsg).Render(r.Context(), w)
	}

	// Check if talk is scheduled
	if talk.Status != models.TalkStatusScheduled {
		renderAttendance("Can only manage attendance for scheduled talks", "")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get status from form
	status := models.AttendanceStatus(r.FormValue("status"))
	if status != models.AttendanceStatusConfirmed && status != models.AttendanceStatusDeclined {
		renderAttendance("Invalid attendance status", "")
		return
	}

	// Check if we're updating an existing attendance record
	existingAttendance, err := h.attendanceRepo.FindByIDs(r.Context(), talkID, user.ID)
	if err == nil {
		existingAttendance.Status = status
		if err := h.attendanceRepo.Update(r.Context(), existingAttendance); err != nil {
			http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
			return
		}
	} else {
		attendance := &models.Attendance{
			TalkID: talkID,
			UserID: user.ID,
			Status: status,
		}
		if err := h.attendanceRepo.Create(r.Context(), attendance); err != nil {
			http.Error(w, "Failed to create attendance record", http.StatusInternalServerError)
			return
		}
	}

	// Render the attendance partial with success message
	renderAttendance("", "Attendance updated successfully")
}

// HandleProvideFeedback handles providing feedback for a talk
func (h *TalkHandler) HandleProvideFeedback(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get talk ID from URL parameters
	talkIDStr := chi.URLParam(r, "id")
	talkID, err := strconv.Atoi(talkIDStr)
	if err != nil {
		http.Error(w, "Invalid talk ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Keep redirect for unauthenticated
		return
	}

	// Get talk from repository
	talk, err := h.talkRepo.FindByID(r.Context(), talkID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Re-fetch attendance record for current state
	attendance, err := h.attendanceRepo.FindByIDs(r.Context(), talkID, user.ID)

	// Reusable function to render the feedback form partial
	renderFeedbackForm := func(errorMsg, successMsg string) {
		w.Header().Set("Content-Type", "text/html")
		templates._talkFeedbackForm(user, talk, attendance, errorMsg, successMsg).Render(r.Context(), w)
	}

	// Check if talk is completed
	if talk.Status != models.TalkStatusCompleted {
		renderFeedbackForm("Can only provide feedback for completed talks", "")
		return
	}

	// Check if user attended the talk (attendance might be nil if not found)
	if err != nil || attendance.Status != models.AttendanceStatusAttended {
		renderFeedbackForm("You must have attended the talk to provide feedback", "")
		return
	}

	// Check if feedback already exists
	if attendance.Feedback != "" {
		renderFeedbackForm("You have already provided feedback for this talk", "")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get feedback from form
	feedback := r.FormValue("feedback")
	if feedback == "" {
		renderFeedbackForm("Feedback cannot be empty", "")
		return
	}

	// Update attendance record with feedback
	attendance.Feedback = feedback
	if err := h.attendanceRepo.Update(r.Context(), attendance); err != nil {
		http.Error(w, "Failed to save feedback", http.StatusInternalServerError)
		return
	}

	// Render the feedback form partial with success message
	renderFeedbackForm("", "Feedback submitted successfully")
}

// HandleAddResource handles adding a resource to a talk
func (h *TalkHandler) HandleAddResource(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get talk ID from URL parameters
	talkIDStr := chi.URLParam(r, "id")
	talkID, err := strconv.Atoi(talkIDStr)
	if err != nil {
		http.Error(w, "Invalid talk ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get talk from repository
	talk, err := h.talkRepo.FindByID(r.Context(), talkID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Reusable function to render the resources partial
	renderResourcesPartial := func(errorMsg, successMsg string) {
		// Re-fetch resources
		resources, err := h.resourceRepo.FindByTalkID(r.Context(), talk.ID)
		if err != nil {
			http.Error(w, "Failed to fetch resources", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		templates._talkResourcesSection(user, talk, resources, errorMsg, successMsg).Render(r.Context(), w)
	}

	// Check if user is the speaker
	if talk.SpeakerID != user.ID {
		if helpers.IsHtmxRequest(r) {
			renderResourcesPartial("Only the speaker can add resources", "")
		} else {
			http.Redirect(w, r, "/talks/"+talkIDStr+"?error=Only the speaker can add resources to a talk", http.StatusSeeOther)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get resource data from form
	title := r.FormValue("title")
	url := r.FormValue("url")
	resourceType := models.ResourceType(r.FormValue("type"))

	// Validate input
	if title == "" || url == "" {
		if helpers.IsHtmxRequest(r) {
			renderResourcesPartial("Title and URL are required", "")
		} else {
			http.Redirect(w, r, "/talks/"+talkIDStr+"?error=Title and URL are required", http.StatusSeeOther)
		}
		return
	}

	// Create new resource
	resource := &models.Resource{
		TalkID: talkID,
		Title:  title,
		URL:    url,
		Type:   resourceType,
	}

	if err := h.resourceRepo.Create(r.Context(), resource); err != nil {
		http.Error(w, "Failed to create resource", http.StatusInternalServerError)
		return
	}

	if helpers.IsHtmxRequest(r) {
		renderResourcesPartial("", "Resource added successfully")
	} else {
		// Redirect back to talk page for non-HTMX requests
		http.Redirect(w, r, "/talks/"+talkIDStr+"?success=Resource added successfully", http.StatusSeeOther)
	}
}

// HandleDeleteResource handles deleting a resource from a talk
func (h *TalkHandler) HandleDeleteResource(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get talk and resource IDs from URL parameters
	talkIDStr := chi.URLParam(r, "id")
	talkID, err := strconv.Atoi(talkIDStr)
	if err != nil {
		http.Error(w, "Invalid talk ID", http.StatusBadRequest)
		return
	}

	resourceIDStr := chi.URLParam(r, "resourceId")
	resourceID, err := strconv.Atoi(resourceIDStr)
	if err != nil {
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get talk from repository
	talk, err := h.talkRepo.FindByID(r.Context(), talkID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Reusable function to render the resources partial
	renderResourcesPartial := func(errorMsg, successMsg string) {
		// Re-fetch resources
		resources, err := h.resourceRepo.FindByTalkID(r.Context(), talk.ID)
		if err != nil {
			http.Error(w, "Failed to fetch resources", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		templates._talkResourcesSection(user, talk, resources, errorMsg, successMsg).Render(r.Context(), w)
	}

	// Check if user is the speaker
	if talk.SpeakerID != user.ID {
		if helpers.IsHtmxRequest(r) {
			renderResourcesPartial("Only the speaker can delete resources", "")
		} else {
			http.Redirect(w, r, "/talks/"+talkIDStr+"?error=Only the speaker can delete resources from a talk", http.StatusSeeOther)
		}
		return
	}

	// Delete the resource
	if err := h.resourceRepo.Delete(r.Context(), resourceID); err != nil {
		http.Error(w, "Failed to delete resource", http.StatusInternalServerError)
		return
	}

	if helpers.IsHtmxRequest(r) {
		renderResourcesPartial("", "Resource deleted successfully")
	} else {
		// Redirect back to talk page for non-HTMX requests
		http.Redirect(w, r, "/talks/"+talkIDStr+"?success=Resource deleted successfully", http.StatusSeeOther)
	}
}
