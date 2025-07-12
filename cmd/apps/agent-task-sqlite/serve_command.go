package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

//go:embed templates static
var embedFS embed.FS

type ServeCommand struct {
	*cmds.CommandDescription
	db *sqlx.DB
}

type ServeSettings struct {
	Port int `glazed.parameter:"port"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &ServeCommand{}

type DashboardData struct {
	Projects    []ProjectView  `json:"projects"`
	Agents      []AgentView    `json:"agents"`
	Tasks       []TaskView     `json:"tasks"`
	RecentTasks []TaskView     `json:"recent_tasks"`
	Stats       DashboardStats `json:"stats"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ProjectView struct {
	ID                int       `json:"id" db:"id"`
	Slug              string    `json:"slug" db:"slug"`
	Name              string    `json:"name" db:"name"`
	Description       string    `json:"description" db:"description"`
	ConciseGuidelines *string   `json:"concise_guidelines" db:"concise_guidelines"`
	FullGuidelines    *string   `json:"full_guidelines" db:"full_guidelines"`
	TaskCount         int       `json:"task_count" db:"task_count"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type AgentView struct {
	ID                 int       `json:"id" db:"id"`
	Slug               string    `json:"slug" db:"slug"`
	Name               string    `json:"name" db:"name"`
	Description        string    `json:"description" db:"description"`
	CurrentProjectID   *int      `json:"current_project_id" db:"current_project_id"`
	CurrentProjectName *string   `json:"current_project_name" db:"current_project_name"`
	CurrentTaskID      *int      `json:"current_task_id" db:"current_task_id"`
	CurrentTaskSlug    *string   `json:"current_task_slug" db:"current_task_slug"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type TaskView struct {
	ID              int        `json:"id" db:"id"`
	Slug            string     `json:"slug" db:"slug"`
	ProjectID       int        `json:"project_id" db:"project_id"`
	ProjectName     string     `json:"project_name" db:"project_name"`
	ProjectSlug     string     `json:"project_slug" db:"project_slug"`
	AgentID         *int       `json:"agent_id" db:"agent_id"`
	AgentName       *string    `json:"agent_name" db:"agent_name"`
	Type            string     `json:"type" db:"type"`
	Status          string     `json:"status" db:"status"`
	Instructions    string     `json:"instructions" db:"instructions"`
	CompletionNotes *string    `json:"completion_notes" db:"completion_notes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	StartedAt       *time.Time `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
}

type DashboardStats struct {
	TotalProjects   int `json:"total_projects"`
	TotalAgents     int `json:"total_agents"`
	TotalTasks      int `json:"total_tasks"`
	PendingTasks    int `json:"pending_tasks"`
	InProgressTasks int `json:"in_progress_tasks"`
	CompletedTasks  int `json:"completed_tasks"`
	FailedTasks     int `json:"failed_tasks"`
}

// New types for extended functionality
type TaskDetailData struct {
	Task         TaskDetailView   `json:"task"`
	Dependencies []TaskView       `json:"dependencies"`
	Dependents   []TaskView       `json:"dependents"`
	Locations    []LocationView   `json:"locations"`
	Reports      []ReportView     `json:"reports"`
	Notes        []NoteView       `json:"notes"`
	Steps        []StepView       `json:"steps"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type TaskDetailView struct {
	ID              int        `json:"id" db:"id"`
	Slug            string     `json:"slug" db:"slug"`
	ProjectID       int        `json:"project_id" db:"project_id"`
	ProjectName     string     `json:"project_name" db:"project_name"`
	ProjectSlug     string     `json:"project_slug" db:"project_slug"`
	AgentID         *int       `json:"agent_id" db:"agent_id"`
	AgentName       *string    `json:"agent_name" db:"agent_name"`
	AgentSlug       *string    `json:"agent_slug" db:"agent_slug"`
	Type            string     `json:"type" db:"type"`
	Status          string     `json:"status" db:"status"`
	Instructions    string     `json:"instructions" db:"instructions"`
	CompletionNotes *string    `json:"completion_notes" db:"completion_notes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	StartedAt       *time.Time `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
}

type LocationView struct {
	ID          int       `json:"id" db:"id"`
	TaskID      int       `json:"task_id" db:"task_id"`
	Location    string    `json:"location" db:"location"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ReportView struct {
	ID        int            `json:"id" db:"id"`
	TaskID    int            `json:"task_id" db:"task_id"`
	Content   string         `json:"content" db:"content"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	Locations []LocationView `json:"locations,omitempty"`
}

type NoteView struct {
	ID        int       `json:"id" db:"id"`
	TaskID    int       `json:"task_id" db:"task_id"`
	Type      string    `json:"type" db:"type"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type StepView struct {
	ID        int       `json:"id" db:"id"`
	TaskID    int       `json:"task_id" db:"task_id"`
	StepType  string    `json:"step_type" db:"step_type"`
	Details   string    `json:"details" db:"details"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ReportsPageData struct {
	Projects        []ProjectView    `json:"projects"`
	Tasks           []TaskView       `json:"tasks"`
	Stats           ReportStats      `json:"stats"`
	Filters         ReportFilters    `json:"filters"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type ReportStats struct {
	TotalTasks      int                `json:"total_tasks"`
	TasksByStatus   map[string]int     `json:"tasks_by_status"`
	TasksByProject  map[string]int     `json:"tasks_by_project"`
	TasksByType     map[string]int     `json:"tasks_by_type"`
	AvgCompletion   float64            `json:"avg_completion_time"`
	Timeline        []TimelinePoint    `json:"timeline"`
}

type TimelinePoint struct {
	Date      string `json:"date"`
	Completed int    `json:"completed"`
	Started   int    `json:"started"`
	Created   int    `json:"created"`
}

type ReportFilters struct {
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Limit     string `json:"limit"`
}

func NewServeCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting (minimal for HTTP server)
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"serve",
		cmds.WithShort("Start HTTP server for task management dashboard"),
		cmds.WithLong("Start an HTTP server that provides a web interface for monitoring the task management system. Shows projects, agents, tasks, and their current status."),
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"port",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Port to serve on"),
				parameters.WithDefault(8080),
			),
		),
		// Add glazed layer
		cmds.WithLayersList(glazedLayer),
	)

	return &ServeCommand{
		CommandDescription: cmdDesc,
	}, nil
}

func (s *ServeCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	settings := &ServeSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()
	s.db = db

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/dashboard", s.handleDashboard).Methods("GET")
	api.HandleFunc("/projects", s.handleProjects).Methods("GET")
	api.HandleFunc("/agents", s.handleAgents).Methods("GET")
	api.HandleFunc("/tasks", s.handleTasks).Methods("GET")
	api.HandleFunc("/task/{id}", s.handleTaskDetail).Methods("GET")
	api.HandleFunc("/reports", s.handleReportsData).Methods("GET")

	// Static files
	staticFS, err := fs.Sub(embedFS, "static")
	if err != nil {
		return errors.Wrap(err, "failed to create static file system")
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Web pages
	r.HandleFunc("/", s.handleIndex).Methods("GET")
	r.HandleFunc("/task/{id}", s.handleTaskDetailPage).Methods("GET")
	r.HandleFunc("/task/{project}/{task}", s.handleTaskDetailPageBySlug).Methods("GET")
	r.HandleFunc("/reports", s.handleReportsPage).Methods("GET")
	r.HandleFunc("/project/{project}/reports", s.handleProjectReportsPage).Methods("GET")

	// Start server
	log.Info().Int("port", settings.Port).Msg("Starting HTTP server")
	return http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), r)
}

func (s *ServeCommand) handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := s.getDashboardData()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dashboard data")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// For simplicity, we'll just send a basic HTML with placeholder data
	// In production, you'd use html/template or similar
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Agent Task Management Dashboard</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css" rel="stylesheet">
    <style>
        .card-metric { font-size: 2rem; font-weight: bold; }
        .status-badge { font-size: 0.75rem; }
        .refresh-indicator { position: fixed; top: 20px; right: 20px; z-index: 1000; }
    </style>
</head>
<body class="bg-light">
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a href="/" class="navbar-brand">
                <i class="bi bi-cpu"></i> Agent Task Management
            </a>
            <div class="d-flex align-items-center">
                <a href="/reports" class="btn btn-outline-light me-3">
                    <i class="bi bi-bar-chart"></i> Reports
                </a>
                <span class="badge bg-success refresh-indicator" id="refreshIndicator">
                    <i class="bi bi-arrow-clockwise"></i> Auto-refresh: 30s
                </span>
            </div>
        </div>
    </nav>

    <div class="container my-4">
        <!-- Dashboard Stats -->
        <div class="row mb-4">
            <div class="col-md-3">
                <div class="card text-center border-primary">
                    <div class="card-body">
                        <i class="bi bi-folder text-primary"></i>
                        <div class="card-metric text-primary">%d</div>
                        <div class="card-title">Projects</div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card text-center border-info">
                    <div class="card-body">
                        <i class="bi bi-robot text-info"></i>
                        <div class="card-metric text-info">%d</div>
                        <div class="card-title">Agents</div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card text-center border-warning">
                    <div class="card-body">
                        <i class="bi bi-list-task text-warning"></i>
                        <div class="card-metric text-warning">%d</div>
                        <div class="card-title">In Progress</div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card text-center border-success">
                    <div class="card-body">
                        <i class="bi bi-check-circle text-success"></i>
                        <div class="card-metric text-success">%d</div>
                        <div class="card-title">Completed</div>
                    </div>
                </div>
            </div>
        </div>

        <div id="content">
            <p class="text-center">Loading dashboard data... <i class="bi bi-arrow-clockwise"></i></p>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        let refreshCounter = 30;
        
        function loadDashboard() {
            fetch('/api/dashboard')
                .then(response => response.json())
                .then(data => updateDashboard(data))
                .catch(error => console.error('Error loading dashboard:', error));
        }
        
        function updateDashboard(data) {
            const content = document.getElementById('content');
            content.innerHTML = buildDashboardHTML(data);
        }
        
        function buildDashboardHTML(data) {
            let html = '<div class="row mb-4">';
            
            // Agents Status
            html += '<div class="col-12"><div class="card"><div class="card-header">';
            html += '<h5 class="mb-0"><i class="bi bi-robot"></i> Agent Status</h5></div>';
            html += '<div class="card-body"><div class="row">';
            
            data.agents.forEach(agent => {
                html += '<div class="col-md-6 col-lg-4 mb-3">';
                html += '<div class="card border-light"><div class="card-body">';
                html += '<h6 class="card-title">' + agent.name + '</h6>';
                
                if (agent.current_task_id) {
                    html += '<span class="badge bg-warning status-badge">';
                    html += '<i class="bi bi-play-circle"></i> Working on ' + agent.current_task_slug + '</span>';
                    html += '<small class="text-muted d-block mt-1">Project: ' + agent.current_project_name + '</small>';
                } else {
                    html += '<span class="badge bg-secondary status-badge">';
                    html += '<i class="bi bi-pause-circle"></i> Idle</span>';
                }
                
                html += '</div></div></div>';
            });
            
            html += '</div></div></div></div>';
            
            // Recent Tasks
            html += '<div class="col-12"><div class="card"><div class="card-header">';
            html += '<h5 class="mb-0"><i class="bi bi-clock-history"></i> Recent Tasks</h5></div>';
            html += '<div class="card-body"><div class="table-responsive">';
            html += '<table class="table table-hover"><thead><tr>';
            html += '<th>Task</th><th>Project</th><th>Agent</th><th>Status</th><th>Type</th><th>Created</th>';
            html += '</tr></thead><tbody>';
            
            data.recent_tasks.forEach(task => {
                html += '<tr><td><a href="/task/' + encodeURIComponent(task.project_slug) + '/' + encodeURIComponent(task.slug) + '" class="text-decoration-none"><strong>#' + task.id + ': ' + task.slug + '</strong></a><br>';
                html += '<small class="text-muted">' + (task.instructions.length > 80 ? task.instructions.substring(0, 80) + '...' : task.instructions) + '</small></td>';
                html += '<td><a href="/project/' + encodeURIComponent(task.project_slug) + '/reports" class="text-decoration-none">' + task.project_name + '</a></td>';
                html += '<td>' + (task.agent_name || '<em>Unassigned</em>') + '</td>';
                html += '<td>';
                
                switch(task.status) {
                    case 'pending': html += '<span class="badge bg-secondary">Pending</span>'; break;
                    case 'in_progress': html += '<span class="badge bg-warning">In Progress</span>'; break;
                    case 'completed': html += '<span class="badge bg-success">Completed</span>'; break;
                    case 'failed': html += '<span class="badge bg-danger">Failed</span>'; break;
                }
                
                html += '</td><td>';
                html += task.type === 'gather_information' ? '<i class="bi bi-search"></i> Gather' : '<i class="bi bi-lightbulb"></i> Analysis';
                html += '</td><td>' + new Date(task.created_at).toLocaleDateString() + '</td></tr>';
            });
            
            html += '</tbody></table></div></div></div></div>';
            html += '</div>';
            
            return html;
        }
        
        function updateRefreshIndicator() {
            const indicator = document.getElementById('refreshIndicator');
            if (refreshCounter > 0) {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise"></i> Auto-refresh: ' + refreshCounter + 's';
                refreshCounter--;
            } else {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise spinning"></i> Refreshing...';
                refreshCounter = 30;
                loadDashboard();
            }
        }
        
        // Load initial data
        loadDashboard();
        
        // Update every second
        setInterval(updateRefreshIndicator, 1000);
        
        // Add spinning animation for refresh icon
        const style = document.createElement('style');
        style.textContent = '.spinning { animation: spin 1s linear infinite; } @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }';
        document.head.appendChild(style);
    </script>
</body>
</html>`, data.Stats.TotalProjects, data.Stats.TotalAgents, data.Stats.InProgressTasks, data.Stats.CompletedTasks)))
}

