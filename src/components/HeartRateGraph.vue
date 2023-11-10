<script setup lang="ts">
import LineGraph from './LineGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <LineGraph
      :name="'Heart Rate'"
      :icon="'fa-heart-pulse'"
      :selected-field="'heartRate'"
      :records="records"
      :graph-records="graphRecords"
      :summary="summary"
      :color="'red'"
      :y-label="'HR (bpm)'"
      :unit="'bpm'"
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
