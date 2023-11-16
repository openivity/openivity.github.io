<template>
  <!-- input -->
  <div class="summary mt-4" v-if="isActivityFileReady">
    <div class="px-3 pt-3 pb-4">
      <div class="time-created pb-2">
        <i class="fa-solid fa-clock"></i>&nbsp;
        {{ formatCreatedDate(activityFiles![0].creator?.timeCreated!, activityFiles![0].timezone) }}
        &nbsp;|
        {{
          activityFiles![0].timezone > 0
            ? '+' + activityFiles![0].timezone
            : activityFiles![0].timezone
        }}
      </div>
      <div class="row" style="font-size: 1em; color: var(--bs-heading-color)">
        <div class="col text-start">
          <span> Sport: {{ summary.sport ?? 'Unknown' }} </span>
        </div>
        <div class="col text-end">
          <i class="fa-solid fa-microchip"></i>
          {{ activityFiles![0]?.creator?.name }}
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
            <div class="summary-value fs-6">{{ summary.totalAscent ?? '-' }} m</div>
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
          <div>
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
import { ActivityFile } from '@/spec/activity'
import { Summary } from '@/spec/summary'
import {
  GMTString,
  secondsToDHMS,
  toTimezone,
  toTimezoneDate,
  toTimezoneDateString
} from '@/toolkit/date'
import { avg, max, sum } from '@/toolkit/number'
export default {
  props: {
    isActivityFileReady: Boolean,
    activityFiles: {
      type: Array<ActivityFile>,
      default: []
    }
  },
  computed: {
    summary(): Summary {
      const summary = new Summary()
      let sports: Set<string> = new Set()

      for (let i = 0; i < this.activityFiles.length; i++) {
        let gapSeconds: number = 0
        if (i > 0) {
          let lastTimestamp: Date = new Date()
          const lastTimezone = this.activityFiles[i - 1].timezone
          for (let j = this.activityFiles[i - 1].records.length - 1; j >= 0; j--) {
            const records = this.activityFiles[i - 1].records
            if (records[j].timestamp == null) continue
            lastTimestamp = toTimezoneDate(records[j].timestamp!, lastTimezone)
            break
          }

          let currentTimestamp: Date = lastTimestamp
          const currentTimezone = this.activityFiles[i].timezone
          for (let j = 0; j < this.activityFiles[i].records.length; j++) {
            const records = this.activityFiles[i].records
            if (records[j].timestamp == null) continue
            currentTimestamp = toTimezoneDate(records[j].timestamp!, currentTimezone)
            break
          }

          gapSeconds = (currentTimestamp!.getTime() - lastTimestamp!.getTime()) / 1000
        }

        for (let j = 0; j < this.activityFiles![i].sessions!.length; j++) {
          const session = this.activityFiles![i].sessions![j]

          session.sport = session.sport ? session.sport : 'Unknown'
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
        }
      }

      summary.sport = Array.from(sports).join(', ')

      this.$emit('summary', summary)
      return summary
    }
  },
  methods: {
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
    toTimezoneDateString: toTimezoneDateString,
    GMTString: GMTString,
    secondsToDHMS: secondsToDHMS,
    avg: avg,
    max: max,
    sum: sum
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
