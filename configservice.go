package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AppConfig holds app-level settings persisted to config.json.
type AppConfig struct {
	// GamesFolder is the root folder the user picked; each immediate subfolder
	// is treated as a game.
	GamesFolder string `json:"gamesFolder"`
}

// ConfigService owns config.json. It is fully decoupled from other services.
type ConfigService struct {
	mu sync.RWMutex
}

// configPath returns the absolute path to config.json under the data dir.
func (s *ConfigService) configPath() (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// GetConfig returns the current app config, or a zero-value config if none has
// been saved yet.
func (s *ConfigService) GetConfig() (AppConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path, err := s.configPath()
	if err != nil {
		return AppConfig{}, err
	}
	var cfg AppConfig
	if err := readJSON(path, &cfg); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return AppConfig{}, nil // first run: no config yet
		}
		return AppConfig{}, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}

// SetGamesFolder validates that path exists and is a directory, then persists it
// as the games root and returns the updated config.
func (s *ConfigService) SetGamesFolder(path string) (AppConfig, error) {
	info, err := os.Stat(path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("games folder %q: %w", path, err)
	}
	if !info.IsDir() {
		return AppConfig{}, fmt.Errorf("games folder %q is not a directory", path)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cfgPath, err := s.configPath()
	if err != nil {
		return AppConfig{}, err
	}
	cfg := AppConfig{GamesFolder: path}
	if err := writeJSONAtomic(cfgPath, cfg); err != nil {
		return AppConfig{}, fmt.Errorf("save config: %w", err)
	}
	return cfg, nil
}

// PickGamesFolder opens the native folder-selection dialog and returns the
// chosen path. It does not persist; the frontend calls SetGamesFolder next. An
// empty string with nil error means the user cancelled.
func (s *ConfigService) PickGamesFolder() (string, error) {
	path, err := application.Get().Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("Select your games folder").
		PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("open folder dialog: %w", err)
	}
	return path, nil
}
