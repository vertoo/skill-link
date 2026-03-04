package core

import (
	"encoding/json"
	"os"
)

// GlobalConfig represents the settings in ~/.skill-link/config.json
type GlobalConfig struct {
	Editor string `json:"editor,omitempty"`
}

// SkillInstallType represents how a skill is mapped locally
type SkillInstallType string

const (
	InstallTypeSymlink SkillInstallType = "symlink"
	InstallTypeCopy    SkillInstallType = "copy"
)

// LocalSkill represents a single mapped skill in .skill-link-lock.json
type LocalSkill struct {
	Type SkillInstallType `json:"type"`
	Hash string           `json:"hash,omitempty"` // Used only when Type is "copy"
}

// LocalManifest represents the .skill-link-lock.json structure
type LocalManifest struct {
	Skills map[string]LocalSkill `json:"skills"`
}

// LoadGlobalConfig reads the global configuration file
func LoadGlobalConfig(path string) (*GlobalConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &GlobalConfig{}, nil
		}
		return nil, err
	}

	var config GlobalConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SaveGlobalConfig writes the global configuration to file
func SaveGlobalConfig(path string, config *GlobalConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadLocalManifest reads the local project's manifest
func LoadLocalManifest(path string) (*LocalManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &LocalManifest{Skills: make(map[string]LocalSkill)}, nil
		}
		return nil, err
	}

	var manifest LocalManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	if manifest.Skills == nil {
		manifest.Skills = make(map[string]LocalSkill)
	}

	return &manifest, nil
}

// SaveLocalManifest writes the local manifest back to disk
func SaveLocalManifest(path string, manifest *LocalManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
