package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func ClearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}
