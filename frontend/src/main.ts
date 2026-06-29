import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './assets/index.css'
import { useTheme } from './composables/useTheme'

// Apply the persisted theme before mount to avoid a flash of the wrong theme.
useTheme()

const app = createApp(App)
app.use(createPinia())
app.mount('#app')
