package util

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TimeTickMsg string

func StartTimeTickLoop(program *tea.Program) tea.Msg {
	for {
		currentTime := time.Now().Format(time.TimeOnly) + " Uhr"

		program.Send(TimeTickMsg(currentTime))
		time.Sleep(time.Duration(time.Second * 1))
	}
}
