package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vertoo/skill-link/internal/core"
)

func RunRestore(args []string) {
	manifestPath, err := core.GetLocalManifestPath()
	if err != nil {
		fmt.Println("Error getting manifest path:", err)
		os.Exit(1)
	}

	manifest, err := core.LoadLocalManifest(manifestPath)
	if err != nil {
		fmt.Println("Error loading manifest:", err)
		os.Exit(1)
	}

	globalSkillsDir, _ := core.GetGlobalSkillsDir()
	localSkillsDir, _ := core.GetLocalSkillsDir()

	for skillName, skillData := range manifest.Skills {
		globalPath := filepath.Join(globalSkillsDir, skillName)
		localPath := filepath.Join(localSkillsDir, skillName)

		if _, err := os.Stat(globalPath); os.IsNotExist(err) {
			fmt.Printf("Warning: Global skill '%s' not found, skipping.\n", skillName)
			continue
		}

		if skillData.Type == core.InstallTypeSymlink {
			if err := core.CreateSymlink(globalPath, localPath); err != nil {
				fmt.Printf("Error restoring symlink for '%s': %v\n", skillName, err)
			} else {
				fmt.Printf("Restored symlink for '%s'\n", skillName)
			}
		} else if skillData.Type == core.InstallTypeCopy {
			if err := core.CopyDir(globalPath, localPath); err != nil {
				fmt.Printf("Error restoring copy for '%s': %v\n", skillName, err)
			} else {
				newHash, _ := core.GenerateSkillHash(localPath)
				skillData.Hash = newHash
				manifest.Skills[skillName] = skillData
				fmt.Printf("Restored copy for '%s'\n", skillName)
			}
		}
	}

	core.SaveLocalManifest(manifestPath, manifest)
	fmt.Println("Restore complete.")
}
