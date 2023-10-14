<script setup lang="ts">
import TheMap from './TheMap.vue'
import ElevationGraphPlot from './ElevationGraphPlot.vue'
import TheNavigator from './TheNavigator.vue'
import Loading from './Loading.vue'
</script>

<template>
  <div class="container-fluid text-center">
    <Transition>
      <Loading v-show="loading"></Loading>
    </Transition>
    <div :class="['row h-100', !isActivityFileReady ? 'align-items-center' : '']">
      <!-- left navigator -->
      <div
        :class="['navigator', 'col-12', isActivityFileReady ? 'col-xl-3 col-md-5' : 'col-md-12']"
      >
        <div class="header">
          <h2 class="title">Open Activity</h2>
          <p class="h6 subtitle">by Openivity</p>
        </div>
        <TheNavigator
          :activityFiles="activityFiles"
          :timezoneOffsetHour="timezoneOffsetHour"
          :isWebAssemblySupported="isWebAssemblySupported"
        />
      </div>
      <!-- map & graph -->
      <div class="col-12 col-md-7 col-xl-9 map-container" v-show="isActivityFileReady">
        <div class="row">
          <!-- the map -->
          <div :class="['col-12', 'map']">
            <TheMap
              :geojsons="geojsons"
              :activity-files="activityFiles"
              :timezoneOffsetHoursList="timezoneOffsetHoursList"
              ref="theMap"
            />
          </div>
          <!-- the bottom graph -->
          <div :class="['col-12', 'bottom-info-button']">
            <nav class="position-relative bottom-info-nav">
              <div class="nav nav-tabs" role="tablist">
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
</template>

<script lang="ts">
import { GeoJSON } from 'ol/format'
import { ActivityFile } from '@/spec/activity'
import { LinierRegression, Point } from '@/toolkit/linier-regression'

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
      begin: 0
      //
    }
  },
  computed: {
    isActivityFileReady: function () {
      return this.activityFiles && this.activityFiles.length > 0
    }
  },
  watch: {
    activityFiles: {
      async handler(activityFiles: ActivityFile[]) {
        const tzOffsetHours = new Array<number>()
        for (let i = 0; i < activityFiles.length; i++) {
          if (!activityFiles[i].activity?.timestamp || !activityFiles[i].activity?.localDateTime)
            return

          const localDateTime = new Date(activityFiles[i].activity!.localDateTime!)
          const timestamp = new Date(activityFiles[i].activity!.timestamp!)
          const tzOffsetMillis = localDateTime.getTime() - timestamp.getTime()
          const tzOffsetHour = Math.floor(tzOffsetMillis / 1000 / 3600)

          tzOffsetHours.push(tzOffsetHour)
        }

        this.timezoneOffsetHoursList = tzOffsetHours
        this.timezoneOffsetHour = tzOffsetHours[0]
        console.log('timezone offsets:', tzOffsetHours)
      }
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

      const geoJSONList = new Array<GeoJSON>()
      for (let i = 0; i < decodeResults.length; i++) {
        geoJSONList.push(decodeResults[i].feature)
      }

      this.geojsons = geoJSONList
      const activityFileList = new Array<ActivityFile>()

      const samplingDistance = 100
      for (let i = 0; i < decodeResults.length; i++) {
        this.calculateGradePercentage(decodeResults[i].activityFile, samplingDistance)
        activityFileList.push(decodeResults[i].activityFile)
      }
      this.activityFiles = activityFileList
      this.loading = false
    },

    calculateGradePercentage(activityFile: ActivityFile, samplingDistance: number) {
      const n = activityFile.records.length

      for (let i = 0; i < n; i++) {
        const points = new Array<Point>()
        let distance: number = 0
        let j = i

        points.push(new Point(activityFile.records[i].distance!, activityFile.records[i].altitude!))
        while (distance < samplingDistance && j > 0) {
          distance += activityFile.records[j].distance! - activityFile.records[j - 1].distance!
          points.push(
            new Point(activityFile.records[j].distance!, activityFile.records[j].altitude!)
          )
          j--
        }

        if (points.length == 1 && i - 1 > 0) {
          // In case the 'samplingDistance' value is less than the actual distance being sampled in the file,
          // add the previous record's value as the pivot.
          points.push(
            new Point(activityFile.records[i - 1].distance!, activityFile.records[i - 1].altitude!)
          )
        }

        const currentAltitude = activityFile.records[i].altitude!
        const linierRegression = new LinierRegression()
        linierRegression.train(points)
        const pivotAltitude = linierRegression.predictY(points[points.length - 1].X) // smoothing the value.

        const rise = currentAltitude - pivotAltitude
        const run = Math.sqrt(Math.pow(distance, 2) - Math.pow(rise, 2))
        const grade = (rise / run) * 100

        activityFile.records[i].grade = grade
      }
    },
    elevationOnRecord(record: any) {
      //@ts-ignore
      this.$refs.theMap.showPopUpRecord(record)
    }
  },
  mounted() {
    document.getElementById('fileInput')?.addEventListener('change', this.fileInputEventListener)

    this.decodeWorker.onmessage = this.decodeWorkerOnMessage
  }
}
</script>

<style>
.v-enter-active,
.v-leave-active {
  transition: opacity 100ms ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

.container-fluid {
  height: 100vh;
}

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
