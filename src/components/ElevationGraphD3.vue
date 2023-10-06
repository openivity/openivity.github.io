<template>
  <div id="chart">
    <div class="elevation">Elevation</div>
    <svg></svg>
    <div id="tooltip">
      <div id="popupElev" class="ol-popup">
        <div class="popup-content">
          <div>
            <span>Time:</span>
            <span>
              {{
                popupRecord.timestamp
                  ? toTimezoneDateString(popupRecord.timestamp, timezoneOffsetHours)
                  : '-'
              }}
            </span>
          </div>
          <div>
            <span>Altitude:</span>
            <span>{{ popupRecord.altitude ? popupRecord.altitude?.toFixed(2) : '-' }}</span>
            <span>&nbsp;masl</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ActivityFile, Record } from '@/spec/activity'
import * as d3 from 'd3'
import { toStringHDMS } from 'ol/coordinate.js'
import { toTimezoneDateString } from '@/toolkit/date'

export default {
  props: {
    activityFile: ActivityFile,
    timezoneOffsetHours: Number
  },
  data() {
    return {
      popupRecord: new Record()
    }
  },
  watch: {
    activityFile: function () {
      const data = this.activityFile?.records || []
      this.renderChart(data)
    }
  },
  methods: {
    toStringHDMS: toStringHDMS,
    toTimezoneDateString: toTimezoneDateString,

    renderChart(data: Array<Record>) {
      const width = 928
      const height = 500
      const marginTop = 20
      const marginRight = 20
      const marginBottom = 30
      const marginLeft = 40

      // Accessors
      // const parseDate = d3.timeParse('%Y-%m-%d %H:%M')
      // const xAccessor = (d) => parseDate(d.timestamp)
      const xAccessor = (d: Record) => {
        return new Date(d.timestamp)
      }
      const yAccessor = (d: Record) => d.altitude

      const max = d3.max(data, (d) => d.altitude) || 0
      const min = d3.min(data, (d) => d.altitude) || 0
      const yMinDom = d3.min([min - min * 2, 0]) || 0
      const yMaxDom = max + max * 2

      // Create the scales.
      const x = d3
        .scaleUtc()
        .domain(d3.extent(data, xAccessor))
        .rangeRound([marginLeft, width - marginRight])

      const y = d3
        .scaleLinear()
        .domain([yMinDom, yMaxDom])
        .nice()
        .rangeRound([height - marginBottom, marginTop])

      // Create the path generator.
      const line = d3
        .area()
        .curve(d3.curveBasis)
        .x((d) => {
          return x(new Date(d.timestamp))
        })
        .y((d) => {
          return y(d.altitude)
        })

      // Create the path generator for the area.
      const area = d3
        .area()
        .curve(d3.curveBasis)
        .x((d) => x(new Date(d.timestamp)))
        .y0(y(yMinDom)) // Set the baseline for the area chart to y=0
        .y1((d) => y(d.altitude))

      // Create the SVG container.
      const svg = d3
        .select(this.$el)
        .select('svg')
        .attr('width', width)
        .attr('height', height)
        .attr('viewBox', [0, 0, width, height])
        .attr('style', 'max-width: 100%; height: auto;')
      // .style('background-color', 'white')

      // Append the axes.
      svg
        .append('g')
        .attr('transform', `translate(0,${height - marginBottom})`)
        .call(
          d3
            .axisBottom(x)
            .ticks(width / 80)
            .tickSizeOuter(0)
        )
        .call((g) => g.select('.domain').remove())

      svg
        .append('g')
        .attr('transform', `translate(${marginLeft},0)`)
        .call(d3.axisLeft(y))
        .call((g) => g.select('.domain').remove())
        .call((g) => g.select('.tick:last-of-type textchart').append('tspan').text(data.y))

      // Create the grid.
      svg
        .append('g')
        .attr('stroke', 'currentColor')
        .attr('stroke-opacity', 0.1)
        .call((g) =>
          g
            .append('g')
            .selectAll('line')
            .data(x.ticks())
            .join('line')
            .attr('x1', (d) => 0.5 + x(d))
            .attr('x2', (d) => 0.5 + x(d))
            .attr('y1', marginTop)
            .attr('y2', height - marginBottom)
        )
        .call((g) =>
          g
            .append('g')
            .selectAll('line')
            .data(y.ticks())
            .join('line')
            .attr('y1', (d) => 0.5 + y(d))
            .attr('y2', (d) => 0.5 + y(d))
            .attr('x1', marginLeft)
            .attr('x2', width - marginRight)
        )

      svg
        .append('path')
        .datum(data)
        .attr('fill', 'none')
        .attr('stroke', 'yellow')
        .attr('stroke-width', 8)
        .attr('stroke-linejoin', 'round')
        .attr('stroke-linecap', 'round')
        .attr('d', line)
        .attr('stroke-opacity', 0.8) // Adjust the opacity as needed

      // Append the area below the line with the same fill color (blue).
      svg
        .append('path')
        .datum(data)
        .attr('fill', '#2A303F') // Set the fill color to blue
        .attr('d', area)

      // Tooltip
      const tooltip = d3.select('#tooltip')
      const tooltipDot = svg
        .append('circle')
        .attr('r', 5)
        .attr('fill', 'none')
        .attr('stroke', '#970000')
        .attr('stroke-width', 2)
        .style('opacity', 0)
        .style('pointer-events', 'none')

      // Hover line.
      var hoverLineGroup = svg.append('g').attr('class', 'hover-line')
      var hoverLine = hoverLineGroup
        .append('line')
        .attr('x1', 10)
        .attr('x2', 10)
        .attr('y1', y(yMinDom))
        .attr('y2', y(yMaxDom))

      // Hide hover line by default.
      hoverLine.style('opacity', 1e-6)

      svg
        .append('rect')
        .attr('width', width)
        .attr('height', height)
        .style('opacity', 0)
        .on('touchmouse mousemove', (event) => {
          const mousePos = d3.pointer(event, this)
          // x coordinate stored in mousePos index 0
          const date = x.invert(mousePos[0])

          // Custom Bisector - left, center, right
          const dateBisector = d3.bisector(xAccessor).left
          const bisectionIndex = dateBisector(data, date)
          // math.max prevents negative index reference error
          const hoveredIndexData = data[Math.max(0, bisectionIndex - 1)]
          if (!hoveredIndexData) return

          this.popupRecord = hoveredIndexData

          // Update Image
          tooltipDot
            .style('opacity', 1)
            .attr('cx', x(xAccessor(hoveredIndexData)))
            .attr('cy', y(yAccessor(hoveredIndexData)))
            .raise()

          tooltip
            .style('display', 'block')
            .style('left', `${event.pageX}px`)
            .style('top', `${event.pageY}px`)
          // .style('left', `${y(xAccessor(hoveredIndexData))}px`)
          // .style('top', `${x(yAccessor(hoveredIndexData)) - 50}px`)

          // tooltip.select('.altitude').text(`${yAccessor(hoveredIndexData)}`)
          // tooltip.select('.timestamp').text(`${xAccessor(hoveredIndexData)}`)

          hoverLine
            .attr('x1', x(xAccessor(hoveredIndexData)))
            .attr('x2', x(xAccessor(hoveredIndexData)))
            .style('opacity', 1)
        })
        .on('mouseleave', function () {
          tooltipDot.style('opacity', 0)
          tooltip.style('display', 'none')

          hoverLine.style('opacity', 1e-6)
        })

      return Object.assign(svg.node())
    }
  },
  mounted() {
    this.renderChart(this.activityFile?.records || [])
  }
}
</script>

<style>
#tooltip {
  border: 1px solid #ccc;
  position: absolute;
  padding: 10px;
  background-color: #fff;
  display: none;
  pointer-events: none;
}

.hover-line {
  stroke: #970000;
  fill: none;
  stroke-width: 2px;
}
</style>
