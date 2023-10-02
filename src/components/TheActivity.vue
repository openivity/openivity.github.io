<script setup lang="ts">
import TheMap from './TheMap.vue'
import ElevationGraph from './ElevationGraph.vue'
import TheNavigator from './TheNavigator.vue'
</script>

<template>
  <div class="container">
    <div class="map">
      <TheMap
        :geojsons="geojsons"
        :activityFiles="activityFiles"
        :timezoneOffsetHoursList="timezoneOffsetHoursList"
      />
      <ElevationGraph />
    </div>
    <div class="navigator">
      <div class="header"><h2 class="title">Open Activity</h2></div>
      <TheNavigator :activityFiles="activityFiles" :timezoneOffsetHours="timezoneOffsetHours" />
    </div>
  </div>
</template>

<script lang="ts">
const isWebAssemblySupported =
  typeof WebAssembly === 'object' && typeof WebAssembly.instantiateStreaming === 'function'

if (isWebAssemblySupported == false) {
  alert('Sorry, it appears that your browser does not support WebAssembly :(')
}

import '@/assets/wasm/wasm_exec.js'
import { ref, watch } from 'vue'
import { GeoJSON } from 'ol/format'
import { ActivityFile } from '@/spec/activity'

const geojsons = ref(new Array<GeoJSON>())
const activityFiles = ref(new Array<ActivityFile>())
const timezoneOffsetHours = ref(0)
const timezoneOffsetHoursList = ref(new Array<Number>())
const byteArrays = ref(new Array<Uint8Array>())

watch(activityFiles, (activityFiles: Array<ActivityFile>) => {
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

const go = new Go()

class Result {
  err: string
  decodeResults: Array<DecodeResult>

  constructor(json?: any) {
    const casted = json as Result

    this.err = casted?.err
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

const wasmUrl = 'wasm/fitsvc.wasm'

WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject).then((wasm) => {
  go.run(wasm.instance)

  watch(byteArrays, (values: Array<Uint8Array>) => {
    const begin = new Date().getTime()

    //@ts-ignore
    const rawResult = decode(values)
    const result = rawResult as Result
    if (result.err != '') {
      alert(`decode return with err: ${result.err}`)
      return
    }

    const decodeResults = new Array<DecodeResult>()
    for (let i = 0; i < result.decodeResults.length; i++) {
      decodeResults[i] = new DecodeResult(result.decodeResults[i])
    }

    console.log('js: e2e decode took: ', new Date().getTime() - begin, 'ms')

    const gjs = new Array<GeoJSON>()
    for (let i = 0; i < decodeResults.length; i++) {
      gjs.push(decodeResults[i].feature)
    }

    geojsons.value = gjs
    const afs = new Array<ActivityFile>()
    for (let i = 0; i < decodeResults.length; i++) {
      afs.push(decodeResults[i].activityFile)
    }
    activityFiles.value = afs
  })

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

    Promise.all(promisers).then((arr) => (byteArrays.value = arr))
  })
})
</script>

<style>
.container {
  max-width: 1280px;
  width: 100vw;
  height: 100vh;
  display: grid;
  grid-template-columns: 70% 30%;
}

.map {
  height: 80vh;
  grid-column: 1;
  grid-row: 1;
}

.navigator {
  overflow: auto;
  grid-column: 2;
  grid-row: 1;
  padding: 1rem;
}

.header {
  text-align: center;
}

@media (pointer: coarse) {
  /* mobile device */

  .container {
    overflow: visible;
    height: unset;
  }

  .map {
    height: 350px;
    width: 100vw;
    grid-column: 1;
    grid-row: 2;
  }

  .navigator {
    width: 100vw;
    overflow: unset;
    grid-column: 1;
    grid-row: 3;
    padding: 0;
  }

  .header {
    margin: 10px auto;
  }
}

@media (pointer: fine), (pointer: none) {
  /* desktop */
}

@media (pointer: fine) and (any-pointer: coarse) {
  /* touch desktop */
}
</style>
