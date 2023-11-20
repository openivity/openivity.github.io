<script setup lang="ts">
import AreaGraph, { Detail } from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Cadence'"
      :icon="'fa-rotate'"
      :record-field="'cadence'"
      :records="records"
      :details="details"
      :graph-records="graphRecords"
      :color="'darkslateblue'"
      :y-label="'Cad. (rpm)'"
      :unit="'rpm'"
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
          title: 'Avg Cadence',
          value: this.summary.avgCadence?.toFixed(0) ?? '0'
        }),
        new Detail({
          title: 'Max Cadence',
          value: this.summary.maxCadence?.toFixed(0) ?? '0'
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
