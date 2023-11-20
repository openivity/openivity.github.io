<template>
  <div style="display: flex; flex-direction: column; height: 100%">
    <div class="graph-detail">
      <div class="detail">
        <span><i class="fa-solid fa-road"></i></span>
        <span class="detail-value"> {{ hoveredResult ? hoveredResult.distance ?? '-' : '-' }}</span>
      </div>
      <div class="detail">
        <span><i class="fa-solid fa-hourglass-half"></i></span>
        <span class="detail-value"> {{ hoveredResult ? hoveredResult.duration ?? '-' : '-' }}</span>
      </div>
      <div class="detail">
        <span><i class="fa-solid fa-mountain"></i></span>
        <span class="detail-value">
          {{ hoveredResult ? hoveredResult.altitude ?? '-' : '-' }} masl</span
        >
      </div>
      <div class="detail">
        <span><i class="fa-solid fa-angle-left"></i></span>
        <span class="detail-value">{{ hoveredResult ? hoveredResult.grade ?? '-' : '-' }}%</span>
      </div>
    </div>
    <div class="graph" style="flex-grow: 1">
      <PlotFigure :options="options" defer :onRender="plotRendered"></PlotFigure>
    </div>
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
      firstData: Record,
      hovered: Record,
      cumulativeData: [],
      yMin: Number.MAX_VALUE,
      yMax: 0,
      plotSizeWidth: 0,
      plotSizeHeight: 0
    }
  },
  watch: {
    activityFiles: {
      handler(activityFiles) {
        let data: Record[] = []
        let lastDistance: Number = 0

        activityFiles?.forEach((activityFile: ActivityFile, activityIndex: number) => {
          if (activityFile.records.length > 0) {
            activityFile.records.map((d, i) => {
              d.totalDistance = lastDistance + d.distance
              if (typeof d.altitude === 'number' && d.altitude > this.yMax) this.yMax = d.altitude
              if (typeof d.altitude === 'number' && d.altitude < this.yMin) this.yMin = d.altitude
              d.activityIndex = activityIndex
              d.recordIndex = i
            })

            lastDistance = activityFile.records[activityFile.records.length - 1].totalDistance
            data = data.concat(activityFile.records)
          }
        })
        this.firstData = data[0]
        this.cumulativeData = data
      }
    }
  },
  computed: {
    hoveredResult: function () {
      if (!this.hovered || !this.firstData) return null

      const diff = DateTime.fromISO(this.hovered.timestamp ?? '').diff(
        DateTime.fromISO(this.firstData.timestamp ?? '')
      )
      return {
        distance: distanceToHuman(this.hovered.totalDistance ?? 0, 2),
        duration: toHuman(diff, 'seconds', {
          unitDisplay: 'short'
        }),
        altitude: (this.hovered.altitude ?? 0).toFixed(2),
        grade: Math.round(this.hovered.grade ?? 0)
      }
    },
    options: function () {
      // window
      const k = (this.cumulativeData.length < 50 ? 50 : this.cumulativeData.length) / 50

      return {
        width: this.plotSizeWidth <= 0 ? undefined : this.plotSizeWidth,
        height: this.plotSizeHeight <= 0 ? undefined : this.plotSizeHeight,
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
            this.cumulativeData,
            Plot.windowY(k, {
              x: this.x,
              y: this.y,
              z: null,
              fill: '#2A303F',
              curve: 'basis',
              y1: this.yMin
            })
          ),
          Plot.lineY(
            this.cumulativeData,
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
          Plot.ruleX(
            this.cumulativeData,
            Plot.pointerX({ x: this.x, py: this.y, stroke: '#f15a22' })
          ),
          Plot.dot(
            this.cumulativeData,
            Plot.pointerX({ x: this.x, y: this.y, r: 10, stroke: '#f15a22' })
          )
          // Plot.text(
          //   this.cumulativeData,
          //   Plot.pointerX({
          //     // px: this.x,
          //     // py: this.y,
          //     dy: 5,
          //     fontSize: 12,
          //     frameAnchor: 'top',
          //     fontVariant: 'tabular-nums'
          //   })
          // )
        ]
      }
    }
  },
  methods: {
    distanceToHuman: distanceToHuman,
    toHuman: toHuman,
    getOptions() {},
    plotRendered(plot: (SVGSVGElement | HTMLElement) & Plot.Plot) {
      this.onResize()
      plot.addEventListener('input', () => {
        this.hovered = plot.value
        this.$emit('record', plot.value)
      })
    },
    onResize() {
      if (!this.$el) return

      let width = this.$el.offsetWidth
      let height = this.$el.offsetHeight - this.$el.querySelector('.graph-detail').offsetHeight - 1

      if (this.plotSizeWidth != width) this.plotSizeWidth = width
      if (this.plotSizeHeight != height) this.plotSizeHeight = height
      console.log(this.plotSizeWidth, this.plotSizeHeight)
    }
  },
  updated() {},
  mounted() {
    this.$nextTick(() => {
      window.addEventListener('resize', this.onResize)
      this.onResize()
    })
  },
  unmounted() {
    window.removeEventListener('resize', this.onResize)
  }
}
</script>

<style lang="scss" scoped>
::v-deep svg {
  background-color: transparent !important;
}
.graph-detail {
  display: grid;
  grid-template-columns: 1fr 1.5fr 1fr 0.7fr;
  color: var(--color-title);
  padding-bottom: 5px;

  .detail-value {
    text-align: center;
  }

  span {
    font-size: 0.7em;
    font-weight: bold;
  }
}
.fa-solid {
  width: 25px;
  text-align: center;
}
</style>
./legacy/PlotFigure