package tui

import (
	"os"
	"path/filepath"

	"github.com/vertoo/skill-link/internal/core"
)

type SkillStatus int

const (
	StatusNotInstalled SkillStatus = iota
	StatusSymlink
	StatusCopyMatch
	StatusCopyMismatch
)

type SkillItem struct {
	Name   string
	Status SkillStatus
}

func loadSkills() ([]SkillItem, error) {
	globalSkillsDir, err := core.GetGlobalSkillsDir()
	if err != nil {
		return nil, err
	}
	os.MkdirAll(globalSkillsDir, 0755)

	entries, err := os.ReadDir(globalSkillsDir)
	if err != nil {
		return nil, err
	}

	manifestPath, _ := core.GetLocalManifestPath()
	manifest, _ := core.LoadLocalManifest(manifestPath)
	localDir, _ := core.GetLocalSkillsDir()

	var skills []SkillItem
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		status := StatusNotInstalled

		if localData, exists := manifest.Skills[name]; exists {
			if localData.Type == core.InstallTypeSymlink {
				status = StatusSymlink
			} else {
				currentHash, err := core.GenerateSkillHash(filepath.Join(localDir, name))
				if err == nil && currentHash == localData.Hash {
					status = StatusCopyMatch
				} else {
					status = StatusCopyMismatch
				}
			}
		}

		skills = append(skills, SkillItem{
			Name:   name,
			Status: status,
		})
	}
	return skills, nil
}
