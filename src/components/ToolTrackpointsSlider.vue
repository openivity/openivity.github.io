<script setup lang="ts">
import { Marker, ToolMode } from '@/spec/activity-service'
</script>
<template>
  <div>
    <div
      class="row m-0"
      style="cursor: pointer"
      data-bs-toggle="collapse"
      :data-bs-target="`#${name}-target`"
      aria-expanded="false"
      :aria-controls="`${name}-target`"
    >
      <div class="text-start p-0">
        <label class="pe-1">{{ title }}</label>
        <i class="fa-regular fa-circle-question" title="Show or Hide Help Text"></i>
      </div>
    </div>
    <div class="collapse show" :id="`${name}-target`">
      <p>
        {{ helpText }}
      </p>
    </div>
    <div class="col-12 pt-2">
      <div>
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
      <div v-if="selected?.value == true">
        <div class="pt-2" v-for="(marker, index) in markers" :key="index">
          <div class="row px-1">
            <div class="d-flex">
              <div class="col">
                <label class="sub-label"
                  >Session {{ index + 1 }}: {{ sessions[index].sport }}
                </label>
              </div>
              <div class="col-auto text-end fs-legend">
                <span>
                  {{ legendLabel }}
                  {{ (distanceByMarkers[index].delta / 1000).toFixed(2) }}
                  km
                </span>
                <span>
                  ({{ (lastDistance(index) / 1000).toFixed(2) }}
                  km)
                </span>
              </div>
            </div>
            <div class="col text-start"><p>Distance from the Start</p></div>
            <div class="col text-end">
              <span>
                {{ (distanceByMarkers[index].start / 1000).toFixed(2) }}
                km
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
              <span
                >{{ ((lastDistance(index) - distanceByMarkers[index].end) / 1000).toFixed(2) }}
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
  </div>
</template>
<script lang="ts">
import type { Session } from '@/spec/activity'
import type { PropType } from 'vue'

class DistanceByMarker {
  start: number = 0
  end: number = 0
  delta: number = 0
}

export default {
  props: {
    name: { type: String, required: true },
    toolMode: { type: Number as PropType<ToolMode>, required: true },
    sessions: { type: Array<Session>, required: true },
    dataSource: { type: Array<Object>, required: true },
    title: String,
    helpText: String,
    legendLabel: String
  },
  data() {
    return {
      selected: null as unknown as any,
      markers: Array<Marker>()
    }
  },
  computed: {
    distanceByMarkers(): DistanceByMarker[] {
      const distanceByMarkers = new Array<DistanceByMarker>()
      for (let i = 0; i < this.markers.length; i++) {
        const distanceByMarker = new DistanceByMarker()

        // Find nearest record with non-nil distance relative to marker's startN
        for (let j = this.markers[i].startN; j < this.sessions[i].records.length; j++) {
          const rec = this.sessions[i].records[j]
          if (rec.distance != null) {
            distanceByMarker.start = rec.distance
            break
          }
        }

        // Find nearest record with non-nil distance relative to marker's endN
        for (let j = this.markers[i].endN; j >= 0; j--) {
          const rec = this.sessions[i].records[j]
          if (rec.distance != null) {
            distanceByMarker.end = rec.distance
            break
          }
        }

        distanceByMarker.delta = distanceByMarker.end - distanceByMarker.start

        distanceByMarkers.push(distanceByMarker)
      }

      return distanceByMarkers
    }
  },
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
    lastDistance(sessionIndex: number): number {
      for (let i = this.sessions[sessionIndex].records.length - 1; i >= 0; i--) {
        const rec = this.sessions[sessionIndex].records[i]
        if (rec.distance != null) return rec.distance
      }
      return 0
    },
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
