# [AppName]

[One-line description of what this app does.]

## Tech Stack

- **Backend**: Go 1.25, Wails 3 (v3.0.0-alpha.74)
- **Frontend**: Vue 3 + TypeScript, Vite 5, Pinia
- **UI**: shadcn-vue (New York style, Tailwind CSS v4, lucide icons)
- **Storage**: JSON files in user-chosen folders
- **Build**: Task runner (go-task), Taskfile.yml

## Project Structure

```
├── main.go                     # App entry point, window setup, service registration
├── appconfig.go                # AppConfig struct
├── appconfigservice.go         # App-level settings (recents, preferences)
├── go.mod / Taskfile.yml
├── plans/                      # Implementation plan docs (gitignored)
├── frontend/
│   ├── bindings/               # Auto-generated Go↔TS bindings (do not edit)
│   │   └── [appname]/
│   ├── src/
│   │   ├── main.ts / App.vue
│   │   ├── components/
│   │   ├── stores/             # Pinia stores — one per domain
│   │   └── composables / lib
└── build/
```

## Rules

See `.claude/rules/` for detailed guidance:

- `architecture.md` — data flow, service design, error handling, bindings
- `build.md` — commands, codegen, package manager
- `conventions.md` — Go, frontend, UI conventions
