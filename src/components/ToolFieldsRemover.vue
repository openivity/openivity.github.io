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
import type { Record, Session } from '@/spec/activity'
import { ToolMode } from '@/spec/activity-service'
import type { PropType } from 'vue'
</script>
<template>
  <div>
    <div
      class="row m-0"
      style="cursor: pointer"
      data-bs-toggle="collapse"
      data-bs-target="#fieldsRemovalTarget"
      aria-expanded="false"
      aria-controls="fieldsRemovalTarget"
    >
      <div class="text-start p-0">
        <label class="pe-1">Remove Fields</label>
        <i class="fa-regular fa-circle-question" title="Show or Hide Help Text"></i>
      </div>
    </div>
    <div class="collapse show" id="fieldsRemovalTarget">
      <p>Select any field you wish to remove from the entire trackpoints.</p>
    </div>
    <div v-for="(item, index) in dataSource" :key="index">
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          :id="item.value"
          :value="item.value"
          :disabled="toolMode == ToolMode.Unknown || !hasField(item.value as keyof Record)"
          v-model="selectedFields"
        />
        <label class="form-check-label" style="color: var(--color-text)" :for="item.value">
          {{ item.label }}
        </label>
      </div>
    </div>
    <div v-if="noAvailableFields"><p>(No available fields to be removed.)</p></div>
  </div>
</template>
<script lang="ts">
export default {
  props: {
    toolMode: { type: Number as PropType<ToolMode>, required: true },
    sessions: { type: Array<Session>, required: true }
  },
  data() {
    return {
      dataSource: [
        { value: 'distance', label: 'Distance' },
        { value: 'speed', label: 'Speed' },
        { value: 'cadence', label: 'Cadence' },
        { value: 'heartRate', label: 'Heart Rate' },
        { value: 'power', label: 'Power' },
        { value: 'temperature', label: 'Temperature' }
      ],
      selectedFields: new Array<String>()
    }
  },
  computed: {
    noAvailableFields(): boolean {
      for (const val of this.dataSource) {
        if (this.hasField(val.value as keyof Record)) return false
      }
      return true
    }
  },
  watch: {
    sessions: {
      handler() {
        this.selectedFields = []
      }
    },
    selectedFields: {
      handler() {
        this.$emit('selectedFields', this.selectedFields)
      }
    }
  },
  methods: {
    hasField<T extends keyof Record>(field: T): boolean {
      for (const ses of this.sessions) {
        for (const rec of ses.records) {
          if (rec[field] != null) return true
        }
      }
      return false
    },
    isSelectable(option: any) {
      return option.heading != true
    }
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
