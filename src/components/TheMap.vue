<template>
  <div class="map-container position-relative w-100 h-100">
    <div v-if="features?.length == 0" class="no-map-data">No map data</div>
    <div v-else class="map" ref="map"></div>
    <div id="popup" class="ol-popup">
      <div class="popup-content" v-if="isIndexed">
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
        <div v-show="hasPace">
          <span :style="{ width: titleWidth + 'px' }">Pace:</span>
          <span>{{ popupRecord.pace ? formatPace(popupRecord.pace) : '-' }}</span>
          <span>&nbsp;/km</span>
        </div>
        <div v-show="hasCadence">
          <span :style="{ width: titleWidth + 'px' }">Cadence:</span>
          <span>{{ popupRecord.cadence ? popupRecord.cadence : '-' }}</span>
          <span>&nbsp;rpm</span>
        </div>
        <div v-show="hasHeartRate">
          <span :style="{ width: titleWidth + 'px' }">Heart Rate:</span>
          <span>{{ popupRecord.heartRate ? popupRecord.heartRate : '-' }}</span>
          <span>&nbsp;bpm</span>
        </div>
        <div v-show="hasPower">
          <span :style="{ width: titleWidth + 'px' }">Power:</span>
          <span>{{ popupRecord.power ? popupRecord.power : '-' }}</span>
          <span>&nbsp;watts</span>
        </div>
        <div v-show="hasTemperature">
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
      <div class="popup-content" v-else>
        <div class="d-flex justify-content-center">
          <div class="spinner-border spinner-border-sm" role="status">
            <span class="sr-only">Indexing...</span>
          </div>
          <div class="px-1">Indexing...</div>
        </div>
      </div>
    </div>
    <div class="options position-absolute d-inline-flex d-none">
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
import concealPointIcon from '@/assets/map/eye-off-green.svg'
import trimPointIcon from '@/assets/map/crop-red.svg'
import 'ol/ol.css'

import { Record, Session } from '@/spec/activity'
import { toTimezoneDateString } from '@/toolkit/date'
import { formatPace } from '@/toolkit/pace'
import { around } from 'geokdbush-tk'
import KDBush from 'kdbush'
import { Feature, MapBrowserEvent, Overlay } from 'ol'
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
import { shallowRef } from 'vue'
import { MULTIPLE, NONE } from '@/components/TheSummary.vue'
import { GeoJSON } from 'ol/format'
import { Marker } from '@/spec/activity-service'

// shallowRef
const kdbush = shallowRef(new KDBush(0))
const kdbushIndexToRecordMapping = shallowRef([] as String[])

// do not require reactivity
const emptyRecord = new Record()

let map: OlMap

const hiddenStyle = () => new Style({ stroke: new Stroke({ color: 'rgba(0, 0, 0, 0)', width: 0 }) })
const concealStyle = () => [
  new Style({ stroke: new Stroke({ color: '#FFFFFF', width: 6 }), zIndex: -1 }), // outliner
  new Style({ stroke: new Stroke({ color: '#2ecc71', width: 4 }), zIndex: 0 })
]
const trimStyle = () => [
  new Style({ stroke: new Stroke({ color: '#FFFFFF', width: 6 }), zIndex: -1 }), // outliner
  new Style({ stroke: new Stroke({ color: '#f50f30', width: 4 }), zIndex: 0 })
]

const tileLayer = new TileLayer({ source: new OSM() })
const routeVecLayer = new VectorImageLayer({
  source: new VectorSource({ features: [] as Feature[] }),
  visible: true,
  style: [
    new Style({ stroke: new Stroke({ color: '#FFFFFF', width: 6 }), zIndex: -1 }), // outliner
    new Style({ stroke: new Stroke({ color: '#34495e', width: 4 }) })
  ]
})

const concealVecLayer = new VectorImageLayer({
  source: new VectorSource({ features: [] as Feature[] }),
  visible: true,
  style: hiddenStyle()
})
const trimVecLayer = new VectorImageLayer({
  source: new VectorSource({ features: [] as Feature[] }),
  visible: true,
  style: hiddenStyle()
})
const pointLayer = new VectorImageLayer({
  source: new VectorSource({ features: [] as Feature[] }),
  visible: true
})

