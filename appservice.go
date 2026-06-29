package main

import "github.com/wailsapp/wails/v3/pkg/application"

// AppService exposes app-level actions to the frontend (e.g. the Exit button).
type AppService struct{}

// Quit terminates the application.
func (s *AppService) Quit() {
	application.Get().Quit()
}
