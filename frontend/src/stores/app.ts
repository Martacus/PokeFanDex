import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const error = ref<string | null>(null)

  function clearError() {
    error.value = null
  }

  return { error, clearError }
})
