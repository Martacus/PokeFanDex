import { ref, watchEffect, type Ref } from 'vue'
import { GameService } from '../../bindings/simplelauncher'

// Module-level cache so the same file path is only read/encoded once.
const cache = new Map<string, string>()

const DEFAULT_COVER = '/default-cover.svg'

/**
 * Resolves a backend image file path to a displayable URL (base64 data URL).
 * Falls back to the bundled default cover when the path is empty or unreadable.
 */
export function useImageUrl(path: Ref<string>) {
  const url = ref(DEFAULT_COVER)

  watchEffect(async () => {
    const p = path.value
    if (!p) {
      url.value = DEFAULT_COVER
      return
    }
    const cached = cache.get(p)
    if (cached) {
      url.value = cached
      return
    }
    try {
      const dataUrl = await GameService.GetImageDataURL(p)
      const resolved = dataUrl || DEFAULT_COVER
      if (dataUrl) cache.set(p, dataUrl)
      url.value = resolved
    } catch {
      url.value = DEFAULT_COVER
    }
  })

  return url
}
