<script setup lang="ts">
import AreaGraph from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Speed'"
      :icon="'fa-gauge-high'"
      :record-field="'speed'"
      :avg="summary?.avgSpeed ? ((summary?.avgSpeed! * 3600) / 1000).toFixed(2) : '0.00'"
      :max="summary?.maxSpeed ? ((summary?.maxSpeed! * 3600) / 1000).toFixed(2) : '0.00'"
      :records="records"
      :graph-records="graphRecords"
      :summary="summary"
      :color="'lightgreen'"
      :pointer-color="'#7f8c8d'"
      :y-label="'Spd. (km/h)'"
      :unit="'km/h'"
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
