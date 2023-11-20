<script setup lang="ts">
import { Summary } from '@/spec/summary'
import CadenceGraph from '@/components/CadenceGraph.vue'
import ElevationGraph from '@/components/ElevationGraph.vue'
import HeartRateGraph from '@/components/HeartRateGraph.vue'
import PowerGraph from '@/components/PowerGraph.vue'
import SpeedGraph from '@/components/SpeedGraph.vue'
import TemperatureGraph from '@/components/TemperatureGraph.vue'
import TheLap from '@/components/TheLap.vue'
import TheLoading from '@/components/TheLoading.vue'
import TheMap from '@/components/TheMap.vue'
import TheNavigatorInput from '@/components/TheNavigatorInput.vue'
import TheSession from '@/components/TheSession.vue'
import TheSummary from '@/components/TheSummary.vue'
</script>

<template>
  <div>
    <Transition>
      <TheLoading v-if="loading"></TheLoading>
    </Transition>
    <div class="activity container-fluid text-center flex-container">
      <div :class="['row', !isActivityFileReady ? 'align-items-center' : '']">
        <div
          id="left"
          :class="[
            'sidebar',
            'col-12',
            'px-0',
            isActivityFileReady
              ? 'col-xxl-3 col-xl-4 col-lg-5 d-md-block float-sm-right'
              : 'col-md-12 landing'
          ]"
        >
          <div>
            <div :class="[isActivityFileReady ? 'default-border-top' : '']">
              <div class="pt-4 pb-2">
                <h5 class="title">Open Activity</h5>
                <div style="font-size: 0.9em">
                  Your data stays in your computer: 100% client-side power.
                </div>
              </div>
              <TheNavigatorInput
                :isActivityFileReady="isActivityFileReady"
                :isWebAssemblySupported="isWebAssemblySupported"
              >
              </TheNavigatorInput>
              <div style="font-size: 0.8em" class="pt-1">Supported files: *.fit, *.gpx, *.tcx</div>
            </div>
            <div class="mb-3" v-if="isActivityFileReady">
              <TheSummary
                :activityFiles="activityFiles"
                :is-activity-file-ready="isActivityFileReady"
                v-on:summary="receiveSummary"
              />
            </div>
            <!-- Tab -->
            <div v-if="isActivityFileReady">
              <!-- Tab Toggler -->
              <ul class="nav nav-tabs ps-2" id="menu" role="tablist">
                <li class="nav-item" role="presentation">
                  <button
                    class="nav-link flat-green active"
                    id="tab1-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#tab1-body"
                    type="button"
                    role="tab"
                    aria-controls="tab1-body"
                    aria-selected="true"
                  >
                    Analysis
                  </button>
                </li>
                <li class="nav-item" role="presentation">
                  <button
                    class="nav-link flat-green"
                    id="tab2-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#tab2-body"
                    type="button"
                    role="tab"
                    aria-controls="tab2-body"
                    aria-selected="false"
                  >
                    Sessions
                  </button>
                </li>
                <li class="nav-item" role="presentation">
                  <button
                    class="nav-link flat-green"
                    id="tab3-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#tab3-body"
                    type="button"
                    role="tab"
                    aria-controls="tab3-body"
                    aria-selected="false"
                  >
                    Laps
                  </button>
                </li>
                <li class="nav-item" role="presentation">
                  <button
                    class="nav-link flat-green"
                    id="tab4-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#tab4-body"
                    type="button"
                    role="tab"
                    aria-controls="tab4-body"
                    aria-selected="false"
                  >
                    Tools
                  </button>
                </li>
              </ul>
              <!-- Tab Toggler Ends -->
              <!-- Tab Content -->
              <div class="tab-content">
                <div
                  class="tab-pane fade show active"
                  id="tab1-body"
                  role="tabpanel"
                  aria-labelledby="tab1-tab"
                >
                  <div class="graph" v-if="hasSpeed">
                    <SpeedGraph
                      :records="combinedRecords"
                      :graph-records="graphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      v-on:hoveredRecord="onHoveredRecord"
                    ></SpeedGraph>
                  </div>
                  <div class="graph" v-if="hasCadence">
                    <CadenceGraph
                      :records="combinedRecords"
                      :graph-records="graphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      v-on:hoveredRecord="onHoveredRecord"
                    ></CadenceGraph>
                  </div>
                  <div class="graph" v-if="hasHeartRate">
                    <HeartRateGraph
                      :records="combinedRecords"
                      :graph-records="graphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      v-on:hoveredRecord="onHoveredRecord"
                    ></HeartRateGraph>
                  </div>
                  <div class="graph" v-if="hasPower">
                    <PowerGraph
                      :records="combinedRecords"
                      :graph-records="graphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      v-on:hoveredRecord="onHoveredRecord"
                    ></PowerGraph>
                  </div>
                  <div class="graph" v-if="hasTemperature">
                    <TemperatureGraph
                      :records="combinedRecords"
                      :graph-records="graphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      v-on:hoveredRecord="onHoveredRecord"
                    ></TemperatureGraph>
                  </div>
                </div>
                <div
                  class="tab-pane fade"
                  id="tab2-body"
                  role="tabpanel"
                  aria-labelledby="tab2-tab"
                >
                  <div v-if="combinedSessions.length > 0">
                    <TheSession :sessions="combinedSessions" />
                  </div>
                </div>
                <div
                  class="tab-pane fade"
                  id="tab3-body"
                  role="tabpanel"
                  aria-labelledby="tab2-tab"
                >
                  <div v-if="combinedLaps.length > 0">
                    <TheLap :laps="combinedLaps" />
                  </div>
                </div>
                <div
                  class="tab-pane fade"
                  id="tab4-body"
                  role="tabpanel"
                  aria-labelledby="tab3-tab"
                >
                  <div class="pt-3">Coming Soon</div>
                </div>
              </div>
              <!-- Tab Content Ends -->
            </div>
            <!-- Tab Ends -->
            <span class="footer pt-3">
              <span class="mx-1">
                <i class="fa-solid fa-copyright fa-rotate-180"></i> {{ new Date().getFullYear() }}
              </span>
              <span class="mx-1">
                <a href="http://github.com/openivity/openivity.github.io" target="_blank">
                  <i class="fa-brands fa-github"></i> Code
                </a>
              </span>
              <div class="mx-1 pt-1">Openivity's Open Source Project</div>
            </span>
          </div>
        </div>

        <div
          id="right"
          class="col-12 col-xxl-9 col-xl-8 col-lg-7 map-container"
          v-if="isActivityFileReady"
        >
          <div class="row">
            <div :class="['col-12', 'map']">
              <TheMap
                :features="features"
                :activity-files="activityFiles"
                :hasCadence="hasCadence"
                :hasHeartRate="hasHeartRate"
                :hasPower="hasPower"
                :hasTemperature="hasTemperature"
                :received-record="hoveredRecord"
                v-on:hoveredRecord="onHoveredRecord"
                ref="theMap"
              />
            </div>
            <div class="col-12 elevation-section">
              <ElevationGraph
                :name="'elev'"
                :has-altitude="hasAltitude"
                :summary="summary"
                :records="combinedRecords"
                :graph-records="graphRecords"
                :received-record="hoveredRecord"
                v-on:hoveredRecord="onHoveredRecord"
              ></ElevationGraph>
            </div>
          </div>
        </div>
        <!-- map & graph end -->
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ActivityFile, Lap, Record, Session } from '@/spec/activity'
import * as d3 from 'd3'
import type { Feature } from 'ol'
import { GeoJSON } from 'ol/format'

