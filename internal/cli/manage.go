package cli

import (
	"fmt"
	"os"

	"github.com/vertoo/skill-link/internal/tui"
)

func RunManage(args []string) {
	if err := tui.StartApp(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
