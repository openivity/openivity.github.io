<script setup lang="ts">
import { Duration } from 'luxon'
</script>

<template>
  <div class="container">
    <div class="row collapsible">
      <div
        class="col text-start"
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
                  {{ formatElev(splitSummary.totalAscend, splitSummary.totalDescend) }}
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
  overallDistance: number = 0 // record.totalDistance

  prosen: number = 0
  pace: number = 0
  isLeftover: boolean = false

  totalDistance: number = 0
  totalDuration: number = 0
  overallDuration: number = 0
  totalAscend: number = 0
  totalDescend: number = 0
  totalRecord: number = 0
  totalHeartRate: number = 0
  totalHeartRateRecord: number = 0
  lastRecord: Record | null = null
}

interface SplitProgress {
  prevRecord: Record
  loopDistance: number
  currentDuration: number
  maxPace: number
  totalHeartRate: number
  totalHeartRateRecord: number
  summarized: boolean
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
      const progress: SplitProgress = {
        prevRecord: emptyRecord,
        loopDistance: 0,
        currentDuration: 0,
        maxPace: 0,
        totalHeartRate: 0,
        totalHeartRateRecord: 0,
        summarized: false
      }

      for (const session of sessions) {
        if (session.records == null) continue
        if (!sessionHasPace(session)) continue

        for (const record of session.records) {
          progress.summarized = false

          if (record.timestamp != null && progress.prevRecord.timestamp != null) {
            const deltaTime =
              new Date(record.timestamp).valueOf() -
              new Date(progress.prevRecord.timestamp).valueOf()
            progress.currentDuration += deltaTime
          }

          // elev delta
          const deltaElev = (record.altitude ?? 0) - (progress.prevRecord.altitude ?? 0)
          if (deltaElev >= 0) {
            splitSummary.totalAscend += deltaElev
          } else {
            splitSummary.totalDescend -= Math.abs(deltaElev)
          }

          if (record.heartRate != null) {
            progress.totalHeartRate += record.heartRate
            progress.totalHeartRateRecord++
          }
          // // test Random HR
          // progress.totalHeartRate += Math.floor(Math.random() * (1 + 200 - 90)) + 90
          // progress.totalHeartRateRecord++

          // split by distance
          if ((record.distance ?? 0) - progress.loopDistance >= this.splitByDistanceInMeter) {
            progress.loopDistance = record.distance ?? 0

            this.recordSplitSummary(splitSummary, record, progress)
            progress.maxPace =
              splitSummary.pace > progress.maxPace ? splitSummary.pace : progress.maxPace

            this.summaries.push(splitSummary)
            splitSummary = new SplitSummary()

            // Reset local process var.
            progress.summarized = true
            progress.currentDuration = 0
            progress.totalHeartRate = 0
            progress.totalHeartRateRecord = 0
          }

          progress.prevRecord = record
        }
      }

      // last split by distance (if any)
      if (!progress.summarized && progress.prevRecord != emptyRecord) {
        progress.loopDistance = progress.prevRecord.distance ?? 0

        this.recordSplitSummary(splitSummary, progress.prevRecord, progress)

        // // omit merge to last split
        if (splitSummary.totalDistance < this.omitByMeter && this.summaries.length > 0) {
          const lastSplitSummary = this.summaries[this.summaries.length - 1]
          lastSplitSummary.overallDistance = splitSummary.overallDistance
          lastSplitSummary.totalDistance += splitSummary.totalDistance
          lastSplitSummary.totalAscend += splitSummary.totalAscend
          lastSplitSummary.totalDescend += splitSummary.totalDescend
          lastSplitSummary.totalDuration += splitSummary.totalDuration
          lastSplitSummary.totalHeartRate += splitSummary.totalHeartRate
          lastSplitSummary.totalHeartRateRecord += splitSummary.totalHeartRateRecord
          lastSplitSummary.totalRecord += splitSummary.totalRecord
          lastSplitSummary.lastRecord = splitSummary.lastRecord
          lastSplitSummary.pace =
            lastSplitSummary.totalDuration /
            (lastSplitSummary.totalDistance == 0 ? 1 : lastSplitSummary.totalDistance)
          lastSplitSummary.isLeftover = false // omit, not leftover
        } else {
          splitSummary.isLeftover = true // flag as leftover
          this.summaries.push(splitSummary)
        }

        const lastSplitSummary = this.summaries[this.summaries.length - 1]
        progress.maxPace =
          lastSplitSummary.pace > progress.maxPace ? lastSplitSummary.pace : progress.maxPace
      }

      // Calculate percentage of pace from max pace
      for (const [_, summary] of this.summaries.entries()) {
        const percentage = (summary.pace / progress.maxPace) * 100
        summary.prosen = percentage
      }

      this.isLoading = false
      console.timeEnd('Splits')
    },
    recordSplitSummary(splitSummary: SplitSummary, lastRecord: Record, progress: SplitProgress) {
      splitSummary.totalDistance =
        (lastRecord.distance ?? 0) -
        (this.summaries.length > 0 ? this.summaries[this.summaries.length - 1].overallDistance : 0)
      splitSummary.totalDuration = progress.currentDuration
      splitSummary.totalHeartRate = progress.totalHeartRate
      splitSummary.totalHeartRateRecord = progress.totalHeartRateRecord
      splitSummary.overallDistance = lastRecord.distance ?? 0
      splitSummary.pace =
        splitSummary.totalDuration /
        (splitSummary.totalDistance == 0 ? 1 : splitSummary.totalDistance)
      splitSummary.lastRecord = lastRecord
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
