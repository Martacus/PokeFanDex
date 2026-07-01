<script setup lang="ts">
import { computed, toRef } from 'vue'
import { Play, Pencil, Clock } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { useGamesStore } from '@/stores/games'
import { useImageUrl } from '@/composables/useImageUrl'
import { relativeTime, formatPlaytime } from '@/lib/datetime'
import type { Game } from '../../bindings/simplelauncher'

const props = defineProps<{ game: Game }>()
const emit = defineEmits<{ edit: [game: Game] }>()

const games = useGamesStore()
const coverPath = computed(() => props.game.coverPath)
const coverUrl = useImageUrl(toRef(coverPath))

const lastPlayed = computed(() => relativeTime(props.game.lastPlayed as unknown as string | null))
const playtime = computed(() => formatPlaytime(props.game.playtimeSeconds))
</script>

<template>
  <div
    class="group flex flex-col overflow-hidden rounded-xl border bg-card shadow-sm transition hover:border-ring/40 hover:shadow-md"
  >
    <!-- Cover with overlaid title + playtime badge -->
    <div class="relative aspect-[3/4] overflow-hidden bg-muted">
      <img
        :src="coverUrl"
        :alt="game.name"
        class="size-full object-cover transition duration-300 group-hover:scale-105"
      />

      <span
        v-if="playtime"
        class="absolute right-2 top-2 flex items-center gap-1 rounded-full bg-black/65 px-2 py-1 text-xs font-medium text-white backdrop-blur-sm"
      >
        <Clock class="size-3" />
        {{ playtime }}
      </span>

      <div
        class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/85 via-black/45 to-transparent p-3 pt-10"
      >
        <p class="truncate text-sm font-semibold text-white" :title="game.name">{{ game.name }}</p>
        <p class="truncate text-xs text-white/70">
          {{ lastPlayed ? `Played ${lastPlayed}` : 'Never played' }}
        </p>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex gap-2 p-2.5">
      <Button class="flex-1 gap-2" size="sm" @click="games.launch(game.id)">
        <Play class="size-4" />
        Play
      </Button>
      <Button variant="outline" size="sm" aria-label="Edit" @click="emit('edit', game)">
        <Pencil class="size-4" />
      </Button>
    </div>
  </div>
</template>
