<template>
  <div class="col-12 h-100 pt-2">
    <div class="row">
      <div class="col">
        <h6 class="pt-1" style="text-align: left">
          Elevation
          <i class="fa-solid fa-solid fa-mountain"></i>
        </h6>
      </div>
      <div class="col text-md-end text-sm-center elevation-hover">
        <span class="text-start">
          <i class="fa-solid fa-hourglass-half"></i>
          {{
            recordView.timestamp
              ? formatDuration(new Date(recordView.timestamp!).getTime() - begin.getTime())
              : '0 s'
          }}
        </span>
        <span class="text-end">
          <i class="fa-solid fa-road"></i>
          {{
            typeof recordView.distance === 'number'
              ? (recordView.distance! / 1000).toFixed(2)
              : '0.00'
          }}
          km
        </span>
        <span class="text-end">
          <i class="fa-solid fa-mountain"></i>
          {{ recordView.altitude?.toFixed(0) ?? '0' }} m
        </span>
        <span class="text-end">
          <i class="fa-solid fa-arrow-up-right-dots fa-sm"></i>
          {{ recordView.grade ? Math.round(recordView.grade) : '0' }} %
        </span>
      </div>
    </div>
    <div :ref="name"></div>
  </div>
</template>
<script lang="ts">
import { Record } from '@/spec/activity'
import { Summary } from '@/spec/summary'
import { formatMillis as formatDuration } from '@/toolkit/date'
import * as d3 from 'd3'

export default {
  props: {
    name: {
      type: String,
      required: true
    },
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
    summary: Summary,
    receivedRecord: Record
  },
  data() {
    return {
      begin: new Date(),
      elapsed: 0, // ms
      hoveredRecord: new Record(),
      recordView: new Record(),
      previousWindowWidth: 0,
      xScale: d3.scaleLinear(),
      color: d3
        .scaleSequential()
        .domain([-45, 45])
        .interpolator(
          d3.interpolateRgbBasis([
            'darkgreen',
            'darkgreen',
            'darkgreen',
            'darkgreen',
            'darkgreen',
            'limegreen',
            'limegreen',
            'lemonchiffon',
            'lemonchiffon',
            'red',
            'red',
            'darkred',
            'darkred',
            'darkred',
            'darkred',
            'darkred'
          ])
        )
    }
  },
  watch: {
    graphRecords: {
      async handler() {
        this.$nextTick(() => {
          requestAnimationFrame(() => this.renderGraph())
        })
      }
    },
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
      handler(record: Record) {
        this.recordView = record
        this.$emit('hoveredRecord', record)
      }
    }
  },
  methods: {
    formatDuration: formatDuration,
    renderGraph() {
      const graphRecords = this.graphRecords
      if (graphRecords.length == 0) return

      this.hoveredRecord = new Record()
      this.begin = d3.min(graphRecords.map((d) => new Date(d.timestamp!)))!

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

      const yScale = d3
        .scaleLinear()
        .domain(d3.extent(graphRecords, (d) => d.altitude) as Number[])
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

      //   svg
      //     .append('text')
      //     .attr('class', 'y-axis-label')
      //     .attr('x', 0)
      //     .attr('y', marginTop - 15)
      //     .style('fill', 'currentColor')
      //     .style('font-size', '0.9em')
      //     .text('↑ Alt. (m)')

      // Add Summary
      svg
        .append('foreignObject')
        .attr('x', 0)
        .attr('y', -5)
        .attr('width', width)
        .attr('height', 100)
        .attr('class', 'text-start')
        .style('font-size', '0.9em')
        .style('width', '100%').html(`
            <span>↑ Alt. (m)&nbsp;</span>
            <span>
                <i class="fa-solid fa-solid fa-mountain"></i>
                ${this.summary?.maxAltitude?.toFixed(0)} m&nbsp;
            </span>
            <span>
                <i class="fa-solid fa-arrow-trend-up"></i>
                ${this.summary?.totalAscent?.toFixed(0)} m&nbsp;
            </span>
            <span>
                <i class="fa-solid fa-arrow-trend-down"></i>
                ${this.summary?.totalDescent?.toFixed(0)} m&nbsp;
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

      // Add Area
      const area = d3
        .area<Record>()
        .curve(d3.curveBasis)
        .x((d: Record) => xScale(d.distance! / 1000) as number)
        .y0(yScale(d3.min(yScale.domain())!))
        .y1((d) => yScale(d.altitude ?? 0))

      svg
        .append('g')
        .append('path')
        .datum(graphRecords)
        .transition()
        .attr('fill', '#222222')
        .style('opacity', 1)
        .attr('d', area)

      const linearGradientId = `linearGradient-${this.name}`

      svg
        .append('linearGradient')
        .attr('id', linearGradientId)
        .attr('gradientUnits', 'userSpaceOnUse')
        .attr('x1', 0)
        .attr('x2', width)
        .selectAll('stop')
        .data(graphRecords)
        .join('stop')
        .attr('offset', (d) => xScale(d.distance! / 1000) / width)
        .attr('stop-color', (d) => this.color(d.grade))

      // Add Line
      const line = d3
        .line<Record>()
        .curve(d3.curveBasis)
        .x((d: Record) => xScale(d.distance! / 1000) as number)
        .y((d) => yScale(d.altitude ?? 0))

      svg
        .append('g')
        .append('path')
        .datum(graphRecords)
        .transition()
        .attr('fill', 'none')
        .attr('stroke', `url(#${linearGradientId})`)
        .attr('stroke-width', 3)
        .attr('d', line)

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
            .attr('stroke', '#78e08f')
            .attr('stroke-width', 1.5)
        })
        .call((g) => {
          g.append('polygon').attr('points', '0,30 -5,20 5,20').attr('fill', '#78e08f')
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
        // let counter = 0 // TODO: remove after debug (currently it's the fastest lookup)
        if (xScale(this.records[lookupIndex].distance! / 1000) <= px) {
          for (let i = lookupIndex; i < this.records.length; i++) {
            const delta = Math.abs(px - xScale(this.records[i].distance! / 1000)!)
            if (delta > dx) break
            nearestRecord = this.records[i]
            dx = delta
            // counter++
          }
          //   console.debug(`look forward for ${counter} records`)
        } else {
          for (let i = lookupIndex; i >= 0; i--) {
            const delta = Math.abs(px - xScale(this.records[i].distance! / 1000)!)
            if (delta > dx) break
            nearestRecord = this.records[i]
            dx = delta
            // counter++
          }
          //   console.debug(`look backward for ${counter} records`)
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
.elevation-hover {
  color: var(--color-text);
  flex: 1 0 0%;
  font-size: 0.9em;
  padding-bottom: 10px;
}

.elevation-hover > span {
  padding-left: 5px;
  padding-right: 5px;
  display: inline-block;
}

@media (pointer: coarse) {
  /* mobile device */

  .elevation-hover {
    display: flex;
    flex: unset;
  }

  .elevation-hover > span:nth-child(1) {
    width: 30% !important;
  }

  .elevation-hover > span:nth-child(4) {
    width: 20% !important;
  }

  .elevation-hover > span {
    padding-left: unset;
    width: 25%;
  }
}
</style>
