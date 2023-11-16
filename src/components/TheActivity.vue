<script setup lang="ts">
import type { Feature } from 'ol'
import ElevationGraphPlot from './ElevationGraphPlot.vue'
import TheLoading from './TheLoading.vue'
import TheMap from './TheMap.vue'
import TheNavigator from './TheNavigator.vue'
import TheNavigatorInput from './TheNavigatorInput.vue'
</script>

<template>
  <div>
    <Transition>
      <TheLoading v-show="loading"></TheLoading>
    </Transition>
    <!-- mobile summary -->
    <div
      class="offcanvas offcanvas-start"
      tabindex="-1"
      id="offcanvasSummary"
      aria-labelledby="offcanvasSummaryLabel"
    >
      <div class="offcanvas-header">
        <h5 class="offcanvas-title title" id="offcanvasSummaryLabel">Open Activity</h5>
        <button
          type="button"
          class="btn-close"
          data-bs-dismiss="offcanvas"
          aria-label="Close"
        ></button>
      </div>
      <div class="offcanvas-body">
        <TheNavigatorInput :isWebAssemblySupported="isWebAssemblySupported" id="fileInputMobile">
        </TheNavigatorInput>
        <TheNavigator :activityFiles="activityFiles" :timezoneOffsetHour="timezoneOffsetHour" />
      </div>
    </div>
    <div class="activity container-fluid text-center">
      <div :class="['row h-100', !isActivityFileReady ? 'align-items-center' : '']">
        <!-- navigator but desktop -->
        <div
          :class="[
            'navigator',
            'col-12',
            isActivityFileReady ? 'col-xl-3 col-md-5 d-none d-md-block' : 'col-md-12'
          ]"
        >
          <div class="header pt-5 pb-1">
            <h2 class="title">Open Activity</h2>
          </div>
          <TheNavigatorInput :isWebAssemblySupported="isWebAssemblySupported"> </TheNavigatorInput>
          <TheNavigator :activityFiles="activityFiles" :timezoneOffsetHour="timezoneOffsetHour" />
        </div>
        <!-- map & graph -->
        <div class="col-12 col-md-7 col-xl-9 map-container" v-show="isActivityFileReady">
          <div class="row">
            <!-- the map -->
            <div :class="['col-12', 'map']">
              <TheMap
                :features="features"
                :activity-files="activityFiles"
                :timezoneOffsetHoursList="timezoneOffsetHoursList"
                :hasCadence="hasCadence"
                :hasHeartRate="hasHeartRate"
                :hasPower="hasPower"
                :hasTemperature="hasTemperature"
                ref="theMap"
              />
            </div>
            <!-- the bottom graph -->
            <div :class="['col-12', 'bottom-info-button']">
              <nav class="position-relative bottom-info-nav">
                <div class="nav nav-tabs" role="tablist">
                  <!-- show only on mobile -->
                  <button
                    class="nav-link d-md-none"
                    type="button"
                    role="tab"
                    data-bs-toggle="offcanvas"
                    data-bs-target="#offcanvasSummary"
                    aria-controls="offcanvasSummary"
                    aria-label="Toggle Summary"
                  >
                    <i class="fa-solid fa-bar"></i> Summary
                  </button>
                  <button
                    class="nav-link active"
                    id="nav-elevation-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#nav-graph"
                    type="button"
                    role="tab"
                    aria-controls="nav-graph"
                    aria-selected="true"
                  >
                    Elevation
                  </button>
                </div>
              </nav>
            </div>
            <div :class="['col-12', 'bottom-info']">
              <div class="tab-content h-100">
                <div
                  :class="['tab-pane fade show active h-100', 'graph']"
                  id="nav-graph"
                  role="tabpanel"
                  aria-labelledby="nav-elevation-tab"
                  tabindex="0"
                >
                  <ElevationGraphPlot
                    :activityFiles="activityFiles"
                    :activityTimezoneOffset="timezoneOffsetHoursList"
                    v-on:record="elevationOnRecord"
                  />
                </div>
                <div
                  :class="['tab-pane fade h-100', 'graph']"
                  id="nav-hr"
                  role="tabpanel"
                  aria-labelledby="nav-hr-tab"
                  tabindex="0"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ActivityFile, Record } from '@/spec/activity'
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
  feature: any
  activityFile: ActivityFile
  err: string

  constructor(json?: any) {
    const casted = json as DecodeResult

    this.feature = casted?.feature
    this.activityFile = new ActivityFile(casted?.activityFile)
    this.err = casted?.err
  }
}

