<script setup lang="ts">
import HeartRateZoneBar from './HeartRateZoneBar.vue'
</script>

<template>
  <div class="container pt-2 pb-3">
    <div class="row">
      <div
        class="col text-start collapsible"
        style="cursor: pointer"
        data-bs-toggle="collapse"
        data-bs-target="#hrzone-graph-content"
        aria-expanded="false"
        aria-controls="hrzone-graph-content"
      >
        <h6 class="pt-1 mb-0 title">
          <i class="fa-solid fa-caret-right collapse-indicator"></i>
          Heart Rate Zone
        </h6>
      </div>
      <div class="col-auto text-end">
        <div class="input-group input-group-sm">
          <span class="input-group-text ps-1">Max HR</span>
          <input
            type="text"
            class="form-control form-control-sm text-end"
            placeholder="-"
            inputmode="numeric"
            maxlength="3"
            style="width: calc(1.35em + 3ch)"
            :value="maxHr"
            @change="maxHrOnChange(maxHr, $event)"
            :readonly="isLoading"
          />
          <span class="input-group-text pe-1"
            >bpm
            <i
              class="fa-solid fa-question-circle ps-1 mt-lg-1"
              data-bs-toggle="tooltip"
              data-bs-html="true"
              data-bs-custom-class="openivity-tooltip"
              data-bs-title="
            <span>Common formula: </span><br />
            <b>Max HR = 220 - Age</b>
          "
              >&nbsp;
            </i>
          </span>
        </div>
      </div>
    </div>
    <div class="row collapse show" id="hrzone-graph-content">
      <div class="col-12 pt-2" v-for="hrZone in hrZones" :key="hrZone.zone">
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
    </div>
  </div>
</template>

<script lang="ts">
import { Record, Session } from '@/spec/activity'
import { Tooltip } from 'bootstrap'

