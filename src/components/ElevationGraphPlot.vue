<template>
  <div>
    <PlotFigure :options="options" defer :onRender="plotRendered"> </PlotFigure>
  </div>
</template>

<script lang="ts">
import { ActivityFile, Record } from '@/spec/activity'
import * as Plot from '@observablehq/plot'
import PlotFigure from './PlotFigure'
import { DateTime } from 'luxon'
import { toHuman } from '@/toolkit/date'
import { distanceToHuman } from '@/toolkit/distance'

export default {
  components: {
    PlotFigure
  },
  props: {
    activityFiles: Array<ActivityFile>,
    activityTimezoneOffset: Array<Number>
  },
  data() {
    return {
      x: (d: Record) => d.totalDistance,
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
    options: function () {
      // combine and get total distance
      let data: Record[] = []
      let lastDistance: Number = 0
      let yMax: number = 0
      let yMin: number = Number.MAX_VALUE

      this.activityFiles?.forEach((activityFile: ActivityFile, activityIndex: number) => {
        if (activityFile.records.length > 0) {
          activityFile.records.map((d, i) => {
            d.totalDistance = lastDistance + d.distance
            if (d.altitude > yMax) yMax = d.altitude
            if (d.altitude < yMin) yMin = d.altitude
            d.activityIndex = activityIndex
            d.recordIndex = i
          })
          lastDistance = activityFile.records[activityFile.records.length - 1].totalDistance
          data = data.concat(activityFile.records)
        }
      })

      // window
      const k = (data.length < 50 ? 50 : data.length) / 50

      return {
        x: {
          grid: true,
          label: 'Distance (km)',
          transform: (d) => d / 1000
        },
        y: {
          grid: true,
          label: 'Altitude (m)',
          nice: true
        },
        color: {
          type: 'diverging',
          scheme: 'RdYlGn',
          reverse: true,
          pivot: 0,
          symmetric: true
        },
        marks: [
          Plot.areaY(
            data,
            Plot.windowY(k, {
              x: this.x,
              y: this.y,
              z: null,
              fill: '#2A303F',
              curve: 'basis',
              y1: yMin
            })
          ),
          Plot.lineY(
            data,
            Plot.windowY(k, {
              x: this.x,
              y: this.y,
              z: null,
              stroke: 'grade',
              curve: 'basis',
              strokeWidth: 4,
              strokeOpacity: 0.5
            })
          ),
          Plot.ruleX(data, Plot.pointerX({ x: this.x, py: this.y, stroke: '#f15a22' })),
          Plot.dot(data, Plot.pointerX({ x: this.x, y: this.y, r: 10, stroke: '#f15a22' })),
          Plot.text(
            data,
            Plot.pointerX({
              px: this.x,
              py: this.y,
              dy: 5,
              fontSize: 12,
              frameAnchor: 'top',
              fontVariant: 'tabular-nums',
              text: (d) => {
                const diff = DateTime.fromISO(d.timestamp).diff(DateTime.fromISO(data[0].timestamp))
                const text = [
                  `Distance ${distanceToHuman(d.totalDistance, 2)}`,
                  `Duration ${toHuman(diff, 'seconds', {
                    unitDisplay: 'short'
                  })}`,
                  `Altitude ${d.altitude.toFixed(2)}`,
                  `Grade ${Math.round(d.grade)}%`
                ].join('  ')
                return text
              }
            })
          )
        ]
      }
    }
  },
  methods: {
    distanceToHuman: distanceToHuman,
    toHuman: toHuman,
    getOptions() {},
    plotRendered(plot: (SVGSVGElement | HTMLElement) & Plot.Plot) {
      plot.addEventListener('input', () => {
        this.$emit('record', plot.value)
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