func (s *ServeCommand) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := s.getDashboardData()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dashboard data")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *ServeCommand) handleProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.getProjects()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get projects")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (s *ServeCommand) handleAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := s.getAgents()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get agents")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (s *ServeCommand) handleTasks(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	status := r.URL.Query().Get("status")
	projectID := r.URL.Query().Get("project_id")
	limit := r.URL.Query().Get("limit")

	tasks, err := s.getTasks(status, projectID, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get tasks")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (s *ServeCommand) getDashboardData() (*DashboardData, error) {
	projects, err := s.getProjects()
	if err != nil {
		return nil, err
	}

	agents, err := s.getAgents()
	if err != nil {
		return nil, err
	}

	tasks, err := s.getTasks("", "", "20")
	if err != nil {
		return nil, err
	}

	stats, err := s.getStats()
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		Projects:    projects,
		Agents:      agents,
		Tasks:       tasks,
		RecentTasks: tasks, // Same for now
		Stats:       stats,
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *ServeCommand) getProjects() ([]ProjectView, error) {
	query := `
		SELECT p.id, p.slug, p.name, p.description, p.concise_guidelines, p.full_guidelines, p.created_at,
		       COUNT(t.id) as task_count
		FROM projects p
		LEFT JOIN tasks t ON p.id = t.project_id
		GROUP BY p.id, p.slug, p.name, p.description, p.concise_guidelines, p.full_guidelines, p.created_at
		ORDER BY p.created_at DESC
	`

	var projects []ProjectView
	err := s.db.Select(&projects, query)
	return projects, err
}