const view = new View({
  center: [0, 0],
  zoom: 1,
  enableRotation: false,
  projection: 'EPSG:4326' // WGS84: World Geodetic System 1984
})

const newIcon = (src: string, scale: number = 1): Icon => {
  return new Icon({ crossOrigin: 'anonymous', src: src, scale: scale })
}

const startingPointStyle = new Style({ image: newIcon(startingPointIcon) })
const destinationPointStyle = new Style({ image: newIcon(destinationPointIcon) })
const concealPointStyle = new Style({ image: newIcon(concealPointIcon, 0.075) })
const trimPointStyle = new Style({ image: newIcon(trimPointIcon, 0.075) })

let popupOverlay = new Overlay({})
let zoomToExtent = new ZoomToExtent()

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
    selectedSessions: {
      type: Array<Session>,
      required: true
    },
    selectSession: {
      type: Number,
      required: true
    },

    // Display Tools Marker
    toolConcealActive: {
      type: Boolean,
      default: false
    },
    toolTrimActive: {
      type: Boolean,
      default: false
    },
    toolConcealMarkers: {
      type: Array<Marker>,
      required: true
    },
    toolTrimMarkers: {
      type: Array<Marker>,
      required: true
    },

    receivedRecord: Record,
    receivedRecordFreeze: Boolean,
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
      popupTimezoneOffsetHours: 0,
      debug: false
    }
  },
  watch: {
    features: {
      handler(features: Array<Feature>) {
        this.$nextTick(() => {
          map.setTarget(this.$refs.map as HTMLElement)
          requestAnimationFrame(() => this.updateMapSource(features))
        })
      }
    },
    sessions: {
      handler() {
        this.kdbushIndexing(this.sessions)
        this.createConcealFeatures(this.sessions)
        this.createTrimFeatures(this.sessions)
      }
    },
    selectSession: {
      handler() {
        this.showConcealFeaturesHandler()
        this.showTrimFeaturesHandler()
      }
    },
    toolConcealActive: {
      handler() {
        this.showConcealFeaturesHandler()
      }
    },
    toolTrimActive: {
      handler() {
        this.showTrimFeaturesHandler()
      }
    },
    receivedRecord: {
      handler(record: Record) {
        this.popupRecord = record
        if (JSON.stringify(record) == JSON.stringify(new Record())) {
          popupOverlay.setPosition(undefined)
          return
        }
        popupOverlay.setPosition([record.positionLong!, record.positionLat!])
      }
    },
    receivedRecordFreeze: {
      handler(freeze: Boolean) {
        this.popupFreeze = freeze
      }
    },
    hoveredRecord: {
      handler(record: Record) {
        this.popupRecord = record
        this.$emit('hoveredRecord', record)
      }
    },
    popupFreeze: {
      handler(freeze: Boolean) {
        this.$emit('hoveredRecordFreeze', freeze)
      }
    },

    // Tool things
    toolConcealMarkers: {
      handler(markers: Array<Marker>) {
        markers.forEach((m, i) => this.setConceal(i, m.startN, m.endN))
      },
      deep: true
    },
    toolTrimMarkers: {
      handler(markers: Array<Marker>) {
        markers.forEach((m, i) => this.setTrim(i, m.startN, m.endN))
      },
      deep: true
    }
  },
  expose: ['showPopUpRecord'],
  methods: {
    toStringHDMS: toStringHDMS,
    toTimezoneDateString: toTimezoneDateString,
    formatPace: formatPace,

    // prepare empty hidden feature for all session (whenever session is created)
    createConcealFeatures(sessions: Array<Session>) {
      if (sessions == null) sessions = this.sessions

      const source = concealVecLayer.getSource()!
      source.clear()

      const pointSource = pointLayer.getSource()!
      pointSource.getFeatures().forEach((v) => {
        if (v.getId()?.toString().startsWith('conceal-')) pointSource.removeFeature(v)
      })

      sessions.forEach((_, sessionIndex) => {
        ;['concealStart', 'concealEnd'].forEach((v) => {
          const feat = new GeoJSON().readFeature({
            id: `conceal-lineString-${v}-${sessionIndex}`,
            type: 'Feature',
            style: undefined
          }) as Feature
          source.addFeature(feat)
        })

        const startingPoint = new Feature(new Point([0, 0]))
        startingPoint.setStyle(undefined)
        startingPoint.setId(`conceal-startPartLastPoint-${sessionIndex}`)

        const destinationPoint = new Feature(new Point([0, 0]))
        destinationPoint.setStyle(undefined)
        destinationPoint.setId(`conceal-endPartFirstPoint-${sessionIndex}`)

        pointSource.addFeatures([startingPoint, destinationPoint])
      })
    },

    // show Conceal features based on selected session
    showConcealFeaturesHandler() {
      if (!this.toolConcealActive) this.showConcealFeatures(NONE)
      else this.showConcealFeatures(this.selectSession)
    },

    showConcealFeatures(index: number) {
      const source = concealVecLayer.getSource()!
      const pointSource = pointLayer.getSource()!

      if (index == NONE) {
        // NONE
        source.getFeatures().forEach((f) => f.setStyle(undefined))
        pointSource.getFeatures().forEach((f) => {
          if (f.getId()?.toString().startsWith('conceal-')) f.setStyle(undefined)
        })
      } else if (index == MULTIPLE) {
        // MULTIPLE
        source.getFeatures().forEach((f) => f.setStyle(concealStyle))
        pointSource.getFeatures().forEach((f) => {
          if (f.getId()?.toString().startsWith('conceal-')) f.setStyle(concealPointStyle as Style)
        })
      } else {
        // DEFAULT
        source.getFeatures().forEach((f) => f.setStyle(undefined))
        pointSource.getFeatures().forEach((f) => {
          if (f.getId()?.toString().startsWith('conceal-')) f.setStyle(undefined)
        })
        ;['concealStart', 'concealEnd'].forEach((v) => {
          ;(
            source.getFeatureById(`conceal-lineString-${v}-${index}`) as Feature<LineString> | null
          )?.setStyle(concealStyle)
        })
        pointSource.getFeatures().forEach((f) => {
          if (
            f.getId()?.toString().startsWith('conceal-') &&
            f.getId()?.toString().endsWith(`-${index}`)
          )
            f.setStyle(concealPointStyle as Style)
        })
      }
    },
    // update feat which side to hide, 0-showStartIndex and showEndIndex-0
    // add/remove coordinates based on trim index record
    setConceal(sessionIndex: number, showStartIndex: number, showEndIndex: number) {
      if (sessionIndex >= this.sessions.length || sessionIndex < 0) return

      const startPartCoords = this.sessions[sessionIndex].records
        .slice(0, showStartIndex + 1)
        .reduce<Array<Array<number>>>((result, r) => {
          if (r.positionLong != null && r.positionLat != null)
            result.push([r.positionLong, r.positionLat])
          return result
        }, [])
      const startPartLineString = new LineString(startPartCoords)
      const startPartLastCoord = startPartLineString.getLastCoordinate()

      const source = concealVecLayer.getSource()!
      const pointSource = pointLayer.getSource()!

      ;(
        source.getFeatureById(
          `conceal-lineString-concealStart-${sessionIndex}`
        ) as Feature<Geometry> | null
      )?.setGeometry(startPartLineString)
      ;(
        (
          pointSource.getFeatureById(
            `conceal-startPartLastPoint-${sessionIndex}`
          ) as Feature<Geometry> | null
        )?.getGeometry() as Point
      ).setCoordinates(startPartLastCoord)

      const endPartCoords = this.sessions[sessionIndex].records
        .slice(showEndIndex, this.sessions[sessionIndex].records.length)
        .reduce<Array<Array<number>>>((result, r) => {
          if (r.positionLong != null && r.positionLat != null)
            result.push([r.positionLong, r.positionLat])
          return result
        }, [])
      const endPartLineString = new LineString(endPartCoords)
      const endPartFirstCoord = endPartLineString.getFirstCoordinate()

      ;(
        source.getFeatureById(
          `conceal-lineString-concealEnd-${sessionIndex}`
        ) as Feature<Geometry> | null
      )?.setGeometry(endPartLineString)
      ;(
        (
          pointSource.getFeatureById(
            `conceal-endPartFirstPoint-${sessionIndex}`
          ) as Feature<Geometry> | null
        )?.getGeometry() as Point
      ).setCoordinates(endPartFirstCoord)
    },

    // prepare empty hidden feature for all session (whenever session is created)
    createTrimFeatures(sessions: Array<Session>) {
      if (sessions == null) sessions = this.sessions

      const source = trimVecLayer.getSource()!
      source.clear()

      const pointSource = pointLayer.getSource()!
      pointSource.getFeatures().forEach((v) => {
        if (v.getId()?.toString().startsWith('trim-')) pointSource.removeFeature(v)
      })

      sessions.forEach((_, sessionIndex) => {
        ;['trimStart', 'trimEnd'].forEach((v) => {
          const feat = new GeoJSON().readFeature({
            id: `trim-lineString-${v}-${sessionIndex}`,
            type: 'Feature',
            style: undefined
          }) as Feature
          source.addFeature(feat)
        })

        const startPartLastPoint = new Feature(new Point([0, 0]))
        startPartLastPoint.setStyle(undefined)
        startPartLastPoint.setId(`trim-startPartLastPoint-${sessionIndex}`)

        const endPartFirstPoint = new Feature(new Point([0, 0]))
        endPartFirstPoint.setStyle(undefined)
        endPartFirstPoint.setId(`trim-endPartFirstPoint-${sessionIndex}`)

        pointSource.addFeatures([startPartLastPoint, endPartFirstPoint])
      })
    },

    // show Conceal features based on selected session
    showTrimFeaturesHandler() {
      if (!this.toolTrimActive) this.showTrimFeatures(NONE)
      else this.showTrimFeatures(this.selectSession)
    },

    showTrimFeatures(index: number) {
      const source = trimVecLayer.getSource()!
      const pointSource = pointLayer.getSource()!

      if (index == NONE) {
        // NONE
        source.getFeatures().forEach((f) => f.setStyle(undefined))
        pointSource.getFeatures().forEach((f) => {
          if (f.getId()?.toString().startsWith('trim-')) f.setStyle(undefined)
        })
      } else if (index == MULTIPLE) {
        // MULTIPLE
        source.getFeatures().forEach((f) => f.setStyle(trimStyle))
        pointSource.getFeatures().forEach((f) => {
          if (f.getId()?.toString().startsWith('trim-')) f.setStyle(trimPointStyle as Style)
        })
      } else {
        // DEFAULT
        source.getFeatures().forEach((f) => f.setStyle(undefined))
        pointSource.getFeatures().forEach((f) => {
          if (f.getId()?.toString().startsWith('trim-')) f.setStyle(undefined)
        })
        ;['trimStart', 'trimEnd'].forEach((v) => {
          ;(
            source.getFeatureById(`trim-lineString-${v}-${index}`) as Feature<LineString> | null
          )?.setStyle(trimStyle)
        })
        pointSource.getFeatures().forEach((f) => {
          if (
            f.getId()?.toString().startsWith('trim-') &&
            f.getId()?.toString().endsWith(`-${index}`)
          )
            f.setStyle(trimPointStyle as Style)
        })
      }
    },

    // update feat which side to hide, 0-showStartIndex and showEndIndex-0
    // add/remove coordinates based on trim index record
    setTrim(sessionIndex: number, showStartIndex: number, showEndIndex: number) {
      if (sessionIndex >= this.sessions.length || sessionIndex < 0) return

      const startPartCoords = this.sessions[sessionIndex].records
        .slice(0, showStartIndex + 1)
        .reduce<Array<Array<number>>>((result, r) => {
          if (r.positionLong != null && r.positionLat != null)
            result.push([r.positionLong, r.positionLat])
          return result
        }, [])
      const startPartLineString = new LineString(startPartCoords)
      const startPartLastCoords = startPartLineString.getLastCoordinate()

      const source = trimVecLayer.getSource()!
      const pointSource = pointLayer.getSource()!

      ;(
        source.getFeatureById(
          `trim-lineString-trimStart-${sessionIndex}`
        ) as Feature<Geometry> | null
      )?.setGeometry(startPartLineString)
      ;(
        (
          pointSource.getFeatureById(
            `trim-startPartLastPoint-${sessionIndex}`
          ) as Feature<Geometry> | null
        )?.getGeometry() as Point
      ).setCoordinates(startPartLastCoords)

      const endPartCoords = this.sessions[sessionIndex].records
        .slice(showEndIndex, this.sessions[sessionIndex].records.length)
        .reduce<Array<Array<number>>>((result, r) => {
          if (r.positionLong != null && r.positionLat != null)
            result.push([r.positionLong, r.positionLat])
          return result
        }, [])
      const endPartLineString = new LineString(endPartCoords)
      const endPartFirstCoord = endPartLineString.getFirstCoordinate()

      ;(
        source.getFeatureById(`trim-lineString-trimEnd-${sessionIndex}`) as Feature<Geometry> | null
      )?.setGeometry(endPartLineString)
      ;(
        (
          pointSource.getFeatureById(
            `trim-endPartFirstPoint-${sessionIndex}`
          ) as Feature<Geometry> | null
        )?.getGeometry() as Point
      ).setCoordinates(endPartFirstCoord)
    },

    updateMapSource(features: Feature[]) {
      popupOverlay.setPosition(undefined)
      const source = routeVecLayer.getSource()!
      source.clear()

      const pointSource = pointLayer.getSource()!
      pointSource.getFeatures().forEach((v) => {
        if (
          v.getId()?.toString().startsWith('startingPoint') ||
          v.getId()?.toString().startsWith('destinationPoint')
        )
          pointSource.removeFeature(v)
      })

      if (features.length == 0) {
        return
      }

      source.addFeatures(features)

      // source.removeFeature(this.overlayFeatures)

      // this.overlayFeatures.length = 0
      // features.forEach((feature) => {
      //   const f = feature.clone()

      //   f.setStyle(lineStringStyle)
      //   this.overlayFeatures.push(f)
      // })
      // source.addFeatures(this.overlayFeatures)

      const pointFeatures = new Array<Feature>()
      for (let i = 0; i < features.length; i++) {
        const geometry = features[i]?.getGeometry() as SimpleGeometry
        const [, sessionIndex] = (features[i].getId() as string).split('-').map((v) => parseInt(v))

        const startingPoint = new Feature(new Point(geometry.getFirstCoordinate()))
        startingPoint.setStyle(startingPointStyle as Style)
        startingPoint.setId(`startingPoint-${sessionIndex}`)
        pointFeatures.push(startingPoint)

        const destinationPoint = new Feature(new Point(geometry.getLastCoordinate()))
        destinationPoint.setStyle(destinationPointStyle as Style)
        destinationPoint.setId(`destinationPoint-${sessionIndex}`)
        pointFeatures.push(destinationPoint)
      }

      pointSource.addFeatures(pointFeatures)
      this.updateExtent()
    },

    async kdbushIndexing(sessions: Array<Session>) {
      console.time('KDtree Indexing')

      let totalRecords = 0
      for (let i = 0; i < sessions.length; i++) {
        totalRecords += sessions[i].records.length
      }

      kdbush.value = new KDBush(totalRecords)
      kdbushIndexToRecordMapping.value = []
      for (let i = 0; i < sessions.length; i++) {
        const records = sessions[i].records
        for (let j = 0; j < records.length; j++) {
          const rec = records[j]
          if (rec.positionLat == null || rec.positionLong == null) {
            kdbush.value.add(0, 0) // only to increment index to match kdbush's size.
            continue
          }
          const index = kdbush.value.add(rec.positionLong, rec.positionLat)
          kdbushIndexToRecordMapping.value[index] = `${i}-${j}`
        }
      }
      kdbush.value.finish()
      console.timeEnd('KDtree Indexing')
    },

    findNearestRecord(coordinate: Coordinate): Record {
      const nearestIndex = around(kdbush.value, coordinate[0], coordinate[1], 1)
      if (nearestIndex.length == 0) return emptyRecord

      const [sessionIndex, recordIndex] = kdbushIndexToRecordMapping.value[nearestIndex[0]]
        .split('-')
        .map((v) => parseInt(v))

      if (this.selectedSessions.length == 1) return this.selectedSessions[0].records[recordIndex]
      return this.selectedSessions[sessionIndex].records[recordIndex]
    },

    lineStringFeatureListener(e: MapBrowserEvent<any>) {
      if (e.type == 'singleclick')
        this.popupFreeze = !this.popupFreeze && popupOverlay.getPosition() != undefined
      if (this.popupFreeze == true && popupOverlay.getPosition() != undefined) return

      const features = map.getFeaturesAtPixel(e.pixel, {
        hitTolerance: 10,
        layerFilter: function (layer) {
          return layer === routeVecLayer
        }
      })
      const feature = features.find((feature) => feature.getGeometry() instanceof LineString)

      if (!feature) {
        this.hoveredRecord = emptyRecord
        popupOverlay.setPosition(undefined)
        return
      }

      if (!this.isIndexed) {
        popupOverlay.setPosition(e.coordinate)
        this.popupFreeze = false
        return
      }

      if (this.debug) console.time('KDTree')
      this.hoveredRecord = this.findNearestRecord(e.coordinate)
      if (this.debug) console.timeEnd('KDTree')

      popupOverlay.setPosition([
        this.hoveredRecord.positionLong,
        this.hoveredRecord.positionLat
      ] as Coordinate)
      this.popupFreeze = e.type == 'singleclick'
    },

    showPopUpRecord(record: Record) {
      popupOverlay.setPosition(undefined)
      if (!record) return

      this.popupRecord = record
      popupOverlay.setPosition([
        this.popupRecord.positionLong,
        this.popupRecord.positionLat
      ] as Coordinate)
      this.popupFreeze = true
    },

    updateExtent() {
      this.$nextTick(() => {
        const extent = routeVecLayer.getSource()!.getExtent()
        if (isEmpty(extent)) return
        map.getView().fit(extent, { padding: [50, 50, 50, 50] })

        map.removeControl(zoomToExtent)
        zoomToExtent = new ZoomToExtent({ extent: map.getView().getViewStateAndExtent().extent })
        map.addControl(zoomToExtent)
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
    isIndexed(): Boolean {
      return kdbush.value._finished
    }
  },
  mounted() {
    const maximizeIcon = document.createElement('i')
    maximizeIcon.setAttribute('class', 'fa-solid fa-expand')

    const minimizeIcon = document.createElement('i')
    minimizeIcon.setAttribute('class', 'fa-solid fa-compress')

    popupOverlay.setElement(document.getElementById('popup')!)

    map = new OlMap({
      overlays: [popupOverlay],
      controls: defaultControls(),
      layers: [tileLayer, routeVecLayer, concealVecLayer, trimVecLayer, pointLayer],
      view: view
    })

    map.setTarget(this.$refs.map as HTMLElement)

    map.addControl(new FullScreen({ label: maximizeIcon, labelActive: minimizeIcon }))
    map.addControl(new ScaleLine())

    map.on('change:size', this.updateExtent)
    map.on('pointermove', this.lineStringFeatureListener)
    map.on('singleclick', this.lineStringFeatureListener)

    this.updateMapSource(this.features)

    this.kdbushIndexing(this.sessions)
    this.createConcealFeatures(this.sessions)
    this.createTrimFeatures(this.sessions)
  },
  unmounted() {
    map.dispose()
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

.no-map-data {
  font-size: 1.25rem;
}
</style>
