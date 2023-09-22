<script setup lang="ts">
import { FileId, Session } from '@/assets/activity'
import moment from 'moment'
</script>

<template>
  <div class="header"><h2 class="title">Open Activity</h2></div>
  <div class="input">
    <input type="file" id="fileInput" accept=".fit" />
  </div>
  <div class="analysis">
    <h3 class="section-title">Summary</h3>
    <div class="summary">
      <div class="summary-grid" v-for="session in sessions">
        <div class="summary-item">
          <div class="summary-title">Sport</div>
          <div class="summary-value">
            {{ session.sport.charAt(0).toUpperCase() + session.sport.slice(1) }}
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Total Elapsed Time</div>
          <div class="summary-value">
            {{ new Date(session.totalElapsedTime * 1000).toISOString().slice(11, 19) }}
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Total Moving Time</div>
          <div class="summary-value">
            {{ new Date(session.totalMovingTime * 1000).toISOString().slice(11, 19) }}
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Total Distance</div>
          <div class="summary-value">{{ (session.totalDistance / 1000).toFixed(2) }} km</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Total Ascent</div>
          <div class="summary-value">{{ session.totalAscent }} m</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Total Descent</div>
          <div class="summary-value">{{ session.totalDescent }} m</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Total Cycles</div>
          <div class="summary-value">{{ session.totalCycles }}</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Total Calories</div>
          <div class="summary-value">{{ session.totalCalories }} Cal</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Avg Speed</div>
          <div class="summary-value">{{ ((session.avgSpeed * 3600) / 1000).toFixed(2) }} km/h</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Max Speed</div>
          <div class="summary-value">{{ ((session.maxSpeed * 3600) / 1000).toFixed(2) }} km/h</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Avg Heart Rate</div>
          <div class="summary-value">{{ session.avgHeartRate }} bpm</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Max Heart Rate</div>
          <div class="summary-value">{{ session.maxHeartRate }} bpm</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Avg Cadence</div>
          <div class="summary-value">{{ session.avgCadence }} rpm</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Max Cadence</div>
          <div class="summary-value">{{ session.maxCadence }} rpm</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Avg Power</div>
          <div class="summary-value">{{ session.avgPower }}</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Max Power</div>
          <div class="summary-value">{{ session.maxPower }}</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Avg Temperature</div>
          <div class="summary-value">{{ session.avgTemperature }} °C</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Max Temperature</div>
          <div class="summary-value">{{ session.maxTemperature }} °C</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Avg Altitude</div>
          <div class="summary-value">{{ session.avgAltitude.toFixed(2) }} m</div>
        </div>
        <div class="summary-item">
          <div class="summary-title">Max Altitude</div>
          <div class="summary-value">{{ session.maxAltitude.toFixed(2) }} m</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    sessions: Array<Session>
  },
  data() {
    return {
      // distance: Session[],
    }
  },
  watch: {},
  methods: {},
  mounted() {}
}
</script>

<style>
.title {
  color: var(--color-title);
  font-weight: bold;
}

.header {
  text-align: center;
}

.input {
  color: var(--color-title);
  margin-top: 1rem;
}

.analysis {
  color: var(--color-title);
  margin-top: 1rem;
}

.section-title {
  text-align: center;
}

.summary-title {
  font-size: 9px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr); /* Two columns */
  gap: 10px; /* Adjust the gap between columns as needed */
}

.summary-item {
  padding: 10px;
  text-align: center;
}
</style>
