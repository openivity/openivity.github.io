<script setup lang="ts">
import AreaGraph, { Detail } from './AreaGraph.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2">
    <AreaGraph
      :name="'Power'"
      :icon="'fa-bolt-lightning'"
      :record-field="'power'"
      :records="records"
      :details="details"
      :graph-records="graphRecords"
      :color="'darkslategray'"
      :y-label="'Pwr. (W)'"
      :unit="'W'"
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
          title: 'Avg Power',
          value: this.summary.avgPower?.toFixed(0) ?? '0'
        }),
        new Detail({
          title: 'Max Power',
          value: this.summary.maxPower?.toFixed(0) ?? '0'
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
