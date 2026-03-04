package core

import (
	"os"
	"path/filepath"
)

// GetGlobalDir returns the absolute path to ~/.agents
func GetGlobalDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agents"), nil
}

// GetGlobalSkillsDir returns the absolute path to ~/.agents/skills
func GetGlobalSkillsDir() (string, error) {
	globalDir, err := GetGlobalDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(globalDir, "skills"), nil
}

// GetLocalDir returns the directory representing the agent's context, optionally falling back to .agent
func GetLocalDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".agent"), nil
}

// GetLocalSkillsDir returns the path where local skills are resolved
func GetLocalSkillsDir() (string, error) {
	localDir, err := GetLocalDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, "skills"), nil
}

// GetLocalManifestPath returns the path to the local .skill-link-lock.json
func GetLocalManifestPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".skill-link-lock.json"), nil
}
