<script setup lang="ts">
import { Gamepad2, Settings, LogOut } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { useAppStore, type View } from '@/stores/app'
import { AppService } from '../../../bindings/simplelauncher'

const app = useAppStore()

const items: { view: View; label: string; icon: any }[] = [
  { view: 'games', label: 'Games', icon: Gamepad2 },
  { view: 'config', label: 'Configuration', icon: Settings },
]

function exit() {
  AppService.Quit()
}
</script>

<template>
  <aside class="flex h-full w-56 shrink-0 flex-col border-r bg-sidebar text-sidebar-foreground">
    <div class="px-4 py-5">
      <span class="text-lg font-semibold tracking-tight">Simple Launcher</span>
    </div>

    <nav class="flex flex-1 flex-col gap-1 px-2">
      <Button
        v-for="item in items"
        :key="item.view"
        :variant="app.currentView === item.view ? 'secondary' : 'ghost'"
        class="justify-start gap-2"
        @click="app.setView(item.view)"
      >
        <component :is="item.icon" class="size-4" />
        {{ item.label }}
      </Button>
    </nav>

    <div class="p-2">
      <Button variant="ghost" class="w-full justify-start gap-2 text-destructive" @click="exit">
        <LogOut class="size-4" />
        Exit
      </Button>
    </div>
  </aside>
</template>
