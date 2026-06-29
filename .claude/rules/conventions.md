# Conventions

## Go

- One service per domain concern, registered in `main.go`
- Each service owns its own `sync.RWMutex` — no shared state
- No service-to-service calls
- Return descriptive errors — no silent failures
- All file writes use temp file + `os.Rename` (atomic)

## Frontend / Pinia

- **Cross-store calls**: call `useOtherStore()` inside the action body, not at module top-level (avoids circular import issues)
- Pinia stores are the orchestration layer — they coordinate multi-service workflows
- Components never call Go bindings directly; always go through a store

## UI

- Use shadcn-vue components wherever possible
- shadcn-vue style: New York, Tailwind CSS v4, lucide icons
- Add new components via `pnpm dlx shadcn-vue@latest add <component>`

## Plans

Plans live in `plans/` — read the relevant plan file before implementing any non-trivial feature.
