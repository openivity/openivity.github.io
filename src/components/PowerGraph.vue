<script setup lang="ts">
import AreaGraph from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Power'"
      :icon="'fa-bolt-lightning'"
      :record-field="'power'"
      :avg="summary?.avgPower?.toFixed(0) ?? 0"
      :max="summary?.maxPower?.toFixed(0) ?? 0"
      :records="records"
      :graph-records="graphRecords"
      :summary="summary"
      :color="'darkslategray'"
      :y-label="'Pwr. (watts)'"
      :unit="'w'"
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
