<script setup lang="ts">
import CadenceGraph from '@/components/CadenceGraph.vue'
import ElevationGraph from '@/components/ElevationGraph.vue'
import HeartRateGraph from '@/components/HeartRateGraph.vue'
import HeartRateZoneGraph from '@/components/HeartRateZoneGraph.vue'
import PaceGraph from '@/components/PaceGraph.vue'
import PowerGraph from '@/components/PowerGraph.vue'
import SpeedGraph from '@/components/SpeedGraph.vue'
import SplitPaceGraph from './SplitPaceGraph.vue'
import TemperatureGraph from '@/components/TemperatureGraph.vue'
import TheLap from '@/components/TheLap.vue'
import TheLoading from '@/components/TheLoading.vue'
import TheMap from '@/components/TheMap.vue'
import TheNavigatorInput from '@/components/TheNavigatorInput.vue'
import TheSession from '@/components/TheSession.vue'
import TheSummary, { MULTIPLE, NONE } from '@/components/TheSummary.vue'
import { Summary } from '@/spec/summary'
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
              <Transition>
                <div style="font-size: 0.8em" class="pt-1">
                  <span v-if="isActivityServiceReady"> Supported files: *.fit, *.gpx, *.tcx </span>
                  <span v-else>
                    Instantiating WebAssembly <i class="fas fa-spinner fa-spin"></i>
                  </span>
                </div>
              </Transition>
            </div>
            <div class="mb-3" v-if="isActivityFileReady">
              <TheSummary
                :sessions="sessions"
                :selected-sessions="selectedSessions"
                :is-activity-file-ready="isActivityFileReady"
                v-on:summary="onReceiveSummary"
                v-on:session-selected="onSessionSelected"
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
                  <div class="graph" v-if="hasHeartRate">
                    <HeartRateZoneGraph
                      :selected-session="selectedSessions"
                      :age="20"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                    ></HeartRateZoneGraph>
                  </div>
                  <div class="graph" v-if="hasPace">
                    <SplitPaceGraph
                      :selected-session="selectedSessions"
                      :age="20"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                    ></SplitPaceGraph>
                  </div>
                  <div class="graph" v-if="hasPace">
                    <PaceGraph
                      :records="selectedRecords"
                      :graph-records="selectedGraphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                      v-on:hoveredRecord="onHoveredRecord"
                      v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
                    ></PaceGraph>
                  </div>
                  <div class="graph" v-if="hasSpeed">
                    <SpeedGraph
                      :records="selectedRecords"
                      :graph-records="selectedGraphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                      v-on:hoveredRecord="onHoveredRecord"
                      v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
                    ></SpeedGraph>
                  </div>
                  <div class="graph" v-if="hasCadence">
                    <CadenceGraph
                      :records="selectedRecords"
                      :graph-records="selectedGraphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                      v-on:hoveredRecord="onHoveredRecord"
                      v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
                    ></CadenceGraph>
                  </div>
                  <div class="graph" v-if="hasHeartRate">
                    <HeartRateGraph
                      :records="selectedRecords"
                      :graph-records="selectedGraphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                      v-on:hoveredRecord="onHoveredRecord"
                      v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
                    ></HeartRateGraph>
                  </div>
                  <div class="graph" v-if="hasPower">
                    <PowerGraph
                      :records="selectedRecords"
                      :graph-records="selectedGraphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                      v-on:hoveredRecord="onHoveredRecord"
                      v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
                    ></PowerGraph>
                  </div>
                  <div class="graph" v-if="hasTemperature">
                    <TemperatureGraph
                      :records="selectedRecords"
                      :graph-records="selectedGraphRecords"
                      :summary="summary"
                      :received-record="hoveredRecord"
                      :received-record-freeze="hoveredRecordFreeze"
                      v-on:hoveredRecord="onHoveredRecord"
                      v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
                    ></TemperatureGraph>
                  </div>
                </div>
                <div
                  class="tab-pane fade"
                  id="tab2-body"
                  role="tabpanel"
                  aria-labelledby="tab2-tab"
                >
                  <div v-if="selectedSessions.length > 0">
                    <TheSession :sessions="selectedSessions" />
                  </div>
                </div>
                <div
                  class="tab-pane fade"
                  id="tab3-body"
                  role="tabpanel"
                  aria-labelledby="tab2-tab"
                >
                  <div v-if="selectedLaps.length > 0">
                    <TheLap :laps="selectedLaps" />
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
            <div class="col-12 map pe-0">
              <TheMap
                :sessions="sessions"
                :selected-sessions="selectedSessions"
                :features="selectedFeatures"
                :hasPace="hasPace"
                :hasCadence="hasCadence"
                :hasHeartRate="hasHeartRate"
                :hasPower="hasPower"
                :hasTemperature="hasTemperature"
                :received-record="hoveredRecord"
                :received-record-freeze="hoveredRecordFreeze"
                v-on:hoveredRecord="onHoveredRecord"
                v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
                ref="theMap"
              />
            </div>
            <div class="col-12 elevation-section">
              <ElevationGraph
                :name="'elev'"
                :has-altitude="hasAltitude"
                :summary="summary"
                :records="selectedRecords"
                :graph-records="selectedGraphRecords"
                :received-record="hoveredRecord"
                :received-record-freeze="hoveredRecordFreeze"
                v-on:hoveredRecord="onHoveredRecord"
                v-on:hoveredRecordFreeze="onHoveredRecordFreeze"
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
import { ActivityFile, Lap, Record, SPORT_GENERIC, Session } from '@/spec/activity'
import { DecodeResult, EncodeResult, ManufacturerListResult } from '@/spec/activity-service'
import { Feature } from 'ol'
import { GeoJSON } from 'ol/format'
import { shallowRef } from 'vue'

