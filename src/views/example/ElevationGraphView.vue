<script setup lang="ts">
import { Record } from '@/spec/activity'
import ElevationGraph from '../../components/ElevationGraph.vue'
</script>

<template>
  <main>
    <h4>Example Elevation Graph</h4>
    <div style="height: 400px">
      <ElevationGraph :name="'example'" :graph-records="records" :records="records" />
    </div>
  </main>
</template>

<script lang="ts">
export default {
  data() {
    return {
      records: Array<Record>()
    }
  },
  mounted() {
    this.records = [
      new Record({ timestamp: '2023-09-20T08:48:29Z', distance: 0, altitude: 90, grade: 0 }),
      new Record({
        timestamp: '2023-09-20T08:49:29Z',
        distance: 10,
        altitude: 95.35,
        grade: -15
      }),
      new Record({ timestamp: '2023-09-20T08:50:29Z', distance: 20, altitude: 91, grade: 5 }),
      new Record({
        timestamp: '2023-09-20T08:51:29Z',
        distance: 30,
        altitude: 95.35,
        grade: 0
      }),
      new Record({ timestamp: '2023-09-20T08:52:29Z', distance: 40, altitude: 93, grade: 10 }),
      new Record({ timestamp: '2023-09-20T08:53:29Z', distance: 50, altitude: 95, grade: 15 }),
      new Record({ timestamp: '2023-09-20T08:54:29Z', distance: 60, altitude: 100, grade: 30 }),
      new Record({ timestamp: '2023-09-20T08:55:29Z', distance: 70, altitude: 150, grade: 45 }),
      new Record({
        timestamp: '2023-09-20T08:56:29Z',
        distance: 80,
        altitude: 100.81,
        grade: 10
      }),
      new Record({ timestamp: '2023-09-20T08:57:29Z', distance: 90, altitude: 150, grade: 0 }),
      new Record({
        timestamp: '2023-09-20T08:58:29Z',
        distance: 100,
        altitude: 103.92,
        grade: 0
      }),
      new Record({
        timestamp: '2023-09-20T08:59:29Z',
        distance: 110,
        altitude: 140,
        grade: -10
      }),
      new Record({
        timestamp: '2023-09-20T09:01:29Z',
        distance: 120,
        altitude: 120,
        grade: -20
      }),
      new Record({
        timestamp: '2023-09-20T09:02:29Z',
        distance: 130,
        altitude: 90,
        grade: -30
      }),
      new Record({
        timestamp: '2023-09-20T09:03:29Z',
        distance: 140,
        altitude: 60,
        grade: -45
      }),
      new Record({
        timestamp: '2023-09-20T09:04:29Z',
        distance: 150,
        altitude: 60,
        grade: 0
      }),
      new Record({
        timestamp: '2023-09-20T09:05:29Z',
        distance: 160,
        altitude: 59,
        grade: -5
      }),
      new Record({
        timestamp: '2023-09-20T09:06:29Z',
        distance: 170,
        altitude: 58,
        grade: -5
      }),
      new Record({
        timestamp: '2023-09-20T09:06:29Z',
        distance: 180,
        altitude: 55,
        grade: -10
      }),
      new Record({
        timestamp: '2023-09-20T09:06:29Z',
        distance: 190,
        altitude: 55,
        grade: 0
      })
    ]
  }
}
</script>