export default {
  data() {
    return {
      decodeWorker: new Worker(new URL('@/workers/fitsvc-decode.ts', import.meta.url), {
        type: 'module'
      }),
      geojsons: new Array<GeoJSON>(),
      activityFiles: new Array<ActivityFile>(),
      timezoneOffsetHour: 0,
      timezoneOffsetHoursList: new Array<Number>(),
      loading: false,
      begin: 0,
      features: new Array<Feature>()
    }
  },
  computed: {
    isActivityFileReady: function () {
      return this.activityFiles && this.activityFiles.length > 0
    },
    combinedRecords: function (): Record[] {
      return this.activityFiles?.flatMap((act) => act.records)
    },
    hasCadence(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        const cad = this.combinedRecords[i].cadence
        if (typeof cad === 'number' && cad != 255 && cad != 0) return true
      }
      return false
    },
    hasHeartRate(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        const hr = this.combinedRecords[i].heartRate
        if (typeof hr === 'number' && hr != 255 && hr != 0) return true
      }
      return false
    },
    hasPower(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        const pwr = this.combinedRecords[i].power
        if (typeof pwr === 'number' && pwr != 65535 && pwr != 0) return true
      }
      return false
    },
    hasTemperature(): boolean {
      for (let i = 0; i < this.combinedRecords.length; i++) {
        const temp = this.combinedRecords[i].temperature
        if (typeof temp === 'number' && temp != 127) return true
      }
      return false
    }
  },
  watch: {
    activityFiles: {
      async handler(activityFiles: ActivityFile[]) {
        const tzOffsetHours = new Array<number>()
        for (let i = 0; i < activityFiles.length; i++) {
          tzOffsetHours.push(activityFiles[i].timezone)
        }

        this.timezoneOffsetHoursList = tzOffsetHours
        this.timezoneOffsetHour = tzOffsetHours[0]
        console.log('timezone offsets:', tzOffsetHours)
      }
    }
  },
  methods: {
    createFeatures(activityFiles: ActivityFile[]): Feature[] {
      const features: Feature[] = []
      activityFiles.forEach((act, i) => {
        const coordinates: number[][] = []
        act.records.forEach((d) => {
          if (d.positionLong == null || d.positionLat == null) return
          coordinates.push([d.positionLong!, d.positionLat!])
        })
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
        this.begin = new Date().getTime()
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

      const decodeResults = new Array<DecodeResult>()
      for (let i = 0; i < result.decodeResults.length; i++) {
        decodeResults[i] = new DecodeResult(result.decodeResults[i])
      }

      const totalDuration = new Date().getTime() - this.begin
      console.group('Elapsed')
      console.log('Decode took:\t\t', result.took, 'ms')
      console.log('Interop wasm to js:\t', totalDuration - result.took, 'ms')
      console.log('Total elapsed:\t\t', totalDuration, 'ms')
      console.groupEnd()

      const activityFileList = new Array<ActivityFile>()

      for (let i = 0; i < decodeResults.length; i++) {
        let records = decodeResults[i].activityFile.records
        records = this.smoothingElevation(records, 30)
        records = this.calculateGrade(records, 100)
        decodeResults[i].activityFile.records = records
        activityFileList.push(decodeResults[i].activityFile)
      }

      this.activityFiles = activityFileList
      this.features = this.createFeatures(activityFileList)

      this.loading = false
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
    elevationOnRecord(record: any) {
      //@ts-ignore
      this.$refs.theMap.showPopUpRecord(record)
    }
  },
  mounted() {
    document.getElementById('fileInput')?.addEventListener('change', this.fileInputEventListener)
    document
      .getElementById('fileInputMobile')
      ?.addEventListener('change', this.fileInputEventListener)

    this.decodeWorker.onmessage = this.decodeWorkerOnMessage
  }
}
</script>

<style scoped lang="scss">
// animation
.v-enter-active,
.v-leave-active {
  transition: opacity 100ms ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

.activity {
  height: 100vh;
}

// summary sidebar
.offcanvas.offcanvas-start {
  width: 100%;
}

// app
.map-container {
  height: 100vh;
}

.map {
  height: 70vh;
}

.bottom-info-button {
  height: 0;
}

.bottom-info {
  height: 30vh;
}

.bottom-info-nav {
  top: -42px;
  height: 0;
}

.bottom-info-nav .nav-link {
  background-color: var(--bs-body-bg);
  color: var(--bs-body-color);
}
.bottom-info-nav .nav-link.active {
  background-color: var(--color-title);
}

/* 
.container {
  display: grid;
  grid-template-columns: 40% 60%;
  height: 100vh;
  position: relative;
  width: 100vw;
}

.map {
  grid-column: 2;
  grid-row: 1;
  height: 100vh;
}

.navigator-center {
  grid-column: 1;
  grid-column-end: 3;
  margin: auto;
}

.navigator-left {
  grid-column: 1;
  margin: 0 auto;
  padding-top: 10px;
}

.navigator {
  grid-row: 1;
  overflow: auto;
  width: 100%;
}

.header {
  text-align: center;
  margin: 10px auto;
} */

@media (pointer: coarse) {
  /* mobile device */

  /* .container {
    overflow: visible;
    height: unset;
  }

  .map {
    grid-column: 1;
    grid-row: 2;
    height: 65vh;
    width: 100vw;
  }

  .navigator {
    grid-column: 1;
    grid-row: 3;
    margin: auto;
    overflow: unset;
    width: 100vw;
  } */
}

@media (pointer: fine), (pointer: none) {
  /* desktop */
}

@media (pointer: fine) and (any-pointer: coarse) {
  /* touch desktop */
}
</style>
