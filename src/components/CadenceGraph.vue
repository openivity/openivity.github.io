<!-- Copyright (C) 2023 Openivity

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. -->

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
    receivedRecordFreeze: Boolean
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
    },
    onHoveredRecordFreeze(freeze: Boolean) {
      this.$emit('hoveredRecordFreeze', freeze)
    }
  }
}
</script>