const isWebAssemblySupported =
  typeof WebAssembly === 'object' && typeof WebAssembly.instantiateStreaming === 'function'

if (isWebAssemblySupported == false) {
  alert('Sorry, it appears that your browser does not support WebAssembly :(')
}

// shallowRef
const sessions = shallowRef(new Array<Session>())
const activities = shallowRef(new Array<ActivityFile>())
const combinedRecords = shallowRef(new Array<Record>())
const combinedSessions = shallowRef(new Array<Session>())
const combinedLaps = shallowRef(new Array<Lap>())
const combinedFeatures = shallowRef(new Array<Feature>())
const combinedGraphRecords = shallowRef(new Array<Record>())
const featureBySession = shallowRef(new Map<number, Feature>())
const graphRecordsBySession = shallowRef(new Map<number, Array<Record>>())
const selectedRecords = shallowRef(new Array<Record>())
const selectedSessions = shallowRef(new Array<Session>())
const selectedLaps = shallowRef(new Array<Lap>())
const selectedFeatures = shallowRef(new Array<Feature>())
const selectedGraphRecords = shallowRef(new Array<Record>())
const manufacturerList = shallowRef(new ManufacturerListResult())

export default {
  data() {
    return {
      loading: false,
      activityService: new Worker(new URL('@/workers/activity-service.ts', import.meta.url), {
        type: 'module'
      }),
      sessionSelected: NONE,
      summary: new Summary(),
      hoveredRecord: new Record(),
      hoveredRecordFreeze: new Boolean(),
      isActivityServiceReady: false
    }
  },
  computed: {
    isActivityFileReady: function () {
      return activities.value.length > 0
    },
    hasPace(): boolean {
      for (let i = 0; i < selectedSessions.value.length; i++) {
        const ses = selectedSessions.value[i]
        switch (selectedSessions.value[i].sport) {
          case 'Hiking':
          case 'Walking':
          case 'Running':
          case 'Swimming':
          case 'Transition':
          case SPORT_GENERIC:
            for (let j = 0; j < ses.records.length; j++) {
              const rec = ses.records[j]
              if (rec.pace != null) return true
            }
        }
      }
      return false
    },
    hasAltitude(): boolean {
      for (let i = 0; i < selectedRecords.value.length; i++) {
        if (selectedRecords.value[i].altitude != null) return true
      }
      return false
    },
    hasSpeed(): boolean {
      for (let i = 0; i < selectedRecords.value.length; i++) {
        if (selectedRecords.value[i].speed) return true
      }
      return false
    },
    hasCadence(): boolean {
      for (let i = 0; i < selectedRecords.value.length; i++) {
        if (selectedRecords.value[i].cadence) return true
      }
      return false
    },
    hasHeartRate(): boolean {
      for (let i = 0; i < selectedRecords.value.length; i++) {
        if (selectedRecords.value[i].heartRate) return true
      }
      return false
    },
    hasPower(): boolean {
      for (let i = 0; i < selectedRecords.value.length; i++) {
        if (selectedRecords.value[i].power) return true
      }
      return false
    },
    hasTemperature(): boolean {
      for (let i = 0; i < selectedRecords.value.length; i++) {
        if (selectedRecords.value[i].temperature != null) return true
      }
      return false
    }
  },
  methods: {
    getExtention(s: string): string {
      if (s == '') return ''
      const parts = s.split('.')
      return parts[parts.length - 1].toLocaleLowerCase()
    },
    getExtentionIdentifier(ext: string): number {
      if (ext == 'fit') return 1
      if (ext == 'gpx') return 2
      if (ext == 'tcx') return 3
      return 0
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
              reject('no file selected')
              return
            }

            const ext = this.getExtention(selectedFile.name)
            const extentionIdentifier = this.getExtentionIdentifier(ext)
            if (extentionIdentifier == 0) {
              reject(`file '${selectedFile.name}' (type: ${ext}) is not supported`)
              return
            }

            const reader = new FileReader()
            reader.onload = (e: ProgressEvent<FileReader>) => {
              const fileData = e.target!.result as ArrayBuffer
              resolve(new Uint8Array([extentionIdentifier, ...new Uint8Array(fileData)]))
            }
            reader.onerror = reject
            reader.readAsArrayBuffer(selectedFile as File)
          })
        )
      }

      Promise.all(promisers)
        .then((arr) => {
          this.loading = true
          this.activityService.postMessage({ type: 'decode', input: arr })
        })
        .catch((e: string) => {
          console.log(e)
          alert(e)
        })
    },
    activityServiceOnMessage(e: MessageEvent) {
      const [type, result, elapsed] = [e.data.type, e.data.result, e.data.elapsed]
      switch (type) {
        case 'isReady':
          this.isActivityServiceReady = true
          break
        case 'decode':
          this.decodeHandler(result, elapsed)
          break
        case 'encode':
          this.encodeHandler(result, elapsed)
          break
        case 'manufacturerList':
          this.manufacturerListHandler(result)
          break
      }
    },
    decodeHandler(result: DecodeResult, elapsed: number) {
      result = new DecodeResult(result)
      if (result.err != null) {
        console.error(`Decode: ${result.err}`)
        alert(`Decode: ${result.err}`)
        this.loading = false
        return
      }

      // Instrumenting...
      console.group('Decoding:')
      console.group('Spent on WASM:')
      console.debug('Decode took:\t\t', result.decodeTook, 'ms')
      console.debug('Serialization took:\t', result.serializationTook, 'ms')
      console.debug('Total elapsed:\t\t', result.totalElapsed, 'ms')
      console.groupEnd()
      console.debug('Interop WASM to JS:\t', elapsed - result.totalElapsed, 'ms')
      console.debug('Total elapsed:\t\t', elapsed, 'ms')
      console.groupEnd()

      requestAnimationFrame(() => {
        this.preprocessing(result)
        this.scrollTop()
      })
    },
    encodeHandler(result: EncodeResult, elapsed: number) {
      result = new EncodeResult(result)
      if (result.err != null) {
        console.error(`Encode: ${result.err}`)
        alert(`Encode: ${result.err}`)
        this.loading = false
        return
      }
    },
    manufacturerListHandler(result: ManufacturerListResult) {
      manufacturerList.value = new ManufacturerListResult(result)
    },
    preprocessing(result: DecodeResult) {
      console.time('Preprocessing')

      activities.value = result.activities

      combinedSessions.value = [] // clone of sessions (record's distance will be accumulated)
      sessions.value = activities.value.flatMap((act) => {
        for (let i = 0; i < act.sessions.length; i++) {
          const ses = act.sessions[i]
          ses.creatorName = act.creator.name
          ses.timeCreated = act.creator.timeCreated
          ses.timezone = act.timezone

          const clonedSes = { ...ses }
          clonedSes.records = []
          for (let j = 0; j < ses.records.length; j++) {
            clonedSes.records.push({ ...ses.records[j] })
          }
          combinedSessions.value.push(clonedSes)
        }
        return act.sessions
      })

      let prevSessionlastDistance = 0
      let lastDistance = 0

      combinedLaps.value = []
      combinedRecords.value = []
      combinedFeatures.value = []
      combinedGraphRecords.value = []
      featureBySession.value.clear()
      graphRecordsBySession.value.clear()

      for (let i = 0; i < combinedSessions.value.length; i++) {
        const ses = combinedSessions.value[i]
        for (let j = 0; j < ses.records.length; j++) {
          if (ses.records[j].distance == null) continue
          ses.records[j].distance! += prevSessionlastDistance
          lastDistance = ses.records[j].distance!
        }
        prevSessionlastDistance = lastDistance

        combinedLaps.value.push(...ses.laps)
        combinedRecords.value.push(...ses.records)

        const feature = this.createFeature(ses, i)
        if (feature != null) {
          combinedFeatures.value.push(feature)
          featureBySession.value.set(i, feature)
        }
      }

      combinedGraphRecords.value = this.summarizeRecords(combinedRecords.value)

      const n = sessions.value.length
      this.selectSession(n == 0 ? NONE : n == 1 ? 0 : MULTIPLE)

      setTimeout(() => (this.loading = false), 200)
      console.timeEnd('Preprocessing')
    },
    createFeature(session: Session, sessionIndex: number): Feature | null {
      let coordinates: number[][] = []
      for (let i = 0; i < session.records.length; i++) {
        const rec = session.records[i]

        if (rec.positionLong == null || rec.positionLat == null) continue
        coordinates.push([rec.positionLong!, rec.positionLat!])
      }

      if (coordinates.length == 0) return null

      return new GeoJSON().readFeature({
        id: `lineString-${sessionIndex}`,
        type: 'Feature',
        geometry: {
          type: 'LineString',
          coordinates: coordinates
        }
      })
    },
    startTime(records: Record[]): string | null {
      let timestamp: string | null = null
      for (let i = 0; i < records.length; i++) {
        if (timestamp != null) break
        timestamp = records[i].timestamp
      }
      return timestamp
    },
    endTime(records: Record[]): string | null {
      let timestamp: string | null = null
      for (let i = records.length - 1; i >= 0; i--) {
        if (timestamp != null) break
        timestamp = records[i].timestamp
      }
      return timestamp
    },
    hasDistance(records: Record[]): boolean {
      for (let i = 0; i < records.length; i++) {
        if (records[i].distance != null && records[i].distance! > 0) return true
      }
      return false
    },
    lastDistance(records: Record[]): number {
      let d: number = 0
      for (let i = records.length - 1; i >= 0; i--) {
        const rec = records[i]
        if (d > 0) break
        if (rec.distance != null) d = records[i].distance!
      }
      return d
    },
    summarizeRecords(records: Record[]): Record[] {
      // There is a limit on the amount of data that can be visually displayed in a graph, so will summarize it based on distance.
      // As a cyclist, the longer the distance, the less likely we are to be concerned about smaller distances.
      // If the distance between data exceeds this scale, then no data will be scaled.
      const percentage = 0.05 / 100 // 0.05% e.g distance is 1km, we summarize every 0.5m.

      const summarizeByDistance = this.hasDistance(records)
      let threshold: number = 0 // the unit value is either distance meters or duration seconds depends on conditions below:
      if (summarizeByDistance) {
        const ld = this.lastDistance(records)
        threshold = ld * percentage
      } else {
        const startTime = this.startTime(records)
        const endTime = this.endTime(records)
        const dur = new Date(endTime!).getTime() - new Date(startTime!).getTime()
        threshold = dur * percentage
      }

      let timestamps: string[] = []
      let distances: number[] = []
      let positionLats: number[] = []
      let positionLongs: number[] = []
      let altitudes: number[] = []
      let cadences: number[] = []
      let heartRates: number[] = []
      let speeds: number[] = []
      let powers: number[] = []
      let temps: number[] = []
      let grades: number[] = []
      let paces: number[] = []

      const summarized: Record[] = []
      let curIndex = 0
      for (let i = 0; i < records.length; i++) {
        const rec = records[i]
        const cur = records[curIndex]

        let delta: number = 0
        if (summarizeByDistance) {
          delta = (rec.distance ?? 0) - (cur.distance ?? 0)
        } else {
          delta = new Date(rec.timestamp!).getTime() - new Date(cur.timestamp!).getTime()
          if (cur.timestamp == null) curIndex = i // move cursor
        }

        if (delta > threshold) {
          let record = new Record()

          if (timestamps.length > 0) record.timestamp = timestamps[0]
          if (distances.length > 0) record.distance = distances[0]
          if (positionLats.length > 0) record.positionLat = positionLats[0]
          if (positionLongs.length > 0) record.positionLong = positionLongs[0]

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
          if (paces.length > 0)
            record.pace = paces.reduce((sum, val) => sum + val, 0) / paces.length

          summarized.push(record)

          // Reset array
          timestamps = []
          distances = []
          positionLats = []
          positionLongs = []
          altitudes = []
          cadences = []
          heartRates = []
          speeds = []
          powers = []
          temps = []
          grades = []
          paces = []

          curIndex = i
        }

        if (records[i].timestamp != null) timestamps.push(records[i].timestamp!)
        if (records[i].distance != null) distances.push(records[i].distance!)
        if (records[i].positionLat != null) positionLats.push(records[i].positionLat!)
        if (records[i].positionLong != null) positionLongs.push(records[i].positionLong!)
        if (records[i].altitude != null) altitudes.push(records[i].altitude!)
        if (records[i].cadence != null) cadences.push(records[i].cadence!)
        if (records[i].heartRate != null) heartRates.push(records[i].heartRate!)
        if (records[i].speed != null) speeds.push(records[i].speed!)
        if (records[i].power != null) powers.push(records[i].power!)
        if (records[i].temperature != null) temps.push(records[i].temperature!)
        if (records[i].grade != null) grades.push(records[i].grade!)
        if (records[i].pace != null) paces.push(records[i].pace!)
      }

      return summarized
    },
    selectSession(index: number) {
      switch (index) {
        case NONE:
          selectedRecords.value = []
          selectedSessions.value = []
          selectedLaps.value = []
          selectedFeatures.value = []
          selectedGraphRecords.value = []
          break
        case MULTIPLE:
          selectedRecords.value = combinedRecords.value
          selectedSessions.value = combinedSessions.value
          selectedLaps.value = combinedLaps.value
          selectedFeatures.value = combinedFeatures.value
          selectedGraphRecords.value = combinedGraphRecords.value
          break
        default: {
          if (index > sessions.value.length - 1) return
          selectedRecords.value = sessions.value[index].records
          selectedSessions.value = [sessions.value[index]]
          selectedLaps.value = sessions.value[index].laps
          const feature = featureBySession.value.get(index)
          selectedFeatures.value = feature ? [feature] : []

          if (graphRecordsBySession.value.get(index) == undefined) {
            const ses = sessions.value[index]
            const graphRecords = this.summarizeRecords(ses.records)
            graphRecordsBySession.value.set(index, graphRecords)
          }

          selectedGraphRecords.value = graphRecordsBySession.value.get(index) as []
          break
        }
      }
    },
    scrollTop() {
      document.body.scrollTop = 0 // For Safari
      document.documentElement.scrollTop = 0 // For Chrome, Firefox, IE and Opera
    },
    onReceiveSummary(summary: Summary) {
      this.summary = summary
    },
    onHoveredRecord(record: Record) {
      this.hoveredRecord = new Record(record)
    },
    onHoveredRecordFreeze(freeze: Boolean) {
      this.hoveredRecordFreeze = freeze
    },
    onSessionSelected(sessionSelected: number) {
      this.sessionSelected = sessionSelected
      this.selectSession(sessionSelected)
    }
  },
  mounted() {
    document.getElementById('fileInput')?.addEventListener('change', this.fileInputEventListener)
    this.activityService.onmessage = this.activityServiceOnMessage
    this.activityService.postMessage({ type: 'isReady' })
    this.activityService.postMessage({ type: 'manufacturerList' })
    this.selectSession(this.sessionSelected)
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
  padding-left: 20px;
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
