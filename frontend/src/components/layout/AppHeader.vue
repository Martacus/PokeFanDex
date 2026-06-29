<script setup lang="ts">
import { RefreshCw } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { useAppStore } from '@/stores/app'
import { useGamesStore } from '@/stores/games'

const app = useAppStore()
const games = useGamesStore()

const titles: Record<string, string> = {
  games: 'Games',
  config: 'Configuration',
}
</script>

<template>
  <header class="flex h-14 shrink-0 items-center justify-between border-b px-6">
    <h1 class="text-base font-medium">{{ titles[app.currentView] }}</h1>
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
  </header>
</template>
