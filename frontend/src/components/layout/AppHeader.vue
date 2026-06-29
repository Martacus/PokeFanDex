<script setup lang="ts">
import { RefreshCw, Sun, Moon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { useAppStore } from '@/stores/app'
import { useGamesStore } from '@/stores/games'
import { useTheme } from '@/composables/useTheme'

const app = useAppStore()
const games = useGamesStore()
const { theme, toggleTheme } = useTheme()

const titles: Record<string, string> = {
  games: 'Games',
  config: 'Configuration',
}
</script>

<template>
  <header class="flex h-14 shrink-0 items-center justify-between border-b px-6">
    <h1 class="text-base font-medium">{{ titles[app.currentView] }}</h1>
    <div class="flex items-center gap-2">
      <Button
        variant="outline"
        size="sm"
        class="gap-2"
        :disabled="games.scanning"
        @click="games.refresh()"
      >
        <RefreshCw class="size-4" :class="{ 'animate-spin': games.scanning }" />
        {{ games.scanning ? 'Scanning…' : 'Refresh' }}
      </Button>
      <Button
        variant="outline"
        size="icon"
        :aria-label="theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'"
        @click="toggleTheme"
      >
        <Sun v-if="theme === 'dark'" class="size-4" />
        <Moon v-else class="size-4" />
      </Button>
    </div>
  </header>
</template>
