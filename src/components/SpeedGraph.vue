<script setup lang="ts">
import LineGraph from './LineGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <LineGraph
      :name="'Speed'"
      :icon="'fa-gauge-high'"
      :record-field="'speed'"
      :avg="summary?.avgSpeed?.toFixed(2) ?? 0"
      :max="summary?.maxSpeed?.toFixed(2) ?? 0"
      :records="records"
      :graph-records="graphRecords"
      :summary="summary"
      :color="'lightgreen'"
      :pointer-color="'#7f8c8d'"
      :y-label="'Spd. (km/h)'"
      :unit="'km/h'"
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
