import { DateTime, Duration, type ToHumanDurationOptions } from 'luxon'

export function toTimezoneDateString(
  s?: string | null,
  timezoneOffsetHours?: Number | null
): string {
  if (!s) return ''
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

export function secondsToDHMS(seconds: number): string {
  const days = Math.floor(seconds / (60 * 60 * 24))
  seconds -= days * (60 * 60 * 24)

  const hours = Math.floor(seconds / (60 * 60))
  seconds -= hours * (60 * 60)

  const minutes = Math.floor(seconds / 60)
  seconds -= minutes * 60

  if (days > 0) {
    return `${String(days).padStart(2, '0')}:${String(hours).padStart(2, '0')}:${String(
      minutes
    ).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
  }

  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(
    seconds
  ).padStart(2, '0')}`
}

export function toHuman(
  dur: Duration,
  smallestUnit = 'seconds',
  opts?: ToHumanDurationOptions
): string {
  const units = ['years', 'months', 'days', 'hours', 'minutes', 'seconds', 'milliseconds']
  const smallestIdx = units.indexOf(smallestUnit)
  const entries = Object.entries(
    dur
      .shiftTo(...units)
      .normalize()
      .toObject()
  ).filter(([_unit, amount], idx) => amount > 0 && idx <= smallestIdx)
  const dur2 = Duration.fromObject(
    entries.length === 0 ? { [smallestUnit]: 0 } : Object.fromEntries(entries)
  )
  return dur2.toHuman(opts)
}

/**
 * Returns a string representation of milliseconds in this format "1d 2h 51m 22s".
 * Zero value will be omitted (expect seconds).
 */
export function formatMillis(d: number): String {
  if (d <= 0) {
    return '0s'
  }

  const days = Math.floor(d / (1000 * 60 * 60 * 24))
  const hours = Math.floor((d % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((d % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((d % (1000 * 60)) / 1000)

  let formattedTime = ''
  if (days > 0) {
    formattedTime += days + 'd '
  }
  if (hours > 0) {
    formattedTime += hours + 'h '
  }
  if (minutes > 0) {
    formattedTime += minutes + 'm '
  }
  if (seconds > 0) {
    formattedTime += seconds + 's'
  }

  return formattedTime
}
