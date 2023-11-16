<template>
  <!-- input -->
  <div class="navigator">
    <div class="navigator-info">
      <div class="time-created" v-if="activityFiles && activityFiles.length != 0">
        Created on:
        {{ toTimezoneDateString(activityFiles[0]?.creator?.timeCreated, timezoneOffsetHour) }}
        {{ GMTString(timezoneOffsetHour) }}
      </div>
      <div class="manufacturer" v-if="activityFiles && activityFiles.length != 0">
        Device: {{ activityFiles[0]?.creator?.name }}
      </div>
    </div>
    <div class="navigator-summary analysis">
      <div class="summary">
        <h4 class="section-title" v-if="activityFiles?.length != 0">Summary</h4>
        <div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.sport">
              <div class="summary-title">Sport</div>
              <div class="summary-value">{{ summary.sport }}</div>
            </div>
            <div class="summary-item" v-if="summary.subSport">
              <div class="summary-title">Sub Sport</div>
              <div class="summary-value">{{ summary.subSport }}</div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.totalMovingTime">
              <div class="summary-title">Total Moving Time</div>
              <div class="summary-value">
                <i class="fa-solid fa-hourglass-half"></i>
                {{ secondsToDHMS(summary.totalMovingTime) }}
              </div>
            </div>
            <div class="summary-item" v-if="summary.totalElapsedTime">
              <div class="summary-title">Total Elapsed Time</div>
              <div class="summary-value">
                <i class="fa-solid fa-hourglass-end"></i>
                {{ secondsToDHMS(summary.totalElapsedTime) }}
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.totalDistance">
              <div class="summary-title">Total Distance</div>
              <div class="summary-value">
                <i class="fa-solid fa-road"></i> {{ (summary.totalDistance / 1000)?.toFixed(2) }} km
              </div>
            </div>
            <div class="summary-item" v-if="summary.totalCalories">
              <div class="summary-title">Total Calories</div>
              <div class="summary-value">
                <i class="fa-solid fa-droplet"></i>
                {{ summary.totalCalories.toLocaleString() }} Cal
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.avgSpeed">
              <div class="summary-title">Avg Speed</div>
              <div class="summary-value">
                <i class="fa-solid fa-gauge"></i>
                {{ ((summary.avgSpeed * 3600) / 1000)?.toFixed(2) }} km/h
              </div>
            </div>
            <div class="summary-item" v-if="summary.maxSpeed">
              <div class="summary-title">Max Speed</div>
              <div class="summary-value">
                <i class="fa-solid fa-gauge-high"></i>
                {{ ((summary.maxSpeed * 3600) / 1000)?.toFixed(2) }} km/h
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.avgHeartRate">
              <div class="summary-title">Avg Heart Rate</div>
              <div class="summary-value">
                <i class="fa-solid fa-heart-pulse"></i> {{ Math.round(summary.avgHeartRate) }} bpm
              </div>
            </div>
            <div class="summary-item" v-if="summary.maxHeartRate">
              <div class="summary-title">Max Heart Rate</div>
              <div class="summary-value">
                <i class="fa-solid fa-heart-pulse"></i> {{ summary.maxHeartRate }} bpm
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.avgCadence">
              <div class="summary-title">Avg Cadence</div>
              <div class="summary-value">
                <i class="fa-solid fa-rotate"></i> {{ Math.round(summary.avgCadence) }} rpm
              </div>
            </div>
            <div class="summary-item" v-if="summary.maxCadence">
              <div class="summary-title">Max Cadence</div>
              <div class="summary-value">
                <i class="fa-solid fa-rotate"></i> {{ summary.maxCadence }} rpm
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.avgPower">
              <div class="summary-title">Avg Power</div>
              <div class="summary-value">
                <i class="fa-solid fa-bolt-lightning"></i> {{ Math.round(summary.avgPower) }}
              </div>
            </div>
            <div class="summary-item" v-if="summary.maxPower">
              <div class="summary-title">Max Power</div>
              <div class="summary-value">
                <i class="fa-solid fa-bolt-lightning"></i> {{ summary.maxPower }}
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.avgTemperature">
              <div class="summary-title">Avg Temperature</div>
              <div class="summary-value">
                <i class="fa-solid fa-temperature-low"></i>
                {{ Math.round(summary.avgTemperature) }} °C
              </div>
            </div>
            <div class="summary-item" v-if="summary.maxTemperature">
              <div class="summary-title">Max Temperature</div>
              <div class="summary-value">
                <i class="fa-solid fa-temperature-high"></i> {{ summary.maxTemperature }} °C
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.totalAscent">
              <div class="summary-title">Total Ascent</div>
              <div class="summary-value">
                <i class="fa-solid fa-arrow-trend-up"></i> {{ summary.totalAscent }} m
              </div>
            </div>
            <div class="summary-item" v-if="summary.totalDescent">
              <div class="summary-title">Total Descent</div>
              <div class="summary-value">
                <i class="fa-solid fa-arrow-trend-down"></i> {{ summary.totalDescent }} m
              </div>
            </div>
          </div>
          <div class="summary-grid">
            <div class="summary-item" v-if="summary.avgAltitude">
              <div class="summary-title">Avg Altitude</div>
              <div class="summary-value">
                <i class="fa-solid fa-mountain"></i> {{ summary.avgAltitude.toFixed(2) }} m
              </div>
            </div>
            <div class="summary-item" v-if="summary.maxAltitude">
              <div class="summary-title">Max Altitude</div>
              <div class="summary-value">
                <i class="fa-solid fa-mountain"></i> {{ summary.maxAltitude.toFixed(2) }} m
              </div>
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
          summary.totalMovingTime = sum(summary.totalMovingTime, session.totalMovingTime)
          summary.totalElapsedTime = sum(summary.totalElapsedTime, session.totalElapsedTime)
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
  color: var(--color-title);
  text-align: center;
  font-size: 10px;
}

.manufacturer {
  color: var(--color-text);
  text-align: center;
  font-size: 12px;
}

.analysis {
  color: var(--color-title);
  margin-top: 10px;
  margin-bottom: 30px;
}

.section-title {
  text-align: center;
}

.summary {
  max-width: 300px;
  margin: auto;
}

.summary-title {
  font-size: 10px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr); /* Two columns */
}

.summary-item {
  padding: 10px;
  text-align: center;
}
</style>
