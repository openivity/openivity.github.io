<template>
  <div
    id="map"
    style="width: 100%; height: 100%"
    v-show="activityFiles && activityFiles.length > 0"
  ></div>
  <div id="popup" class="ol-popup">
    <div class="popup-content">
      <div>
        <span>Time:</span>
        <span>
          {{
            popupRecord.timestamp
              ? toTimezoneDateString(popupRecord.timestamp, popupTimezoneOffsetHours)
              : '-'
          }}
        </span>
      </div>
      <div>
        <span>Distance:</span>
        <span>{{ popupRecord.distance ? (popupRecord.distance / 1000).toFixed(2) : '-' }}</span>
        <span>&nbsp;km</span>
      </div>
      <div>
        <span>Speed:</span>
        <span>{{ popupRecord.speed ? ((popupRecord.speed * 3600) / 1000).toFixed(2) : '-' }}</span>
        <span>&nbsp;km/h</span>
      </div>
      <div>
        <span>Cadence:</span>
        <span>{{ popupRecord.cadence ? popupRecord.cadence : '-' }}</span>
        <span>&nbsp;rpm</span>
      </div>
      <div>
        <span>Heart Rate:</span>
        <span>{{ popupRecord.heartRate ? popupRecord.heartRate : '-' }}</span>
        <span>&nbsp;bpm</span>
      </div>
      <div>
        <span>Power:</span>
        <span>{{ popupRecord.power ? popupRecord.power : '-' }}</span>
        <span>&nbsp;watts</span>
      </div>
      <div>
        <span>Altitude:</span>
        <span>{{ popupRecord.altitude ? popupRecord.altitude?.toFixed(2) : '-' }}</span>
        <span>&nbsp;masl</span>
      </div>
      <div>
        <span>Location:</span>
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
</template>

<script lang="ts">
import 'ol/ol.css'
import destinationPointIcon from '@/assets/map/destination-point.svg'
import startingPointIcon from '@/assets/map/starting-point.svg'

import OlMap from 'ol/Map'
import View from 'ol/View'
import { GeoJSON } from 'ol/format'
import TileLayer from 'ol/layer/Tile'
import VectorImageLayer from 'ol/layer/VectorImage'
import OSM from 'ol/source/OSM'
import VectorSource from 'ol/source/Vector'
import { Icon, Stroke, Style } from 'ol/style'
import { Point, Geometry, SimpleGeometry, LineString } from 'ol/geom'
import { FullScreen, ScaleLine, ZoomToExtent, defaults as defaultControls } from 'ol/control.js'
import { Feature, MapBrowserEvent, Overlay } from 'ol'
import type { Coordinate } from 'ol/coordinate'
import { ActivityFile, Record } from '@/spec/activity'
import { toStringHDMS } from 'ol/coordinate.js'
import { toTimezoneDateString } from '@/toolkit/date'
import type { Extent } from 'ol/extent'

const maximizeIcon = document.createElement('i')
const minimizeIcon = document.createElement('i')
maximizeIcon.setAttribute('class', 'fa-solid fa-expand')
minimizeIcon.setAttribute('class', 'fa-solid fa-compress')

