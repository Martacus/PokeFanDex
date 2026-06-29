<script setup lang="ts">
import { computed } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import CoverCandidate from '@/components/CoverCandidate.vue'
import { useGamesStore } from '@/stores/games'

const games = useGamesStore()

// Resolve one game at a time: always show the first pending choice.
const current = computed(() => games.pendingCoverChoices[0] ?? null)

async function choose(png: string) {
  if (!current.value) return
  const gameId = current.value.gameId
  await games.setCover(gameId, png)
  games.dismissCoverChoice(gameId)
}

function skip() {
  if (current.value) games.dismissCoverChoice(current.value.gameId)
}
</script>

<template>
  <Dialog :open="!!current" @update:open="(v) => !v && skip()">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Choose a cover image</DialogTitle>
        <DialogDescription>
          <template v-if="current">
            Pick a cover for <span class="font-medium">{{ current.name }}</span>, or skip to use the
            default.
          </template>
        </DialogDescription>
      </DialogHeader>

      <div v-if="current" class="grid max-h-80 grid-cols-3 gap-3 overflow-y-auto py-2">
        <CoverCandidate
          v-for="png in current.pngs"
          :key="png"
          :path="png"
          @select="choose(png)"
        />
      </div>

      <DialogFooter>
        <Button variant="outline" @click="skip">Skip</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
