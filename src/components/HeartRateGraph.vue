<script setup lang="ts">
import AreaGraph, { Detail } from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Heart Rate'"
      :icon="'fa-heart-pulse'"
      :record-field="'heartRate'"
      :records="records"
      :details="details"
      :graph-records="graphRecords"
      :color="'red'"
      :y-label="'HR (bpm)'"
      :unit="'bpm'"
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
          title: 'Avg Heart Rate',
          value: this.summary.avgHeartRate?.toFixed(0) ?? '0'
        }),
        new Detail({
          title: 'Max Heart Rate',
          value: this.summary.maxHeartRate?.toFixed(0) ?? '0'
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
