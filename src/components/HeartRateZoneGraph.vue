<script setup lang="ts">
import HeartRateZoneBar from './HeartRateZoneBar.vue'
</script>

<template>
  <div class="col-12 h-100 pt-2" v-for="hrZone in hrZones" :key="hrZone.zone">
    <HeartRateZoneBar
      :zone="hrZone.zone"
      :zoneSub="hrZone.zoneSub"
      :minmax="hrZone.minmax"
      :timeInSecond="hrZone.timeInSecond"
      :prosen="hrZone.prosen"
      :progress-class="hrZone.progressClass"
      :progress-text="hrZone.showProgressText"
      :isLoading="isLoading"
    ></HeartRateZoneBar>
  </div>
</template>

<script lang="ts">
import { Record, Session } from '@/spec/activity'

export default {
  props: {
    // The record, it don't need to know it's single or combined
    selectedSession: {
      type: Array<Session>,
      required: true,
      default: []
    },
    age: {
      type: Number,
      required: true,
      default: 30
    },
    receivedRecord: Record,
    receivedRecordFreeze: Boolean
  },
  components: {
    HeartRateZoneBar
  },
  data() {
    return {
      isLoading: false,
      hrZones: [
        {
          zone: 'Zone 5',
          zoneSub: 'Maximum',
          minmax: [174, Infinity],
          timeInSecond: 0,
          prosen: 0,
          progressClass: ['bg-danger'],
          showProgressText: false
        },
        {
          zone: 'Zone 4',
          zoneSub: 'Threshold',
          minmax: [155, 174],
          timeInSecond: 0,
          prosen: 0,
          progressClass: ['bg-warning', 'text-dark'],
          showProgressText: false
        },
        {
          zone: 'Zone 3',
          zoneSub: 'Aaerobic',
          minmax: [136, 154],
          timeInSecond: 0,
          prosen: 0,
          progressClass: ['bg-success'],
          showProgressText: false
        },
        {
          zone: 'Zone 2',
          zoneSub: 'Easy',
          minmax: [116, 135],
          timeInSecond: 0,
          prosen: 0,
          progressClass: ['bg-primary'],
          showProgressText: false
        },
        {
          zone: 'Zone 1',
          zoneSub: 'Warm Up',
          minmax: [97, 115],
          timeInSecond: 0,
          prosen: 0,
          progressClass: ['bg-secondary'],
          showProgressText: false
        }
      ]
    }
  },
  watch: {
    selectedSession: {
      handler(sessions: Array<Session>) {
        this.summarizedGraph(sessions)
      }
    },
    age: {
      handler() {
        this.summarizedGraph(this.selectedSession)
      }
    }
  },
  computed: {
    // details(): Detail[] {
    //   return [
    //     new Detail({
    //       title: 'Avg Heart Rate',
    //       value: this.summary.avgHeartRate?.toFixed(0) ?? '0'
    //     }),
    //     new Detail({
    //       title: 'Max Heart Rate',
    //       value: this.summary.maxHeartRate?.toFixed(0) ?? '0'
    //     })
    //   ]
    // }
  },
  methods: {
    summarizedGraph(sessions: Array<Session>) {
      this.isLoading = true
      console.time('HR Zone Process')

      // Example
      // let zones = {
      //   zone1: { min: 110, max: 120 },
      //   zone2: { min: 121, max: 130 },
      //   zone3: { min: 131, max: 140 },
      //   zone4: { min: 141, max: 150 },
      //   zone5: { min: 151, max: 160 }
      // }
      // const data = [
      //   { timestamp: '2022-07-10T11:36:00Z', heartRate: 150 },
      //   { timestamp: '2022-07-10T11:36:40Z', heartRate: 120 }
      // ]
      // // Heart Rate Zone Breakdown:
      // // zone4: 12.90 seconds (32.26%)
      // // zone3: 12.90 seconds (32.26%)
      // // zone2: 12.90 seconds (32.26%)
      // // zone1: 1.29 seconds (3.23%)
      // // Total Time: 40.00 seconds
      //
      // const data = [
      //   { timestamp: '2022-07-10T11:36:00Z', heartRate: 120 },
      //   { timestamp: '2022-07-10T11:36:40Z', heartRate: 150 }
      // ]
      // // Heart Rate Zone Breakdown:
      // // zone1: 1.29 seconds (3.23%)
      // // zone2: 12.90 seconds (32.26%)
      // // zone3: 12.90 seconds (32.26%)
      // // zone4: 12.90 seconds (32.26%)
      // // Total Time: 40.00 seconds

      /* Re-adjust HR Zone */
      const maxHr = 220 - (this.age <= 1 ? 1 : this.age >= 150 ? 150 : this.age)

      this.hrZones.forEach((zone) => {
        zone.prosen = 0
        zone.timeInSecond = 0
      })

      this.hrZones[0].minmax = [Math.floor(maxHr * 0.9), Infinity]
      this.hrZones[1].minmax = [Math.floor(maxHr * 0.8), this.hrZones[0].minmax[0] - 1]
      this.hrZones[2].minmax = [Math.floor(maxHr * 0.7), this.hrZones[1].minmax[0] - 1]
      this.hrZones[3].minmax = [Math.floor(maxHr * 0.6), this.hrZones[2].minmax[0] - 1]
      this.hrZones[4].minmax = [Math.floor(maxHr * 0.5), this.hrZones[3].minmax[0] - 1]

      /* Categorized HR Zone */
      // Initialize variables for tracking total time in each zone
      if (!sessions) {
        console.timeEnd('HR Zone Process')
        this.isLoading = false
        return
      }
      let zoneTotals: number[] = []
      let totalSeconds = 0

      // Process each data point and calculate heart rate zone and total time
      sessions.forEach((session) => {
        if (!session.records) return
        for (let i = 0; i < session.records.length - 1; i++) {
          const entry = session.records[i]
          const nextEntry = session.records[i + 1]

          const hrZoneIndex = this.getHeartRateZoneIndex(entry.heartRate || 0)
          const nextHrZoneIndex = this.getHeartRateZoneIndex(nextEntry.heartRate || 0)

          if (entry.timestamp == null && nextEntry.timestamp == null) continue

          const timestamp1 = new Date(entry.timestamp || nextEntry.timestamp)
          const timestamp2 = new Date(nextEntry.timestamp || nextEntry.timestamp)
          const secondsDiff: any = (timestamp2 - timestamp1) / 1000

          totalSeconds += secondsDiff

          if (!zoneTotals[hrZoneIndex]) {
            zoneTotals[hrZoneIndex] = 0
          }

          // If there is a transition between heart rate zones, distribute time proportionally
          if (hrZoneIndex !== nextHrZoneIndex) {
            // calculate the bpm step
            const zonesInvolved = this.determineZonesInvolved(hrZoneIndex, nextHrZoneIndex)
            const totalSteps = zonesInvolved.reduce((total, zoneIndex) => {
              const start = Math.min(
                Math.max(this.hrZones[zoneIndex].minmax[0], entry.heartRate || 0),
                this.hrZones[zoneIndex].minmax[1]
              )
              const end = Math.min(
                Math.max(this.hrZones[zoneIndex].minmax[0], nextEntry.heartRate || 0),
                this.hrZones[zoneIndex].minmax[1]
              )
              return 1 + total + (Math.max(end, start) - Math.min(end, start))
            }, 0)

            // calculate the fraction or delta
            for (let j = 0; j < zonesInvolved.length; j++) {
              const zIndex = zonesInvolved[j]
              zoneTotals[zIndex] = zoneTotals[zIndex] || 0
              const start = Math.min(
                Math.max(this.hrZones[zIndex].minmax[0], entry.heartRate || 0),
                this.hrZones[zIndex].minmax[1]
              )
              const end = Math.min(
                Math.max(this.hrZones[zIndex].minmax[0], nextEntry.heartRate || 0),
                this.hrZones[zIndex].minmax[1]
              )
              const fraction = 1 + (Math.max(end, start) - Math.min(end, start))
              zoneTotals[zIndex] += secondsDiff * (fraction / totalSteps)
            }
          } else {
            // If the heart rate zone remains the same, accumulate time in the current zone
            zoneTotals[hrZoneIndex] += secondsDiff
          }
        }
      })

      // Calculate percentage of time in each zone and assign to hr zone
      const zonePercentages: any = {}
      for (const [zoneIndex, zoneSeconds] of zoneTotals.entries()) {
        const percentage = (zoneSeconds / totalSeconds) * 100
        zonePercentages[zoneIndex] = percentage

        if (this.hrZones[zoneIndex]) {
          this.hrZones[zoneIndex].prosen = percentage || 0
          this.hrZones[zoneIndex].timeInSecond = zoneSeconds
        }
      }

      // Debugging
      console.log('> Heart Rate Zone Breakdown:')
      for (const [zoneIndex, zoneSeconds] of zoneTotals.entries()) {
        console.log(
          ` > ${this.hrZones[zoneIndex]?.zone}: ${zoneSeconds?.toFixed(
            2
          )} seconds (${zonePercentages[zoneIndex].toFixed(2)}%)`
        )
      }

      console.log(`> Total Time: ${totalSeconds.toFixed(2)} seconds`)

      console.timeEnd('HR Zone Process')
      this.isLoading = false
    },
    // get hr this.hrZones index based on hr
    getHeartRateZoneIndex(heartRate: number) {
      for (const [i, d] of this.hrZones.entries()) {
        if (heartRate >= d.minmax[0] && heartRate <= d.minmax[1]) {
          return i
        }
      }
      return -1
    },
    // get hrzones involved between calculate transition 2 data hr
    determineZonesInvolved(startIndex: number, endIndex: number) {
      if (startIndex === -1 || endIndex === -1) {
        return []
      }

      const direction = startIndex < endIndex ? 1 : -1
      const zonesInvolved = []

      for (let i = startIndex; i !== endIndex + direction; i += direction) {
        zonesInvolved.push(i)
      }

      return zonesInvolved
    }
  },
  mounted() {
    this.summarizedGraph(this.selectedSession)
  }
}
</script>
<style lang="scss" scoped>
.progress {
  height: 15px;
  margin-top: 5px;
}
.hr-time {
  font-weight: bold;
}
.hr-prosen {
}
</style>
