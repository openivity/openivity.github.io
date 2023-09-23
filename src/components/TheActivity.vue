<script setup lang="ts">
import TheMap from './TheMap.vue'
import ElevationGraph from './ElevationGraph.vue'
import TheNavigator from './TheNavigator.vue'
</script>

<template>
  <div class="container">
    <div class="map">
      <TheMap :geojson="geojson" />
      <ElevationGraph />
    </div>
    <div class="navigator">
      <TheNavigator :sessions="sessions" />
    </div>
  </div>
</template>

<script lang="ts">
import { ref } from 'vue'
import { GeoJSON } from 'ol/format'
import '@/assets/wasm_exec.js'
import { Session } from '@/spec/activity'

const geojson = ref(new GeoJSON())
const sessions = ref(new Array<Session>())

const go = new Go()

declare class Result {
  err: string
  feature: any
  activityFile: any
}

WebAssembly.instantiateStreaming(fetch('wasm/fitsvc.wasm'), go.importObject).then((result) => {
  go.run(result.instance)

  document.getElementById('fileInput')?.addEventListener('change', (e) => {
    const fileInput = e.target as HTMLInputElement
    const selectedFile = (fileInput.files as FileList)[0]
    if (!selectedFile) {
      return
    }

    const reader = new FileReader()

    reader.onload = (e: ProgressEvent<FileReader>) => {
      const fileData = e.target!.result as ArrayBuffer
      const byteArray = new Uint8Array(fileData)

      // @ts-ignore
      const res: Result = JSON.parse(decode(byteArray))
      if (res.err) {
        alert(res.err)
        return
      }
      geojson.value = res.feature
      sessions.value = res.activityFile.sessions
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
