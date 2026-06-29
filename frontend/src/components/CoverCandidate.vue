<script setup lang="ts">
import { toRef } from 'vue'
import { useImageUrl } from '@/composables/useImageUrl'

const props = defineProps<{ path: string }>()
defineEmits<{ select: [] }>()

const url = useImageUrl(toRef(props, 'path'))
</script>

<template>
  <button
    type="button"
    class="overflow-hidden rounded-md border transition hover:ring-2 hover:ring-ring"
    @click="$emit('select')"
  >
    <img :src="url" class="aspect-[3/4] w-full object-cover" alt="Cover candidate" />
  </button>
</template>