const isWebAssemblySupported =
  typeof WebAssembly === 'object' && typeof WebAssembly.instantiateStreaming === 'function'

if (isWebAssemblySupported == false) {
  alert('Sorry, it appears that your browser does not support WebAssembly :(')
}

class Result {
  err: string
  took: number
  decodeResults: Array<DecodeResult>

  constructor(json?: any) {
    const casted = json as Result

    this.err = casted?.err
    this.took = casted?.took
    this.decodeResults = casted?.decodeResults
  }
}

class DecodeResult {
  activityFile: ActivityFile
  err: string

  constructor(json?: any) {
    const casted = json as DecodeResult

    this.activityFile = new ActivityFile(casted?.activityFile)
    this.err = casted?.err
  }
}

export default {
  data() {
    return {
      loading: false,
      decodeWorker: new Worker(new URL('@/workers/fitsvc-decode.ts', import.meta.url), {
        type: 'module'
      }),
      decodeBeginTimestamp: 0,
      activityFiles: new Array<ActivityFile>(),
      features: new Array<Feature>(),
      combinedRecords: new Array<Record>(),
      combinedSessions: new Array<Session>(),
      combinedLaps: new Array<Lap>(),
      graphRecords: new Array<Record>(),
      summary: new Summary(),
      hoveredRecord: new Record()
    }
  },
  computed: {
    isActivityFileReady: function () {
      return this.activityFiles.length > 0
    },
    hasAltitude(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        if (this.combinedRecords[i].altitude != null) return true
      }
      return false
    },
    hasSpeed(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        if (this.combinedRecords[i].speed) return true
      }
      return false
    },
    hasCadence(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        if (this.combinedRecords[i].cadence) return true
      }
      return false
    },
    hasHeartRate(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        if (this.combinedRecords[i].heartRate) return true
      }
      return false
    },
    hasPower(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        if (this.combinedRecords[i].power) return true
      }
      return false
    },
    hasTemperature(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        if (this.combinedRecords[i].temperature != null) return true
      }
      return false
    }
  },
  methods: {
    fileInputEventListener(e: Event) {
      const fileInput = e.target as HTMLInputElement
      if (!fileInput.files) return

      let promisers = new Array<Promise<Uint8Array>>()
      for (let i = 0; i < fileInput.files?.length!; i++) {
        promisers.push(
          new Promise<Uint8Array>((resolve, reject) => {
            const selectedFile = (fileInput.files as FileList)[i]
            if (!selectedFile) {
              return
            }

            const reader = new FileReader()
            reader.onload = (e: ProgressEvent<FileReader>) => {
              const fileData = e.target!.result as ArrayBuffer
              resolve(new Uint8Array(fileData))
            }
            reader.onerror = reject
            reader.readAsArrayBuffer(selectedFile as File)
          })
        )
      }

      Promise.all(promisers).then((arr) => {
        this.decodeBeginTimestamp = new Date().getTime()
        this.loading = true
        this.decodeWorker.postMessage(arr)
      })
    },
    decodeWorkerOnMessage(e: MessageEvent) {
      const result = e.data as Result
      if (result.err != '') {
        console.error(`decode return with err: ${result.err}`)
        alert(`decode return with err: ${result.err}`)
        this.loading = false
        return
      }

      const totalDuration = new Date().getTime() - this.decodeBeginTimestamp
      console.group('Decoding')
      console.debug('Decode took:\t\t', result.took, 'ms')
      console.debug('Interop WASM to JS:\t', totalDuration - result.took, 'ms')
      console.debug('Total elapsed:\t\t', totalDuration, 'ms')
      console.groupEnd()

      requestAnimationFrame(() => {
        this.preprocessingResults(result.decodeResults)
        this.scrollTop()
      })
    },
    async preprocessingResults(decodeResults: DecodeResult[]) {
      const activityFiles = new Array<ActivityFile>()

      const begin = new Date()

      let lastDistance = 0
      for (let i = 0; i < decodeResults.length; i++) {
        let records = decodeResults[i].activityFile.records

        const filteredRecords: Record[] = []
        for (let i = 0; i < records.length; i++) {
          const rec = records[i]
          if (typeof rec.distance !== 'number') continue
          rec.distance! += lastDistance
          filteredRecords.push(rec)
        }

        records = filteredRecords
        lastDistance = d3.max(records.map((d2) => d2.distance!))!

        records = this.smoothingElevation(records, 30)
        records = this.calculateGrade(records, 100)

        decodeResults[i].activityFile.records = records

        activityFiles.push(decodeResults[i].activityFile)
      }

      this.activityFiles = activityFiles
      this.features = this.createFeatures(activityFiles)
      this.combinedRecords = this.activityFiles.flatMap((d) => d.records)
      this.createLapsIfNotExist(activityFiles)
      this.combinedLaps = this.activityFiles.flatMap((d) => d.laps)
      this.combinedSessions = this.activityFiles.flatMap((d) => d.sessions)
      this.graphRecords = this.summarizeRecords(this.combinedRecords, 100)

      setTimeout(() => (this.loading = false), 200)

      const elapsed = new Date().getTime() - begin.getTime()
      console.debug('Preprocessing Results:\t\t', elapsed, 'ms')
    },
    createLapsIfNotExist(activityFiles: ActivityFile[]) {
      activityFiles.forEach((d) => {
        if (d.laps.length > 0) return
        const laps: Lap[] = []
        d.sessions.forEach((ses) => laps.push(new Lap(ses)))
        d.laps = laps
      })
    },
    createFeatures(activityFiles: ActivityFile[]): Feature[] {
      const features: Feature[] = []
      activityFiles.forEach((act, i) => {
        const coordinates: number[][] = []
        act.records.forEach((d) => {
          if (d.positionLong == null || d.positionLat == null) return
          coordinates.push([d.positionLong!, d.positionLat!])
        })

        if (coordinates.length == 0) return

        const feature = new GeoJSON().readFeature({
          id: 'lineString-' + i,
          type: 'Feature',
          geometry: {
            type: 'LineString',
            coordinates: coordinates
          }
        })
        features.push(feature)
      })
      return features
    },
    smoothingElevation(records: Record[], distance: number): Record[] {
      // Simple Moving Average
      for (let i = 1; i < records.length; i++) {
        if (records[i].distance! < distance) continue
        if (typeof records[i].altitude !== 'number') continue

        let sum = 0
        let counter = 0
        for (let j = i; j >= 0; j--) {
          if (typeof records[j].altitude !== 'number') continue

          const d = records[i].distance! - records[j].distance!
          if (d > distance) break
          sum += records[j].altitude!
          counter++
        }

        records[i].altitude = sum / counter
      }
      return records
    },
    calculateGrade(records: Record[], distance: number): Record[] {
      for (let i = 0; i < records.length - 1; i++) {
        if (typeof records[i].altitude !== 'number') continue

        let run = 0
        let rise = 0
        for (let j = i + 1; j < records.length; j++) {
          if (typeof records[j].altitude !== 'number') continue

          const d = records[j].distance! - records[i].distance!
          if (d > distance) break
          rise = records[j].altitude! - records[i].altitude!
          run = d
        }

        records[i].grade = (rise / run) * 100
      }

      return records
    },
    summarizeRecords(records: Record[], distance: number): Record[] {
      if (records.length < 1000) return records

      let altitudes: number[] = []
      let cadences: number[] = []
      let heartRates: number[] = []
      let speeds: number[] = []
      let powers: number[] = []
      let temps: number[] = []
      let grades: number[] = []

      const summarized: Record[] = []
      let cur = 0
      for (let i = 0; i < records.length; i++) {
        const d = records[i].distance! - records[cur].distance!
        if (d > distance) {
          let record = new Record({
            timestamp: records[i].timestamp,
            distance: records[i].distance,
            positionLat: records[i].positionLat,
            positionLong: records[i].positionLong
          })

          if (altitudes.length > 0)
            record.altitude = altitudes.reduce((sum, val) => sum + val, 0) / altitudes.length
          if (cadences.length > 0)
            record.cadence = cadences.reduce((sum, val) => sum + val, 0) / cadences.length
          if (grades.length > 0)
            record.grade = grades.reduce((sum, val) => sum + val, 0) / grades.length
          if (heartRates.length > 0)
            record.heartRate = heartRates.reduce((sum, val) => sum + val, 0) / heartRates.length
          if (powers.length > 0)
            record.power = powers.reduce((sum, val) => sum + val, 0) / powers.length
          if (speeds.length > 0)
            record.speed = speeds.reduce((sum, val) => sum + val, 0) / speeds.length

          if (temps.length > 0)
            record.temperature = temps.reduce((sum, val) => sum + val, 0) / temps.length

          summarized.push(record)

          // Reset array
          altitudes = []
          cadences = []
          heartRates = []
          speeds = []
          powers = []
          temps = []
          grades = []

          cur = i + 1
        }

        if (records[i].altitude) {
          altitudes.push(records[i].altitude!)
        }
        if (records[i].cadence) {
          cadences.push(records[i].cadence!)
        }
        if (records[i].heartRate) {
          heartRates.push(records[i].heartRate!)
        }
        if (records[i].speed) {
          speeds.push(records[i].speed!)
        }
        if (records[i].power) {
          powers.push(records[i].power!)
        }
        if (records[i].temperature) {
          temps.push(records[i].temperature!)
        }
        if (records[i].grade) {
          grades.push(records[i].grade!)
        }
      }

      return summarized
    },
    scrollTop() {
      document.body.scrollTop = 0 // For Safari
      document.documentElement.scrollTop = 0 // For Chrome, Firefox, IE and Opera
    },
    receiveSummary(summary: Summary) {
      this.summary = summary
    },
    onHoveredRecord(record: Record) {
      this.hoveredRecord = new Record(record)
    }
  },
  mounted() {
    document.getElementById('fileInput')?.addEventListener('change', this.fileInputEventListener)
    this.decodeWorker.onmessage = this.decodeWorkerOnMessage
  }
}
</script>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 100ms ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

