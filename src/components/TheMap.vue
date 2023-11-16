<template>
  <div class="map-container w-100 h-100">
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
    </div>
  </div>
</template>

<script lang="ts">
import destinationPointIcon from '@/assets/map/destination-point.svg'
import startingPointIcon from '@/assets/map/starting-point.svg'
import 'ol/ol.css'

import { ActivityFile, Record } from '@/spec/activity'
import { toTimezoneDateString } from '@/toolkit/date'
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
    activityFiles: Array<ActivityFile>,
    receivedRecord: Record,
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
      pointer: new Feature()
    }
  },
  watch: {
    features: {
      handler(features: Array<Feature>) {
        this.$nextTick(() => {
          this.map.setTarget(this.$refs.map as HTMLElement)
          requestAnimationFrame(() => this.updateMapSource(features))
        })
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
      const destinationPoints = new Array<Feature>()
      for (let i = 0; i < features.length; i++) {
        const startingPoint = new Feature(
          new Point((features[i]?.getGeometry() as SimpleGeometry).getFirstCoordinate())
        )
        startingPoint.setStyle(
          new Style({
            image: new Icon({ crossOrigin: 'anonymous', src: startingPointIcon, scale: 1 })
          })
        )
        startingPoint.setId('startingPoint-' + i)
        startingPoints.push(startingPoint)

        const destinationPoint = new Feature(
          new Point((features[i]?.getGeometry() as SimpleGeometry).getLastCoordinate())
        )
        destinationPoint.setStyle(
          new Style({
            image: new Icon({ crossOrigin: 'anonymous', src: destinationPointIcon, scale: 1 })
          })
        )
        destinationPoint.setId('destinationPoint-' + i)
        destinationPoints.push(destinationPoint)
      }

      source.addFeatures(startingPoints)
      source.addFeatures(destinationPoints)

      this.zoomToExtent = new ZoomToExtent({ extent: view.getViewStateAndExtent().extent })
    },

    findNearestRecord(featureId: string, coordinate: Coordinate): Record {
      const index = Number(featureId.split('-')[1])
      this.popupActivityIndex = index
      let nearestRecord: Record = new Record()
      let nearestEuclidean: number = Number.MAX_VALUE

      this.popupTimezoneOffsetHours = this.activityFiles![index].timezone
      this.activityFiles![index].records?.forEach((record) => {
        if (!record.positionLong || !record.positionLat) return
        const euclidean = Math.abs(
          Math.sqrt(
            Math.pow(record.positionLong - coordinate[0], 2) +
              Math.pow(record.positionLat - coordinate[1], 2)
          )
        )
        if (euclidean < nearestEuclidean) {
          nearestEuclidean = euclidean
          nearestRecord = record
        }
      })

      return nearestRecord
    },

    lineStringFeatureListener(e: MapBrowserEvent<any>) {
      if (e.type == 'singleclick')
        this.popupFreeze = !this.popupFreeze && this.popupOverlay.getPosition() != undefined
      if (!this.popupFreeze) {
        this.hoveredRecord = new Record()
        this.popupOverlay.setPosition(undefined)
      }
      if (this.popupFreeze == true && this.popupOverlay.getPosition() != undefined) return

      this.map.forEachFeatureAtPixel(
        e.pixel,
        (feature) => {
          if (!(feature.getGeometry() instanceof LineString)) return

          this.hoveredRecord = this.findNearestRecord(feature.getId() as string, e.coordinate)
          this.popupOverlay.setPosition([
            this.hoveredRecord.positionLong,
            this.hoveredRecord.positionLat
          ] as Coordinate)
          this.popupFreeze = e.type == 'singleclick'
        },
        { hitTolerance: 10 }
      )
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
    }
  },
  mounted() {
    this.popupOverlay = new Overlay({ element: document.getElementById('popup')! })
    this.map.addOverlay(this.popupOverlay as Overlay)

    this.map.addLayer(this.vec as VectorImageLayer<VectorSource<Geometry>>)
    this.map.setTarget(this.$refs.map as HTMLElement)

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
<style scoped>
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
