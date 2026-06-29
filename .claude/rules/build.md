# Build & Commands

## Task Commands

| Command | Description |
|---------|-------------|
| `task dev` | Dev mode with hot-reload |
| `task build` | Build for current platform |
| `task run` | Run the built application |
| `task package` | Create installer (NSIS on Windows) |

## Codegen & Verification

- **Regen bindings** after any Go service change: `wails3 generate bindings -ts -d frontend/bindings`
- **Go build check**: `go build -v .` (ignore `build/ios` errors — pre-existing Wails platform stubs)
- **Vite build check**: `cd frontend && pnpm vite build`

## Package Manager

Use **pnpm** — `pnpm add` / `pnpm remove`. Never use npm.

Add shadcn-vue components via: `pnpm dlx shadcn-vue@latest add <component>`
