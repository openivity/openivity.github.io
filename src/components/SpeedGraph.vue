<script setup lang="ts">
import AreaGraph, { Detail } from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Speed'"
      :icon="'fa-gauge-high'"
      :record-field="'speed'"
      :records="records"
      :details="details"
      :graph-records="graphRecords"
      :color="'lightgreen'"
      :pointer-color="'#7f8c8d'"
      :y-label="'Spd. (km/h)'"
      :unit="'km/h'"
      :received-record="receivedRecord"
      :received-record-freeze="receivedRecordFreeze"
      v-on:hoveredRecord="onHoveredRecord"
      v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
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
    receivedRecord: Record,
    receivedRecordFreeze: Boolean,
  },
  computed: {
    details(): Detail[] {
      return [
        new Detail({
          title: 'Avg Speed',
          value: this.summary?.avgSpeed
            ? ((this.summary?.avgSpeed! * 3600) / 1000).toFixed(2)
            : '0.00'
        }),
        new Detail({
          title: 'Max Speed',
          value: this.summary?.maxSpeed
            ? ((this.summary?.maxSpeed! * 3600) / 1000).toFixed(2)
            : '0.00'
        })
      ]
    }
  },
  methods: {
    onHoveredRecord(record: Record) {
      this.$emit('hoveredRecord', record)
    },
    onHoveredRecordFreeze(freeze: Boolean) {
      this.$emit('hoveredRecordFreeze', freeze)
    }
  }
}
</script>
