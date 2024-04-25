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

<script setup lang="ts">
import { Duration } from 'luxon'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <div class="row">
      <div class="col d-flex align-items-end pb-1">
        <h6 class="zone-label fw-bolder text-body-secondary mb-0">{{ zone }}</h6>
        <span class="zone-label small text-body-secondary ms-2 me-1">{{ validMinMax }} bpm</span>
        <span class="zone-label small text-body-secondary">&#183; {{ zoneSub }}</span>
      </div>
    </div>
    <div class="row">
      <div class="col d-flex flex-row">
        <div
          class="progress flex-fill"
          role="progressbar"
          aria-label=""
          aria-valuenow="{{ validProsen }}"
          aria-valuemin="0"
          aria-valuemax="100"
        >
          <div
            :class="[
              'progress-bar',
              ...progressClass,
              isLoading ? 'progress-bar-striped progress-bar-animated' : ''
            ]"
            :style="{ width: `${isLoading ? 100 : validProsen.toFixed(0)}%` }"
          >
            <span v-if="progressText"> {{ validProsen.toFixed(1) }}% </span>
          </div>
        </div>
        <small class="text-body-secondary text-end hr-time">{{ formattedTime }}</small>
        <small class="text-body-secondary text-end hr-prosen">{{ validProsen.toFixed(1) }}%</small>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    zone: {
      type: String,
      default: 'Zone 1'
    },
    minmax: {
      type: Array<Number>,
      required: true
    },
    zoneSub: {
      type: String,
      default: 'Maximum'
    },
    timeInSecond: {
      type: Number,
      default: 0
    },
    prosen: {
      type: Number,
      default: 0
    },
    progressText: {
      type: Boolean
    },
    progressClass: {
      type: Array<String>,
      default: []
    },
    isLoading: {
      type: Boolean
    }
  },
  data() {
    return {}
  },
  watch: {},
  computed: {
    formattedTime(): String {
      if (this.timeInSecond >= 60 * 60)
        return Duration.fromMillis(this.timeInSecond * 1000).toFormat('h:mm:ss')
      else return Duration.fromMillis(this.timeInSecond * 1000).toFormat('mm:ss')
    },
    validProsen(): Number {
      return this.prosen >= 0 && this.prosen <= 100 ? this.prosen : this.prosen > 100 ? 100 : 0 // invalid number will be 0
    },
    validMinMax(): String {
      if (this.max == Infinity) return `> ${this.min}`
      else if (this.min == Infinity) return `< ${this.max}`
      return `${this.min} - ${this.max}`
    },
    min(): Number {
      if (this.minmax[0] == Infinity) return Infinity
      return this.minmax[0] || 0
    },
    max(): Number {
      if (this.minmax[1] == Infinity) return Infinity
      return this.minmax[1] || 0
    }
  },
  methods: {},
  mounted() {}
}
</script>
<style lang="scss" scoped>
.progress {
  height: 15px;
  margin-top: 5px;
}
.hr-time {
  width: 60px;
  font-weight: bold;
}
.hr-prosen {
  width: 45px;
}

.zone-label {
  line-height: 80%;
  align-self: baseline;
}
</style>