func (s *ServeCommand) getAgents() ([]AgentView, error) {
	query := `
		SELECT a.id, a.slug, a.name, a.description, a.current_project_id, a.current_task_id, a.created_at,
		       p.name as current_project_name,
		       t.slug as current_task_slug
		FROM agents a
		LEFT JOIN projects p ON a.current_project_id = p.id
		LEFT JOIN tasks t ON a.current_task_id = t.id
		ORDER BY a.created_at DESC
	`

	var agents []AgentView
	err := s.db.Select(&agents, query)
	if err != nil {
		return nil, err
	}

	// Set status based on current task
	for i := range agents {
		if agents[i].CurrentTaskID != nil {
			agents[i].Status = "working"
		} else {
			agents[i].Status = "idle"
		}
	}

	return agents, nil
}

func (s *ServeCommand) getTasks(status, projectID, limit string) ([]TaskView, error) {
	query := `
		SELECT t.id, t.slug, t.project_id, t.agent_id, t.type, t.status, 
		       t.instructions, t.completion_notes, t.created_at, t.started_at, t.completed_at,
		       p.name as project_name, p.slug as project_slug,
		       a.name as agent_name
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		LEFT JOIN agents a ON t.agent_id = a.id
		WHERE 1=1
	`

	args := []interface{}{}

	if status != "" {
		query += " AND t.status = ?"
		args = append(args, status)
	}

	if projectID != "" {
		query += " AND t.project_id = ?"
		pid, err := strconv.Atoi(projectID)
		if err != nil {
			return nil, errors.Wrap(err, "invalid project_id")
		}
		args = append(args, pid)
	}

	query += " ORDER BY t.created_at DESC"

	if limit != "" {
		query += " LIMIT ?"
		lim, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errors.Wrap(err, "invalid limit")
		}
		args = append(args, lim)
	}

	var tasks []TaskView
	err := s.db.Select(&tasks, query, args...)
	return tasks, err
}

func (s *ServeCommand) getStats() (DashboardStats, error) {
	var stats DashboardStats

	// Get project count
	err := s.db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&stats.TotalProjects)
	if err != nil {
		return stats, err
	}

	// Get agent count
	err = s.db.QueryRow("SELECT COUNT(*) FROM agents").Scan(&stats.TotalAgents)
	if err != nil {
		return stats, err
	}

	// Get task counts by status
	err = s.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&stats.TotalTasks)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'pending'").Scan(&stats.PendingTasks)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'in_progress'").Scan(&stats.InProgressTasks)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'completed'").Scan(&stats.CompletedTasks)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'failed'").Scan(&stats.FailedTasks)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

// New handler functions for task detail and reports

func (s *ServeCommand) handleTaskDetailPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	data, err := s.getTaskDetailData(taskID)
	if err != nil {
		log.Error().Err(err).Str("task_id", taskID).Msg("Failed to get task detail data")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	s.renderTaskDetailPage(w, data)
}

