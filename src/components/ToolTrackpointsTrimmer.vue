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
      :name="'trimmer'"
      :title="title"
      :help-text="helpText"
      :data-source="dataSource"
      :sessions="sessions"
      :tool-mode="toolMode"
      :legend-label="'New dist.'"
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
