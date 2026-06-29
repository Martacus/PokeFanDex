package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// exeExtensions are the executable file types a scan recognizes. Only ".exe"
// for now; more can be added here later.
var exeExtensions = []string{".exe"}

// imageExtensions are the cover-image file types a scan recognizes.
var imageExtensions = []string{".png", ".jpg", ".jpeg", ".webp"}

// Game is a single launchable fan game. Persisted as games/{id}.json.
type Game struct {
	ID         string `json:"id"`
	Name       string `json:"name"`       // defaults to folder name, editable
	FolderPath string `json:"folderPath"` // absolute path to the game subfolder
	ExePath    string `json:"exePath"`    // absolute path to the chosen executable
	CoverPath  string `json:"coverPath"`  // absolute path to chosen cover, or "" for default
}

// GameCoverOptions lists candidate cover images discovered for a newly added
// game so the frontend can prompt the user to pick one.
type GameCoverOptions struct {
	GameID string   `json:"gameId"`
	Name   string   `json:"name"`
	Images []string `json:"images"` // candidate image paths in the folder
}

// ScanResult is returned by Scan so the frontend can drive follow-up prompts.
type ScanResult struct {
	Added            []Game             `json:"added"`            // newly discovered games
	NeedsCoverChoice []GameCoverOptions `json:"needsCoverChoice"` // games with candidate covers
	Missing          []Game             `json:"missing"`          // known games whose folder/exe is gone
}

// GameService owns games/*.json and covers/. Fully decoupled from other
// services; the games root is passed in by the frontend.
type GameService struct {
	mu sync.RWMutex
}

// gamesDir returns the directory holding per-game json files.
func (s *GameService) gamesDir() (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "games"), nil
}

// ListGames reads every games/*.json, sorted by name (case-insensitive).
func (s *GameService) ListGames() ([]Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.listGames()
}

// listGames is the lock-free implementation, callable from methods that already
// hold the mutex.
func (s *GameService) listGames() ([]Game, error) {
	dir, err := s.gamesDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []Game{}, nil // nothing scanned yet
		}
		return nil, fmt.Errorf("read games dir: %w", err)
	}

	games := make([]Game, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		var g Game
		if err := readJSON(filepath.Join(dir, e.Name()), &g); err != nil {
			return nil, fmt.Errorf("read game %q: %w", e.Name(), err)
		}
		games = append(games, g)
	}
	sort.Slice(games, func(i, j int) bool {
		return strings.ToLower(games[i].Name) < strings.ToLower(games[j].Name)
	})
	return games, nil
}

// Scan walks each immediate subfolder of root. New folders with an executable
// become games; known games whose folder/exe is gone are reported as Missing.
// Scan never deletes — the frontend resolves Missing via Keep/Delete.
func (s *GameService) Scan(root string) (ScanResult, error) {
	if root == "" {
		return ScanResult{}, errors.New("no games folder configured")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.listGames()
	if err != nil {
		return ScanResult{}, err
	}
	byFolder := make(map[string]Game, len(existing))
	for _, g := range existing {
		byFolder[g.FolderPath] = g
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return ScanResult{}, fmt.Errorf("read games root %q: %w", root, err)
	}

	result := ScanResult{Added: []Game{}, NeedsCoverChoice: []GameCoverOptions{}, Missing: []Game{}}
	seen := make(map[string]bool)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		folder := filepath.Join(root, e.Name())
		seen[folder] = true

		if known, ok := byFolder[folder]; ok {
			// Already known; leave user edits untouched, but try to fill in a
			// cover if it still doesn't have a valid one.
			if s.assignCover(&known, folder, &result) {
				if err := s.saveGame(known); err != nil {
					return ScanResult{}, err
				}
			}
			continue
		}

		exe := findExecutable(folder, e.Name())
		if exe == "" {
			continue // no executable found; not a game (yet)
		}

		g := Game{
			ID:         gameID(folder),
			Name:       e.Name(),
			FolderPath: folder,
			ExePath:    exe,
		}
		// A single image is auto-assigned; multiple are queued for the user.
		s.assignCover(&g, folder, &result)
		if err := s.saveGame(g); err != nil {
			return ScanResult{}, err
		}
		result.Added = append(result.Added, g)
	}

	// Any known game whose folder is no longer present is reported as missing.
	for _, g := range existing {
		if !seen[g.FolderPath] {
			if _, err := os.Stat(g.FolderPath); errors.Is(err, fs.ErrNotExist) {
				result.Missing = append(result.Missing, g)
			}
		}
	}

	return result, nil
}

// SetCover points a game's cover at the given image path and persists it. An
// empty path clears the cover (falls back to the default).
func (s *GameService) SetCover(gameID, imagePath string) (Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	g, err := s.loadGame(gameID)
	if err != nil {
		return Game{}, err
	}
	if imagePath != "" {
		if _, err := os.Stat(imagePath); err != nil {
			return Game{}, fmt.Errorf("cover image %q: %w", imagePath, err)
		}
	}
	g.CoverPath = imagePath
	if err := s.saveGame(g); err != nil {
		return Game{}, err
	}
	return g, nil
}

// PickCoverImage opens a native file dialog filtered to image types and returns
// the chosen path (empty string if cancelled). The frontend then calls SetCover.
func (s *GameService) PickCoverImage() (string, error) {
	path, err := application.Get().Dialog.OpenFile().
		CanChooseFiles(true).
		CanChooseDirectories(false).
		SetTitle("Select a cover image").
		AddFilter("Images", "*.png;*.jpg;*.jpeg;*.webp").
		PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("open image dialog: %w", err)
	}
	return path, nil
}

