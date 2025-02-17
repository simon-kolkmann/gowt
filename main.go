package main

import (
	"flag"
	"fmt"
	"gowt/store"
	"gowt/util"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

type Flags struct {
	v *bool
}

func main() {
	flags := setupAndParseFlags()

	if *flags.v {
		fmt.Printf("Version: %s\n", version)
		return
	}

	p := tea.NewProgram(NewApp(), tea.WithAltScreen())
	store.Init()

	go util.StartTimeTickLoop(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func setupAndParseFlags() Flags {
	var flags Flags

	flags.v = flag.Bool("v", false, "print version information")

	flag.Parse()

	return flags
}
