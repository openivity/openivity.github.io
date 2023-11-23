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
      <div class="col text-start">
        <h6 class="title">
          <i class="fa-solid fa-caret-down when-opened"></i>
          <i class="fa-solid fa-caret-right when-closed"></i>
          {{ name }}
          <i :class="['fa-solid', icon]"></i>
        </h6>
      </div>
      <div class="col text-end">
        <span>
          <span class="pe-1 fs-6">
            {{ formatUnit(scaleUnit(recordView[recordField as keyof Record] as number)) }}
          </span>
          <span>{{ unit }}</span>
        </span>
      </div>
    </div>
    <div class="collapse show" v-bind:id="graphContainerName">
      <div :ref="graphName"></div>
      <div class="graph-summary pt-1 pb-1">
        <div class="row" v-for="(detail, index) in details" v-bind:key="index">
          <span class="col px-0 text-start">
            <span style="font-size: 0.9em">{{ detail.title }}</span>
          </span>
          <span class="col px-0 text-end">
            <span class="fs-6 pe-1">{{ detail.value }}</span>
            <span style="font-size: 0.9em">{{ unit }}</span>
          </span>
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
      xScale: d3.scaleLinear()
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
        pointer.style('opacity', 1)
        pointer.attr('transform', `translate(${this.xScale(record.distance! / 1000)}, 0)`)
      }
    },
    receivedRecordFreeze: {
      handler(freeze: Boolean) {
        this.hoveredRecordFreeze = freeze
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
    graphName(): string {
      return this.name.toLowerCase().replace(/\s/g, '-') + '-graph'
    },
    graphContainerName(): string {
      return this.graphName + '-container'
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
      if (graphRecords.length == 0) return

      const marginTop = 25
      const marginRight = 10
      const marginBottom = 25
      const marginLeft = 45

      const width = this.elemWidth
      const height = 230 - marginBottom

      const xTicks = width > 720 ? 10 : 5
      const yTicks = 3

      const container = d3.select(this.$refs[`${this.graphName}`] as HTMLElement)

      container
        .select('svg')
        .on('pointerdown', null)
        .on('pointermove', null)
        .on('mouseleave', null)
        .remove()

      const svg = container
        .append('svg')
        .attr('width', width)
        .attr('height', height)
        .style('touch-action', 'none')
        .style('user-select', 'none')
        .style('--webkit-user-select', 'none') /* Safari */
        .style('--ms-user-select', 'none') /* IE 10 and IE 11 */

      // Creating Scales
      const xScale = d3
        .scaleLinear()
        .domain(d3.extent(graphRecords, (d) => d.distance! / 1000) as Number[])
        .range([marginLeft, width - marginRight])

      this.xScale = xScale

      const altitudeScale = d3
        .scaleLinear()
        .domain(d3.extent(graphRecords, (d) => d.altitude) as Number[])
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

      // Add X-Y Axis Label
      svg
        .append('text')
        .attr('class', 'x-axis-label')
        .attr('x', width - marginRight - 70)
        .attr('y', marginTop - 15)
        .style('fill', 'currentColor')
        .style('font-size', '0.9em')
        .text('Dist. (km) →')

      svg
        .append('text')
        .attr('class', 'y-axis-label')
        .attr('x', 0)
        .attr('y', marginTop - 15)
        .style('fill', 'currentColor')
        .style('font-size', '0.9em')
        .text(`↑ ${this.yLabel}`)

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
            .data(xScale.ticks(xTicks))
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

      // Add Altitude Area
      const altitudeArea = d3
        .area<Record>()
        .curve(d3.curveBasis)
        .x((d: Record) => xScale(d.distance! / 1000) as number)
        .y0(altitudeScale(d3.min(altitudeScale.domain())!))
        .y1((d) => altitudeScale(d.altitude ?? 0))

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
        .x((d: Record) => xScale(d.distance! / 1000) as number)
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
          g.append('polygon').attr('points', '0,30 -5,20 5,20').attr('fill', `${this.pointerColor}`)
        })

      // Add Events
      const pointerListener = (e: Event) => {
        if (e.type == 'pointerdown')
          this.hoveredRecordFreeze = !this.hoveredRecordFreeze && this.hoveredRecord != new Record()
        if (this.hoveredRecordFreeze == true && this.hoveredRecord != new Record()) return

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

        pointer.attr('transform', `translate(${px}, 0)`)

        const pointerPercentage = (px - xMin) / (xMax - xMin)
        const lookupIndex = Math.round(pointerPercentage * (this.records.length - 1))

        let nearestRecord: Record = new Record()
        let dx = Number.MAX_VALUE
        if (xScale(this.records[lookupIndex].distance! / 1000) <= px) {
          for (let i = lookupIndex; i < this.records.length; i++) {
            const delta = Math.abs(px - xScale(this.records[i].distance! / 1000)!)
            if (delta > dx) break
            nearestRecord = this.records[i]
            dx = delta
          }
        } else {
          for (let i = lookupIndex; i >= 0; i--) {
            const delta = Math.abs(px - xScale(this.records[i].distance! / 1000)!)
            if (delta > dx) break
            nearestRecord = this.records[i]
            dx = delta
          }
        }

        this.hoveredRecord = nearestRecord
        this.hoveredRecordFreeze = this.hoveredRecordFreeze && e.type == 'pointerdown'
      }

      svg.on('pointerdown', pointerListener)
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
  display: inline-block;
}

.line-graph-hover {
  font-size: 0.9em;
}
.graph-summary {
  margin-left: 45px;
  padding-left: 0px;
  padding-right: 10px;
  margin-right: 15px;
}

.graph-summary div {
  padding-bottom: 1px;
  margin-bottom: 3px;
}
.graph-summary div:nth-child(odd) {
  box-shadow: 0px 0.5px grey;
}

.collapsed > div > h6 > .fa-caret-down,
:not(.collapsed) > div > h6 > .fa-caret-right {
  display: none;
}
</style>