const defaultHr = 193

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
        // Please don't remove, it's for reference if using Zone 5A/B/C
        // {
        //   zone: 'Zone 5C',
        //   zoneSub: 'Maximum',
        //   minmax: [179, Infinity],
        //   timeInSecond: 0,
        //   prosen: 0,
        //   progressClass: ['bg-danger'],
        //   showProgressText: false
        // },
        // {
        //   zone: 'Zone 5B',
        //   zoneSub: 'Maximum',
        //   minmax: [176, 179],
        //   timeInSecond: 0,
        //   prosen: 0,
        //   progressClass: ['bg-danger'],
        //   showProgressText: false
        // },
        // {
        //   zone: 'Zone 5A',
        //   zoneSub: 'Maximum',
        //   minmax: [174, 176],
        //   timeInSecond: 0,
        //   prosen: 0,
        //   progressClass: ['bg-danger'],
        //   showProgressText: false
        // },
        {
          zone: 'Zone 5',
          zoneSub: 'Maximum',
          minmax: [174, 176],
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
      ],
      maxHr: defaultHr
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
  computed: {},
  methods: {
    maxHrOnChange(oldHr: number, event: Event) {
      this.$nextTick(() => {
        this.maxHr = this.validateHr(parseFloat((event.target as HTMLInputElement).value), oldHr)
        this.$forceUpdate()

        if (oldHr == this.maxHr) return
        this.summarizedGraph(this.selectedSession)
      })
    },
    validateHr(value: number, defaultValue: number) {
      return value >= 50 && value <= 300 ? value : defaultValue
    },
    summarizedGraph(sessions: Array<Session>) {
      this.isLoading = true
      console.time('HR Zone Process')

      /* Re-adjust HR Zone */
      // this.maxHr = this.hrConst - (this.age <= 1 ? 1 : this.age >= 150 ? 150 : this.age)

      this.hrZones.forEach((zone) => {
        zone.prosen = 0
        zone.timeInSecond = 0
      })

      // Please don't remove, it's for reference if using Zone 5A/B/C
      // this.hrZones[0].minmax = [Math.floor(this.maxHr * 0.96), Infinity]
      // this.hrZones[1].minmax = [Math.floor(this.maxHr * 0.93), this.hrZones[0].minmax[0]]
      // this.hrZones[2].minmax = [Math.floor(this.maxHr * 0.915), this.hrZones[1].minmax[0]]
      // this.hrZones[3].minmax = [Math.floor(this.maxHr * 0.85), this.hrZones[2].minmax[0]]
      // this.hrZones[4].minmax = [Math.floor(this.maxHr * 0.8), this.hrZones[3].minmax[0]]
      // this.hrZones[5].minmax = [Math.floor(this.maxHr * 0.745), this.hrZones[4].minmax[0]]
      // this.hrZones[6].minmax = [Math.floor(this.maxHr * 0), this.hrZones[5].minmax[0]]
      this.hrZones[0].minmax = [Math.floor(this.maxHr * 0.9), Infinity]
      this.hrZones[1].minmax = [Math.floor(this.maxHr * 0.8), this.hrZones[0].minmax[0] - 1]
      this.hrZones[2].minmax = [Math.floor(this.maxHr * 0.7), this.hrZones[1].minmax[0] - 1]
      this.hrZones[3].minmax = [Math.floor(this.maxHr * 0.6), this.hrZones[2].minmax[0] - 1]
      this.hrZones[4].minmax = [Math.floor(this.maxHr * 0.5), this.hrZones[3].minmax[0] - 1]

      if (!sessions) {
        console.timeEnd('HR Zone Process')
        this.isLoading = false
        return
      }

      // categorize...
      // this.categorizedByDiffData(sessions, false)
      // this.categorizedByDiffData(sessions, true)
      this.categorizedByTransition(sessions)

      console.timeEnd('HR Zone Process')
      this.isLoading = false
    },
    /* Categorized HR Zone ver. Diff */
    categorizedByDiffData(sessions: Array<Session>, useNextDataAsHrZone: Boolean) {
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
      //   { timestamp: '2022-07-10T11:36:20Z', heartRate: 120 }
      // ]
      // // Heart Rate Zone Breakdown:
      // // zone4: 20s
      // // zone3: 0s
      // // zone2: 0s
      // // zone1: 0s
      // // Total Time: 20.00 seconds
      // // Reverse if useNextDataAsHrZone = true

      // Initialize variables for tracking total time in each zone
      let zoneTotals: number[] = []
      let totalSeconds = 0

      // Process each data point and calculate heart rate zone and total time
      sessions.forEach((session) => {
        if (session.records == null) return
        for (let i = 0; i < session.records.length - 1; i++) {
          const entry = session.records[i]
          if (entry == null || entry.heartRate == null) continue

          let { nextEntry, nextIndex } = this.getNextValidEntry(session, entry, i)
          i = nextIndex // skip loop to latest valid entry

          const hrZoneIndex = this.getHeartRateZoneIndex(entry.heartRate)
          const nextHrZoneIndex = this.getHeartRateZoneIndex(nextEntry.heartRate ?? 0) // should be valid, but tslinter can't check

          if (entry.timestamp == null || nextEntry.timestamp == null) continue

          const timestamp1 = new Date(entry.timestamp)
          const timestamp2 = new Date(nextEntry.timestamp)
          let secondsDiff: number = (timestamp2.valueOf() - timestamp1.valueOf()) / 1000
          if (secondsDiff > 30 || secondsDiff < 0) secondsDiff = 1
          if (entry == nextEntry) secondsDiff = 1

          totalSeconds += secondsDiff

          let selectedHrZoneIndex = hrZoneIndex
          if (useNextDataAsHrZone) selectedHrZoneIndex = nextHrZoneIndex

          if (!zoneTotals[selectedHrZoneIndex]) {
            zoneTotals[selectedHrZoneIndex] = 0
          }

          zoneTotals[selectedHrZoneIndex] += secondsDiff
        }
      })

      // Calculate percentage of time in each zone and assign to hr zone
      const zonePercentages: any = {}
      for (const [zoneIndex, zoneSeconds] of zoneTotals.entries()) {
        if (this.hrZones[zoneIndex]) {
          const percentage = (zoneSeconds / totalSeconds) * 100
          zonePercentages[zoneIndex] = percentage

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
    },
    /* Categorized HR Zone ver. Transition Zone */
    categorizedByTransition(sessions: Array<Session>) {
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

      // Initialize variables for tracking total time in each zone
      let zoneTotals: number[] = []
      let totalSeconds = 0

      // Process each data point and calculate heart rate zone and total time
      // TODO optimize calculation
      sessions.forEach((session) => {
        if (session.records == null) return
        console.time('totalSteps')
        for (let i = 0; i < session.records.length - 1; i++) {
          const entry = session.records[i]
          if (entry == null || entry.heartRate == null) continue

          let { nextEntry, nextIndex } = this.getNextValidEntry(session, entry, i)
          i = nextIndex // skip loop to latest valid entry

          const hrZoneIndex = this.getHeartRateZoneIndex(entry.heartRate)
          const nextHrZoneIndex = this.getHeartRateZoneIndex(nextEntry.heartRate ?? 0) // should be valid, but tslinter can't check

          if (entry.timestamp == null || nextEntry.timestamp == null) continue

          const timestamp1 = new Date(entry.timestamp)
          const timestamp2 = new Date(nextEntry.timestamp)
          let secondsDiff: number = (timestamp2.valueOf() - timestamp1.valueOf()) / 1000
          if (secondsDiff > 30 || secondsDiff < 0) secondsDiff = 1
          if (entry == nextEntry) secondsDiff = 1

          totalSeconds += secondsDiff

          if (!zoneTotals[hrZoneIndex]) {
            zoneTotals[hrZoneIndex] = 0
          }

          // If there is a transition between heart rate zones, distribute time proportionally
          if (hrZoneIndex !== nextHrZoneIndex) {
            // calculate the bpm step
            const zonesInvolved = this.determineZonesInvolved(hrZoneIndex, nextHrZoneIndex)
            const totalSteps = zonesInvolved.reduce((total, zoneIndex) => {
              const start =
                zoneIndex == -1
                  ? entry.heartRate ?? 0
                  : Math.min(
                      Math.max(this.hrZones[zoneIndex].minmax[0], entry.heartRate ?? 0),
                      this.hrZones[zoneIndex].minmax[1]
                    )
              const end =
                zoneIndex == -1
                  ? nextEntry.heartRate ?? 0
                  : Math.min(
                      Math.max(this.hrZones[zoneIndex].minmax[0], nextEntry.heartRate ?? 0),
                      this.hrZones[zoneIndex].minmax[1]
                    )
              return 1 + total + (Math.max(end, start) - Math.min(end, start))
            }, 0)

            // calculate the fraction or delta
            for (let j = 0; j < zonesInvolved.length; j++) {
              const zIndex = zonesInvolved[j]
              zoneTotals[zIndex] = zoneTotals[zIndex] || 0
              const start =
                zIndex == -1
                  ? entry.heartRate ?? 0
                  : Math.min(
                      Math.max(this.hrZones[zIndex].minmax[0], entry.heartRate || 0),
                      this.hrZones[zIndex].minmax[1]
                    )
              const end =
                zIndex == -1
                  ? nextEntry.heartRate ?? 0
                  : Math.min(
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
        console.timeEnd('totalSteps')
      })

      // Calculate percentage of time in each zone and assign to hr zone
      const zonePercentages: any = {}
      const invalidTotalSeconds = zoneTotals[-1] ?? 0
      for (const [zoneIndex, zoneSeconds] of zoneTotals.entries()) {
        if (this.hrZones[zoneIndex]) {
          const percentage = (zoneSeconds / (totalSeconds - invalidTotalSeconds)) * 100
          zonePercentages[zoneIndex] = percentage

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
    },
    getNextValidEntry(session: Session, currentEntry: Record, currentIndex: number) {
      // findout next record with valid HR
      for (let index = currentIndex + 1; index < session.records.length; index++) {
        const r = session.records[index]
        // r.heartRate = (Math.floor(Math.random() * (10 + 1)) + 1) % 2 == 0 ? r.heartRate : 55 // Test random null HR
        if (r.heartRate != null) return { nextEntry: r, nextIndex: index }
      }
      // no next entry, use current entry as last comparator
      return { nextEntry: currentEntry, nextIndex: session.records.length - 1 }
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
      const direction = startIndex < endIndex ? 1 : -1
      const zonesInvolved = []

      if (startIndex == -1) zonesInvolved.push(-1)
      for (let i = startIndex; i !== endIndex + direction; i += direction) {
        zonesInvolved.push(i)
      }
      if (endIndex == -1) zonesInvolved.push(-1)

      return zonesInvolved
    }
  },
  mounted() {
    this.summarizedGraph(this.selectedSession)
    new Tooltip(document.body, {
      selector: "[data-bs-toggle='tooltip']"
    })
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
  font-weight: normal;
}
</style>
