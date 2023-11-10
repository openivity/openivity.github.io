<script setup lang="ts">
import LineGraph from './LineGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <LineGraph
      :name="'Cadence'"
      :icon="'fa-rotate'"
      :record-field="'cadence'"
      :avg="summary?.avgCadence?.toFixed(0) ?? 0"
      :max="summary?.maxCadence?.toFixed(0) ?? 0"
      :records="records"
      :graph-records="graphRecords"
      :summary="summary"
      :color="'darkslateblue'"
      :y-label="'Cad. (rpm)'"
      :unit="'rpm'"
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