// assignCover gives g a cover from images in folder when it lacks a valid one.
// A single image is assigned directly; multiple images are queued in result for
// the user to choose. Returns true if g.CoverPath was modified.
func (s *GameService) assignCover(g *Game, folder string, result *ScanResult) bool {
	if g.CoverPath != "" {
		if _, err := os.Stat(g.CoverPath); err == nil {
			return false // still has a valid cover
		}
	}
	images := findImages(folder)
	switch len(images) {
	case 0:
		return false
	case 1:
		g.CoverPath = images[0]
		return true
	default:
		result.NeedsCoverChoice = append(result.NeedsCoverChoice, GameCoverOptions{
			GameID: g.ID,
			Name:   g.Name,
			Images: images,
		})
		return false
	}
}

// ListFolderImages returns the image files directly inside a game's folder, so
// the edit dialog can offer covers from the folder. gameID identifies the game.
func (s *GameService) ListFolderImages(gameID string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	g, err := s.loadGame(gameID)
	if err != nil {
		return nil, err
	}
	return findImages(g.FolderPath), nil
}

// UpdateGame persists user edits to an existing game (name, exe, cover).
func (s *GameService) UpdateGame(game Game) (Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if game.ID == "" {
		return Game{}, errors.New("game id is required")
	}
	if _, err := s.loadGame(game.ID); err != nil {
		return Game{}, err
	}
	if err := s.saveGame(game); err != nil {
		return Game{}, err
	}
	return game, nil
}

// LaunchGame starts the game's executable detached, leaving the launcher open.
func (s *GameService) LaunchGame(gameID string) error {
	s.mu.RLock()
	g, err := s.loadGame(gameID)
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	if g.ExePath == "" {
		return fmt.Errorf("game %q has no executable set", g.Name)
	}
	if _, err := os.Stat(g.ExePath); err != nil {
		return fmt.Errorf("executable %q: %w", g.ExePath, err)
	}

	cmd := exec.Command(g.ExePath)
	cmd.Dir = g.FolderPath // run from the game's folder so relative assets resolve
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launch %q: %w", g.Name, err)
	}
	// Release so the game keeps running independently of the launcher.
	return cmd.Process.Release()
}

// RemoveGame deletes the game's json and any copied cover. It does not touch the
// game's own folder on disk.
func (s *GameService) RemoveGame(gameID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	gamesDir, err := s.gamesDir()
	if err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(gamesDir, gameID+".json")); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove game %q: %w", gameID, err)
	}
	return nil
}

// GetImageDataURL reads an image file and returns it as a base64 data URL so the
// webview can display covers/candidates that live outside the web root. Returns
// "" for an empty path (caller falls back to the default cover).
func (s *GameService) GetImageDataURL(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read image %q: %w", path, err)
	}
	mime := "image/png"
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".gif":
		mime = "image/gif"
	case ".webp":
		mime = "image/webp"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

// --- internal helpers ---

// loadGame reads a single game by id.
func (s *GameService) loadGame(id string) (Game, error) {
	dir, err := s.gamesDir()
	if err != nil {
		return Game{}, err
	}
	var g Game
	if err := readJSON(filepath.Join(dir, id+".json"), &g); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Game{}, fmt.Errorf("game %q not found", id)
		}
		return Game{}, fmt.Errorf("read game %q: %w", id, err)
	}
	return g, nil
}

// saveGame persists a game atomically.
func (s *GameService) saveGame(g Game) error {
	dir, err := s.gamesDir()
	if err != nil {
		return err
	}
	return writeJSONAtomic(filepath.Join(dir, g.ID+".json"), g)
}

var slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// gameID derives a stable id from the folder: a slug of the folder name plus a
// short hash of the absolute path (so identically named folders don't collide).
func gameID(folder string) string {
	base := strings.ToLower(filepath.Base(folder))
	slug := strings.Trim(slugNonAlnum.ReplaceAllString(base, "-"), "-")
	if slug == "" {
		slug = "game"
	}
	sum := sha1.Sum([]byte(folder))
	return slug + "-" + hex.EncodeToString(sum[:])[:8]
}

// findExecutable picks the best executable in folder. If exactly one matches an
// exe extension it's used; with several, prefer one whose name matches the
// folder, otherwise the first alphabetically. Returns "" if none found.
func findExecutable(folder, folderName string) string {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return ""
	}
	var exes []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isExecutable(e.Name()) {
			exes = append(exes, e.Name())
		}
	}
	if len(exes) == 0 {
		return ""
	}
	sort.Strings(exes)
	if len(exes) == 1 {
		return filepath.Join(folder, exes[0])
	}
	wanted := strings.ToLower(folderName)
	for _, name := range exes {
		stem := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))
		if stem == wanted {
			return filepath.Join(folder, name)
		}
	}
	return filepath.Join(folder, exes[0])
}

// isExecutable reports whether name has a recognized executable extension.
func isExecutable(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	for _, e := range exeExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

// findImages returns absolute paths of recognized image files directly in
// folder, sorted.
func findImages(folder string) []string {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil
	}
	var images []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isImage(e.Name()) {
			images = append(images, filepath.Join(folder, e.Name()))
		}
	}
	sort.Strings(images)
	return images
}

// isImage reports whether name has a recognized image extension.
func isImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	for _, e := range imageExtensions {
		if ext == e {
			return true
		}
	}
	return false
}
