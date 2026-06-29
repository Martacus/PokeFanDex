import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ConfigService } from '../../bindings/pokefanlauncher'

export const useConfigStore = defineStore('config', () => {
  const gamesFolder = ref('')
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function load() {
    error.value = null
    try {
      const cfg = await ConfigService.GetConfig()
      gamesFolder.value = cfg.gamesFolder
    } catch (e) {
      error.value = String(e)
    }
  }

  // Opens the native folder picker and, if the user chose a folder, persists it.
  // Returns the chosen path (empty string if cancelled).
  async function pickFolder(): Promise<string> {
    error.value = null
    loading.value = true
    try {
      const path = await ConfigService.PickGamesFolder()
      if (path) {
        const cfg = await ConfigService.SetGamesFolder(path)
        gamesFolder.value = cfg.gamesFolder
      }
      return path
    } catch (e) {
      error.value = String(e)
      return ''
    } finally {
      loading.value = false
    }
  }

  return { gamesFolder, loading, error, load, pickFolder }
})
