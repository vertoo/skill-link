package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vertoo/skill-link/internal/core"
)

func RunCreate(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing skill name.")
		fmt.Println("Usage: skill-link create <nazwa-skilla>")
		os.Exit(1)
	}
	skillName := args[0]

	globalSkillsDir, err := core.GetGlobalSkillsDir()
	if err != nil {
		fmt.Println("Error getting global directory:", err)
		os.Exit(1)
	}

	targetDir := filepath.Join(globalSkillsDir, skillName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Println("Failed to create skill directory:", err)
		os.Exit(1)
	}

	skillMdPath := filepath.Join(targetDir, "SKILL.md")
	if _, err := os.Stat(skillMdPath); os.IsNotExist(err) {
		template := `# Skill: ` + skillName + `

## Description
Provide a brief description of the skill here.

## Triggers
- trigger1
- trigger2

## Instructions
1. Step one
2. Step two
`
		if err := os.WriteFile(skillMdPath, []byte(template), 0644); err != nil {
			fmt.Println("Failed to generate SKILL.md:", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully created skill '%s' at %s\n", skillName, targetDir)
	} else {
		fmt.Println("Skill already exists.")
	}
}
