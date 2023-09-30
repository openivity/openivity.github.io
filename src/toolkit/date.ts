import { DateTime } from 'luxon'

export function toTimezoneDateString(s: string, timezoneOffsetHours?: number): string {
  let d = DateTime.fromISO(s)
  if (timezoneOffsetHours) d = d.setZone(`UTC+${timezoneOffsetHours}`)

  return d.toLocaleString({
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    hour12: false,
    minute: '2-digit',
    second: '2-digit'
  })
}

export function GMTString(timezoneOffsetHours?: number): string {
  if (!timezoneOffsetHours) return 'UTC'
  let s = 'GMT'
  if (timezoneOffsetHours > 0) s += '+'
  return s + timezoneOffsetHours.toString()
}
