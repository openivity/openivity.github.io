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
        data-bs-target="#hrzone-graph-content"
        aria-expanded="false"
        aria-controls="hrzone-graph-content"
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
    <div class="row collapse show" id="hrzone-graph-content">
      <div class="col-12 pt-2">
        <table class="table table-sm table-">
          <thead>
            <tr>
              <th scope="col" class="small col-2 text-start">KM</th>
              <th scope="col" class="small col-2 text-start">Pace</th>
              <th scope="col"></th>
              <th scope="col" class="small col-2 text-end">Elev</th>
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
                  {{ formatPace(splitSummary.pace) }}
                </td>
                <td class="">
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
                </td>
                <td class="small text-end">
                  {{ formatElev(splitSummary.totalAscend, splitSummary.totalDescend) }}
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

  totalAscend: number = 0
  totalDescend: number = 0

  totalRecord: number = 0
  lastRecord: Record | null = null
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
      ]
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
    formatPace(pace: number): String {
      if (pace >= 60 * 60) return Duration.fromMillis(pace * 1000).toFormat('h:mm:ss')
      else return Duration.fromMillis(pace * 1000).toFormat('mm:ss')
    },
    // TODO optimize calculate & performance
    summarize(sessions: Array<Session>) {
      console.time('Splits')
      this.isLoading = true
      this.summaries.length = 0
      let splitSummary = new SplitSummary()

      // local process var.
      let _prevRecord: Record = emptyRecord
      let _loopDistance = 0
      let _currentDuration = 0
      let _maxPace = 0
      let _summarized = false

      for (const session of sessions) {
        if (session.records == null) continue
        if (!sessionHasPace(session)) continue
        
        for (const record of session.records) {
          _summarized = false

          if (record.timestamp != null && _prevRecord.timestamp != null) {
            const deltaTime =
              new Date(record.timestamp).valueOf() - new Date(_prevRecord.timestamp).valueOf()
            _currentDuration += deltaTime
          }

          // elev delta
          const deltaElev = (record.altitude ?? 0) - (_prevRecord.altitude ?? 0)
          if (deltaElev >= 0) {
            splitSummary.totalAscend += deltaElev
          } else {
            splitSummary.totalDescend -= Math.abs(deltaElev)
          }

          // split by distance
          if ((record.distance ?? 0) - _loopDistance >= this.splitByDistanceInMeter) {
            _loopDistance = record.distance ?? 0

            this.recordSplitSummary(splitSummary, record, _currentDuration)

            _maxPace = splitSummary.pace > _maxPace ? splitSummary.pace : _maxPace

            this.summaries.push(splitSummary)
            splitSummary = new SplitSummary()

            // Reset local process var.
            _summarized = true
            _currentDuration = 0
          }

          _prevRecord = record
        }
      }

      // last split by distance (if any)
      if (!_summarized && _prevRecord != emptyRecord) {
        _loopDistance = _prevRecord.distance ?? 0

        this.recordSplitSummary(splitSummary, _prevRecord, _currentDuration)

        splitSummary.isLeftover = true // flag as leftover

        _maxPace = splitSummary.pace > _maxPace ? splitSummary.pace : _maxPace

        this.summaries.push(splitSummary)
      }

      // Calculate percentage of pace from max pace
      for (const [_, summary] of this.summaries.entries()) {
        const percentage = (summary.pace / _maxPace) * 100
        summary.prosen = percentage
      }

      this.isLoading = false
      console.timeEnd('Splits')
    },
    recordSplitSummary(splitSummary: SplitSummary, lastRecord: Record, _currentDuration: number) {
      splitSummary.totalDistance =
        (lastRecord.distance ?? 0) -
        (this.summaries.length > 0 ? this.summaries[this.summaries.length - 1].overallDistance : 0)
      splitSummary.totalDuration = _currentDuration
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
.progress {
  border-radius: 0px;
  height: 20px;
  margin-top: 2px;
}

tbody td {
  border: none;
}
.table-sm > :not(caption) > * > * {
  padding: 0.1rem 0rem;
}
</style>
