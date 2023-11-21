<template>
  <div class="map-container position-relative w-100 h-100">
    <div v-if="features?.length == 0">No map data</div>
    <div v-else class="map" ref="map"></div>
    <div id="popup" class="ol-popup">
      <div class="popup-content">
        <div>
          <span :style="{ width: titleWidth + 'px' }">Time:</span>
          <span>
            {{
              popupRecord.timestamp
                ? toTimezoneDateString(popupRecord.timestamp, popupTimezoneOffsetHours)
                : '-'
            }}
          </span>
        </div>
        <div>
          <span :style="{ width: titleWidth + 'px' }">Distance:</span>
          <span>{{ popupRecord.distance ? (popupRecord.distance / 1000).toFixed(2) : '-' }}</span>
          <span>&nbsp;km</span>
        </div>
        <div>
          <span :style="{ width: titleWidth + 'px' }">Speed:</span>
          <span>{{
            popupRecord.speed ? ((popupRecord.speed * 3600) / 1000).toFixed(2) : '-'
          }}</span>
          <span>&nbsp;km/h</span>
        </div>
        <div v-if="hasPace">
          <span :style="{ width: titleWidth + 'px' }">Pace:</span>
          <span>{{ popupRecord.pace ? formatPace(popupRecord.pace) : '-' }}</span>
          <span>&nbsp;/km</span>
        </div>
        <div v-if="hasCadence">
          <span :style="{ width: titleWidth + 'px' }">Cadence:</span>
          <span>{{ popupRecord.cadence ? popupRecord.cadence : '-' }}</span>
          <span>&nbsp;rpm</span>
        </div>
        <div v-if="hasHeartRate">
          <span :style="{ width: titleWidth + 'px' }">Heart Rate:</span>
          <span>{{ popupRecord.heartRate ? popupRecord.heartRate : '-' }}</span>
          <span>&nbsp;bpm</span>
        </div>
        <div v-if="hasPower">
          <span :style="{ width: titleWidth + 'px' }">Power:</span>
          <span>{{ popupRecord.power ? popupRecord.power : '-' }}</span>
          <span>&nbsp;watts</span>
        </div>
        <div v-if="hasTemperature">
          <span :style="{ width: titleWidth + 'px' }">Temperature:</span>
          <span>{{ popupRecord.temperature ? popupRecord.temperature : '-' }}</span>
          <span>&nbsp;°C</span>
        </div>
        <div style="display: grid">
          <div style="grid-column: 1">
            <span :style="{ width: titleWidth + 'px' }">Altitude:</span>
            <span>{{ popupRecord.altitude ? popupRecord.altitude?.toFixed(2) : '-' }}</span>
            <span>&nbsp;masl</span>
          </div>
          <div style="grid-column: 2">
            <span>(Grade:&nbsp;</span>
            <span>{{ popupRecord.grade ? Math.round(popupRecord.grade) : '0' }}</span>
            <span>&nbsp;%)</span>
          </div>
        </div>

        <div>
          <span :style="{ width: titleWidth + 'px' }">Location:</span>
          <span>
            {{
              popupRecord.positionLong && popupRecord.positionLat
                ? toStringHDMS([popupRecord.positionLong, popupRecord.positionLat] as Coordinate)
                : '-'
            }}
          </span>
        </div>
      </div>
    </div>
    <div class="options position-absolute d-inline-flex">
      <div class="form-control-sm">
        <select
          class="custom-select custom-select-sm"
          aria-label="Hover Method"
          v-model="searchMethod"
          :disabled="kdIndexing"
        >
          <option value="standard" selected>Standard</option>
          <option value="kdtree">KD Tree</option>
        </select>
      </div>
      <div class="form-control-sm form-check">
        <input class="form-check-input" type="checkbox" v-model="debug" id="flexCheckDefault" />
        <label class="form-check-label" for="flexCheckDefault"> Debug </label>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import destinationPointIcon from '@/assets/map/destination-point.svg'
import startingPointIcon from '@/assets/map/starting-point.svg'
import 'ol/ol.css'

import { Record, Session } from '@/spec/activity'
import { toTimezoneDateString } from '@/toolkit/date'
import { formatPace } from '@/toolkit/pace'
import { Feature, MapBrowserEvent, Overlay } from 'ol'
import type { FeatureLike } from 'ol/Feature'
import OlMap from 'ol/Map'
import View from 'ol/View'
import { FullScreen, ScaleLine, ZoomToExtent, defaults as defaultControls } from 'ol/control.js'
import type { Coordinate } from 'ol/coordinate'
import { toStringHDMS } from 'ol/coordinate.js'
import { isEmpty } from 'ol/extent'
import { Geometry, LineString, Point, SimpleGeometry } from 'ol/geom'
import TileLayer from 'ol/layer/Tile'
import VectorImageLayer from 'ol/layer/VectorImage'
import OSM from 'ol/source/OSM'
import VectorSource from 'ol/source/Vector'
import { Icon, Stroke, Style } from 'ol/style'
import KDBush from 'kdbush'
import { around } from 'geokdbush-tk'
import { DateTime } from 'luxon'

