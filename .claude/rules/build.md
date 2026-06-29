# Build & Commands

## Commands

| Command | Description |
|---------|-------------|
| `wails3 dev` | Dev mode with hot-reload (also: `task dev`) |
| `wails3 build` | Build for current platform (also: `task build`) |
| `wails3 package` | Create installer (NSIS on Windows, also: `task package`) |

## Codegen & Verification

- **Regen bindings** after any Go service change: `wails3 generate bindings -ts -d frontend/bindings`
- **Go build check**: `go build -v .` (ignore `build/ios` errors — pre-existing Wails platform stubs)
- **Vite build check**: `cd frontend && pnpm vite build`

## Package Manager

Use **pnpm** — `pnpm add` / `pnpm remove`. Never use npm.

Add shadcn-vue components via: `pnpm dlx shadcn-vue@latest add <component>`