export default {
  props: {
    geojsons: Array<GeoJSON>,
    activityFiles: Array<ActivityFile>,
    timezoneOffsetHoursList: Array<Number>
  },
  data() {
    return {
      popupFreeze: new Boolean(),
      popupRecord: new Record(),
      popupTimezoneOffsetHours: new Number(0),
      popupOverlay: new Overlay({}),
      vec: new VectorImageLayer({
        source: new VectorSource({
          features: []
        }),
        visible: true,
        style: new Style({
          stroke: new Stroke({
            // color: [232, 65, 24, 1.0],
            color: [60, 99, 130, 1.0],
            width: 3
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
          enableRotation: false,
          projection: 'EPSG:4326' // WGS84: World Geodetic System 1984
        })
      }),
      zoomToExtent: new ZoomToExtent()
    }
  },
  watch: {
    geojsons: {
      handler(geojsons: Array<GeoJSON>) {
        this.updateMapSource(geojsons)
      }
    },
    zoomToExtent: {
      handler(newValue: ZoomToExtent, oldValue: ZoomToExtent) {
        this.map.removeControl(oldValue)
        this.map.addControl(newValue)
      }
    }
  },
  methods: {
    toStringHDMS: toStringHDMS,
    toTimezoneDateString: toTimezoneDateString,

    updateMapSource(geojsons: Array<GeoJSON>) {
      this.popupOverlay.setPosition(undefined)

      const view = this.map.getView()
      const source = this.vec.getSource()!
      const features = new Array<Feature>()

      for (let i = 0; i < geojsons.length; i++) {
        const feature = new GeoJSON().readFeature(geojsons[i])
        feature.setId('lineString-' + i)
        features.push(feature)
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
            image: new Icon({ crossOrigin: 'anonymous', src: startingPointIcon, scale: 0.8 })
          })
        )
        startingPoint.setId('startingPoint-' + i)
        startingPoints.push(startingPoint)

        const destinationPoint = new Feature(
          new Point((features[i]?.getGeometry() as SimpleGeometry).getLastCoordinate())
        )
        destinationPoint.setStyle(
          new Style({
            image: new Icon({ crossOrigin: 'anonymous', src: destinationPointIcon, scale: 0.8 })
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
      let nearestRecord: Record = new Record()
      let nearestEuclidean: number = Number.MAX_VALUE

      this.popupTimezoneOffsetHours = this.timezoneOffsetHoursList![index] as Number
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

    lineStringFeatureListener(type: string, e: MapBrowserEvent<any>) {
      if (type == 'singleclick')
        this.popupFreeze = !this.popupFreeze && this.popupOverlay.getPosition() != undefined
      if (!this.popupFreeze) this.popupOverlay.setPosition(undefined)
      if (this.popupFreeze == true && this.popupOverlay.getPosition() != undefined) return

      this.map.forEachFeatureAtPixel(
        e.pixel,
        (feature) => {
          if (!(feature.getGeometry() instanceof LineString)) return
          this.popupRecord = this.findNearestRecord(feature.getId() as string, e.coordinate)
          this.popupOverlay.setPosition([
            this.popupRecord.positionLong,
            this.popupRecord.positionLat
          ] as Coordinate)
          this.popupFreeze = type == 'singleclick'
        },
        { hitTolerance: 10 }
      )
    },
    updateExtent() {
      this.map.getView().fit(this.vec.getSource()!.getExtent(), { padding: [50, 50, 50, 50] })
      this.zoomToExtent = new ZoomToExtent({
        extent: this.map.getView().getViewStateAndExtent().extent
      })
    }
  },
  mounted() {
    this.popupOverlay = new Overlay({ element: document.getElementById('popup')! })
    this.map.addOverlay(this.popupOverlay as Overlay)

    this.map.addLayer(this.vec as VectorImageLayer<VectorSource<Geometry>>)
    this.map.setTarget('map')

    this.map.addControl(new FullScreen({ label: maximizeIcon, labelActive: minimizeIcon }))
    this.map.addControl(new ScaleLine())

    this.map.once('precompose', () => this.updateExtent()) // init

    this.map.on('change:size', () => this.updateExtent())
    this.map.on('pointermove', (e) => this.lineStringFeatureListener(e.type, e))
    this.map.on('singleclick', (e) => this.lineStringFeatureListener(e.type, e))
  }
}
</script>
<style>
.ol-popup {
  position: absolute;
  color: #000;
  background-color: white;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
  padding: 15px;
  border-radius: 10px;
  border: 1px solid #cccccc;
  bottom: 12px;
  left: -50px;
  min-width: 250px;
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
  border-top-color: white;
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
  font-size: 10px;
}

.popup-content span {
  display: inline-block;
}

.popup-content span:nth-child(1) {
  width: 60px;
}

.popup-content span:nth-child(2) {
  font-weight: 500;
}
</style>
