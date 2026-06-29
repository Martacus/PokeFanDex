import { defineStore } from 'pinia'
import { ref } from 'vue'

export type View = 'games' | 'config'

export const useAppStore = defineStore('app', () => {
  const currentView = ref<View>('games')

  function setView(view: View) {
    currentView.value = view
  }

  return { currentView, setView }
})
