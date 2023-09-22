<template>
  <div id="map" style="width: 100%; height: 100%"></div>
</template>

<script lang="ts">
import 'ol/ol.css'

import Map from 'ol/Map'
import View from 'ol/View'
import { GeoJSON } from 'ol/format'
import TileLayer from 'ol/layer/Tile'
import VectorImageLayer from 'ol/layer/VectorImage'
import OSM from 'ol/source/OSM'
import VectorSource from 'ol/source/Vector'
import { Stroke, Style } from 'ol/style'
import type { Geometry } from 'ol/geom'

export default {
  props: {
    geojson: {}
  },
  data() {
    return {
      vec: new VectorImageLayer({
        source: new VectorSource({
          features: []
        }),
        visible: true,
        style: new Style({
          stroke: new Stroke({
            // color: [211, 84, 0, 1.0],
            color: [241, 90, 34, 1.0],
            width: 5
          })
        })
      }) as VectorImageLayer<VectorSource<Geometry>>,
      map: new Map({
        layers: [
          new TileLayer({
            source: new OSM()
          })
        ],
        view: new View({
          enableRotation: false,
          projection: 'EPSG:4326' // WGS84: World Geodetic System 1984
        })
      })
    }
  },
  watch: {
    geojson: {
      handler(geojson: JSON) {
        this.updateSource(geojson)
      }
    }
  },
  methods: {
    updateSource(geojson: JSON) {
      const view = this.map.getView()
      const source = this.vec.getSource()!

      const features = new GeoJSON().readFeatures(geojson)

      source.clear()
      source.addFeatures(features)

      view.fit(source.getExtent(), {
        padding: [50, 50, 50, 50]
      })
    }
  },
  mounted() {
    let vec = this.vec as VectorImageLayer<VectorSource<Geometry>>
    this.map.addLayer(vec)
    this.map.setTarget('map')
  }
}
</script>
