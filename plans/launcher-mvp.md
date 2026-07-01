# Simple Launcher — MVP Design

A desktop launcher for your games and apps. The user points the app at a
folder; the app scans subfolders for game executables, builds a library, and
lets the user launch each game from a tile with cover art, name, and Play/Edit
controls.

## Goals (MVP)

- Static left **menu**: Games, Configuration, and a pinned **Exit** at the bottom.
- **Header bar** with the app name (logo later) and a **Refresh** button that
  rescans the games folder for anything new.
- **Configuration screen**: pick the games root folder; trigger a scan that
  categorizes subfolders and finds executables (`.exe` for now, more later).
- **Games screen**: grid of game tiles. Each tile = cover image, name above,
  **Play** button below, **Edit** button beside Play.
- Cover art: if `.png` files are found in a game's subfolder, prompt the user to
  pick one as the preview. Otherwise fall back to a bundled default cover.
- Everything persists as **JSON under the OS appdata directory**.

## Non-Goals (for now)

- Logo asset (placeholder text name for now).
- Non-`.exe` executable types (design leaves room; only `.exe` implemented).
- Download/install of games, auto-updates, online catalogs.

---

## Architecture

Follows `.claude/rules/architecture.md`:

```
Vue components → Pinia stores → Wails bindings → Go services → JSON on disk
```

- Go services own all filesystem I/O. Frontend never touches disk.
- Atomic writes everywhere: temp file + `os.Rename`.
- One service per domain, each with its own `sync.RWMutex`, no service-to-service
  calls. Pinia stores orchestrate cross-service workflows.

### Storage location

All data lives under the per-user appdata dir, resolved in Go via
`os.UserConfigDir()` (on Windows → `%AppData%`):

```
%AppData%/Simple Launcher/
├── config.json                 # app config (games root folder, preferences)
├── games/
│   ├── {gameId}.json           # one file per game
│   └── ...
└── covers/
    └── {gameId}.png            # copied/selected cover art (optional)
```

`{gameId}` is a stable ID derived from the game's subfolder name (slug + short
hash) so rescans match existing entries instead of duplicating them.

---

## Go Services

### 1. `ConfigService` (`configservice.go`)

Owns `config.json`. App-level settings.

```go
type AppConfig struct {
    GamesFolder string `json:"gamesFolder"` // root folder the user picked
}
```

Methods:
- `GetConfig() (AppConfig, error)`
- `SetGamesFolder(path string) (AppConfig, error)` — validates the path exists,
  persists, returns updated config.
- `PickGamesFolder() (string, error)` — opens the native Wails folder dialog,
  returns the chosen path (does not persist; frontend then calls SetGamesFolder).

### 2. `GameService` (`gameservice.go`)

Owns `games/*.json` and `covers/`. The library + scanning.

```go
type Game struct {
    ID         string `json:"id"`
    Name       string `json:"name"`        // defaults to folder name, editable
    FolderPath string `json:"folderPath"`  // absolute path to the game subfolder
    ExePath    string `json:"exePath"`     // absolute path to the chosen executable
    CoverPath  string `json:"coverPath"`   // path under covers/, or "" for default
}

// Returned by a scan so the frontend can prompt for cover choices etc.
type ScanResult struct {
    Added            []Game             `json:"added"`            // newly discovered games
    NeedsCoverChoice []GameCoverOptions `json:"needsCoverChoice"`
    Missing          []Game             `json:"missing"`          // known games whose folder/exe is gone
}

type GameCoverOptions struct {
    GameID  string   `json:"gameId"`
    Pngs    []string `json:"pngs"`        // candidate .png paths in the folder
}
```

Methods:
- `ListGames() ([]Game, error)` — read all `games/*.json`.
- `Scan() (ScanResult, error)` — walk each immediate subfolder of the configured
  games root. For each subfolder:
  - find the first/likely `.exe` (heuristic: single exe → use it; multiple →
    pick by name match to folder, leave editable).
  - if the folder is already a known game (by ID), skip / update missing fields.
  - if new, create a `Game`, persist it, and if `.png` files exist, add the
    folder to `NeedsCoverChoice` for the frontend to prompt on.
  - Scan never deletes on its own. Any known game whose folder (or exe) no longer
    exists is returned in `Missing` so the frontend can ask the user, per game,
    to **Keep** or **Delete** it.
- `SetCover(gameId, pngPath string) (Game, error)` — copy chosen png into
  `covers/{gameId}.png`, update the game.
