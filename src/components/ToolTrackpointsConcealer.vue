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
      v-on:selected="onSelected"
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
    },
    onSelected(data: any) {
      this.$emit('active', data?.value)
    }
  },
  mounted() {}
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
