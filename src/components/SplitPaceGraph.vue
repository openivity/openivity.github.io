<script setup lang="ts">
import { Duration } from 'luxon'
</script>

<template>
  <div class="container">
    <div class="row">
      <div
        class="col text-start collapsible"
        style="cursor: pointer"
        data-bs-toggle="collapse"
        data-bs-target="#split-pace-graph-content"
        aria-expanded="false"
        aria-controls="split-pace-graph-content"
      >
        <h6 class="pt-1 mb-0 title">
          <i class="fa-solid fa-caret-right collapse-indicator"></i>
          Splits
        </h6>
      </div>
      <div class="col-auto text-end">
        <div class="row g-0">
          <label for="splitDistance" class="col-auto col-form-label col-form-label-sm me-2"
            >Split By</label
          >
          <div class="col-auto">
            <select
              class="form-select form-select-sm"
              name="splitDistance"
              id="splitDistance"
              v-model="splitByDistanceInMeter"
            >
              <option v-for="option in splitByOptions" :key="option.value" :value="option.value">
                {{ option.text }}
              </option>
            </select>
          </div>
        </div>
      </div>
    </div>
    <div class="row collapse show" id="split-pace-graph-content">
      <div class="col-12 pt-2">
        <table class="table table-sm table-">
          <thead>
            <tr>
              <th scope="col" class="small col-2 text-start">KM</th>
              <th scope="col" class="small col-2 text-start">Pace</th>
              <th scope="col"></th>
              <th scope="col" class="small col-2 text-end">Elev</th>
              <th scope="col" class="small col-2 text-end">HR</th>
            </tr>
          </thead>
          <tbody class="table-group-divider">
            <template v-for="(splitSummary, i) in summaries" :key="i">
              <tr :class="[isRowBold(splitSummary) ? 'table-active' : '']">
                <td scope="row" class="fw-bold small text-start">
                  <template v-if="splitSummary.isLeftover">
                    {{ (splitSummary.totalDistance / 1000).toFixed(1) }}
                  </template>
                  <template v-else>
                    {{ (splitSummary.overallDistance / 1000).toFixed(0) }}
                  </template>
                </td>
                <td class="small text-start">
                  {{ formatPaceTime(splitSummary.pace) }}
                </td>
                <td class="position-relative">
                  <div
                    class="progress"
                    role="progressbar"
                    :aria-valuenow="isLoading ? 100 : formatProsen(splitSummary.prosen).toFixed(0)"
                    aria-valuemin="0"
                    aria-valuemax="100"
                  >
                    <div
                      :class="[
                        'progress-bar bg-primary',
                        isLoading ? 'progress-bar-striped progress-bar-animated' : ''
                      ]"
                      :style="{
                        width: `${isLoading ? 100 : formatProsen(splitSummary.prosen).toFixed(0)}%`
                      }"
                    ></div>
                  </div>
                  <div class="position-absolute top-50 end-0 translate-middle-y small me-1 d-none">
                    {{ formatDuration(splitSummary.overallDuration) }}
                  </div>

                  <div class="position-relative"></div>
                </td>
                <td class="small text-end">
                  {{ formatElev(splitSummary.totalAscent, splitSummary.totalDescent) }}
                </td>
                <td class="small text-end">
                  {{ formatAvgHr(splitSummary.totalHeartRate, splitSummary.totalHeartRateRecord) }}
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Record, Session } from '@/spec/activity'
import { sessionHasPace } from '@/toolkit/activity'

const emptyRecord = new Record()
class SplitSummary {
  prosen: number = 0
  pace: number = 0
  isLeftover: boolean = false

  totalDistance: number = 0 // Total distance of current split
  overallDistance: number = 0 // Overall distance from start till this split
  totalDuration: number = 0 // Total duration of current split
  overallDuration: number = 0 // Overall duration from start till this split

  totalAscent: number = 0
  totalDescent: number = 0
  totalRecord: number = 0
  totalHeartRate: number = 0
  totalHeartRateRecord: number = 0
  lastRecord: Record | null = null
}

