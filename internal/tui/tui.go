package tui

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vertoo/skill-link/internal/core"
)

// StartApp starts the Bubble Tea program for the install command
func StartApp() error {
	// Initialize local setup implicitly just in case? Or rely on init?
	// The problem specs say init scans the project. So we assume it was initialized.
	// We could also check if .skill-link-lock.json exists here.

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

// executeAction applies the requested operation (copy/symlink) and updates the local lockfile
func executeAction(skillName string, action actionType) error {
	globalSkillsDir, _ := core.GetGlobalSkillsDir()
	localSkillsDir, _ := core.GetLocalSkillsDir()
	manifestPath, _ := core.GetLocalManifestPath()

	manifest, err := core.LoadLocalManifest(manifestPath)
	if err != nil {
		return err
	}

	globalPath := filepath.Join(globalSkillsDir, skillName)
	localPath := filepath.Join(localSkillsDir, skillName)

	switch action {
	case actionInstallSymlink, actionSwitchToSymlink:
		if err := core.CreateSymlink(globalPath, localPath); err != nil {
			return err
		}
		manifest.Skills[skillName] = core.LocalSkill{Type: core.InstallTypeSymlink}

	case actionInstallCopy, actionUpdateCopy, actionSwitchToCopy:
		if err := core.CopyDir(globalPath, localPath); err != nil {
			return err
		}
		hash, err := core.GenerateSkillHash(localPath)
		if err != nil {
			return fmt.Errorf("copied files but failed to generate hash: %w", err)
		}
		manifest.Skills[skillName] = core.LocalSkill{
			Type: core.InstallTypeCopy,
			Hash: hash,
		}

	case actionUninstall:
		os.RemoveAll(localPath)
		delete(manifest.Skills, skillName)
	}

	return core.SaveLocalManifest(manifestPath, manifest)
}
