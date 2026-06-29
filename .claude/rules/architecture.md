# Architecture

## Data Flow

```
Vue 3 components → Pinia stores → Wails-generated bindings → Go services → JSON files on disk
```

- Go services handle ALL filesystem I/O — frontend never touches disk directly
- Atomic writes: temp file + `os.Rename` to prevent corruption on crash
- Wails auto-generates TypeScript bindings from Go service structs into `frontend/bindings/`

## Service Design

- One Go service per domain concern, registered in `main.go`
- Each service owns its own data and `sync.RWMutex` — no shared state between services
- No service-to-service calls — Go services are fully decoupled
- **Frontend orchestrates**: Pinia stores coordinate cross-service workflows

## Error Handling

Go services return descriptive errors → Wails rejects the promise → Pinia store catches and sets `error` ref → component displays to user. No silent failures.

## Storage

- All file I/O goes through Go services, never from frontend
- Atomic writes: all file mutations use temp file + `os.Rename`
- Store entities as individual JSON files: `<domain-folder>/{id}.json`

## Bindings

- Bindings import path from `src/`: use relative `../../bindings/[appname]`
- `@/` alias won't reach bindings outside `src/` — always use the relative path