func (s *ServeCommand) handleTaskDetailPageBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectSlug := vars["project"]
	taskSlug := vars["task"]

	// Resolve task by project/task slug
	taskID, err := ResolveTaskID(r.Context(), s.db, projectSlug+"/"+taskSlug)
	if err != nil {
		log.Error().Err(err).Str("project", projectSlug).Str("task", taskSlug).Msg("Failed to resolve task ID")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	data, err := s.getTaskDetailData(strconv.Itoa(taskID))
	if err != nil {
		log.Error().Err(err).Int("task_id", taskID).Msg("Failed to get task detail data")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	s.renderTaskDetailPage(w, data)
}

func (s *ServeCommand) handleTaskDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	data, err := s.getTaskDetailData(taskID)
	if err != nil {
		log.Error().Err(err).Str("task_id", taskID).Msg("Failed to get task detail data")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *ServeCommand) handleReportsPage(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	filters := ReportFilters{
		ProjectID: r.URL.Query().Get("project_id"),
		Status:    r.URL.Query().Get("status"),
		Type:      r.URL.Query().Get("type"),
		Limit:     r.URL.Query().Get("limit"),
	}

	if filters.Limit == "" {
		filters.Limit = "50" // Default limit
	}

	data, err := s.getReportsData(filters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get reports data")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	s.renderReportsPage(w, data)
}

func (s *ServeCommand) handleReportsData(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	filters := ReportFilters{
		ProjectID: r.URL.Query().Get("project_id"),
		Status:    r.URL.Query().Get("status"),
		Type:      r.URL.Query().Get("type"),
		Limit:     r.URL.Query().Get("limit"),
	}

	if filters.Limit == "" {
		filters.Limit = "50" // Default limit
	}

	data, err := s.getReportsData(filters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get reports data")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *ServeCommand) handleProjectReportsPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectSlug := vars["project"]

	// Resolve project ID from slug
	projectID, err := s.resolveProjectID(projectSlug)
	if err != nil {
		log.Error().Err(err).Str("project_slug", projectSlug).Msg("Failed to resolve project ID")
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Get query parameters for filtering
	filters := ReportFilters{
		ProjectID: strconv.Itoa(projectID),
		Status:    r.URL.Query().Get("status"),
		Type:      r.URL.Query().Get("type"),
		Limit:     r.URL.Query().Get("limit"),
	}

	if filters.Limit == "" {
		filters.Limit = "50" // Default limit
	}

	data, err := s.getReportsData(filters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get project reports data")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	s.renderProjectReportsPage(w, data, projectSlug)
}

func (s *ServeCommand) getTaskDetailData(taskIDStr string) (*TaskDetailData, error) {
	// Resolve task ID
	var taskID int
	var err error

	if strings.Contains(taskIDStr, "/") {
		taskID, err = ResolveTaskID(context.Background(), s.db, taskIDStr)
		if err != nil {
			return nil, err
		}
	} else {
		taskID, err = strconv.Atoi(taskIDStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid task ID")
		}
	}

	// Get task details
	task, err := s.getTaskDetail(taskID)
	if err != nil {
		return nil, err
	}

	// Get dependencies
	dependencies, err := s.getTaskDependencies(taskID)
	if err != nil {
		return nil, err
	}

	// Get dependents
	dependents, err := s.getTaskDependents(taskID)
	if err != nil {
		return nil, err
	}

	// Get locations
	locations, err := s.getTaskLocations(taskID)
	if err != nil {
		return nil, err
	}

	// Get reports
	reports, err := s.getTaskReports(taskID)
	if err != nil {
		return nil, err
	}

	// Get notes
	notes, err := s.getTaskNotes(taskID)
	if err != nil {
		return nil, err
	}

	// Get steps
	steps, err := s.getTaskSteps(taskID)
	if err != nil {
		return nil, err
	}

	return &TaskDetailData{
		Task:         *task,
		Dependencies: dependencies,
		Dependents:   dependents,
		Locations:    locations,
		Reports:      reports,
		Notes:        notes,
		Steps:        steps,
		UpdatedAt:    time.Now(),
	}, nil
}

func (s *ServeCommand) getTaskDetail(taskID int) (*TaskDetailView, error) {
	query := `
		SELECT t.id, t.slug, t.project_id, t.agent_id, t.type, t.status, 
		       t.instructions, t.completion_notes, t.created_at, t.started_at, t.completed_at,
		       p.name as project_name, p.slug as project_slug,
		       a.name as agent_name, a.slug as agent_slug
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		LEFT JOIN agents a ON t.agent_id = a.id
		WHERE t.id = ?
	`

	var task TaskDetailView
	err := s.db.Get(&task, query, taskID)
	return &task, err
}

func (s *ServeCommand) getTaskDependencies(taskID int) ([]TaskView, error) {
	query := `
		SELECT t.id, t.slug, t.project_id, t.agent_id, t.type, t.status, 
		       t.instructions, t.completion_notes, t.created_at, t.started_at, t.completed_at,
		       p.name as project_name,
		       a.name as agent_name
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		LEFT JOIN agents a ON t.agent_id = a.id
		JOIN task_dependencies td ON t.id = td.parent_task_id
		WHERE td.task_id = ?
		ORDER BY t.created_at DESC
	`

	var tasks []TaskView
	err := s.db.Select(&tasks, query, taskID)
	return tasks, err
}

func (s *ServeCommand) getTaskDependents(taskID int) ([]TaskView, error) {
	query := `
		SELECT t.id, t.slug, t.project_id, t.agent_id, t.type, t.status, 
		       t.instructions, t.completion_notes, t.created_at, t.started_at, t.completed_at,
		       p.name as project_name,
		       a.name as agent_name
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		LEFT JOIN agents a ON t.agent_id = a.id
		JOIN task_dependencies td ON t.id = td.task_id
		WHERE td.parent_task_id = ?
		ORDER BY t.created_at DESC
	`

	var tasks []TaskView
	err := s.db.Select(&tasks, query, taskID)
	return tasks, err
}

func (s *ServeCommand) getTaskLocations(taskID int) ([]LocationView, error) {
	query := `
		SELECT id, task_id, location, description, created_at
		FROM gathered_locations
		WHERE task_id = ?
		ORDER BY created_at DESC
	`

	var locations []LocationView
	err := s.db.Select(&locations, query, taskID)
	return locations, err
}

func (s *ServeCommand) getTaskReports(taskID int) ([]ReportView, error) {
	query := `
		SELECT id, task_id, content, created_at
		FROM reports
		WHERE task_id = ?
		ORDER BY created_at DESC
	`

	var reports []ReportView
	err := s.db.Select(&reports, query, taskID)
	if err != nil {
		return reports, err
	}

	// Get locations for each report
	for i := range reports {
		locations, err := s.getReportLocations(reports[i].ID)
		if err != nil {
			log.Error().Err(err).Int("report_id", reports[i].ID).Msg("Failed to get report locations")
			continue
		}
		reports[i].Locations = locations
	}

	return reports, nil
}

func (s *ServeCommand) getReportLocations(reportID int) ([]LocationView, error) {
	query := `
		SELECT gl.id, gl.task_id, gl.location, gl.description, gl.created_at
		FROM gathered_locations gl
		JOIN report_locations rl ON gl.id = rl.location_id
		WHERE rl.report_id = ?
		ORDER BY gl.created_at DESC
	`

	var locations []LocationView
	err := s.db.Select(&locations, query, reportID)
	return locations, err
}

func (s *ServeCommand) resolveProjectID(projectSlug string) (int, error) {
	var projectID int
	query := "SELECT id FROM projects WHERE slug = ?"
	err := s.db.QueryRow(query, projectSlug).Scan(&projectID)
	return projectID, err
}

func (s *ServeCommand) getTaskNotes(taskID int) ([]NoteView, error) {
	query := `
		SELECT id, task_id, type, content, created_at
		FROM task_notes
		WHERE task_id = ?
		ORDER BY created_at DESC
	`

	var notes []NoteView
	err := s.db.Select(&notes, query, taskID)
	return notes, err
}

func (s *ServeCommand) getTaskSteps(taskID int) ([]StepView, error) {
	query := `
		SELECT id, task_id, step_type, details, created_at
		FROM agent_steps
		WHERE task_id = ?
		ORDER BY created_at ASC
	`

	var steps []StepView
	err := s.db.Select(&steps, query, taskID)
	return steps, err
}

func (s *ServeCommand) getReportsData(filters ReportFilters) (*ReportsPageData, error) {
	// Get projects for filter dropdown
	projects, err := s.getProjects()
	if err != nil {
		return nil, err
	}

	// Get filtered tasks
	tasks, err := s.getTasks(filters.Status, filters.ProjectID, filters.Limit)
	if err != nil {
		return nil, err
	}

	// Get statistics
	stats, err := s.getReportStats(filters)
	if err != nil {
		return nil, err
	}

	return &ReportsPageData{
		Projects:  projects,
		Tasks:     tasks,
		Stats:     stats,
		Filters:   filters,
		UpdatedAt: time.Now(),
	}, nil
}

func (s *ServeCommand) getReportStats(filters ReportFilters) (ReportStats, error) {
	var stats ReportStats

	// Base query for filtering
	baseQuery := "FROM tasks t JOIN projects p ON t.project_id = p.id WHERE 1=1"
	args := []interface{}{}

	if filters.ProjectID != "" {
		baseQuery += " AND t.project_id = ?"
		pid, err := strconv.Atoi(filters.ProjectID)
		if err != nil {
			return stats, errors.Wrap(err, "invalid project_id")
		}
		args = append(args, pid)
	}

	if filters.Status != "" {
		baseQuery += " AND t.status = ?"
		args = append(args, filters.Status)
	}

	if filters.Type != "" {
		baseQuery += " AND t.type = ?"
		args = append(args, filters.Type)
	}

	// Get total tasks
	err := s.db.QueryRow("SELECT COUNT(*) "+baseQuery, args...).Scan(&stats.TotalTasks)
	if err != nil {
		return stats, err
	}

	// Get tasks by status
	stats.TasksByStatus = make(map[string]int)
	statusRows, err := s.db.Query("SELECT t.status, COUNT(*) "+baseQuery+" GROUP BY t.status", args...)
	if err != nil {
		return stats, err
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			return stats, err
		}
		stats.TasksByStatus[status] = count
	}

	// Get tasks by project
	stats.TasksByProject = make(map[string]int)
	projectRows, err := s.db.Query("SELECT p.name, COUNT(*) "+baseQuery+" GROUP BY p.name", args...)
	if err != nil {
		return stats, err
	}
	defer projectRows.Close()

	for projectRows.Next() {
		var project string
		var count int
		if err := projectRows.Scan(&project, &count); err != nil {
			return stats, err
		}
		stats.TasksByProject[project] = count
	}

	// Get tasks by type
	stats.TasksByType = make(map[string]int)
	typeRows, err := s.db.Query("SELECT t.type, COUNT(*) "+baseQuery+" GROUP BY t.type", args...)
	if err != nil {
		return stats, err
	}
	defer typeRows.Close()

	for typeRows.Next() {
		var taskType string
		var count int
		if err := typeRows.Scan(&taskType, &count); err != nil {
			return stats, err
		}
		stats.TasksByType[taskType] = count
	}

	return stats, nil
}

// Simple HTML rendering functions
func (s *ServeCommand) renderTaskDetailPage(w http.ResponseWriter, data *TaskDetailData) {
	agentInfo := ""
	if data.Task.AgentName != nil {
		agentInfo = fmt.Sprintf(`<div class="mt-2"><span class="badge bg-info"><i class="bi bi-robot"></i> %s</span></div>`, *data.Task.AgentName)
	}

	completionNotes := ""
	if data.Task.CompletionNotes != nil && *data.Task.CompletionNotes != "" {
		completionNotes = fmt.Sprintf(`
			<h6>Completion Notes</h6>
			<div class="alert alert-success">
				<pre class="mb-0">%s</pre>
			</div>
		`, *data.Task.CompletionNotes)
	}

	statusBadge := s.getStatusBadge(data.Task.Status)

	// Enhanced task detail page with locations, reports, and notes
	locationsHTML := s.buildLocationsHTML(data.Locations)
	reportsHTML := s.buildReportsHTML(data.Reports)
	notesHTML := s.buildNotesHTML(data.Notes)
	dependenciesHTML := s.buildDependenciesHTML(data.Dependencies)
	dependentsHTML := s.buildDependentsHTML(data.Dependents)

	w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Task #%d: %s - Agent Task Management</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css" rel="stylesheet">
    <style>
        .status-indicator { width: 10px; height: 10px; border-radius: 50%%; display: inline-block; margin-right: 8px; }
        .status-pending { background-color: #6c757d; }
        .status-in_progress { background-color: #ffc107; }
        .status-completed { background-color: #198754; }
        .status-failed { background-color: #dc3545; }
        .refresh-indicator { position: fixed; top: 20px; right: 20px; z-index: 1000; }
    </style>
</head>
<body class="bg-light">
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a href="/" class="navbar-brand"><i class="bi bi-cpu"></i> Agent Task Management</a>
            <div class="d-flex align-items-center">
                <a href="/reports" class="btn btn-outline-light me-3"><i class="bi bi-bar-chart"></i> Reports</a>
                <a href="/project/%s/reports" class="btn btn-outline-light me-3"><i class="bi bi-folder"></i> Project</a>
                <span class="badge bg-success refresh-indicator" id="refreshIndicator">
                    <i class="bi bi-arrow-clockwise"></i> Auto-refresh: 5s
                </span>
            </div>
        </div>
    </nav>
    <div class="container my-4">
        <nav aria-label="breadcrumb" class="mb-4">
            <ol class="breadcrumb">
                <li class="breadcrumb-item"><a href="/">Dashboard</a></li>
                <li class="breadcrumb-item active">%s</li>
            </ol>
        </nav>
        <div class="row mb-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-start">
                        <div>
                            <h4 class="card-title mb-2">
                                <span class="status-indicator status-%s"></span>#%d: %s
                            </h4>
                            <div class="text-muted">
                                <i class="bi bi-folder"></i> %s • <i class="bi bi-tag"></i> %s • <i class="bi bi-calendar"></i> %s
                            </div>
                        </div>
                        <div class="text-end">%s%s</div>
                    </div>
                    <div class="card-body">
                        <h6>Instructions</h6>
                        <p class="mb-3">%s</p>
                        %s
                    </div>
                </div>
            </div>
        </div>
        <!-- Dependencies and Dependents -->
        <div class="row mb-4">
            <div class="col-md-6">%s</div>
            <div class="col-md-6">%s</div>
        </div>
        <!-- Locations and Reports -->
        <div class="row mb-4">
            <div class="col-md-6">%s</div>
            <div class="col-md-6">%s</div>
        </div>
        <!-- Notes -->
        <div class="row mb-4">
            <div class="col-12">%s</div>
        </div>
        <!-- Quick Actions -->
        <div class="row mb-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header">
                        <h6 class="mb-0"><i class="bi bi-lightning"></i> Quick Actions</h6>
                    </div>
                    <div class="card-body">
                        <div class="btn-group me-2" role="group">
                            <button type="button" class="btn btn-outline-primary" onclick="showTakeNoteModal()">
                                <i class="bi bi-sticky"></i> Take Note
                            </button>
                            <button type="button" class="btn btn-outline-success" onclick="showWriteReportModal()" %s>
                                <i class="bi bi-file-text"></i> Write Report
                            </button>
                        </div>
                        <div class="btn-group" role="group">
                            <button type="button" class="btn btn-outline-info" onclick="refreshPage()">
                                <i class="bi bi-arrow-clockwise"></i> Refresh
                            </button>
                            <a href="/project/%s/reports" class="btn btn-outline-secondary">
                                <i class="bi bi-folder"></i> Project Reports
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        let refreshCounter = 5;
        function refreshPage() { window.location.reload(); }
        function updateRefreshIndicator() {
            const indicator = document.getElementById('refreshIndicator');
            if (refreshCounter > 0) {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise"></i> Auto-refresh: ' + refreshCounter + 's';
                refreshCounter--;
            } else {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise spinning"></i> Refreshing...';
                refreshCounter = 5;
                refreshPage();
            }
        }
        setInterval(updateRefreshIndicator, 1000);
        const style = document.createElement('style');
        style.textContent = '.spinning { animation: spin 1s linear infinite; } @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }';
        document.head.appendChild(style);
        
        // Modal functions for quick actions
        function showTakeNoteModal() {
            const taskId = %d;
            const content = prompt('Enter your note:');
            if (content) {
                const noteType = prompt('Note type (observations, bugs, ideas, lessons_learned, issues, notes):', 'notes');
                alert('Note functionality would be implemented with API call:\\n\\nCommand: ./agent-task-sqlite take-note --task=' + taskId + ' --type=' + (noteType || 'notes') + ' --content="' + content + '"');
            }
        }
        
        function showWriteReportModal() {
            const taskId = %d;
            const content = prompt('Enter your completion report:');
            if (content) {
                alert('Report functionality would be implemented with API call:\\n\\nCommand: ./agent-task-sqlite write-completion-report --task=' + taskId + ' --content="' + content + '"');
            }
        }
    </script>
</body>
</html>`,
		data.Task.ID, data.Task.Slug, data.Task.ProjectSlug, data.Task.Slug, data.Task.Status, data.Task.ID, data.Task.Slug,
		data.Task.ProjectName, data.Task.Type, data.Task.CreatedAt.Format("Jan 2, 2006 15:04"),
		statusBadge, agentInfo, data.Task.Instructions, completionNotes,
		dependenciesHTML, dependentsHTML, locationsHTML, reportsHTML, notesHTML,
		s.getReportButtonState(data.Task.Status), data.Task.ProjectSlug, data.Task.ID, data.Task.ID,
	)))
}

func (s *ServeCommand) renderProjectReportsPage(w http.ResponseWriter, data *ReportsPageData, projectSlug string) {
	taskRows := s.buildTaskRowsHTML(data.Tasks)

	w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Project Reports: %s - Agent Task Management</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css" rel="stylesheet">
    <style>.refresh-indicator { position: fixed; top: 20px; right: 20px; z-index: 1000; }</style>
</head>
<body class="bg-light">
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a href="/" class="navbar-brand"><i class="bi bi-cpu"></i> Agent Task Management</a>
            <div class="d-flex align-items-center">
                <a href="/reports" class="btn btn-outline-light me-3"><i class="bi bi-bar-chart"></i> All Reports</a>
                <span class="badge bg-success refresh-indicator" id="refreshIndicator">
                    <i class="bi bi-arrow-clockwise"></i> Auto-refresh: 5s
                </span>
            </div>
        </div>
    </nav>
    <div class="container my-4">
        <nav aria-label="breadcrumb" class="mb-4">
            <ol class="breadcrumb">
                <li class="breadcrumb-item"><a href="/">Dashboard</a></li>
                <li class="breadcrumb-item"><a href="/reports">Reports</a></li>
                <li class="breadcrumb-item active">Project: %s</li>
            </ol>
        </nav>
        <div class="row mb-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header">
                        <h4 class="card-title mb-0"><i class="bi bi-folder"></i> Project Reports: %s</h4>
                    </div>
                    <div class="card-body">
                        <div class="row">
                            <div class="col-md-3"><div class="text-center"><h5 class="text-primary">%d</h5><small>Total Tasks</small></div></div>
                            <div class="col-md-3"><div class="text-center"><h5 class="text-success">%d</h5><small>Completed</small></div></div>
                            <div class="col-md-3"><div class="text-center"><h5 class="text-warning">%d</h5><small>In Progress</small></div></div>
                            <div class="col-md-3"><div class="text-center"><h5 class="text-secondary">%d</h5><small>Pending</small></div></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="row">
            <div class="col-12">
                <div class="card">
                    <div class="card-header"><h5 class="mb-0"><i class="bi bi-list-task"></i> Project Tasks</h5></div>
                    <div class="card-body"><div class="table-responsive">%s</div></div>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        let refreshCounter = 5;
        function refreshPage() { window.location.reload(); }
        function updateRefreshIndicator() {
            const indicator = document.getElementById('refreshIndicator');
            if (refreshCounter > 0) {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise"></i> Auto-refresh: ' + refreshCounter + 's';
                refreshCounter--;
            } else {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise spinning"></i> Refreshing...';
                refreshCounter = 5;
                refreshPage();
            }
        }
        setInterval(updateRefreshIndicator, 1000);
        const style = document.createElement('style');
        style.textContent = '.spinning { animation: spin 1s linear infinite; } @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }';
        document.head.appendChild(style);
    </script>
</body>
</html>`,
		projectSlug, projectSlug, projectSlug,
		data.Stats.TotalTasks,
		data.Stats.TasksByStatus["completed"],
		data.Stats.TasksByStatus["in_progress"],
		data.Stats.TasksByStatus["pending"],
		taskRows,
	)))
}

func (s *ServeCommand) renderReportsPage(w http.ResponseWriter, data *ReportsPageData) {
	taskRows := s.buildTaskRowsHTML(data.Tasks)

	w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reports - Agent Task Management</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css" rel="stylesheet">
    <style>.refresh-indicator { position: fixed; top: 20px; right: 20px; z-index: 1000; }</style>
</head>
<body class="bg-light">
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a href="/" class="navbar-brand"><i class="bi bi-cpu"></i> Agent Task Management</a>
            <div class="d-flex align-items-center">
                <span class="badge bg-success refresh-indicator" id="refreshIndicator">
                    <i class="bi bi-arrow-clockwise"></i> Auto-refresh: 5s
                </span>
            </div>
        </div>
    </nav>
    <div class="container my-4">
        <nav aria-label="breadcrumb" class="mb-4">
            <ol class="breadcrumb">
                <li class="breadcrumb-item"><a href="/">Dashboard</a></li>
                <li class="breadcrumb-item active">Reports</li>
            </ol>
        </nav>
        <div class="row mb-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header"><h4 class="mb-0"><i class="bi bi-bar-chart"></i> Task Reports & Analytics</h4></div>
                    <div class="card-body">
                        <div class="row">
                            <div class="col-md-3"><div class="text-center"><div class="h3 text-primary">%d</div><div class="text-muted">Total Tasks</div></div></div>
                            <div class="col-md-3"><div class="text-center"><div class="h3 text-success">%d</div><div class="text-muted">Completed</div></div></div>
                            <div class="col-md-3"><div class="text-center"><div class="h3 text-warning">%d</div><div class="text-muted">In Progress</div></div></div>
                            <div class="col-md-3"><div class="text-center"><div class="h3 text-secondary">%d</div><div class="text-muted">Pending</div></div></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="row">
            <div class="col-12">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h6 class="mb-0"><i class="bi bi-list-task"></i> All Tasks</h6>
                        <small class="text-muted">%d tasks shown</small>
                    </div>
                    <div class="card-body p-0">
                        <div class="table-responsive">
                            <table class="table table-hover mb-0">
                                <thead class="table-light">
                                    <tr><th>Task</th><th>Project</th><th>Agent</th><th>Status</th><th>Type</th><th>Created</th><th>Actions</th></tr>
                                </thead>
                                <tbody>%s</tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        let refreshCounter = 5;
        function refreshPage() { window.location.reload(); }
        function updateRefreshIndicator() {
            const indicator = document.getElementById('refreshIndicator');
            if (refreshCounter > 0) {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise"></i> Auto-refresh: ' + refreshCounter + 's';
                refreshCounter--;
            } else {
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise spinning"></i> Refreshing...';
                refreshCounter = 5;
                refreshPage();
            }
        }
        setInterval(updateRefreshIndicator, 1000);
        const style = document.createElement('style');
        style.textContent = '.spinning { animation: spin 1s linear infinite; } @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }';
        document.head.appendChild(style);
    </script>
</body>
</html>`,
		data.Stats.TotalTasks, data.Stats.TasksByStatus["completed"],
		data.Stats.TasksByStatus["in_progress"], data.Stats.TasksByStatus["pending"],
		len(data.Tasks), taskRows,
	)))
}

func (s *ServeCommand) buildTaskRowsHTML(tasks []TaskView) string {
	html := ""
	for _, task := range tasks {
		agentName := "<em class=\"text-muted\">Unassigned</em>"
		if task.AgentName != nil {
			agentName = *task.AgentName
		}

		typeIcon := `<i class="bi bi-lightbulb"></i> Analysis`
		if task.Type == "gather_information" {
			typeIcon = `<i class="bi bi-search"></i> Gather`
		}

		instructions := task.Instructions
		if len(instructions) > 60 {
			instructions = instructions[:60] + "..."
		}

		html += fmt.Sprintf(`
			<tr>
				<td>
					<a href="/task/%s/%s" class="text-decoration-none"><strong>#%d: %s</strong></a><br/>
					<small class="text-muted">%s</small>
				</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td><div>%s</div><small class="text-muted">%s</small></td>
				<td><a href="/task/%s/%s" class="btn btn-sm btn-outline-primary"><i class="bi bi-eye"></i></a></td>
			</tr>
		`, task.ProjectSlug, task.Slug, task.ID, task.Slug, instructions,
		   task.ProjectName, agentName, s.getStatusBadge(task.Status), typeIcon,
		   task.CreatedAt.Format("Jan 2, 2006"), task.CreatedAt.Format("15:04"),
		   task.ProjectSlug, task.Slug)
	}
	return html
}

func (s *ServeCommand) getStatusBadge(status string) string {
	switch status {
	case "pending":
		return `<span class="badge bg-secondary">Pending</span>`
	case "in_progress":
		return `<span class="badge bg-warning text-dark">In Progress</span>`
	case "completed":
		return `<span class="badge bg-success">Completed</span>`
	case "failed":
		return `<span class="badge bg-danger">Failed</span>`
	default:
		return fmt.Sprintf(`<span class="badge bg-light text-dark">%s</span>`, status)
	}
}

func (s *ServeCommand) buildDependenciesHTML(dependencies []TaskView) string {
	if len(dependencies) == 0 {
		return `<div class="card">
			<div class="card-header">
				<h6 class="mb-0"><i class="bi bi-arrow-down-circle"></i> Dependencies</h6>
			</div>
			<div class="card-body">
				<p class="text-muted">No dependencies</p>
			</div>
		</div>`
	}

	html := `<div class="card">
		<div class="card-header">
			<h6 class="mb-0"><i class="bi bi-arrow-down-circle"></i> Dependencies</h6>
		</div>
		<div class="card-body">
			<div class="list-group list-group-flush">`

	for _, dep := range dependencies {
		html += fmt.Sprintf(`
			<div class="list-group-item">
				<div class="d-flex justify-content-between align-items-start">
					<div>
						<h6 class="mb-1"><a href="/task/%s/%s" class="text-decoration-none">#%d: %s</a></h6>
						<p class="mb-1 text-muted small">%s</p>
					</div>
					<small>%s</small>
				</div>
			</div>`,
			dep.ProjectSlug, dep.Slug, dep.ID, dep.Slug,
			dep.Instructions,
			s.getStatusBadge(dep.Status))
	}

	html += `</div></div></div>`
	return html
}

func (s *ServeCommand) buildDependentsHTML(dependents []TaskView) string {
	if len(dependents) == 0 {
		return `<div class="card">
			<div class="card-header">
				<h6 class="mb-0"><i class="bi bi-arrow-up-circle"></i> Dependents</h6>
			</div>
			<div class="card-body">
				<p class="text-muted">No dependents</p>
			</div>
		</div>`
	}

	html := `<div class="card">
		<div class="card-header">
			<h6 class="mb-0"><i class="bi bi-arrow-up-circle"></i> Dependents</h6>
		</div>
		<div class="card-body">
			<div class="list-group list-group-flush">`

	for _, dep := range dependents {
		html += fmt.Sprintf(`
			<div class="list-group-item">
				<div class="d-flex justify-content-between align-items-start">
					<div>
						<h6 class="mb-1"><a href="/task/%s/%s" class="text-decoration-none">#%d: %s</a></h6>
						<p class="mb-1 text-muted small">%s</p>
					</div>
					<small>%s</small>
				</div>
			</div>`,
			dep.ProjectSlug, dep.Slug, dep.ID, dep.Slug,
			dep.Instructions,
			s.getStatusBadge(dep.Status))
	}

	html += `</div></div></div>`
	return html
}

func (s *ServeCommand) buildLocationsHTML(locations []LocationView) string {
	if len(locations) == 0 {
		return `<div class="card">
			<div class="card-header">
				<h6 class="mb-0"><i class="bi bi-geo-alt"></i> Code Locations</h6>
			</div>
			<div class="card-body">
				<p class="text-muted">No locations</p>
			</div>
		</div>`
	}

	html := `<div class="card">
		<div class="card-header">
			<h6 class="mb-0"><i class="bi bi-geo-alt"></i> Code Locations (%d)</h6>
		</div>
		<div class="card-body">
			<div class="list-group list-group-flush">`

	for _, loc := range locations {
		html += fmt.Sprintf(`
			<div class="list-group-item">
				<div class="d-flex justify-content-between align-items-start">
					<div>
						<h6 class="mb-1"><code class="text-primary">%s</code></h6>
						<p class="mb-0 text-muted small">%s</p>
					</div>
					<small class="text-muted">%s</small>
				</div>
			</div>`,
			loc.Location, loc.Description, loc.CreatedAt.Format("Jan 2 15:04"))
	}

	html += `</div></div></div>`
	return fmt.Sprintf(html, len(locations))
}

func (s *ServeCommand) buildReportsHTML(reports []ReportView) string {
	if len(reports) == 0 {
		return `<div class="card">
			<div class="card-header">
				<h6 class="mb-0"><i class="bi bi-file-text"></i> Reports</h6>
			</div>
			<div class="card-body">
				<p class="text-muted">No reports</p>
			</div>
		</div>`
	}

	html := `<div class="card">
		<div class="card-header">
			<h6 class="mb-0"><i class="bi bi-file-text"></i> Reports (%d)</h6>
		</div>
		<div class="card-body">`

	for i, report := range reports {
		html += fmt.Sprintf(`
			<div class="mb-3%s">
				<div class="d-flex justify-content-between align-items-start mb-2">
					<h6 class="mb-0">Report #%d</h6>
					<small class="text-muted">%s</small>
				</div>
				<div class="alert alert-info">
					<pre class="mb-0 small">%s</pre>
				</div>`,
			s.getBottomMargin(i, len(reports)-1),
			report.ID, report.CreatedAt.Format("Jan 2 15:04"), report.Content)

		// Add locations for this report
		if len(report.Locations) > 0 {
			html += `<div class="mb-2"><strong>Locations:</strong></div><ul class="list-unstyled ms-3">`
			for _, loc := range report.Locations {
				html += fmt.Sprintf(`<li><code class="text-primary">%s</code> - %s</li>`, loc.Location, loc.Description)
			}
			html += `</ul>`
		}

		html += `</div>`
	}

	html += `</div></div>`
	return fmt.Sprintf(html, len(reports))
}

func (s *ServeCommand) getBottomMargin(current, last int) string {
	if current == last {
		return ""
	}
	return " border-bottom pb-3"
}

func (s *ServeCommand) buildNotesHTML(notes []NoteView) string {
	if len(notes) == 0 {
		return `<div class="card">
			<div class="card-header">
				<h6 class="mb-0"><i class="bi bi-sticky"></i> Task Notes</h6>
			</div>
			<div class="card-body">
				<p class="text-muted">No notes recorded</p>
			</div>
		</div>`
	}

	html := `<div class="card">
		<div class="card-header">
			<h6 class="mb-0"><i class="bi bi-sticky"></i> Task Notes (%d)</h6>
		</div>
		<div class="card-body">`

	for i, note := range notes {
		noteTypeClass := s.getNoteTypeClass(note.Type)
		noteTypeIcon := s.getNoteTypeIcon(note.Type)
		
		html += fmt.Sprintf(`
			<div class="mb-3%s">
				<div class="d-flex justify-content-between align-items-start mb-2">
					<h6 class="mb-0">
						<span class="badge %s">%s %s</span>
					</h6>
					<small class="text-muted">%s</small>
				</div>
				<div class="alert alert-light">
					<pre class="mb-0 small">%s</pre>
				</div>
			</div>`,
			s.getBottomMargin(i, len(notes)-1),
			noteTypeClass, noteTypeIcon, note.Type,
			note.CreatedAt.Format("Jan 2 15:04"),
			note.Content)
	}

	html += `</div></div>`
	return fmt.Sprintf(html, len(notes))
}

func (s *ServeCommand) getNoteTypeClass(noteType string) string {
	switch noteType {
	case "bugs":
		return "bg-danger"
	case "issues":
		return "bg-warning text-dark"
	case "observations":
		return "bg-info"
	case "lessons_learned":
		return "bg-success"
	case "ideas":
		return "bg-primary"
	default:
		return "bg-secondary"
	}
}

func (s *ServeCommand) getNoteTypeIcon(noteType string) string {
	switch noteType {
	case "bugs":
		return "<i class=\"bi bi-bug\"></i>"
	case "issues":
		return "<i class=\"bi bi-exclamation-triangle\"></i>"
	case "observations":
		return "<i class=\"bi bi-eye\"></i>"
	case "lessons_learned":
		return "<i class=\"bi bi-lightbulb\"></i>"
	case "ideas":
		return "<i class=\"bi bi-star\"></i>"
	default:
		return "<i class=\"bi bi-sticky\"></i>"
	}
}

func (s *ServeCommand) getReportButtonState(status string) string {
	if status != "completed" {
		return "disabled title=\"Task must be completed to write a report\""
	}
	return ""
}
