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
    },
    onHoveredRecordFreeze(freeze: Boolean) {
      this.$emit('hoveredRecordFreeze', freeze)
    }
  }
}
</script>