.map-container {
  padding-left: 0;
}

.map {
  height: 73vh;
}

.elevation-section {
  height: 27vh;
}

.sidebar {
  scrollbar-width: thin;
  scrollbar-gutter: stable;
  overflow-y: scroll;
  overflow-x: hidden;
  height: 100vh;
}

.graph {
  padding: 10px;
  border-bottom: 0.5rem solid var(--color-background-soft);
}

.graph:last-child {
  border-bottom: unset !important;
}

.landing {
  display: flex;
  justify-content: center;
  align-items: center;
}

.footer {
  display: inline-block;
  height: 70px;
  font-size: 0.8em;
  color: var(--green-text);
}

.footer a {
  color: var(--green-text);
}

.footer div {
  font-size: 0.9em;
  color: var(--color-text);
}

@media (pointer: coarse) {
  /* mobile device */

  .default-border-top {
    border-top: 0.5rem solid var(--color-background-soft);
  }

  .flex-container {
    display: flex;
    flex-direction: column;
  }

  .sidebar {
    overflow: unset;
    height: 80vh;
  }

  #left {
    order: 2;
  }
  #right {
    order: 1;
  }

  .activity {
    height: 100%;
  }

  .map-container {
    padding-left: 5px;
    height: 100%;
  }

  .elevation-section {
    padding-left: 20px;
  }

  .map {
    padding: 0;
    height: 50vh;
  }

  .elevation-section {
    height: 35vh;
  }
}

@media (pointer: fine), (pointer: none) {
  /* desktop */
}

@media (pointer: fine) and (any-pointer: coarse) {
  /* touch desktop */
}
</style>
