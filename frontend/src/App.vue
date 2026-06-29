<script setup lang="ts">
import { onMounted } from 'vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import AppHeader from '@/components/layout/AppHeader.vue'
import GamesView from '@/views/GamesView.vue'
import ConfigView from '@/views/ConfigView.vue'
import CoverChoiceDialog from '@/components/CoverChoiceDialog.vue'
import MissingGamesDialog from '@/components/MissingGamesDialog.vue'
import { useAppStore } from '@/stores/app'
import { useConfigStore } from '@/stores/config'
import { useGamesStore } from '@/stores/games'

const app = useAppStore()
const config = useConfigStore()
const games = useGamesStore()

onMounted(async () => {
  await config.load()
  await games.load()
})
</script>

<template>
  <div class="flex h-screen w-screen overflow-hidden bg-background text-foreground">
    <AppSidebar />
    <div class="flex min-w-0 flex-1 flex-col">
      <AppHeader />
      <main class="min-h-0 flex-1">
        <GamesView v-show="app.currentView === 'games'" />
        <ConfigView v-show="app.currentView === 'config'" />
      </main>
    </div>

    <CoverChoiceDialog />
    <MissingGamesDialog />
  </div>
</template>
