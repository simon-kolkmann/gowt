package util

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

var Validators = struct {
	Time textinput.ValidateFunc
}{
	Time: func(v string) error {
		_, err := time.Parse(time.TimeOnly, v)
		return err
	},
}
