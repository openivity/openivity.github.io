<script setup lang="ts">
import LineGraph from './LineGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <LineGraph
      :name="'Temperature'"
      :icon="'fa-temperature-low'"
      :record-field="'temperature'"
      :avg="summary?.avgTemperature?.toFixed(0) ?? 0"
      :max="summary?.maxTemperature?.toFixed(0) ?? 0"
      :records="records"
      :graph-records="graphRecords"
      :summary="summary"
      :color="'midnightblue'"
      :y-label="'Temp (°C)'"
      :unit="'°C'"
      :received-record="receivedRecord"
      v-on:hoveredRecord="onHoveredRecord"
    ></LineGraph>
  </div>
</template>

<script lang="ts">
import { Record } from '@/spec/activity'
import { Summary } from '@/spec/summary'

export default {
  props: {
    graphRecords: {
      type: Array<Record>,
      required: true,
      default: []
    },
    records: {
      type: Array<Record>,
      required: true,
      default: []
    },
    color: String,
    summary: Summary,
    receivedRecord: Record
  },
  methods: {
    onHoveredRecord(record: Record) {
      this.$emit('hoveredRecord', record)
    }
  }
}
</script>
