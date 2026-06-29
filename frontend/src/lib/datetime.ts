const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' })

const UNITS: [Intl.RelativeTimeFormatUnit, number][] = [
  ['year', 1000 * 60 * 60 * 24 * 365],
  ['month', 1000 * 60 * 60 * 24 * 30],
  ['week', 1000 * 60 * 60 * 24 * 7],
  ['day', 1000 * 60 * 60 * 24],
  ['hour', 1000 * 60 * 60],
  ['minute', 1000 * 60],
]

/**
 * Formats a backend timestamp (ISO string or null) as a human-friendly relative
 * time, e.g. "2 days ago". Returns null for never-played (null or Go zero time).
 */
export function relativeTime(value: string | null | undefined): string | null {
  if (!value) return null
  const date = new Date(value)
  const ms = date.getTime()
  // Invalid date or Go's zero value (year 1) means "never".
  if (Number.isNaN(ms) || date.getFullYear() <= 1) return null

  const diff = ms - Date.now()
  for (const [unit, unitMs] of UNITS) {
    if (Math.abs(diff) >= unitMs) {
      return rtf.format(Math.round(diff / unitMs), unit)
    }
  }
  return 'just now'
}