const maximizeIcon = document.createElement('i')
const minimizeIcon = document.createElement('i')
maximizeIcon.setAttribute('class', 'fa-solid fa-expand')
minimizeIcon.setAttribute('class', 'fa-solid fa-compress')

export default {
  props: {
    features: {
      type: Array<Feature>,
      required: true
    },
    sessions: {
      type: Array<Session>,
      required: true
    },
    receivedRecord: Record,
    hasPace: Boolean,
    hasCadence: Boolean,
    hasHeartRate: Boolean,
    hasPower: Boolean,
    hasTemperature: Boolean
  },
  data() {
    return {
      popupFreeze: new Boolean(),
      popupRecord: new Record(),
      hoveredRecord: new Record(),
      popupActivityIndex: 0,
      popupTimezoneOffsetHours: 0,
      popupOverlay: new Overlay({}),
      startingPointStyle: new Style({
        image: new Icon({ crossOrigin: 'anonymous', src: startingPointIcon, scale: 1 })
      }),
      destinationPointStyle: new Style({
        image: new Icon({ crossOrigin: 'anonymous', src: destinationPointIcon, scale: 1 })
      }),
      vec: new VectorImageLayer({
        source: new VectorSource({
          features: []
        }),
        visible: true,
        style: new Style({
          stroke: new Stroke({
            // color: [232, 65, 24, 1.0],
            // color: [60, 99, 130, 1],
            color: '#34495e',
            width: 4
          })
        })
      }),
      map: new OlMap({
        controls: defaultControls(),
        layers: [
          new TileLayer({
            source: new OSM()
          })
        ],
        view: new View({
          center: [0, 0],
          zoom: 1,
          enableRotation: false,
          projection: 'EPSG:4326' // WGS84: World Geodetic System 1984
        })
      }),
      zoomToExtent: new ZoomToExtent(),
      pointer: new Feature(),
      debug: false,
      /** @type {'standard' | 'kdtree'} */
      searchMethod: 'standard',
      sessionTimestamp: '',
      kdIndexTimestamp: '',
      kdIndexToRecord: [] as String[],
      kdIndex: new KDBush(0),
      kdIndexing: false
    }
  },
  watch: {
    searchMethod: {
      handler(value: string) {
        if (value == 'kdtree' && !this.kdIndexed) {
          this.indexing_KDtree(this.sessions)
        }
      }
    },

    features: {
      handler(features: Array<Feature>) {
        this.$nextTick(() => {
          this.map.setTarget(this.$refs.map as HTMLElement)
          requestAnimationFrame(() => this.updateMapSource(features))
        })
      }
    },
    sessions: {
      handler() {
        this.sessionTimestamp = DateTime.now().toString()
        if (this.searchMethod == 'kdtree' && !this.kdIndexed) {
          this.indexing_KDtree(this.sessions)
        }
      }
    },
    zoomToExtent: {
      handler(newValue: ZoomToExtent, oldValue: ZoomToExtent) {
        this.map.removeControl(oldValue)
        this.map.addControl(newValue)
      }
    },
    receivedRecord: {
      handler(record: Record) {
        this.popupRecord = record
        if (JSON.stringify(record) == JSON.stringify(new Record())) {
          this.popupOverlay.setPosition(undefined)
          return
        }
        this.popupOverlay.setPosition([record.positionLong!, record.positionLat!])
      }
    },
    hoveredRecord: {
      handler(record: Record) {
        this.popupRecord = record
        this.$emit('hoveredRecord', record)
      }
    }
  },
  expose: ['showPopUpRecord'],
  methods: {
    toStringHDMS: toStringHDMS,
    toTimezoneDateString: toTimezoneDateString,
    formatPace: formatPace,

    updateMapSource(features: Feature[]) {
      this.popupOverlay.setPosition(undefined)

      const view = this.map.getView()
      const source = this.vec.getSource()!

      if (features.length == 0) {
        source.clear()
        return
      }

      source.clear()
      source.addFeatures(features)

      view.fit(source.getExtent(), { padding: [50, 50, 50, 50] })

      const startingPoints = new Array<Feature>()

      let lastSessionIndex = -1
      const destinationPointById = new Map<string, Feature>()
      for (let i = 0; i < features.length; i++) {
        const [, sessionIndex] = (features[i].getId() as string).split('-').map((v) => parseInt(v))

        const geometry = features[i]?.getGeometry() as SimpleGeometry
        if (lastSessionIndex != sessionIndex) {
          const startingPoint = new Feature(new Point(geometry.getFirstCoordinate()))
          startingPoint.setStyle(this.startingPointStyle as Style)
          startingPoint.setId(`startingPoint-${sessionIndex}`)
          startingPoints.push(startingPoint)
        }

        lastSessionIndex = sessionIndex

        const destinationPoint = new Feature(new Point(geometry.getLastCoordinate()))
        destinationPoint.setStyle(this.destinationPointStyle as Style)
        destinationPoint.setId(`destinationPoint-${sessionIndex}`)
        destinationPointById.set(`destinationPoint-${sessionIndex}`, destinationPoint)
      }

      source.addFeatures(startingPoints)
      source.addFeatures(Array.from(destinationPointById, ([, feature]) => feature))

      this.zoomToExtent = new ZoomToExtent({ extent: view.getViewStateAndExtent().extent })
    },
    findNearestRecord(featureId: string, coordinate: Coordinate): Record {
      const debug = this.debug
      if (debug) console.time('Standard')

      const [, sessionIndex, startIndex, endIndex] = featureId.split('-').map((v) => parseInt(v))

      this.popupActivityIndex = sessionIndex
      let nearestRecord: Record = new Record()
      let nearestEuclidean: number = Number.MAX_VALUE

      this.popupTimezoneOffsetHours = this.sessions[sessionIndex].timezone

      for (let i = startIndex; i <= endIndex; i++) {
        const rec = this.sessions[sessionIndex].records[i]
        if (!rec.positionLong || !rec.positionLat) continue
        const euclidean = Math.abs(
          Math.sqrt(
            Math.pow(rec.positionLong - coordinate[0], 2) +
              Math.pow(rec.positionLat - coordinate[1], 2)
          )
        )
        if (euclidean < nearestEuclidean) {
          nearestEuclidean = euclidean
          nearestRecord = rec
        }
      }
      // console.log(coordinate, nearestRecord)
      if (debug) console.timeEnd('Standard')
      return nearestRecord
    },
    // Start Indexing & Transform points to KDBush (1 time only)
    async indexing_KDtree(sessions: Array<Session>) {
      const debug = this.debug
      if (debug) console.time('KDtree Indexing')

      this.kdIndexing = true
      let totalRecords = 0
      sessions.forEach((ses) => {
        totalRecords += ses.records.length
      })
      this.kdIndexToRecord = [] as String[]
      this.kdIndex = new KDBush(totalRecords)
      sessions.forEach((d, sessionIndex) =>
        d.records.forEach((r, recordIndex) => {
          if (r.positionLat != null && r.positionLong != null) {
            this.kdIndexToRecord[
              this.kdIndex.add(r.positionLong, r.positionLat)
            ] = `${sessionIndex}-${recordIndex}`
          } else {
            this.kdIndexToRecord[this.kdIndex.add(0, 0)] = `${sessionIndex}-${recordIndex}`
          }
        })
      )
      this.kdIndex.finish()
      this.kdIndexTimestamp = this.sessionTimestamp
      this.kdIndexing = false
      if (debug) console.timeEnd('KDtree Indexing')
    },
    getClosestPoint_KDtree(coordinate: Coordinate): Record {
      const debug = this.debug
      if (debug) console.time('KDTree')

      const nearestIndex = around(this.kdIndex, coordinate[0], coordinate[1], 1)
      if (nearestIndex.length > 0) {
        const [sessionIndex, recordIndex] = this.kdIndexToRecord[
          nearestIndex[nearestIndex.length - 1]
        ]
          .split('-')
          .map((v) => parseInt(v))

        // console.log(coordinate, this.sessions[sessionIndex].records[recordIndex])
        if (debug) console.timeEnd('KDTree')
        return this.sessions[sessionIndex].records[recordIndex]
      }
      if (debug) console.timeEnd('KDTree')
      return new Record()
    },

    lineStringFeatureListener(e: MapBrowserEvent<any>) {
      if (e.type == 'singleclick')
        this.popupFreeze = !this.popupFreeze && this.popupOverlay.getPosition() != undefined
      if (!this.popupFreeze) {
        this.hoveredRecord = new Record()
        this.popupOverlay.setPosition(undefined)
      }
      if (this.popupFreeze == true && this.popupOverlay.getPosition() != undefined) return

      const features = this.map.getFeaturesAtPixel(e.pixel, { hitTolerance: 10 })
      let feature: FeatureLike | null = null
      for (let i = 0; i < features.length; i++) {
        if (features[i].getGeometry() instanceof LineString) {
          feature = features[i]
          break
        }
      }

      if (feature == null) return

      this.hoveredRecord = (() => {
        if (this.searchMethod == 'kdtree' && this.kdIndexed) {
          return this.getClosestPoint_KDtree(e.coordinate)
        }
        return this.findNearestRecord(feature.getId() as string, e.coordinate)
      })()
      this.popupOverlay.setPosition([
        this.hoveredRecord.positionLong,
        this.hoveredRecord.positionLat
      ] as Coordinate)
      this.popupFreeze = e.type == 'singleclick'
    },

    showPopUpRecord(record: Record) {
      this.popupOverlay.setPosition(undefined)
      if (!record) return

      this.popupRecord = record
      this.popupOverlay.setPosition([
        this.popupRecord.positionLong,
        this.popupRecord.positionLat
      ] as Coordinate)
      this.popupFreeze = true
    },

    updateExtent() {
      this.$nextTick(() => {
        const extent = this.vec.getSource()!.getExtent()
        if (isEmpty(extent)) return
        this.map.getView().fit(extent, { padding: [50, 50, 50, 50] })
        this.zoomToExtent = new ZoomToExtent({
          extent: this.map.getView().getViewStateAndExtent().extent
        })
      })
    }
  },
  computed: {
    titleWidth(): number {
      this.popupRecord // make reactive to this
      let maxWidth = 0
      const titleElements = document.querySelectorAll(
        '.popup-content > div:not([style*="display: none"]) > span:nth-child(1)'
      )
      titleElements.forEach((element: Element) => {
        const width = (element as HTMLElement).innerHTML.length * 6 // assume 1 char require 6 pixels
        maxWidth = Math.max(maxWidth, width)
      })
      return maxWidth
    },
    kdIndexed(): Boolean {
      return this.kdIndexTimestamp != '' && this.kdIndexTimestamp == this.sessionTimestamp
    }
  },
  mounted() {
    this.popupOverlay = new Overlay({ element: document.getElementById('popup')! })
    this.map.setTarget(this.$refs.map as HTMLElement)
    this.map.addOverlay(this.popupOverlay as Overlay)
    this.map.addLayer(this.vec as VectorImageLayer<VectorSource<Geometry>>)
    this.map.addControl(new FullScreen({ label: maximizeIcon, labelActive: minimizeIcon }))
    this.map.addControl(new ScaleLine())
    this.map.on('change:size', this.updateExtent)
    this.map.on('pointermove', this.lineStringFeatureListener)
    this.map.on('singleclick', this.lineStringFeatureListener)
    this.updateMapSource(this.features)
  },
  unmounted() {
    this.map.un('change:size', this.updateExtent)
    this.map.un('pointermove', this.lineStringFeatureListener)
    this.map.un('singleclick', this.lineStringFeatureListener)
  }
}
</script>
<style>
/* override open layers's default style */
.ol-overlay-container {
  position: relative !important;
}
</style>
<style scoped lang="scss">
.map-container {
  display: flex;
  align-items: center;
  justify-content: center;
}

