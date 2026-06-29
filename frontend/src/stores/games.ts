import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GameService } from '../../bindings/pokefanlauncher'
import type { Game, GameCoverOptions } from '../../bindings/pokefanlauncher'
import { useConfigStore } from './config'

export const useGamesStore = defineStore('games', () => {
  const games = ref<Game[]>([])
  const scanning = ref(false)
  const error = ref<string | null>(null)

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

  return {
    games,
    scanning,
    error,
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
