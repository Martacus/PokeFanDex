import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { Events } from '@wailsio/runtime'
import { GameService } from '../../bindings/simplelauncher'
import type { Game, GameCoverOptions } from '../../bindings/simplelauncher'
import { useConfigStore } from './config'

export type SortBy = 'name' | 'lastPlayed'

const SORT_KEY = 'pfl-sort'

function loadSort(): SortBy {
  const saved = localStorage.getItem(SORT_KEY)
  return saved === 'lastPlayed' ? 'lastPlayed' : 'name'
}

// Parse a backend timestamp to epoch ms; never-played (null / Go zero time) → 0.
function playedAt(g: Game): number {
  const v = g.lastPlayed as unknown as string | null
  if (!v) return 0
  const t = new Date(v).getTime()
  return Number.isNaN(t) ? 0 : t
}

export const useGamesStore = defineStore('games', () => {
  const games = ref<Game[]>([])
  const scanning = ref(false)
  const error = ref<string | null>(null)
  const sortBy = ref<SortBy>(loadSort())

  // Games ordered by the active sort: most-recently-played first (never-played
  // last, then alphabetical), or alphabetical by name.
  const sortedGames = computed(() => {
    const list = [...games.value]
    if (sortBy.value === 'lastPlayed') {
      return list.sort((a, b) => {
        const diff = playedAt(b) - playedAt(a)
        return diff !== 0 ? diff : a.name.localeCompare(b.name)
      })
    }
    return list.sort((a, b) => a.name.localeCompare(b.name))
  })

  function setSortBy(value: SortBy) {
    sortBy.value = value
    localStorage.setItem(SORT_KEY, value)
  }

  // Follow-up prompts produced by the most recent scan.
  const pendingCoverChoices = ref<GameCoverOptions[]>([])
  const missingGames = ref<Game[]>([])

  async function load() {
    error.value = null
    try {
      games.value = await GameService.ListGames()
    } catch (e) {
      error.value = String(e)
    }
  }

  // Rescan the configured games folder, reload the library, and surface any
  // cover-choice / missing-game prompts for the UI to resolve.
  async function refresh() {
    const config = useConfigStore()
    if (!config.gamesFolder) {
      error.value = 'No games folder configured. Set one in Configuration.'
      return
    }
    error.value = null
    scanning.value = true
    try {
      const result = await GameService.Scan(config.gamesFolder)
      games.value = await GameService.ListGames()
      pendingCoverChoices.value = result.needsCoverChoice ?? []
      missingGames.value = result.missing ?? []
    } catch (e) {
      error.value = String(e)
    } finally {
      scanning.value = false
    }
  }

  async function launch(gameId: string) {
    error.value = null
    try {
      const updated = await GameService.LaunchGame(gameId)
      // Reflect the new lastPlayed timestamp in the local list.
      const idx = games.value.findIndex((g) => g.id === updated.id)
      if (idx !== -1) games.value[idx] = updated
    } catch (e) {
      error.value = String(e)
    }
  }

  async function update(game: Game) {
    error.value = null
    try {
      await GameService.UpdateGame(game)
      await load()
    } catch (e) {
      error.value = String(e)
    }
  }

  async function setCover(gameId: string, imagePath: string) {
    error.value = null
    try {
      await GameService.SetCover(gameId, imagePath)
      await load()
    } catch (e) {
      error.value = String(e)
    }
  }

  // Opens the native image picker and returns the chosen path (empty if
  // cancelled). Does not persist — the caller decides when to save.
  async function chooseCoverFile(): Promise<string> {
    error.value = null
    try {
      return await GameService.PickCoverImage()
    } catch (e) {
      error.value = String(e)
      return ''
    }
  }

  async function listFolderImages(gameId: string): Promise<string[]> {
    try {
      return await GameService.ListFolderImages(gameId)
    } catch (e) {
      error.value = String(e)
      return []
    }
  }

  async function remove(gameId: string) {
    error.value = null
    try {
      await GameService.RemoveGame(gameId)
      await load()
    } catch (e) {
      error.value = String(e)
    }
  }

  function dismissCoverChoice(gameId: string) {
    pendingCoverChoices.value = pendingCoverChoices.value.filter((c) => c.gameId !== gameId)
  }

  function dismissMissing(gameId: string) {
    missingGames.value = missingGames.value.filter((g) => g.id !== gameId)
  }

  // The backend emits "game:updated" when a play session ends and playtime is
  // added; reload so tiles reflect the new total. Registered once per store.
  Events.On('game:updated', () => {
    load()
  })

  return {
    games,
    scanning,
    error,
    sortBy,
    sortedGames,
    setSortBy,
    pendingCoverChoices,
    missingGames,
    load,
    refresh,
    launch,
    update,
    setCover,
    chooseCoverFile,
    listFolderImages,
    remove,
    dismissCoverChoice,
    dismissMissing,
  }
})
