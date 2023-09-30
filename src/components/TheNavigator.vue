<template>
  <div style="text-align: center">
    <input class="input" type="file" id="fileInput" accept=".fit" />
  </div>
  <div class="time-created" v-if="activityFile?.fileId?.timeCreated">
    Created on:
    {{ toTimezoneDateString(activityFile.fileId.timeCreated, timezoneOffsetHours) }}
    {{ GMTString(timezoneOffsetHours) }}
  </div>
  <div class="manufacturer" v-if="activityFile?.fileId">
    Device: {{ activityFile.fileId.manufacturer }} ({{ activityFile.fileId.product }})
  </div>
  <div class="analysis">
    <div class="summary">
      <h4 class="section-title" v-if="activityFile?.fileId">Summary</h4>
      <div v-for="(session, index) in activityFile?.sessions" v-bind:key="index">
        <div class="summary-grid">
          <div class="summary-item" v-if="session.sport">
            <div class="summary-title">Sport</div>
            <div class="summary-value">{{ session.sport }}</div>
          </div>
          <div class="summary-item" v-if="session.subSport">
            <div class="summary-title">Sub Sport</div>
            <div class="summary-value">{{ session.subSport }}</div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.totalElapsedTime">
            <div class="summary-title">Total Elapsed Time</div>
            <div class="summary-value">
              <i class="fa-solid fa-hourglass-end"></i>
              {{ new Date(session.totalElapsedTime * 1000).toISOString().slice(11, 19) }}
            </div>
          </div>
          <div class="summary-item" v-if="session.totalMovingTime">
            <div class="summary-title">Total Moving Time</div>
            <div class="summary-value">
              <i class="fa-solid fa-hourglass-half"></i>
              {{ new Date(session.totalMovingTime * 1000).toISOString().slice(11, 19) }}
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.totalDistance">
            <div class="summary-title">Total Distance</div>
            <div class="summary-value">
              <i class="fa-solid fa-road"></i> {{ (session.totalDistance / 1000)?.toFixed(2) }} km
            </div>
          </div>
          <div class="summary-item" v-if="session.totalCalories">
            <div class="summary-title">Total Calories</div>
            <div class="summary-value">
              <i class="fa-solid fa-droplet"></i>
              {{ session.totalCalories.toLocaleString() }} Cal
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.avgSpeed">
            <div class="summary-title">Avg Speed</div>
            <div class="summary-value">
              <i class="fa-solid fa-gauge"></i>
              {{ ((session.avgSpeed * 3600) / 1000)?.toFixed(2) }} km/h
            </div>
          </div>
          <div class="summary-item" v-if="session.maxSpeed">
            <div class="summary-title">Max Speed</div>
            <div class="summary-value">
              <i class="fa-solid fa-gauge-high"></i>
              {{ ((session.maxSpeed * 3600) / 1000)?.toFixed(2) }} km/h
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.avgHeartRate">
            <div class="summary-title">Avg Heart Rate</div>
            <div class="summary-value">
              <i class="fa-solid fa-heart-pulse"></i> {{ session.avgHeartRate }} bpm
            </div>
          </div>
          <div class="summary-item" v-if="session.maxHeartRate">
            <div class="summary-title">Max Heart Rate</div>
            <div class="summary-value">
              <i class="fa-solid fa-heart-pulse"></i> {{ session.maxHeartRate }} bpm
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.avgCadence">
            <div class="summary-title">Avg Cadence</div>
            <div class="summary-value">
              <i class="fa-solid fa-rotate"></i> {{ session.avgCadence }} rpm
            </div>
          </div>
          <div class="summary-item" v-if="session.maxCadence">
            <div class="summary-title">Max Cadence</div>
            <div class="summary-value">
              <i class="fa-solid fa-rotate"></i> {{ session.maxCadence }} rpm
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.avgPower">
            <div class="summary-title">Avg Power</div>
            <div class="summary-value">
              {<i class="fa-solid fa-bolt-lightning"></i> { session.avgPower }}
            </div>
          </div>
          <div class="summary-item" v-if="session.maxPower">
            <div class="summary-title">Max Power</div>
            <div class="summary-value">
              <i class="fa-solid fa-bolt-lightning"></i> {{ session.maxPower }}
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.avgTemperature">
            <div class="summary-title">Avg Temperature</div>
            <div class="summary-value">
              <i class="fa-solid fa-temperature-low"></i> {{ session.avgTemperature }} °C
            </div>
          </div>
          <div class="summary-item" v-if="session.maxTemperature">
            <div class="summary-title">Max Temperature</div>
            <div class="summary-value">
              <i class="fa-solid fa-temperature-high"></i> {{ session.maxTemperature }} °C
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.totalAscent">
            <div class="summary-title">Total Ascent</div>
            <div class="summary-value">
              <i class="fa-solid fa-arrow-trend-up"></i> {{ session.totalAscent }} m
            </div>
          </div>
          <div class="summary-item" v-if="session.totalDescent">
            <div class="summary-title">Total Descent</div>
            <div class="summary-value">
              <i class="fa-solid fa-arrow-trend-down"></i> {{ session.totalDescent }} m
            </div>
          </div>
        </div>
        <div class="summary-grid">
          <div class="summary-item" v-if="session.avgAltitude">
            <div class="summary-title">Avg Altitude</div>
            <div class="summary-value">
              <i class="fa-solid fa-mountain"></i> {{ session.avgAltitude.toFixed(2) }} m
            </div>
          </div>
          <div class="summary-item" v-if="session.maxAltitude">
            <div class="summary-title">Max Altitude</div>
            <div class="summary-value">
              <i class="fa-solid fa-mountain"></i> {{ session.maxAltitude.toFixed(2) }} m
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ActivityFile } from '@/spec/activity'
import { GMTString, toTimezoneDateString } from '@/toolkit/date'
export default {
  props: {
    activityFile: ActivityFile,
    timezoneOffsetHours: Number
  },
  methods: {
    toTimezoneDateString: toTimezoneDateString,
    GMTString: GMTString
  }
}
</script>

<style>
.title {
  color: var(--color-title);
  font-weight: bold;
}

.input {
  color: var(--color-title);
  display: inline-block;
  padding: 15px 20px;
  background-color: var(--color-background-mute);
  color: var(--color-text);
  border-radius: 5px;
  cursor: pointer;
  margin: 0 0 5px 0;
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
}

.section-title {
  text-align: center;
}

.summary {
  padding: 0px 20px;
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