.map {
  width: 100%;
  height: 100%;
  /* border: 0.5px solid black; */
}
.options {
  background-color: var(--ol-background-color);
  top: 8px !important;
  left: 35px !important;
  .form-check {
    // background-color: var(--ol-background-color);
    color: #333;
    font-weight: bold;
  }
}
.ol-popup {
  position: absolute;
  color: #000;
  background-color: rgba(255, 255, 255, 0.88);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
  padding: 15px;
  border-radius: 5px;
  border: 1px solid #cccccc;
  bottom: 12px;
  left: -50px;
}

.ol-popup:after,
.ol-popup:before {
  top: 100%;
  border: solid transparent;
  content: ' ';
  height: 0;
  width: 0;
  position: absolute;
  pointer-events: none;
}

.ol-popup:after {
  border-top-color: rgba(255, 255, 255, 0.88);
  border-width: 10px;
  left: 48px;
  margin-left: -10px;
}

.ol-popup:before {
  border-top-color: #cccccc;
  border-width: 11px;
  left: 48px;
  margin-left: -11px;
}

.popup-content {
  text-align: left;
  position: relative;
  font-size: 10px;
}

.popup-content span {
  display: inline-block;
}

.popup-content span:nth-child(2) {
  font-weight: 500;
}
</style>
