<script setup lang="ts">
import { ref, watch, toRef } from 'vue'
import { ImageOff } from '@lucide/vue'
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
import CoverCandidate from '@/components/CoverCandidate.vue'
import { useGamesStore } from '@/stores/games'
import { useImageUrl } from '@/composables/useImageUrl'
import { Game } from '../../bindings/pokefanlauncher'

const props = defineProps<{ game: Game | null }>()
const emit = defineEmits<{ 'update:open': [value: boolean] }>()

const games = useGamesStore()

const name = ref('')
const exePath = ref('')
const coverPath = ref('')
const folderImages = ref<string[]>([])

const previewUrl = useImageUrl(toRef(coverPath))

watch(
  () => props.game,
  async (g) => {
    if (g) {
      name.value = g.name
      exePath.value = g.exePath
      coverPath.value = g.coverPath
      folderImages.value = await games.listFolderImages(g.id)
    }
  },
  { immediate: true },
)

async function chooseFile() {
  const path = await games.chooseCoverFile()
  if (path) coverPath.value = path
}

function close() {
  emit('update:open', false)
}

async function save() {
  if (!props.game) return
  const updated = new Game({
    ...props.game,
    name: name.value.trim() || props.game.name,
    exePath: exePath.value.trim(),
    coverPath: coverPath.value,
  })
  await games.update(updated)
  close()
}
</script>

<template>
  <Dialog :open="!!game" @update:open="(v) => emit('update:open', v)">
    <DialogContent class="max-h-[85vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Edit game</DialogTitle>
        <DialogDescription>Update the name, executable, or cover for this game.</DialogDescription>
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

        <div class="grid gap-2">
          <span class="text-sm font-medium">Cover</span>
          <div class="flex gap-3">
            <img
              :src="previewUrl"
              alt="Cover preview"
              class="aspect-[3/4] w-24 shrink-0 rounded-md border object-cover"
            />
            <div class="flex flex-col gap-2">
              <Button variant="outline" size="sm" @click="chooseFile">Choose file…</Button>
              <Button
                v-if="coverPath"
                variant="ghost"
                size="sm"
                class="gap-2 text-muted-foreground"
                @click="coverPath = ''"
              >
                <ImageOff class="size-4" />
                Use default
              </Button>
            </div>
          </div>

          <div v-if="folderImages.length" class="grid gap-1.5">
            <span class="text-xs text-muted-foreground">From this game's folder</span>
            <div class="grid grid-cols-4 gap-2">
              <CoverCandidate
                v-for="img in folderImages"
                :key="img"
                :path="img"
                @select="coverPath = img"
              />
            </div>
          </div>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="close">Cancel</Button>
        <Button @click="save">Save</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
