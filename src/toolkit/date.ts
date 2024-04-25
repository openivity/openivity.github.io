// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import { DateTime, Duration, type DurationUnit, type ToHumanDurationOptions } from 'luxon'

export function toTimezone(s: string, timezoneOffsetHours: number = 0): DateTime {
  let d = DateTime.fromISO(s)
  let tzStr = ''
  if (timezoneOffsetHours > 0) tzStr = '+' + timezoneOffsetHours?.toString()
  else if (timezoneOffsetHours < 0) tzStr = timezoneOffsetHours?.toString()

  if (timezoneOffsetHours) d = d.setZone(`UTC${tzStr}`)

  return d
}

export function toTimezoneDate(s: string, timezoneOffsetHours: number = 0): Date {
  return toTimezone(s, timezoneOffsetHours).toJSDate()
}

export function toTimezoneDateString(s?: string | null, timezoneOffsetHours: number = 0): string {
  if (!s) s
  return toTimezone(s!, timezoneOffsetHours).toLocaleString({
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
  seconds = Math.round(seconds)

  if (days > 0) {
    return `${String(days).padStart(2, '0')}:${String(hours).padStart(2, '0')}:${String(
      minutes
    ).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
  }

  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(
    seconds
  ).padStart(2, '0')}`
}

/**
 * Returns a string representation using Luxon Duration (e.g., 1 day, 5 hr, 6 min)
 * Zero value will be omitted (except seconds).
 */
export function toHuman(
  dur: Duration,
  smallestUnit: DurationUnit = 'seconds',
  opts?: ToHumanDurationOptions
): string {
  const units: DurationUnit[] = [
    'years',
    'months',
    'days',
    'hours',
    'minutes',
    'seconds',
    'milliseconds'
  ]
  const smallestIdx = units.indexOf(smallestUnit)
  const entries = Object.entries(
    dur
      .shiftTo(...units)
      .normalize()
      .toObject()
  ).filter(([, amount], idx) => amount > 0 && idx <= smallestIdx)
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

/**
 * Format hour to hh:mm:ss with hh as optional "1:00:22", "21:55", "0:00"
 * Zero value will be omitted (expect minute & second).
 */
export function formatMillisToHours(d: number): String {
  if (d <= 0) {
    return '0:00'
  }

  const hours = Math.floor(d / (1000 * 60 * 60))
  const minutes = Math.floor((d % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((d % (1000 * 60)) / 1000)

  let formattedTime = ''
  if (hours > 0) {
    formattedTime += hours + ':'
  }
  if (minutes > 0) {
    formattedTime += minutes + ':'
  }
  if (seconds > 0) {
    formattedTime += seconds + ''
  }

  return formattedTime
}