- `UpdateGame(game Game) (Game, error)` — edit name / exePath / cover.
- `LaunchGame(gameId string) error` — `exec` the game's `ExePath` (detached;
  launcher stays open).
- `RemoveGame(gameId string) error` — delete the json (and cover) for a game.

> Note: `Scan` needs the games root. Per the "no service-to-service calls" rule,
> the **frontend passes the folder** into `Scan(rootPath)` (read from
> ConfigService in the store), rather than GameService calling ConfigService.

---

## Frontend

### Stores (Pinia)

- `stores/config.ts` — wraps ConfigService. Holds `gamesFolder`, `pickFolder()`,
  `setGamesFolder()`.
- `stores/games.ts` — wraps GameService. Holds `games`, `loadGames()`,
  `refresh()` (orchestrates: read folder from config store → `Scan(folder)` →
  reload list → surface any `needsCoverChoice` prompts), `launch()`, `update()`,
  `setCover()`.

### Layout / Components

- `App.vue` — app shell: left `<AppSidebar>` + right column of `<AppHeader>` +
  `<router-view>`-style content (simple view switch via a `currentView` ref in an
  app store is fine for MVP; no vue-router needed yet).
- `components/layout/AppSidebar.vue` — static menu: **Games**, **Configuration**,
  spacer, **Exit** pinned to bottom (Exit calls a Go `Quit()` or
  `application.Quit()` binding).
- `components/layout/AppHeader.vue` — app name text + **Refresh** button (calls
  `gamesStore.refresh()`), with a spinner/disabled state while scanning.
- `views/GamesView.vue` — responsive grid of `<GameTile>`.
- `components/GameTile.vue` — cover image (or default), name above, **Play** +
  **Edit** buttons below.
- `components/EditGameDialog.vue` — shadcn-vue dialog to edit name, exe path,
  and cover.
- `components/ConfigView.vue` — folder picker (shows current `gamesFolder`,
  "Choose folder" button → `configStore.pickFolder()`), and a "Scan now" button.
- `components/CoverChoiceDialog.vue` — shown when a scan returns
  `needsCoverChoice`; lets the user pick one png per game (or skip → default).
- `components/MissingGamesDialog.vue` — shown when a scan returns `missing`;
  lists each game whose folder/exe is gone with **Keep** / **Delete** per game.

### UI building blocks (shadcn-vue)

Add as needed: `button`, `dialog`, `input`, `card`, `scroll-area`, `sonner`
(toasts), `tooltip`. Use lucide icons (`@lucide/vue`): `Gamepad2`, `Settings`,
`LogOut`, `RefreshCw`, `Play`, `Pencil`.

### Default cover

Bundle `frontend/public/default-cover.png`. `GameTile` uses it whenever
`game.coverPath` is empty.

---

## Workflows

**First run / configure**
1. User opens Configuration → "Choose folder" → native dialog →
   `configStore.setGamesFolder(path)`.
2. App offers "Scan now" → `gamesStore.refresh()`.

**Scan / Refresh** (header button or Config "Scan now")
1. Store reads `gamesFolder` from config store.
2. `GameService.Scan(folder)` walks subfolders, adds new games, returns
   `ScanResult`.
3. Store reloads `games`, then if `needsCoverChoice` is non-empty, opens
   `CoverChoiceDialog` to resolve covers one game at a time.
4. If `missing` is non-empty, opens `MissingGamesDialog` listing each game whose
   folder/exe is gone, with **Keep** (leave the entry) or **Delete** (calls
   `RemoveGame`) per game.

**Play**
1. `GameTile` Play → `gamesStore.launch(gameId)` → `GameService.LaunchGame`
   spawns the exe detached. Launcher stays open.

**Edit**
1. `GameTile` Edit → `EditGameDialog` → on save → `gamesStore.update(game)`.

---

## Build / Codegen reminders

- After adding/changing Go services: `wails3 generate bindings -ts -d frontend/bindings`.
- Register `ConfigService` and `GameService` in `main.go` (remove demo
  `GreetService` + the `time` event emitter once the shell is in place).
- Verify: `go build -v .` and `cd frontend && pnpm vite build`.

## Open questions / later

- More executable types beyond `.exe` (config-driven extension list).
- Categorization rules ("categorize" subfolders) — MVP treats each immediate
  subfolder as one game; revisit if nested categories are needed.
- Per-game launch args / working directory.
