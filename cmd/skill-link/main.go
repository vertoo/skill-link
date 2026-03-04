package main

import (
	"fmt"
	"os"

	"github.com/vertoo/skill-link/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cli.RunInit(os.Args[2:])
	case "create":
		cli.RunCreate(os.Args[2:])
	case "manage":
		cli.RunManage(os.Args[2:])
	case "restore":
		cli.RunRestore(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: skill-link <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  init      Initialize a new local agent container")
	fmt.Println("  create    Create a new skill template globally")
	fmt.Println("  manage    Open TUI to install, remove, and manage skills")
	fmt.Println("  restore   Restore skills based on the local lockfile")
}
