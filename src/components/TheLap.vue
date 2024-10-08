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
  <div class="col-12 h-100" v-if="laps?.length > 0">
    <div v-for="(lap, index) in laps" :key="index">
      <div
        class="tab row m-0 mt-2 pt-2 collapsible"
        style="cursor: pointer"
        data-bs-toggle="collapse"
        v-bind:data-bs-target="'#laps-' + index"
        aria-expanded="false"
        v-bind:aria-controls="'laps-' + index"
      >
        <div class="row text-start">
          <div class="col-auto d-inline-block" style="height: 50px">
            <h6 style="text-align: left" class="mb-0">
              <i class="fa-solid fa-caret-right collapse-indicator"></i>
              <span class="px-1">Lap {{ index + 1 }}</span>
            </h6>
            <span>{{ lap.sport }}</span>
          </div>
          <div class="col">
            <div class="row overview-title">Distance</div>
            <div class="row overview-value">
              {{ lap.totalDistance ? (lap.totalDistance / 1000).toFixed(2) : '-' }} km
            </div>
          </div>
          <div class="col">
            <div class="row overview-title">Moving Time</div>
            <div class="row overview-value">
              {{ lap.totalMovingTime ? secondsToDHMS(lap.totalMovingTime!) : '-:-' }}
            </div>
          </div>
          <div class="col">
            <div class="row overview-title">Avg Cadence</div>
            <div class="row overview-value">{{ lap.avgCadence ? lap.avgCadence : '-' }} rpm</div>
          </div>
        </div>
      </div>
      <div class="collapse show text-start pb-3" v-bind:id="'laps-' + index">
        <div class="row m-0">
          <div class="col ps-4 pe-0">
            <div>&nbsp;</div>
            <div class="right-border">
              <div class="detail-title">Speed</div>
              <div class="detail-title">Cadence</div>
              <div class="detail-title">Heart Rate</div>
              <div class="detail-title">Power</div>
            </div>
          </div>
          <div class="col fw-bold text-center px-0">
            <div class="fw-normal">Avg</div>
            <div class="right-border">
              <div class="detail-value">
                {{ lap.avgSpeed ? ((lap.avgSpeed * 3600) / 1000).toFixed(2) : '-' }}
              </div>
              <div class="detail-value">{{ lap.avgCadence ? lap.avgCadence : '-' }}</div>
              <div class="detail-value">{{ lap.avgHeartRate ? lap.avgHeartRate : '-' }}</div>
              <div class="detail-value">{{ lap.avgPower ? lap.avgPower : '-' }}</div>
            </div>
          </div>
          <div class="col fw-bold text-center px-0">
            <div class="fw-normal">Max</div>
            <div class="right-border">
              <div class="detail-value">
                {{ lap.maxSpeed ? ((lap.maxSpeed * 3600) / 1000).toFixed(2) : '-' }}
              </div>
              <div class="detail-value">{{ lap.maxCadence ? lap.maxCadence : '-' }}</div>
              <div class="detail-value">{{ lap.maxHeartRate ? lap.maxHeartRate : '-' }}</div>
              <div class="detail-value">{{ lap.maxPower ? lap.maxPower : '-' }}</div>
            </div>
          </div>
          <div class="col m-0">
            <div>&nbsp;</div>
            <div>
              <div class="ps-2 detail-unit">km/h</div>
              <div class="ps-2 detail-unit">rpm</div>
              <div class="ps-2 detail-unit">bpm</div>
              <div class="ps-2 detail-unit">W</div>
            </div>
          </div>
        </div>
        <div class="calories mt-1">
          <div class="row m-0 pt-1 pb-1">
            <div class="col ps-4 m-0 text-start detail-title">Calories</div>
            <div class="col pe-4 m-0 text-end">
              <span class="fw-bold pe-1">{{
                lap.totalCalories ? lap.totalCalories.toLocaleString() : '-'
              }}</span>
              <span class="detail-unit">Cal</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Lap } from '@/spec/activity'
import { secondsToDHMS } from '@/toolkit/date'

export default {
  props: {
    laps: {
      type: Array<Lap>,
      required: true
    }
  },
  methods: {
    secondsToDHMS: secondsToDHMS
  }
}
</script>
<style scoped>
.tab {
  background-color: var(--color-background-mute);
}

.overview-title {
  font-size: 0.8em;
}

.overview-value {
  font-size: 1, 2em;
  font-weight: 700;
}

.detail-title,
.detail-unit {
  font-size: 0.9em;
}

.detail-value {
  font-size: 1em;
}

.right-border {
  border-right: 1px solid grey;
  height: 90px;
}

.calories {
  background-color: var(--color-background-soft);
}
</style>
