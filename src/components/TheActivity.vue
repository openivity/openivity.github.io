<script setup lang="ts">
import TheMap from './TheMap.vue'
import ElevationGraphPlot from './ElevationGraphPlot.vue'
import TheNavigator from './TheNavigator.vue'
import Loading from './Loading.vue'

import { onMounted } from 'vue'

const isWebAssemblySupported = ref(false)

onMounted(() => {
  isWebAssemblySupported.value =
    typeof WebAssembly === 'object' && typeof WebAssembly.instantiateStreaming === 'function'

  if (isWebAssemblySupported.value == false) {
    alert('Sorry, it appears that your browser does not support WebAssembly :(')
    return
  }

  document.getElementById('fileInput')?.addEventListener('change', (e) => {
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
      begin.value = new Date().getTime()
      loading.value = true
      decodeWorker.postMessage(arr)
    })
  })
})
</script>

<template>
  <div class="container">
    <Transition>
      <Loading v-show="loading"></Loading>
    </Transition>
    <div
      :class="activityFiles.length > 0 ? 'navigator-left' : 'navigator-center'"
      class="navigator"
    >
      <div class="header"><h2 class="title">Open Activity</h2></div>
      <TheNavigator
        :activityFiles="activityFiles"
        :timezoneOffsetHours="timezoneOffsetHours"
        :isWebAssemblySupported="isWebAssemblySupported"
      />
      <div class="graph" v-show="activityFiles && activityFiles.length > 0">
        <ElevationGraphPlot :activityFile="activityFiles[0]" />
      </div>
    </div>
    <div class="map" v-show="activityFiles && activityFiles.length > 0">
      <TheMap
        :geojsons="geojsons"
        :activity-files="activityFiles"
        :timezoneOffsetHoursList="timezoneOffsetHoursList"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { ref, watch } from 'vue'
import { GeoJSON } from 'ol/format'
import { ActivityFile } from '@/spec/activity'
import { LinierRegression, Point } from '@/toolkit/linier-regression'

const decodeWorker = new Worker(new URL('@/workers/fitsvc-decode.ts', import.meta.url), {
  type: 'module'
})

const geojsons = ref(new Array<GeoJSON>())
const activityFiles = ref(new Array<ActivityFile>())
const timezoneOffsetHours = ref(0)
const timezoneOffsetHoursList = ref(new Array<Number>())
const loading = ref(false)
const begin = ref(0)

watch(activityFiles, async (activityFiles: Array<ActivityFile>) => {
  const timezoneOffsetHours = new Array<Number>()
  for (let i = 0; i < activityFiles.length; i++) {
    if (!activityFiles[i].activity?.timestamp || !activityFiles[i].activity?.localDateTime) return

    const localDateTime = new Date(activityFiles[i].activity!.localDateTime!)
    const timestamp = new Date(activityFiles[i].activity!.timestamp!)
    const tzOffsetMillis = localDateTime.getTime() - timestamp.getTime()
    const tzOffsetHours = Math.floor(tzOffsetMillis / 1000 / 3600)

    timezoneOffsetHours.push(tzOffsetHours)
  }

  timezoneOffsetHoursList.value = timezoneOffsetHours
  console.log('timezone offsets:', timezoneOffsetHours)
})

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

decodeWorker.onmessage = (e) => {
  const result = e.data as Result
  if (result.err != '') {
    console.error(`decode return with err: ${result.err}`)
    alert(`decode return with err: ${result.err}`)
    loading.value = false
    return
  }

  const decodeResults = new Array<DecodeResult>()
  for (let i = 0; i < result.decodeResults.length; i++) {
    decodeResults[i] = new DecodeResult(result.decodeResults[i])
  }

  const totalDuration = new Date().getTime() - begin.value
  console.group('Elapsed')
  console.log('Decode took:\t\t', result.took, 'ms')
  console.log('Interop wasm to js:\t', totalDuration - result.took, 'ms')
  console.log('Total elapsed:\t\t', totalDuration, 'ms')
  console.groupEnd()

  const geoJSONList = new Array<GeoJSON>()
  for (let i = 0; i < decodeResults.length; i++) {
    geoJSONList.push(decodeResults[i].feature)
  }

  geojsons.value = geoJSONList
  const activityFileList = new Array<ActivityFile>()

  const samplingDistance = 100
  for (let i = 0; i < decodeResults.length; i++) {
    calculateGradePercentage(decodeResults[i].activityFile, samplingDistance)
    activityFileList.push(decodeResults[i].activityFile)
  }
  activityFiles.value = activityFileList
  loading.value = false
}

function calculateGradePercentage(activityFile: ActivityFile, samplingDistance: number) {
  const n = activityFile.records.length

  for (let i = 0; i < n; i++) {
    const points = new Array<Point>()
    let distance: number = 0
    let j = i

    points.push(new Point(activityFile.records[i].distance!, activityFile.records[i].altitude!))
    while (distance < samplingDistance && j > 0) {
      distance += activityFile.records[j].distance! - activityFile.records[j - 1].distance!
      points.push(new Point(activityFile.records[j].distance!, activityFile.records[j].altitude!))
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
}

@media (pointer: coarse) {
  /* mobile device */

  .container {
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
  }
}

@media (pointer: fine), (pointer: none) {
  /* desktop */
}

@media (pointer: fine) and (any-pointer: coarse) {
  /* touch desktop */
}
</style>
