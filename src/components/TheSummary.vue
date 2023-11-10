<template>
  <!-- input -->
  <div class="navigator">
    <div class="navigator-info">
      <div class="time-created" v-if="activityFiles && activityFiles.length != 0">
        Created on:
        {{ toTimezoneDateString(activityFiles[0]?.fileId?.timeCreated, timezoneOffsetHour) }}
        {{ GMTString(timezoneOffsetHour) }}
      </div>
      <div class="row pt-4" style="font-size: 1em">
        <div class="col text-start" v-if="summary.sport">
          <span>
            Sport: {{ summary.sport }}
            <span v-if="summary.subSport">({{ summary.subSport }})</span>
          </span>
        </div>
        <div class="col text-end" v-if="activityFiles && activityFiles.length != 0">
          Device: {{ activityFiles[0]?.fileId?.manufacturer }} ({{
            activityFiles[0]?.fileId?.product
          }})
        </div>
      </div>
    </div>
    <div class="row pt-2">
      <div class="col" v-if="summary.totalElapsedTime">
        <div class="summary-title">Total Elapsed</div>
        <div class="summary-value">
          <i class="fa-solid fa-hourglass-end"></i>
          {{ secondsToDHMS(summary.totalElapsedTime) }}
        </div>
      </div>
      <div class="col" v-if="summary.totalMovingTime">
        <div class="row">
          <div class="summary-title">Total Moving</div>
          <div class="summary-value">
            <i class="fa-solid fa-hourglass-half"></i>
            {{ secondsToDHMS(summary.totalMovingTime) }}
          </div>
        </div>
      </div>
      <div class="col" v-if="summary.totalDistance">
        <div class="summary-title">Total Distance</div>
        <div class="summary-value">
          <i class="fa-solid fa-road"></i> {{ (summary.totalDistance / 1000)?.toFixed(2) }} km
        </div>
      </div>
      <div class="col" v-if="summary.totalCalories">
        <div class="summary-title">Total Calories</div>
        <div class="summary-value">
          <i class="fa-solid fa-droplet"></i>
          {{ summary.totalCalories.toLocaleString() }} Cal
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ActivityFile } from '@/spec/activity'
import { Summary } from '@/spec/summary'
import { GMTString, secondsToDHMS, toTimezoneDateString } from '@/toolkit/date'
import { avg, max, sum } from '@/toolkit/number'
export default {
  props: {
    activityFiles: Array<ActivityFile>,
    timezoneOffsetHour: Number
  },
  computed: {
    summary(): Summary {
      const summary = new Summary()
      for (let i = 0; i < this.activityFiles!.length; i++) {
        for (let j = 0; j < this.activityFiles![i].sessions!.length; j++) {
          const session = this.activityFiles![i].sessions![j]

          summary.sport = session.sport
          summary.subSport = session.subSport
          summary.totalMovingTime = sum(summary.totalMovingTime, session.totalMovingTime)
          summary.totalElapsedTime = sum(summary.totalElapsedTime, session.totalElapsedTime)
          summary.totalTimerTime = sum(summary.totalTimerTime, session.totalTimerTime)
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

      this.$emit('summary', summary)
      return summary
    }
  },
  methods: {
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
  text-align: center;
  font-size: 10px;
}

.summary-title {
  font-size: 0.8em;
}

.summary-value {
  font-size: 0.8em;
}

.navigator {
  margin-bottom: 10px;
}

.navigator-info {
  color: var(--color-title);
}

@media (pointer: coarse) {
  .summary-title {
    font-size: 0.7em;
  }

  .summary-value {
    font-size: 0.8em;
  }
}
</style>
