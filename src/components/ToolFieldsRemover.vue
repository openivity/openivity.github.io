<script setup lang="ts">
import type { Session } from '@/spec/activity'
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
      <div class="form-check" v-show="show(item.value)">
        <input
          class="form-check-input"
          type="checkbox"
          :id="item.value"
          :value="item.value"
          :disabled="toolMode == ToolMode.Unknown"
          v-model="selectedFields"
        />
        <label class="form-check-label" style="color: var(--color-text)" :for="item.value">
          {{ item.label }}
        </label>
      </div>
    </div>
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
        { value: 'cadence', label: 'Cadence' },
        { value: 'heartRate', label: 'Heart Rate' },
        { value: 'power', label: 'Power' },
        { value: 'temperature', label: 'Temperature' }
      ],
      selectedFields: new Array<String>()
    }
  },
  computed: {
    hasCadence(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.cadence != null) return true
        }
      }
      return false
    },
    hasHeartRate(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.heartRate != null) return true
        }
      }
      return false
    },
    hasPower(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.power != null) return true
        }
      }
      return false
    },
    hasTemperature(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.temperature != null) return true
        }
      }
      return false
    }
  },
  watch: {
    selectedFields: {
      handler(value) {
        this.$emit('selectedFields', this.selectedFields)
      }
    }
  },
  methods: {
    show(value: string): boolean {
      switch (value) {
        case 'cadence':
          return this.hasCadence
        case 'heartRate':
          return this.hasHeartRate
        case 'power':
          return this.hasPower
        case 'temperature':
          return this.hasTemperature
        default:
          return false
      }
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
