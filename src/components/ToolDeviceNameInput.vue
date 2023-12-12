<script setup lang="ts">
import type { ActivityFile } from '@/spec/activity'
import { ToolMode } from '@/spec/activity-service'
import type { PropType } from 'vue'
</script>
<template>
  <div>
    <label>Target Device Name</label>
    <div class="col-12 pt-1">
      <input
        class="form-control form-control-sm"
        v-model="deviceName"
        placeholder="-- Please input device name --"
        :disabled="toolMode == ToolMode.Unknown"
      />
    </div>
  </div>
</template>
<script lang="ts">
export default {
  props: {
    activities: { type: Array<ActivityFile>, required: true },
    toolMode: { type: Number as PropType<ToolMode>, required: true }
  },

  data() {
    return {
      deviceName: ''
    }
  },
  computed: {
    deviceNameFromFile(): string {
      for (let i = 0; i < this.activities.length; i++) {
        if (this.deviceName != '') break
        this.deviceName = this.activities[i].creator.name
      }
      return this.deviceName
    }
  },
  watch: {
    toolMode: {
      handler() {
        this.deviceName = this.deviceNameFromFile
      }
    },
    deviceName: {
      handler(value: string) {
        this.$emit('deviceName', value)
      }
    }
  },
  mounted() {
    this.deviceName = this.deviceNameFromFile
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
