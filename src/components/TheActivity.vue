<script setup lang="ts">
import TheMap from './TheMap.vue'
import ElevationGraph from './ElevationGraph.vue'
import TheNavigator from './TheNavigator.vue'
</script>

<template>
  <div class="container">
    <div class="map">
      <TheMap :geojson="geojson" :slider="slider" :activityFile="activityFile" />
      <input type="range" min="1" max="10000" v-model.number="slider" />
      <ElevationGraph />
    </div>
    <div class="navigator">
      <TheNavigator :activityFile="activityFile" />
    </div>
  </div>
</template>

<script lang="ts">
const isWebAssemblySupported =
  typeof WebAssembly === 'object' && typeof WebAssembly.instantiateStreaming === 'function'

if (isWebAssemblySupported == false) {
  alert('Sorry, it appears that your browser does not support WebAssembly :(')
}

import { ref, watch } from 'vue'
import { GeoJSON } from 'ol/format'
import '@/assets/wasm_exec.js'
import { ActivityFile } from '@/spec/activity'

const geojson = ref(new GeoJSON())
const activityFile = ref(new ActivityFile())
const slider = ref(0)
const byteArray = ref(new Uint8Array())

const go = new Go()

class Result {
  feature: any
  activityFile: ActivityFile
  err: string

  constructor(json?: any) {
    const casted = json as Result

    this.feature = casted?.feature
    this.activityFile = new ActivityFile(casted?.activityFile)
    this.err = casted?.err
  }
}

const wasmUrl = 'wasm/fitsvc.wasm'

WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject).then((wasm) => {
  go.run(wasm.instance)

  watch(byteArray, (value: Uint8Array) => {
    //@ts-ignore
    const rawResult = decode(value)

    const begin = new Date().getTime()
    const result: Result = new Result(rawResult)
    console.log('js: deserialization took: ', new Date().getTime() - begin, 'ms')

    geojson.value = result.feature
    activityFile.value = result.activityFile

    // console.log(lookupRecord)
  })

  document.getElementById('fileInput')?.addEventListener('change', (e) => {
    const fileInput = e.target as HTMLInputElement
    const selectedFile = (fileInput.files as FileList)[0]
    if (!selectedFile) {
      return
    }

    const reader = new FileReader()

    reader.onload = (e: ProgressEvent<FileReader>) => {
      const fileData = e.target!.result as ArrayBuffer
      byteArray.value = new Uint8Array(fileData)
    }

    reader.readAsArrayBuffer(selectedFile as File)
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
</style>
