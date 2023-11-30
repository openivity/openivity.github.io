<script setup lang="ts">
import { Marker, ToolMode } from '@/spec/activity-service'
import ToolTrackpointsSlider from './ToolTrackpointsSlider.vue'
</script>
<template>
  <div>
    <ToolTrackpointsSlider
      :name="'trimmer'"
      :title="title"
      :help-text="helpText"
      :data-source="dataSource"
      :sessions="sessions"
      :tool-mode="toolMode"
      :legend-label="'New dist.'"
      v-on:markers="onMarkers"
    >
    </ToolTrackpointsSlider>
  </div>
</template>
<script lang="ts">
import type { Session } from '@/spec/activity'
import type { PropType } from 'vue'

export default {
  props: {
    toolMode: { type: Number as PropType<ToolMode>, required: true },
    sessions: { type: Array<Session>, required: true }
  },
  data() {
    return {
      title: `Trim Trackpoints`,
      helpText: `Select this option if you want to trim certain trackpoints in your activities. For instance,
        if you finish cycling and transition to another mode of transportation without turning off
        your cyclocomputer, this feature allows you to remove the unwanted trackpoints.`,
      dataSource: [
        { label: 'Do Not Trim Trackpoints', value: false },
        { label: 'Trim Trackpoints', value: true }
      ]
    }
  },
  methods: {
    onMarkers(value: Marker[]) {
      this.$emit('markers', value)
    }
  },
  mounted() {}
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
