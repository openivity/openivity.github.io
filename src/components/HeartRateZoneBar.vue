<script setup lang="ts">
import { Duration } from 'luxon'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <div class="row">
      <div class="col d-flex flex-row pb-1">
        <h6 class="p-0 m-0">{{ zone }}</h6>
        <small class="text-body-secondary ms-2 me-1">{{ validMinMax }} bpm</small>
        <small class="text-body-secondary">&#183; {{ zoneSub }}</small>
      </div>
    </div>
    <div class="row gx-0">
      <div class="col-9">
        <div
          class="progress"
          role="progressbar"
          aria-label="Success example"
          aria-valuenow="{{ validProsen }}"
          aria-valuemin="0"
          aria-valuemax="100"
        >
          <div :class="['progress-bar', ...progressClass]" :style="{ width: `${validProsen}%` }">
            {{ validProsen }}%
          </div>
        </div>
      </div>
      <div class="col-3 d-flex flex-row">
        <small class="text-body-secondary flex-fill text-end hr-time">{{ formattedTime }}</small>
        <small class="text-body-secondary text-end hr-prosen">{{ validProsen }}%</small>
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
    progressClass: {
      type: Array<String>,
      default: []
    }
  },
  data() {},
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
  font-weight: bold;
}
.hr-prosen {
  width: 40px;
}
</style>
