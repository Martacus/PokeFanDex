<script setup lang="ts">
import { ref } from 'vue'
import { Gamepad2 } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import GameTile from '@/components/GameTile.vue'
import EditGameDialog from '@/components/EditGameDialog.vue'
import { useGamesStore } from '@/stores/games'
import { useAppStore } from '@/stores/app'
import type { Game } from '../../bindings/pokefanlauncher'

const games = useGamesStore()
const app = useAppStore()

const editing = ref<Game | null>(null)
</script>

<template>
  <div class="h-full overflow-y-auto p-6">
    <p v-if="games.error" class="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
      {{ games.error }}
    </p>

    <div
      v-if="games.games.length === 0"
      class="flex h-full flex-col items-center justify-center gap-3 text-center text-muted-foreground"
    >
      <Gamepad2 class="size-12 opacity-50" />
      <p>No games yet.</p>
      <p class="text-sm">
        Set your games folder in Configuration, then hit Refresh to scan for games.
      </p>
      <Button variant="outline" size="sm" @click="app.setView('config')">Go to Configuration</Button>
    </div>

    <div
      v-else
      class="grid grid-cols-[repeat(auto-fill,minmax(160px,1fr))] gap-5"
    >
      <GameTile v-for="g in games.games" :key="g.id" :game="g" @edit="editing = $event" />
    </div>

    <EditGameDialog :game="editing" @update:open="(v) => !v && (editing = null)" />
  </div>
</template>
