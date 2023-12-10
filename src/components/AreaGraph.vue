<template>
  <div class="col-12 h-100">
    <div
      class="line-graph-hover row m-0"
      style="cursor: pointer"
      data-bs-toggle="collapse"
      v-bind:data-bs-target="`#${graphContainerName}`"
      aria-expanded="false"
      v-bind:aria-controls="graphContainerName"
    >
      <div class="d-flex pb-1 pe-2">
        <h6 class="title mb-0 col text-start">
          <i class="fa-solid fa-caret-down when-opened"></i>
          <i class="fa-solid fa-caret-right when-closed"></i>
          {{ name }}
          <i :class="['fa-solid', icon]"></i>
        </h6>
        <span class="lh-sm pe-1 fs-6 col text-end">
          {{ formatUnit(scaleUnit(recordView[recordField as keyof Record] as number)) }}
          <span>{{ unit }}</span>
        </span>
      </div>
    </div>
    <div class="collapse show" v-bind:id="graphContainerName">
      <div class="row label mx-0 pt-1 pb-1">
        <span class="col text-start">↑ {{ yLabel }}</span>
        <span class="col text-end">{{ xAxisLabel }} →</span>
      </div>
      <div :ref="graphName"></div>
      <div class="graph-summary pt-2 pb-1 px-0">
        <div class="row mx-0 px-0" v-for="(detail, index) in details" v-bind:key="index">
          <div class="d-flex px-0">
            <span class="col text-start summary-text">{{ detail.title }}</span>
            <span class="col-auto summary-text fs-6 pe-1">{{ detail.value }}</span>
            <span class="col-auto summary-text pe-0">{{ unit }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Record } from '@/spec/activity'
import { formatPace } from '@/toolkit/pace'
import * as d3 from 'd3'

export class Detail {
  title: string = ''
  value: string = ''

  constructor(data?: any) {
    const casted = data as Detail
    this.title = casted?.title
    this.value = casted?.value
  }
}

