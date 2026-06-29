<script setup lang="ts">
import { computed } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useGamesStore } from '@/stores/games'

const games = useGamesStore()

const open = computed(() => games.missingGames.length > 0)

async function remove(gameId: string) {
  await games.remove(gameId)
  games.dismissMissing(gameId)
}

function keep(gameId: string) {
  games.dismissMissing(gameId)
}
</script>

<template>
  <Dialog :open="open" @update:open="(v) => !v && (games.missingGames = [])">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Missing games</DialogTitle>
        <DialogDescription>
          These games' folders were not found during the scan. Keep the entry or delete it from your
          library.
        </DialogDescription>
      </DialogHeader>

      <ul class="flex max-h-80 flex-col gap-2 overflow-y-auto py-2">
        <li
          v-for="g in games.missingGames"
          :key="g.id"
          class="flex items-center justify-between gap-3 rounded-md border p-3"
        >
          <div class="min-w-0">
            <p class="truncate text-sm font-medium">{{ g.name }}</p>
            <p class="truncate text-xs text-muted-foreground">{{ g.folderPath }}</p>
          </div>
          <div class="flex shrink-0 gap-2">
            <Button variant="outline" size="sm" @click="keep(g.id)">Keep</Button>
            <Button variant="destructive" size="sm" @click="remove(g.id)">Delete</Button>
          </div>
        </li>
      </ul>
    </DialogContent>
  </Dialog>
</template>
