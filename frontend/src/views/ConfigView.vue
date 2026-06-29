<script setup lang="ts">
import { FolderOpen, ScanLine } from '@lucide/vue'
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useConfigStore } from '@/stores/config'
import { useGamesStore } from '@/stores/games'

const config = useConfigStore()
const games = useGamesStore()

async function chooseFolder() {
  await config.pickFolder()
}
</script>

<template>
  <div class="h-full overflow-y-auto p-6">
    <Card class="max-w-2xl">
      <CardHeader>
        <CardTitle>Games folder</CardTitle>
        <CardDescription>
          Choose the folder that contains your games. Each subfolder is treated as a game; scanning
          looks for an executable (.exe) inside.
        </CardDescription>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div class="rounded-md border bg-muted/40 px-3 py-2 text-sm">
          <span v-if="config.gamesFolder" class="break-all">{{ config.gamesFolder }}</span>
          <span v-else class="text-muted-foreground">No folder selected yet.</span>
        </div>

        <p v-if="config.error" class="text-sm text-destructive">{{ config.error }}</p>

        <div class="flex gap-2">
          <Button class="gap-2" :disabled="config.loading" @click="chooseFolder">
            <FolderOpen class="size-4" />
            Choose folder
          </Button>
          <Button
            variant="outline"
            class="gap-2"
            :disabled="!config.gamesFolder || games.scanning"
            @click="games.refresh()"
          >
            <ScanLine class="size-4" />
            {{ games.scanning ? 'Scanning…' : 'Scan now' }}
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