export default {
  props: {
    name: {
      type: String,
      required: true
    },
    icon: String,
    recordField: { type: String, required: true },
    graphRecords: {
      type: Array<Record>,
      required: true,
      default: []
    },
    records: {
      type: Array<Record>,
      required: true,
      default: []
    },
    yLabel: { type: String, required: true },
    unit: { type: String, required: true },
    color: { type: String, required: true },
    pointerColor: { type: String, default: '#78e08f' },
    receivedRecord: Record,
    receivedRecordFreeze: Boolean,
    details: Array<Detail>
  },
  data() {
    return {
      elemWidth: 0,
      hoveredRecord: new Record(),
      hoveredRecordFreeze: new Boolean(),
      recordView: new Record(),
      scaleByDistanceOrTimestamp: null as unknown as (rec: Record) => number
    }
  },
  watch: {
    receivedRecord: {
      handler(record: Record) {
        this.recordView = record
        const pointer = d3
          .select(this.$refs[`${this.graphName}`] as HTMLElement)
          .select('g#pointer')

        if (JSON.stringify(record) == JSON.stringify(new Record())) {
          pointer.style('opacity', 0)
          return
        }

        if (this.scaleByDistanceOrTimestamp == null) return

        pointer.style('opacity', 1)
        pointer.attr('transform', `translate(${this.scaleByDistanceOrTimestamp(record)}, 0)`)
      }
    },
    receivedRecordFreeze: {
      handler(freeze: Boolean) {
        this.hoveredRecordFreeze = freeze
        if (!freeze) {
          const pointer = d3
            .select(this.$refs[`${this.graphName}`] as HTMLElement)
            .select('g#pointer')
          pointer.style('opacity', 0)
        }
      }
    },
    hoveredRecord: {
      handler(record: Record) {
        this.recordView = record
        this.$emit('hoveredRecord', record)
      }
    },
    hoveredRecordFreeze: {
      handler(freeze: Boolean) {
        this.$emit('hoveredRecordFreeze', freeze)
      }
    },
    graphRecords: {
      handler() {
        this.$nextTick(() => requestAnimationFrame(() => this.renderGraph()))
      }
    }
  },
  computed: {
    hasDistance(): boolean {
      for (let i = 0; i < this.graphRecords.length; i++) {
        if (this.graphRecords[i].distance) return true
      }
      return false
    },
    graphName(): string {
      return this.name.toLowerCase().replace(/\s/g, '-') + '-graph'
    },
    graphContainerName(): string {
      return this.graphName + '-container'
    },
    xAxisLabel(): string {
      return this.hasDistance ? 'Dist. (km)' : 'Time'
    }
  },
  methods: {
    formatLabel(value: number | { valueOf(): number }): string {
      const val = typeof value === 'number' ? value : value.valueOf()
      if (this.name == 'Pace') {
        return formatPace(val)
      }
      if (this.name == 'Speed') {
        value = (val * 3600) / 1000
      }
      return d3.format('~s')(val)
    },
    renderGraph() {
      const graphRecords = this.graphRecords

      const marginTop = 5
      const marginRight = 10
      const marginBottom = 20
      const marginLeft = 40

      const width = this.elemWidth
      const height = 190 - marginBottom

      const xTicks = width > 720 ? 10 : 5
      const yTicks = 3

      const container = d3.select(this.$refs[`${this.graphName}`] as HTMLElement)

      container
        .select('svg')
        .on('pointerdown', null)
        .on('pointermove', null)
        .on('mouseleave', null)
        .remove()

      if (graphRecords.length == 0) return

      const svg = container
        .append('svg')
        .attr('width', width)
        .attr('height', height)
        .style('touch-action', 'none')
        .style('user-select', 'none')
        .style('--webkit-user-select', 'none') /* Safari */
        .style('--ms-user-select', 'none') /* IE 10 and IE 11 */

      // Creating Scales
      const xScaleDistance = d3
        .scaleLinear()
        .domain(d3.extent(graphRecords, (d) => (d.distance ?? 0) / 1000) as number[])
        .range([marginLeft, width - marginRight])

      const xScaleTimestamp = d3
        .scaleTime()
        .domain(d3.extent(graphRecords, (d) => new Date(d.timestamp!).getTime()) as number[])
        .range([marginLeft, width - marginRight])

      const xScale = this.hasDistance ? xScaleDistance : xScaleTimestamp
      this.scaleByDistanceOrTimestamp = (rec: Record): number => {
        if (this.hasDistance) return xScale((rec.distance ?? 0) / 1000)
        return xScale(new Date(rec.timestamp!).getTime())
      }

      let altitudeExtent = d3.extent(graphRecords, (d) => d.altitude) as number[]
      const altitudeExtentMinSize = 100 // masl
      const currentSize = altitudeExtent[1] - altitudeExtent[0]
      if (currentSize < altitudeExtentMinSize)
        altitudeExtent[1] += altitudeExtentMinSize - currentSize

      const altitudeScale = d3
        .scaleLinear()
        .domain(altitudeExtent)
        .rangeRound([height - marginBottom, marginTop])
        .nice()

      const yExtent = d3.extent(
        graphRecords,
        (d) => this.scaleUnit(d[this.recordField as keyof Record] as number) ?? 0
      ) as number[]

      const yScale = d3
        .scaleLinear()
        .domain(yExtent)
        .rangeRound([height - marginBottom, marginTop])
        .nice()

      if (this.name == 'Pace') yScale.domain([yExtent[1], yExtent[0]])

      // Add X & Y Axis
      svg
        .append('g')
        .style('font-size', '0.8em')
        .attr('transform', `translate(0,${height - marginBottom})`)
        .call(d3.axisBottom(xScale).ticks(xTicks))

      svg
        .append('g')
        .style('font-size', '0.8em')
        .attr('transform', `translate(${marginLeft},0)`)
        .call(d3.axisLeft(yScale).ticks(yTicks).tickFormat(this.formatLabel))

      // Create Grid Lines
      svg
        .append('g')
        .attr('stroke', 'currentColor')
        .attr('stroke-opacity', 0.1)
        .style('stroke', 'lightgray')
        .style('stroke-dasharray', '2,2')
        .call((g) =>
          g
            .append('g')
            .selectAll('line')
            .data(xScale.ticks(xTicks) as number[])
            .join('line')
            .attr('x1', (d) => 0.5 + xScale(d))
            .attr('x2', (d) => 0.5 + xScale(d))
            .attr('y1', marginTop)
            .attr('y2', height - marginBottom)
        )
        .call((g) =>
          g
            .append('g')
            .selectAll('line')
            .data(yScale.ticks(yTicks))
            .join('line')
            .attr('y1', (d) => 0.5 + yScale(d))
            .attr('y2', (d) => 0.5 + yScale(d))
            .attr('x1', marginLeft)
            .attr('x2', width - marginRight)
        )

      const minDomain = d3.min(altitudeScale.domain())!

      // Add Altitude Area
      const altitudeArea = d3
        .area<Record>()
        .curve(d3.curveBasis)
        .x(this.scaleByDistanceOrTimestamp)
        .y0(altitudeScale(minDomain))
        .y1((d) => altitudeScale(d.altitude ?? minDomain))

      svg
        .append('g')
        .append('path')
        .datum(graphRecords)
        .transition()
        .attr('fill', 'lightgrey')
        .style('opacity', 0.9)
        .attr('d', altitudeArea)

      // Add Speed Area
      const area = d3
        .area<Record>()
        .curve(d3.curveBasis)
        .x(this.scaleByDistanceOrTimestamp)
        .y0(yScale(d3.min(yScale.domain())!))
        .y1((d): number => {
          let value = d[this.recordField as keyof Record]
          value = value != null ? value : yScale.domain()[0]
          return yScale(this.scaleUnit(value as number) as number)
        })

      if (this.name == 'Pace') area.y0(yScale(d3.max(yScale.domain())!))

      svg
        .append('g')
        .append('path')
        .datum(graphRecords)
        .transition()
        .attr('fill', `${this.color}`)
        .style('opacity', 0.8)
        .attr('d', area)

      // Add Pointer
      const pointer = svg
        .append('g')
        .attr('id', 'pointer')
        .style('opacity', 0)
        .call((g) => {
          g.append('line')
            .attr('x1', 0)
            .attr('y1', marginTop)
            .attr('x2', 0)
            .attr('y2', height - marginBottom)
            .attr('stroke', `${this.pointerColor}`)
            .attr('stroke-width', 1.5)
        })
        .call((g) => {
          g.append('polygon')
            .attr('points', '0,30 -5,20 5,20')
            .attr('fill', `${this.pointerColor}`)
            .attr('transform', `translate(0, ${-marginTop - 10})`)
        })

      // Add Events
      const pointerListener = (e: Event) => {
        if (e.type == 'pointerdown' && this.hoveredRecordFreeze) this.hoveredRecordFreeze = false
        if (this.hoveredRecordFreeze == true) return

        const [px] = d3.pointer(e)
        const [xMin, xMax] = xScale.range()

        pointer.style('opacity', 1)

        if (px <= xMin) {
          pointer.attr('transform', `translate(${xMin}, 0)`)
          this.hoveredRecord = this.records[0]
          return
        } else if (px >= xMax) {
          pointer.attr('transform', `translate(${xMax}, 0)`)
          this.hoveredRecord = this.records[this.records.length - 1]
          return
        }

        const pointerPercentage = (px - xMin) / (xMax - xMin)
        const lookupIndex = Math.round(pointerPercentage * (this.records.length - 1))

        let nearestRecord: Record = new Record()
        let dx = Number.MAX_VALUE
        if (this.scaleByDistanceOrTimestamp(this.records[lookupIndex]) <= px) {
          for (let i = lookupIndex; i < this.records.length; i++) {
            const delta = Math.abs(px - this.scaleByDistanceOrTimestamp(this.records[i]))
            if (delta > dx) break
            nearestRecord = this.records[i]
            dx = delta
          }
        } else {
          for (let i = lookupIndex; i >= 0; i--) {
            const delta = Math.abs(px - this.scaleByDistanceOrTimestamp(this.records[i]))
            if (delta > dx) break
            nearestRecord = this.records[i]
            dx = delta
          }
        }

        pointer.attr('transform', `translate(${this.scaleByDistanceOrTimestamp(nearestRecord)}, 0)`)

        this.hoveredRecord = nearestRecord
        this.hoveredRecordFreeze = e.type == 'pointerup'
      }

      svg.on('pointerdown', pointerListener)
      svg.on('pointerup', pointerListener)
      svg.on('pointermove', pointerListener, { passive: true })
      svg.on('mouseleave', () => {
        if (this.hoveredRecordFreeze) return

        this.hoveredRecord = new Record()
        pointer.style('opacity', 0)
      })
    },
    formatUnit(value: number | null): string {
      if (this.recordField === 'pace') {
        return value ? formatPace(value) : '-:-'
      }
      return value ? value.toFixed(2) : '0.00'
    },
    scaleUnit(value: number): number | null {
      if (this.recordField === 'speed') {
        value = (value * 3600) / 1000
      }
      if (isNaN(value)) return null
      return value
    },
    onResize() {
      // Ensure DOM is fully re-rendered after resize
      this.$nextTick(() => {
        // Prevent re-render graph when width is zero
        const $el = this.$el as HTMLElement
        if ($el.clientWidth == 0) return

        // Prevent re-render graph when width is not changing
        if ($el.clientWidth == this.elemWidth) return

        this.elemWidth = $el.clientWidth
        this.renderGraph()
      })
    }
  },
  mounted() {
    const $el = this.$el as HTMLElement

    this.$nextTick(() => {
      this.elemWidth = $el.clientWidth
      requestAnimationFrame(() => this.renderGraph())
    })

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && this.elemWidth == 0) {
            this.elemWidth = $el.clientWidth
            this.renderGraph()
          }
        })
      },
      {
        root: null,
        rootMargin: '0px',
        threshold: 0.5 // Trigger the callback when 50% of the element is visible);
      }
    )

    observer.observe($el)

    window.addEventListener('resize', this.onResize)
  },
  unmounted() {
    window.removeEventListener('resize', this.onResize)
  }
}
</script>
<style scoped>
.title {
  text-align: left;
}

.label {
  font-size: 0.9em;
}

.line-graph-hover {
  font-size: 0.9em;
}
.graph-summary {
  margin-left: 40px;
  padding-left: 0px;
  padding-right: 10px;
  margin-right: 10px;
}

.graph-summary div {
  padding-bottom: 1px;
  margin-bottom: 3px;
}
.summary-text {
  font-size: 0.9em;
  line-height: normal;
  align-self: center;
}

.graph-summary > div:nth-child(odd) {
  box-shadow: 0px 0.5px grey;
}

.collapsed > div > h6 > .fa-caret-down,
:not(.collapsed) > div > h6 > .fa-caret-right {
  display: none;
}
</style>
