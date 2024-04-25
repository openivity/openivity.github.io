<!-- Copyright (C) 2023 Openivity

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. -->

<template>
  <!-- input -->
  <div class="summary mt-4" v-if="isActivityFileReady">
    <div class="px-3 pt-3 pb-4">
      <div class="time-created pb-2">
        <i class="fa-solid fa-clock"></i>&nbsp;
        {{ formatCreatedDate(timeCreated!, timezone!) }}
        &nbsp;|
        {{ timezone == 0 ? 'UTC' : timezone > 0 ? '+' + timezone : timezone }}
      </div>
      <div class="row" style="font-size: 1em; color: var(--bs-heading-color)">
        <div class="col text-start fs-8">
          <select
            class="form-select form-select-sm"
            name="sessions"
            id="sessions"
            v-model="sessionSelected"
          >
            <option value="-2" v-if="sessions.length == 0">No Session</option>
            <option value="-1" v-if="sessions.length > 1">Multi-Sessions</option>
            <option v-for="(session, index) in sessions" v-bind:key="index" :value="index">
              {{ sessions.length > 1 ? 'S' + (index + 1) : 'Session' }}: {{ session.sport }}
            </option>
          </select>
        </div>
        <div class="col-auto text-end pt-1">
          <i class="fa-solid fa-microchip pe-1"></i>{{ creatorName }}
        </div>
      </div>
      <div class="row pt-2 text-center">
        <div class="col px-0">
          <div>
            <div class="summary-title">Moving Time</div>
            <div class="summary-value fs-6">
              {{ summary.totalMovingTime ? secondsToDHMS(summary.totalMovingTime) : '-:-' }}
            </div>
          </div>
        </div>
        <div class="col px-0">
          <div>
            <div class="summary-title">Distance</div>
            <div class="summary-value fs-6 lr-border">
              {{ summary.totalDistance ? (summary.totalDistance / 1000)?.toFixed(2) : '-' }} km
            </div>
          </div>
        </div>
        <div class="col px-0">
          <div>
            <div class="summary-title">Elevation Gain</div>
            <div class="summary-value fs-6">{{ summary.totalAscent?.toFixed(0) ?? '-' }} m</div>
          </div>
        </div>
      </div>
      <div class="row pt-2 text-center">
        <div class="col px-0">
          <div>
            <div class="summary-title">Elapsed</div>
            <div class="summary-value fs-6">
              {{ summary.totalElapsedTime ? secondsToDHMS(summary.totalElapsedTime) : '-:-' }}
            </div>
          </div>
        </div>
        <div class="col px-0">
          <div v-if="summary.sport == 'Running'">
            <div class="summary-title">Avg Pace</div>
            <div class="summary-value fs-6 lr-border">
              {{ summary.avgPace ? formatPace(summary.avgPace) : '-:-' }} /km
            </div>
          </div>
          <div v-else>
            <div class="summary-title">Avg Speed</div>
            <div class="summary-value fs-6 lr-border">
              {{ summary.avgSpeed ? ((summary.avgSpeed * 3600) / 1000).toFixed(2) : '-' }} km/h
            </div>
          </div>
        </div>
        <div class="col px-0">
          <div>
            <div class="summary-title">Calories</div>
            <div class="summary-value fs-6">
              {{ summary.totalCalories ? summary.totalCalories.toLocaleString() : '-' }} Cal
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Session } from '@/spec/activity'
import { Summary } from '@/spec/summary'
import { GMTString, secondsToDHMS, toTimezone, toTimezoneDateString } from '@/toolkit/date'
import { avg, max, sum } from '@/toolkit/number'
import { formatPace } from '@/toolkit/pace'

export const NONE: number = -2
export const MULTIPLE: number = -1

