package cmd

import (
	"strings"

	"github.com/rs/zerolog/log"
)

// logWriter is a simple writer that logs output to zerolog
type logWriter struct {
	prefix string
}

// newLogWriter creates a writer that logs to zerolog
func newLogWriter(prefix string) *logWriter {
	return &logWriter{prefix: prefix}
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	// Trim trailing whitespace and only log non-empty lines
	text := strings.TrimSpace(string(p))
	if text != "" {
		log.Debug().Str("source", l.prefix).Msg(text)
	}
	return len(p), nil
}
