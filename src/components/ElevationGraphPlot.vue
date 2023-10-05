<template>
  <PlotFigure :options="options" :mark="marks" defer :onRender="plotRendered"></PlotFigure>
</template>

<script lang="ts">
import { ActivityFile, Record } from '@/spec/activity'
import * as Plot from '@observablehq/plot'
import * as d3 from 'd3'
import PlotFigure from './PlotFigure'

export default {
  components: {
    PlotFigure: PlotFigure
  },
  props: {
    activityFiles: Array<ActivityFile>
  },
  data() {
    return {
      x: (d: Record) => new Date(String(d.timestamp)),
      y: (d: Record) => d.altitude,
      hovered: new Record()
    }
  },
  watch: {
    activityFile: {
      handler() {}
    }
  },
  computed: {
    marks: function () {
      let data: Record[] = []
      this.activityFiles?.forEach((activityFile: ActivityFile) => {
        data = data.concat(activityFile.records)
      })
      return [
        Plot.areaY(data, {
          x: this.x,
          y: this.y,
          z: null,
          fill: '#2A303F'
        }),
        Plot.lineY(
          data,
          Plot.map(
            {
              stroke: Plot.window({ k: 5, reduce: 'difference' })
            },
            {
              x: this.x,
              y: this.y,
              z: null,
              stroke: this.y,
              strokeWidth: 5,
              strokeOpacity: 0.7,
              ariaDescription: 'elevation-line'
            }
          )
        ),
        Plot.ruleX(data, Plot.pointerX({ x: this.x, py: this.y, stroke: 'lawngreen' })),
        Plot.dot(data, Plot.pointerX({ x: this.x, y: this.y, r: 10, stroke: 'lawngreen' })),
        Plot.tip(
          data,
          Plot.pointerX({
            x: this.x,
            y: this.y,
            fill: '#2A303F',
            fontSize: 16,
            fontWeight: 'bolder',
            channels: {
              altitude: {
                label: 'Altitude',
                value: 'altitude'
              },
              speed: {
                label: 'Speed',
                value: 'speed'
              }
            },
            format: {
              y: false,
              x: false
            }
          })
        )
      ]
    },
    options: function () {
      return {
        x: {
          grid: true,
          label: 'time'
        },
        y: { grid: true, label: 'Altitude (m)' },
        color: {
          interpolate: d3.piecewise(d3.interpolateRgb.gamma(2.2), [
            'lawngreen',
            'lawngreen',
            'yellow',
            'yellow',
            'yellow',
            'yellow',
            'yellow',
            'red',
            'red'
          ]),
          domain: [-0.5, 0.5]
        }
      }
    }
  },
  methods: {
    getOptions() {},
    plotRendered(plot: (SVGSVGElement | HTMLElement) & Plot.Plot) {
      plot.addEventListener('input', () => {
        this.$emit('input', plot.value)
      })
    }
  },
  updated() {},
  mounted() {}
}
</script>

<style>
svg {
  background-color: transparent !important;
}
</style>
