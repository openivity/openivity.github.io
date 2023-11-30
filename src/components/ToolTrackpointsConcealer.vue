<script setup lang="ts">
import { Marker, ToolMode } from '@/spec/activity-service'
import ToolTrackpointsSlider from './ToolTrackpointsSlider.vue'
</script>
<template>
  <div>
    <ToolTrackpointsSlider
      :name="'concealer'"
      :title="title"
      :help-text="helpText"
      :data-source="dataSource"
      :sessions="sessions"
      :tool-mode="toolMode"
      :legend-label="'Visible'"
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
      title: `Conceal GPS Positions`,
      helpText: `Select this option if you wish to conceal your GPS positions at a specific start and end distance 
      while maintaining other trackpoints data. It is useful when you want to share the activity on social media, 
      like Strava, and you want to avoid revealing the exact location of your home for security reasons and 
      you wouldn't want your information stored on any platform's server.`,
      dataSource: [
        { label: 'Do Not Conceal My GPS Positions', value: false },
        { label: 'Conceal My GPS positions', value: true }
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
