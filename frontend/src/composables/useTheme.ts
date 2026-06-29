import { ref } from 'vue'

export type Theme = 'light' | 'dark'

const STORAGE_KEY = 'pfl-theme'

// Shared across all callers so the toggle and any other consumer stay in sync.
const theme = ref<Theme>(getInitialTheme())

function getInitialTheme(): Theme {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'light' || saved === 'dark') return saved
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function apply(value: Theme) {
  document.documentElement.classList.toggle('dark', value === 'dark')
}

export function useTheme() {
  function setTheme(value: Theme) {
    theme.value = value
    localStorage.setItem(STORAGE_KEY, value)
    apply(value)
  }

  function toggleTheme() {
    setTheme(theme.value === 'dark' ? 'light' : 'dark')
  }

  // Ensure the DOM reflects the current theme whenever this is used.
  apply(theme.value)

  return { theme, setTheme, toggleTheme }
}
