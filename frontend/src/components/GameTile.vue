<script setup lang="ts">
import { computed, toRef } from 'vue'
import { Play, Pencil } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { useGamesStore } from '@/stores/games'
import { useImageUrl } from '@/composables/useImageUrl'
import type { Game } from '../../bindings/pokefanlauncher'

const props = defineProps<{ game: Game }>()
const emit = defineEmits<{ edit: [game: Game] }>()

const games = useGamesStore()
const coverPath = computed(() => props.game.coverPath)
const coverUrl = useImageUrl(toRef(coverPath))
</script>

<template>
  <div class="flex flex-col gap-2">
    <p class="truncate text-sm font-medium" :title="game.name">{{ game.name }}</p>

    <div class="overflow-hidden rounded-lg border bg-muted">
      <img :src="coverUrl" :alt="game.name" class="aspect-[3/4] w-full object-cover" />
    </div>

    <div class="flex gap-2">
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