class SplitProgress {
  firstRecord: Record = emptyRecord
  prevRecord: Record = emptyRecord
  currentDuration: number = 0
  totalHeartRate: number = 0
  totalHeartRateRecord: number = 0
  maxPace: number = 0
  distance: number = 0
  prevAltitude: number | null = null

  summarized: boolean = false
}

export default {
  props: {
    // The record, it don't need to know it's single or combined
    selectedSession: {
      type: Array<Session>,
      required: true,
      default: []
    },
    receivedRecord: Record,
    receivedRecordFreeze: Boolean
  },
  components: {},
  data() {
    return {
      isLoading: false,
      summaries: Array<SplitSummary>(),
      splitByDistanceInMeter: 1000,
      splitByOptions: [
        { text: '1 KM', value: 1000 },
        { text: '5 KM', value: 5000 },
        { text: '10 KM', value: 10000 }
      ],
      listRowBoldBy: [
        {
          splitBy: 1000,
          boldIns: [5, 10, 21, 42]
        }
      ],
      omitByMeter: 100
    }
  },
  watch: {
    selectedSession: {
      handler(sessions: Array<Session>) {
        this.summarize(sessions)
      }
    },
    splitByDistanceInMeter() {
      this.summarize(this.selectedSession)
    }
  },
  computed: {},
  methods: {
    isRowBold(splitSummary: SplitSummary): boolean {
      const criteria = this.listRowBoldBy.find((criteria) => {
        return this.splitByDistanceInMeter == criteria.splitBy
      })
      if (criteria) {
        let calc = Math.round
        if (splitSummary.isLeftover) {
          calc = Math.ceil
        }

        return criteria.boldIns.includes(calc(splitSummary.overallDistance / 1000))
      }
      return false
    },
    formatProsen(prosen: number): Number {
      return prosen >= 0 && prosen <= 100 ? prosen : prosen > 100 ? 100 : 0 // invalid number will be 0
    },
    formatElev(ascent: number, descent: number): String {
      return Math.round(ascent - descent).toFixed(0)
    },
    formatDuration(t: number): String {
      if (t >= 60 * 60 * 1000) return Duration.fromMillis(t).toFormat('h:mm:ss')
      else return Duration.fromMillis(t).toFormat('mm:ss')
    },
    formatPaceTime(t: number): String {
      if (t >= 60 * 60) return Duration.fromMillis(t * 1000).toFormat('h:mm:ss')
      else return Duration.fromMillis(t * 1000).toFormat('mm:ss')
    },
    formatAvgHr(hrTotal: number, hrRecord: number): String {
      if (hrRecord <= 0) return '-'
      return (hrTotal / hrRecord).toFixed(0)
    },
    // TODO optimize calculate & performance
    summarize(sessions: Array<Session>) {
      console.time('Splits')
      this.isLoading = true
      this.summaries.length = 0
      let splitSummary = new SplitSummary()

      // local process var.
      let progress: SplitProgress = new SplitProgress()
      let lastSplitDistance: number = 0

      for (const session of sessions) {
        if (session.records == null) continue
        if (!sessionHasPace(session)) {
          continue
        }

        progress = new SplitProgress()

        for (const record of session.records) {
          // ignore invalid distance
          if (record.distance == null) continue
          progress.summarized = false

          if (progress.firstRecord == emptyRecord) {
            progress.firstRecord = record
            lastSplitDistance = record.distance
          }

          // distance delta
          let distance = record.distance - progress.distance
          progress.distance += distance <= 0 ? 0 : distance

          if (record.timestamp != null && progress.prevRecord.timestamp != null) {
            const deltaTime =
              new Date(record.timestamp).valueOf() -
              new Date(progress.prevRecord.timestamp).valueOf()
            progress.currentDuration += deltaTime
          }

          // Elevation Gain, compare current altitude vs latest valid altitude
          // 1st record always ignored
          if (record.altitude != null) {
            if (progress.prevAltitude != null) {
              const deltaElev = record.altitude - progress.prevAltitude
              if (deltaElev > 0) {
                splitSummary.totalAscent += deltaElev
              } else {
                splitSummary.totalDescent += Math.abs(deltaElev)
              }
            }

            progress.prevAltitude = record.altitude
          }

          if (record.heartRate != null) {
            progress.totalHeartRate += record.heartRate
            progress.totalHeartRateRecord++
          }

          // // test Random HR
          // progress.totalHeartRate += Math.floor(Math.random() * (1 + 200 - 90)) + 90
          // progress.totalHeartRateRecord++

          // split by distance
          if (progress.distance - lastSplitDistance >= this.splitByDistanceInMeter) {
            splitSummary.totalDistance = progress.distance - lastSplitDistance
            splitSummary.totalDuration = progress.currentDuration
            splitSummary.totalHeartRate = progress.totalHeartRate
            splitSummary.totalHeartRateRecord = progress.totalHeartRateRecord
            splitSummary.overallDistance = progress.distance

            splitSummary.pace =
              splitSummary.totalDuration /
              (splitSummary.totalDistance == 0 ? 1 : splitSummary.totalDistance)
            splitSummary.lastRecord = record

            progress.maxPace =
              splitSummary.pace > progress.maxPace ? splitSummary.pace : progress.maxPace

            this.summaries.push(splitSummary)
            splitSummary = new SplitSummary()

            // Reset local process var.
            progress.summarized = true
            progress.currentDuration = 0
            progress.totalHeartRate = 0
            progress.totalHeartRateRecord = 0
            lastSplitDistance = progress.distance
          }

          progress.prevRecord = record
        }
      }

      // last split by distance (if any)
      if (!progress.summarized && progress.prevRecord != emptyRecord) {
        splitSummary.totalDistance = progress.distance - lastSplitDistance
        splitSummary.totalDuration = progress.currentDuration
        splitSummary.totalHeartRate = progress.totalHeartRate
        splitSummary.totalHeartRateRecord = progress.totalHeartRateRecord
        splitSummary.overallDistance = progress.distance
        splitSummary.pace =
          splitSummary.totalDuration /
          (splitSummary.totalDistance == 0 ? 1 : splitSummary.totalDistance)
        splitSummary.lastRecord = progress.prevRecord

        // // omit merge to last split
        if (splitSummary.totalDistance < this.omitByMeter && this.summaries.length > 0) {
          const lastSplitSummary = this.summaries[this.summaries.length - 1]
          lastSplitSummary.overallDistance = splitSummary.overallDistance // replace
          lastSplitSummary.totalDistance += splitSummary.totalDistance
          lastSplitSummary.totalAscent += splitSummary.totalAscent
          lastSplitSummary.totalDescent += splitSummary.totalDescent
          lastSplitSummary.totalDuration += splitSummary.totalDuration
          lastSplitSummary.totalHeartRate += splitSummary.totalHeartRate
          lastSplitSummary.totalHeartRateRecord += splitSummary.totalHeartRateRecord
          lastSplitSummary.totalRecord += splitSummary.totalRecord
          lastSplitSummary.lastRecord = splitSummary.lastRecord // replace
          lastSplitSummary.pace =
            lastSplitSummary.totalDuration /
            (lastSplitSummary.totalDistance == 0 ? 1 : lastSplitSummary.totalDistance)
          lastSplitSummary.isLeftover = false // omit, not leftover
        } else {
          splitSummary.isLeftover = true // flag as leftover
          this.summaries.push(splitSummary)
        }

        const lastestSplitSummary = this.summaries[this.summaries.length - 1]
        progress.maxPace =
          lastestSplitSummary.pace > progress.maxPace ? lastestSplitSummary.pace : progress.maxPace
      }

      // Calculate percentage of pace from max pace
      for (const [_, summary] of this.summaries.entries()) {
        const percentage = (summary.pace / progress.maxPace) * 100
        summary.prosen = percentage
      }

      this.isLoading = false
      console.timeEnd('Splits')
    }
  },
  mounted() {
    this.summarize(this.selectedSession)
  }
}
</script>
<style lang="scss" scoped>
.progress,
.progress-stacked {
  border-radius: 0px;
  height: 20px;
  margin-top: 2px;
  background-color: var(--bs-progress-bg-focus);
}

tbody td {
  border: none;
}
.table-sm > :not(caption) > * > * {
  padding: 0.1rem 0rem;
}
</style>
