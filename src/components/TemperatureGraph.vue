<script setup lang="ts">
import AreaGraph, { Detail } from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Temperature'"
      :icon="'fa-temperature-low'"
      :record-field="'temperature'"
      :records="records"
      :details="details"
      :graph-records="graphRecords"
      :color="'midnightblue'"
      :y-label="'Temp (°C)'"
      :unit="'°C'"
      :received-record="receivedRecord"
      v-on:hoveredRecord="onHoveredRecord"
    ></AreaGraph>
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
    summary: {
      type: Summary,
      required: true
    },
    receivedRecord: Record
  },
  computed: {
    details(): Detail[] {
      return [
        new Detail({
          title: 'Avg Temperature',
          value: this.summary.avgTemperature?.toFixed(0) ?? '0'
        }),
        new Detail({
          title: 'Max Temperature',
          value: this.summary.maxTemperature?.toFixed(0) ?? '0'
        })
      ]
    }
  },
  methods: {
    onHoveredRecord(record: Record) {
      this.$emit('hoveredRecord', record)
    }
  }
}
</script>
