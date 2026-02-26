package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"railyard/internal/files"
)

// AppConfig is persisted at AppConfigPath() and is used for global configuration
type AppConfig struct {
	ModFolderPath  string `json:"modFolderPath,omitempty"`
	ExecutablePath string `json:"executablePath,omitempty"`
}

// Result of validating the appConfig paths
type ConfigPathValidation struct {
	IsConfigured        bool `json:"isConfigured"`
	ModFolderPathValid  bool `json:"modFolderPathValid"`
	ExecutablePathValid bool `json:"executablePathValid"`
}

// Checks if both required paths have been set in the AppConfig
func (c AppConfig) AreConfigPathsConfigured() bool {
	return strings.TrimSpace(c.ModFolderPath) != "" && strings.TrimSpace(c.ExecutablePath) != ""
}

// Checks whether the AppConfig has been configured and whether or not its specified paths exist on disk
func (c AppConfig) ValidateConfigPaths() (bool, ConfigPathValidation) {
	result := ConfigPathValidation{
		IsConfigured: c.AreConfigPathsConfigured(),
	}

	if !result.IsConfigured {
		return false, result
	}

	modInfo, modErr := os.Stat(c.ModFolderPath)
	result.ModFolderPathValid = modErr == nil && modInfo.IsDir()
	exeInfo, exeErr := os.Stat(c.ExecutablePath)
	result.ExecutablePathValid = exeErr == nil && !exeInfo.IsDir()

	return result.ModFolderPathValid && result.ExecutablePathValid, result
}

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func readAppConfig() (AppConfig, error) {
	return files.ReadJSON[AppConfig](
		AppConfigPath(),
		"app config",
		files.JSONReadOptions{
			AllowMissing: true,
			AllowEmpty:   true,
		},
	)
}

func writeAppConfig(cfg AppConfig) error {
	configPath := AppConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory for %q: %w", configPath, err)
	}

	serialized, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize app config: %w", err)
	}

	if err := os.WriteFile(configPath, serialized, 0644); err != nil {
		return fmt.Errorf("failed to write app config %q: %w", configPath, err)
	}

	return nil
}

// ResolveConfig returns the current config from disk, or empty defaults when missing.
func (s *ConfigService) ResolveConfig() (AppConfig, error) {
	return readAppConfig()
}

// UpdateExecutable updates and persists ExecutablePath.
func (s *ConfigService) UpdateExecutable(executablePath string) (AppConfig, error) {
	cfg, err := readAppConfig()
	if err != nil {
		return AppConfig{}, err
	}

	cfg.ExecutablePath = strings.TrimSpace(executablePath)
	if err := writeAppConfig(cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}

// UpdateModFolder updates and persists ModFolderPath.
func (s *ConfigService) UpdateModFolder(modFolderPath string) (AppConfig, error) {
	cfg, err := readAppConfig()
	if err != nil {
		return AppConfig{}, err
	}

	cfg.ModFolderPath = strings.TrimSpace(modFolderPath)
	if err := writeAppConfig(cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}