export default {
  props: {
    isActivityFileReady: Boolean,
    sessions: {
      type: Array<Session>,
      required: true
    },
    selectedSessions: {
      type: Array<Session>,
      required: true
    }
  },
  data() {
    return {
      sessionSelected: NONE
    }
  },
  watch: {
    sessions: {
      handler(sessions: Session[]) {
        this.updateSelectedSession(sessions)
      }
    },
    sessionSelected: {
      handler(value) {
        this.$emit('sessionSelected', parseInt(value))
      }
    }
  },
  computed: {
    creatorName(): string {
      if (this.selectedSessions.length > 0) return this.selectedSessions[0].creatorName
      return 'Unknown'
    },
    timeCreated(): string {
      if (this.selectedSessions.length > 0) return this.selectedSessions[0].timeCreated!
      return ''
    },
    timezone(): number {
      if (this.selectedSessions.length > 0) return this.selectedSessions[0].timezone
      return 0
    },
    summary(): Summary {
      const summary = new Summary()
      let sports: Set<string> = new Set()

      for (let i = 0; i < this.selectedSessions.length; i++) {
        let gapSeconds: number = 0
        const session = this.selectedSessions[i]
        if (i > 0) {
          const prev = this.selectedSessions[i - 1]

          const prevLastTimestamp = new Date(prev.endTime).getTime()
          const currentFirstTimestamp = new Date(session.startTime).getTime()

          gapSeconds = (currentFirstTimestamp - prevLastTimestamp) / 1000
        }

        sports.add(session.sport)

        summary.totalMovingTime = sum(summary.totalMovingTime, session.totalMovingTime)

        summary.totalElapsedTime = sum(summary.totalElapsedTime, session.totalElapsedTime)
        summary.totalElapsedTime = summary.totalElapsedTime! + gapSeconds

        summary.totalDistance = sum(summary.totalDistance, session.totalDistance)
        summary.totalAscent = sum(summary.totalAscent, session.totalAscent)
        summary.totalDescent = sum(summary.totalDescent, session.totalDescent)
        summary.totalCycles = sum(summary.totalCycles, session.totalCycles)
        summary.totalCalories = sum(summary.totalCalories, session.totalCalories)
        summary.avgSpeed = avg(summary.avgSpeed, session.avgSpeed)
        summary.maxSpeed = max(summary.maxSpeed, session.maxSpeed)
        summary.avgHeartRate = avg(summary.avgHeartRate, session.avgHeartRate)
        summary.maxHeartRate = max(summary.maxHeartRate, session.maxHeartRate)
        summary.avgCadence = avg(summary.avgCadence, session.avgCadence)
        summary.maxCadence = max(summary.maxCadence, session.maxCadence)
        summary.avgPower = avg(summary.avgPower, session.avgPower)
        summary.maxPower = max(summary.maxPower, session.maxPower)
        summary.avgTemperature = avg(summary.avgTemperature, session.avgTemperature)
        summary.maxTemperature = max(summary.maxTemperature, session.maxTemperature)
        summary.avgAltitude = avg(summary.avgAltitude, session.avgAltitude)
        summary.maxAltitude = max(summary.maxAltitude, session.maxAltitude)
        summary.avgPace = avg(summary.avgPace, session.avgPace)
        summary.avgElapsedPace = max(summary.avgElapsedPace, session.avgElapsedPace)
      }

      summary.sport = Array.from(sports).join(', ')

      this.$emit('summary', summary)

      return summary
    }
  },
  methods: {
    avg: avg,
    max: max,
    sum: sum,
    secondsToDHMS: secondsToDHMS,
    formatPace: formatPace,
    toTimezoneDateString: toTimezoneDateString,
    GMTString: GMTString,
    formatCreatedDate(s: string, tz: number): string {
      const d = toTimezone(s, tz)
      return d.setLocale('en-US').toLocaleString({
        localeMatcher: 'best fit',
        weekday: undefined,
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        hour12: true,
        minute: '2-digit',
        second: undefined
      })
    },
    updateSelectedSession(sessions: Session[]) {
      if (sessions.length == 0) {
        this.sessionSelected = NONE
      } else if (sessions.length == 1) {
        this.sessionSelected = 0 // default
      } else {
        this.sessionSelected = MULTIPLE
      }
    }
  },
  mounted() {
    this.updateSelectedSession(this.sessions)
  }
}
</script>

<style scoped>
.title {
  color: var(--color-title);
  font-weight: bold;
}

.time-created {
  color: var(--color-text);
  text-align: left;
  font-size: 0.9em;
}

.summary-title {
  display: inline-block;
  width: 100%;
  font-size: 0.8em;
}

.lr-border {
  border-left: 0.5px solid grey;
  border-right: 0.5px solid grey;
}

.summary-value {
  display: inline-block;
  width: 100%;
}

.summary {
  border-top: 0.5rem solid var(--color-background-soft);
  border-bottom: 0.5rem solid var(--color-background-soft);
}

@media (pointer: coarse) {
  .summary-title {
    font-size: 0.8em;
  }
}
</style>
