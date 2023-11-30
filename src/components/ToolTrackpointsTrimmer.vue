<script setup lang="ts">
import type { Session } from '@/spec/activity'
import { Marker, ToolMode } from '@/spec/activity-service'
</script>
<template>
  <div>
    <div
      class="row m-0"
      style="cursor: pointer"
      data-bs-toggle="collapse"
      data-bs-target="#trimTarget"
      aria-expanded="false"
      aria-controls="trimTarget"
    >
      <div class="text-start p-0">
        <label class="pe-1">Trim Trackpoints</label>
        <i class="fa-regular fa-circle-question" title="Show or Hide Help Text"></i>
      </div>
    </div>
    <div class="collapse show" id="trimTarget">
      <p>
        Select this option if you want to trim certain trackpoints in your activities. For instance,
        if you finish cycling and transition to another mode of transportation without turning off
        your cyclocomputer, this feature allows you to remove the unwanted trackpoints.
      </p>
    </div>
    <div class="col-12 pt-2">
      <v-select
        label="label"
        :options="dataSource"
        :clearable="false"
        :searchable="false"
        v-model="selected"
        :disabled="toolMode == ToolMode.Unknown"
      >
      </v-select>
    </div>
    <div v-if="selected?.value">
      <div class="pt-2" v-for="(marker, index) in markers" :key="index">
        <div class="row px-1">
          <div class="d-flex">
            <div class="col">
              <label class="sub-label">Session {{ index + 1 }}: {{ sessions[index].sport }} </label>
            </div>
            <div class="col-auto text-end fs-legend">
              <span>
                New dist.
                {{
                  (
                    ((sessions[index].records[marker.endN].distance ?? 0) -
                      (sessions[index].records[marker.startN].distance ?? 0)) /
                    1000
                  ).toFixed(2)
                }}
                km
              </span>
              <span>
                ({{
                  (
                    (sessions[index].records[sessions[index].records.length - 1].distance ?? 0) /
                    1000
                  ).toFixed(2)
                }}
                km)
              </span>
            </div>
          </div>
          <div class="col text-start"><p>Distance from the Start</p></div>
          <div class="col text-end">
            <span>
              {{ ((sessions[index].records[marker.startN].distance ?? 0) / 1000).toFixed(2) }} km
            </span>
          </div>
          <div class="ps-2">
            <input
              class="form-range openivity-form-range"
              type="range"
              :min="0"
              :max="sessions[index].records.length - 1"
              v-model="marker.startN"
            />
          </div>
        </div>
        <div class="row px-1">
          <div class="col text-start"><p>Distance from the End</p></div>
          <div class="col text-end">
            <span>
              {{
                (
                  ((sessions[index].records[sessions[index].records.length - 1].distance ?? 0) -
                    (sessions[index].records[marker.endN].distance ?? 0)) /
                  1000
                ).toFixed(2)
              }}
              km
            </span>
          </div>
          <div class="ps-2">
            <input
              class="form-range openivity-form-range"
              type="range"
              :min="0"
              :max="sessions[index].records.length - 1"
              v-model="marker.endN"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import type { PropType } from 'vue'
export default {
  props: {
    toolMode: { type: Number as PropType<ToolMode>, required: true },
    sessions: { type: Array<Session>, required: true }
  },
  data() {
    return {
      dataSource: [
        { label: 'Do Not Trim Trackpoints', value: false },
        { label: 'Trim Trackpoints', value: true }
      ],
      selected: null as unknown as any,
      markers: Array<Marker>()
    }
  },
  computed: {},
  watch: {
    sessions: {
      handler() {
        this.selected = this.dataSource[0]
        this.updateMarkers()
      }
    },
    markers: {
      handler(markers: Marker[]) {
        this.limitMarkers(markers)
        this.$emit('markers', markers)
      },
      deep: true
    },
    selected: {
      handler() {
        this.updateMarkers()
      }
    }
  },
  methods: {
    updateMarkers() {
      if (this.sessions.length == 0) return
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        this.markers[i] = new Marker({ startN: 0, endN: ses.records.length - 1 })
      }
    },
    limitMarkers(markers: Marker[]) {
      for (let i = 0; i < markers.length; i++) {
        const m = markers[i]
        m.startN = parseInt(m.startN as unknown as string)
        m.endN = parseInt(m.endN as unknown as string)

        if (m.startN > m.endN) {
          m.endN = m.startN
        }
      }
    }
  },
  mounted() {
    this.selected = this.dataSource[0]
    this.updateMarkers()
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
