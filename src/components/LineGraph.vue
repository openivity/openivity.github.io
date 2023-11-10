<template>
  <div class="col-12 h-100">
    <div class="line-graph-hover row">
      <div class="col text-start">
        <h6 style="text-align: left; display: inline-block">
          {{ name }}
          <i :class="['fa-solid', icon]"></i>
        </h6>
      </div>
      <div class="col text-end">
        <span>
          {{ scaleUnit(recordView[selectedField as keyof Record] as number)?.toFixed(2) ?? '0.00' }}
          {{ unit }}&nbsp;
        </span>
      </div>
    </div>
    <div :ref="name"></div>
  </div>
</template>
<script lang="ts">
import { Record } from '@/spec/activity'
import { Summary } from '@/spec/summary'
import * as d3 from 'd3'

export default {
  props: {
    name: {
      type: String,
      required: true
    },
    icon: String,
    selectedField: { type: String, required: true },
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
    summary: Summary,
    receivedRecord: Record
  },
  data() {
    return {
      previousWindowWidth: 0,
      hoveredRecord: new Record(),
      recordView: new Record(),
      xScale: d3.scaleLinear()
    }
  },
  watch: {
    receivedRecord: {
      handler(record: Record) {
        this.recordView = record
        const pointer = d3.select(this.$refs[`${this.name}`] as HTMLElement).select('g#pointer')
        if (JSON.stringify(record) == JSON.stringify(new Record())) {
          pointer.style('opacity', 0)
          return
        }
        pointer.style('opacity', 1)
        pointer.attr('transform', `translate(${this.xScale(record.distance! / 1000)}, 0)`)
      }
    },
    hoveredRecord: {
      async handler(record: Record) {
        this.recordView = record
        this.$emit('hoveredRecord', record)
      }
    },
    graphRecords: {
      async handler() {
        this.$nextTick(() => {
          requestAnimationFrame(() => this.renderGraph())
        })
      }
    }
  },
  computed: {
    avg(): string {
      switch (this.selectedField) {
        case 'speed':
          return (this.scaleUnit(this.summary?.avgSpeed!) ?? 0).toFixed(2)
        case 'cadence':
          return (this.summary?.avgCadence ?? 0).toFixed(0)
        case 'heartRate':
          return (this.summary?.avgHeartRate ?? 0).toFixed(0)
        case 'temperature':
          return (this.summary?.avgTemperature ?? 0).toFixed(0)
        case 'power':
          return (this.summary?.avgPower ?? 0).toFixed(0)
        default:
          return '0'
      }
    },
    max(): string {
      switch (this.selectedField) {
        case 'speed':
          return (this.scaleUnit(this.summary?.maxSpeed!) ?? 0).toFixed(2)
        case 'cadence':
          return (this.summary?.maxCadence ?? 0).toFixed(0)
        case 'heartRate':
          return (this.summary?.maxHeartRate ?? 0).toFixed(0)
        case 'temperature':
          return (this.summary?.maxTemperature ?? 0).toFixed(0)
        case 'power':
          return (this.summary?.maxPower ?? 0).toFixed(0)
        default:
          return '0'
      }
    }
  },
  methods: {
    renderGraph() {
      const graphRecords = this.graphRecords
      if (graphRecords.length == 0) return

      const marginTop = 25
      const marginRight = 5
      const marginBottom = 35
      const marginLeft = 35

      const $elem = this.$el as Element
      const width = $elem.clientWidth
      const height = $elem.clientHeight - marginBottom

      const xTicks = width > 720 ? 10 : 5
      const yTicks = 3

      const container = d3.select(this.$refs[`${this.name}`] as HTMLElement)

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

      const yScale = d3
        .scaleLinear()
        .domain(
          d3.extent(
            graphRecords,
            (d) => this.scaleUnit(d[this.selectedField as keyof Record] as number) ?? 0
          ) as Number[]
        )
        .rangeRound([height - marginBottom, marginTop])
        .nice()

      // Add X & Y Axis
      svg
        .append('g')
        .style('font-size', '0.9em')
        .attr('transform', `translate(0,${height - marginBottom})`)
        .call(d3.axisBottom(xScale).ticks(xTicks))

      svg
        .append('g')
        .style('font-size', '0.9em')
        .attr('transform', `translate(${marginLeft},0)`)
        .call(d3.axisLeft(yScale).ticks(yTicks).tickFormat(d3.format('~s')))

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
        .append('foreignObject')
        .attr('x', 0)
        .attr('y', -5)
        .attr('width', width)
        .attr('height', 100)
        .attr('class', 'text-start')
        .style('font-size', '0.9em')
        .style('width', '100%').html(`
            <span>↑ ${this.yLabel}&nbsp;</span>
            <span>
                Avg: ${this.avg} ${this.unit}&nbsp;
            </span>
            <span>
                Max: ${this.max} ${this.unit}&nbsp;
            </span>
            `)

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
        .y1((d) => yScale(this.scaleUnit(d[this.selectedField as keyof Record] as number) ?? 0))

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
        const lookupIndex = Math.round(pointerPercentage * this.records.length - 1)

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
      }

      svg.on('pointerdown', pointerListener)
      svg.on('pointermove', pointerListener, { passive: true })
      svg.on('mouseleave', () => {
        this.hoveredRecord = new Record()
        pointer.style('opacity', 0)
      })
    },
    scaleUnit(value: number): number | null {
      if (this.selectedField === 'speed') {
        value = (value * 3600) / 1000
      }
      if (isNaN(value)) return null
      return value
    },
    onResize() {
      // Ensure DOM is fully re-rendered after resize
      this.$nextTick(() => {
        // Prevent re-render graph when width is not changing
        if (window.innerWidth == this.previousWindowWidth) return
        this.previousWindowWidth = window.innerWidth
        this.renderGraph()
      })
    }
  },
  mounted() {
    this.$nextTick(() => (this.previousWindowWidth = window.innerHeight))
    window.addEventListener('resize', this.onResize)
  },
  unmounted() {
    window.removeEventListener('resize', this.onResize)
  }
}
</script>
<style>
.line-graph-hover {
  font-size: 0.9em;
}
</style>
