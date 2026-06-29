<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useGamesStore } from '@/stores/games'
import { Game } from '../../bindings/pokefanlauncher'

const props = defineProps<{ game: Game | null }>()
const emit = defineEmits<{ 'update:open': [value: boolean] }>()

const games = useGamesStore()

const name = ref('')
const exePath = ref('')

watch(
  () => props.game,
  (g) => {
    if (g) {
      name.value = g.name
      exePath.value = g.exePath
    }
  },
  { immediate: true },
)

function close() {
  emit('update:open', false)
}

async function save() {
  if (!props.game) return
  const updated = new Game({
    ...props.game,
    name: name.value.trim() || props.game.name,
    exePath: exePath.value.trim(),
  })
  await games.update(updated)
  close()
}
</script>

<template>
  <Dialog :open="!!game" @update:open="(v) => emit('update:open', v)">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Edit game</DialogTitle>
        <DialogDescription>Update the name or executable path for this game.</DialogDescription>
      </DialogHeader>

      <div class="grid gap-4 py-2">
        <div class="grid gap-1.5">
          <label class="text-sm font-medium" for="edit-name">Name</label>
          <Input id="edit-name" v-model="name" placeholder="Game name" />
        </div>
        <div class="grid gap-1.5">
          <label class="text-sm font-medium" for="edit-exe">Executable path</label>
          <Input id="edit-exe" v-model="exePath" placeholder="C:\\path\\to\\game.exe" />
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="close">Cancel</Button>
        <Button @click="save">Save</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
