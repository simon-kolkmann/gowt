package util

import (
	"io"
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

var dump io.Writer

func init() {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}
}

func LogMessage(msg tea.Msg) {
	_, isTimeTickMsg := msg.(TimeTickMsg)
	_, isBlinkMsg := msg.(cursor.BlinkMsg)

	if dump != nil && !isTimeTickMsg && !isBlinkMsg {
		spew.Fdump(dump, msg)
	}

}
