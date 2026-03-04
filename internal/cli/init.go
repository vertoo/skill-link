package cli

import (
	"fmt"
	"os"

	"github.com/vertoo/skill-link/internal/core"
)

func RunInit(args []string) {
	localDir, err := core.GetLocalDir()
	if err != nil {
		fmt.Println("Error determining local directory:", err)
		os.Exit(1)
	}

	skillsDir, err := core.GetLocalSkillsDir()
	if err != nil {
		fmt.Println("Error determining local skills directory:", err)
		os.Exit(1)
	}

	// Create directories
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		fmt.Println("Failed to create agent directory structure:", err)
		os.Exit(1)
	}

	// Create empty manifest
	manifestPath, err := core.GetLocalManifestPath()
	if err != nil {
		fmt.Println("Failed to get manifest path:", err)
		os.Exit(1)
	}

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		manifest := &core.LocalManifest{
			Skills: make(map[string]core.LocalSkill),
		}
		if err := core.SaveLocalManifest(manifestPath, manifest); err != nil {
			fmt.Println("Failed to initialize manifest:", err)
			os.Exit(1)
		}
		fmt.Println("Initialized skill-link in", localDir)
		fmt.Println("Created empty manifest at", manifestPath)
	} else {
		fmt.Println("Project is already initialized.")
	}
}
