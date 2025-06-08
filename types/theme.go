package types

type MaterialColor string

const (
	mcAmber500  MaterialColor = "#FFC107"
	mcGreen500  MaterialColor = "#4CAF50"
	mcIndigo500 MaterialColor = "#3F51B5"
	mcRed500    MaterialColor = "#F44336"
	mcWhite     MaterialColor = "#FFFFFF"
)

var Theme = struct {
	Primary string
	Success string
	Error   string
	Text    string
	Warn    string
}{
	Primary: string(mcIndigo500),
	Success: string(mcGreen500),
	Error:   string(mcRed500),
	Text:    string(mcWhite),
	Warn:    string(mcAmber500),
}
