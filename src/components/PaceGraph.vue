<script setup lang="ts">
import { formatPace } from '@/toolkit/pace'
import AreaGraph, { Detail } from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Pace'"
      :icon="'fa-clock-rotate-left'"
      :record-field="'pace'"
      :records="records"
      :details="details"
      :graph-records="graphRecords"
      :summary="summary"
      :color="'dodgerblue'"
      :y-label="'Pace (duration/km)'"
      :unit="'/km'"
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
          title: 'Avg Pace',
          value: this.summary.avgPace ? formatPace(this.summary.avgPace) : '-:-'
        }),
        new Detail({
          title: 'Avg Elapsed Pace',
          value: this.summary.avgElapsedPace ? formatPace(this.summary.avgElapsedPace) : '-:-'
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
