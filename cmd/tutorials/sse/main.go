// main.go
package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Event struct {
	Message string
	Time    string
}

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/events", handleSSE)
	http.HandleFunc("/clear", handleClear)
	http.HandleFunc("/start-sse", handleStartSSE)

	fmt.Println("Server starting at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	log.Printf("SSE connection requested from %s", r.RemoteAddr)

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	log.Printf("SSE headers set for %s", r.RemoteAddr)

	// Create context that will be canceled when client disconnects
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Create error group with context
	g, ctx := errgroup.WithContext(ctx)

	// Create channel for events with buffer to prevent blocking
	events := make(chan Event, 1)
	log.Printf("Starting event generation for %s", r.RemoteAddr)

	// Generate events in separate goroutine
	g.Go(func() error {
		defer func() {
			log.Printf("Closing events channel for %s", r.RemoteAddr)
			close(events)
		}()

		count := 1
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Context done, stopping event generation for %s: %v", r.RemoteAddr, ctx.Err())
				return ctx.Err()
			case <-ticker.C:
				event := Event{
					Message: fmt.Sprintf("Event #%d", count),
					Time:    time.Now().Format("15:04:05"),
				}
				log.Printf("Generated event #%d for %s", count, r.RemoteAddr)

				select {
				case events <- event:
					log.Printf("Sent event #%d to channel for %s", count, r.RemoteAddr)
					count++
				case <-ctx.Done():
					log.Printf("Context done while sending event #%d for %s: %v", count, r.RemoteAddr, ctx.Err())
					return ctx.Err()
				}
			}
		}
	})

	// Send events to client
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Streaming unsupported for %s", r.RemoteAddr)
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	g.Go(func() error {
		log.Printf("Starting event sender for %s", r.RemoteAddr)
		for {
			select {
			case <-ctx.Done():
				log.Printf("Context done in event sender for %s: %v", r.RemoteAddr, ctx.Err())
				return ctx.Err()
			case event, ok := <-events:
				if !ok {
					log.Printf("Event channel closed for %s", r.RemoteAddr)
					return nil
				}

				// Send named event for updates
				updateEvent := fmt.Sprintf("event: Update\ndata: <div class=\"event\"><span class=\"time\">%s</span> %s</div>\n\n",
					event.Time, event.Message)
				_, err := w.Write([]byte(updateEvent))
				if err != nil {
					log.Printf("Error sending update event for %s: %v", r.RemoteAddr, err)
					return err
				}
				log.Printf("Sent update event for %s", r.RemoteAddr)
				flusher.Flush()

				time.Sleep(1 * time.Second)

				// Send unnamed event for status
				statusEvent := fmt.Sprintf("event: TestMessage\ndata: <div class=\"status\">Active - Last update: %s</div>\n\n",
					event.Time)
				_, err = w.Write([]byte(statusEvent))
				if err != nil {
					log.Printf("Error sending status event for %s: %v", r.RemoteAddr, err)
					return err
				}
				log.Printf("Sent status event for %s", r.RemoteAddr)

				flusher.Flush()
				log.Printf("Flushed events to %s", r.RemoteAddr)
			}
		}
	})

	log.Printf("Waiting for event handlers to complete for %s", r.RemoteAddr)
	if err := g.Wait(); err != nil && err != context.Canceled {
		log.Printf("Error in SSE handler for %s: %v", r.RemoteAddr, err)
	}
	log.Printf("SSE handler completed for %s", r.RemoteAddr)
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<div id=\"events\"></div>"))
}

// handleStartSSE returns the SSE container HTML
func handleStartSSE(w http.ResponseWriter, r *http.Request) {
	html := `
		<div class="sse-container" 
			 hx-ext="sse" 
			 sse-connect="/events">
			 Nammed
			<!-- Named events will go here -->
			<div id="events" sse-swap="Update">
				<div class="event">Waiting for first event...</div>
			</div>

			Status

			<!-- Unnamed events (status updates) will go here -->
			<div id="status" sse-swap="TestMessage">
				<div class="status">Connected - Waiting for updates...</div>
			</div>

			Foobar

			<!-- Clear events button -->
			<div class="controls">
				<button class="btn"
						hx-get="/clear"
						hx-target="#events"
						hx-swap="innerHTML">
					Clear Events
				</button>
			</div>
		</div>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
