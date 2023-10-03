<template>
  <PlotFigure :options="options"></PlotFigure>
</template>

<script lang="ts">
import { ActivityFile } from '@/spec/activity'
import * as Plot from '@observablehq/plot'
import * as d3 from 'd3'
import PlotFigure from './PlotFigure'

export default {
  components: {
    PlotFigure: PlotFigure
  },
  props: {
    activityFile: ActivityFile
  },
  data() {
    return {}
  },
  watch: {
    activityFile: {
      handler() {}
    }
  },
  computed: {
    options: function () {
      console.log('options ', this.activityFile?.records)
      const data = this.activityFile?.records || []
      return {
        y: { grid: true },
        color: {
          interpolate: d3.piecewise(d3.interpolateRgb.gamma(2.2), [
            'green',
            'green',
            'yellow',
            'red',
            'red'
          ]),
          domain: [-1, 1]
        },
        marks: [
          Plot.ruleX(
            data,
            Plot.pointerX({ x: (d) => new Date(d.timestamp), py: 'altitude', stroke: 'red' })
          ),
          Plot.dot(
            data,
            Plot.pointerX({ x: (d) => new Date(d.timestamp), y: 'altitude', stroke: 'red' })
          ),
          Plot.areaY(data, { x: (d) => new Date(d.timestamp), y: 'altitude', z: null }),
          Plot.lineY(
            data,
            Plot.map(
              { stroke: Plot.window({ k: 2, reduce: 'difference' }) },
              {
                x: (d) => new Date(d.timestamp),
                y: 'altitude',
                z: null,
                stroke: 'altitude',
                strokeWidth: 5,
                strokeOpacity: 0.7,
                tip: 'x'
              }
            )
          ),
        ]
      }
    }
  },
  methods: {
    getOptions() {}
  },
  mounted() {}
}
</script>

<style>
text {
  color: black;
}
</style>
