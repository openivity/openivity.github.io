<script setup lang="ts">
import { FileType, ToolMode } from '@/spec/activity-service'
import type { PropType } from 'vue'
</script>
<template>
  <div>
    <label class="pe-1">Target File Type</label>
    <div class="col-12 pt-1">
      <v-select
        label="label"
        placeholder="Please select file type"
        :clearable="false"
        :searchable="false"
        :options="dataSource"
        v-model="selected"
        :disabled="toolMode == ToolMode.Unknown"
      >
      </v-select>
      <div class="pt-1">
        <p v-show="selected.value == FileType.FIT">
          FIT is currently the most advanced file format for storing activity data developed by
          Garmin. We strives to comply with the FIT Activity File (FIT_FILE_TYPE = 4) as defined by
          <a href="https://developer.garmin.com/fit" target="_blank" rel="noopener noreferrer"
            >Garmin FIT</a
          >
          that being implemented by
          <a href="https://github.com/muktihari/fit" target="_blank" rel="noopener noreferrer"
            >FIT SDK for Go</a
          >
          .
        </p>
        <div v-show="selected.value == FileType.GPX">
          <p>
            GPX is a widely used XML format for geospacial data developed by Topografix. We follow
            <a
              href="https://www.topografix.com/gpx/1/1/gpx.xsd"
              target="_blank"
              rel="noopener noreferrer"
              >Schema V1.1</a
            >
            with
            <a
              href="https://www8.garmin.com/xmlschemas/TrackPointExtensionv2.xsd"
              target="_blank"
              rel="noopener noreferrer"
              >Garmin Trackpoint Extension V2</a
            >
            which does not support the <strong>Power</strong> data field.
          </p>
        </div>
        <div class="pt-1" v-show="selected.value == FileType.TCX">
          <p>
            The advantage of TCX is its ability to include <strong>Power</strong> information.
            However, it does not support the <strong>Temperature</strong> data field. We follow
            <a
              href="https://www8.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd"
              target="_blank"
              rel="noopener noreferrer"
            >
              Schema V2
            </a>
            with
            <a
              href="https://www8.garmin.com/xmlschemas/ActivityExtensionv2.xsd"
              target="_blank"
              rel="noopener noreferrer"
            >
              Garmin Activity Extension V2
            </a>
            .
          </p>
        </div>
        <div
          class="pt-2"
          v-show="selected.value != FileType.Unsupported && selected.value != FileType.FIT"
        >
          <p>
            Note: If your target platform is Strava, we recommend choosing FIT instead. Strava have
            shifted their support to FIT files. If you choose GPX or TCX, the target device name may
            not map correctly with Strava's device mapping database.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
export class FileTypeOption {
  label: string = '-- Please select file type -- '
  value: FileType = FileType.Unsupported
}

export default {
  props: {
    toolMode: { type: Number as PropType<ToolMode>, required: true }
  },

  data() {
    return {
      selected: new FileTypeOption()
    }
  },
  computed: {
    dataSource(): FileTypeOption[] {
      const dataSource: FileTypeOption[] = [
        { label: 'FIT - Flexible and Interoperable Data Transfer', value: FileType.FIT },
        { label: 'GPX - GPS Exchange Format', value: FileType.GPX },
        { label: 'TCX - Training Center XML', value: FileType.TCX }
      ]
      return dataSource
    }
  },
  watch: {
    toolMode: {
      handler() {
        this.selected = new FileTypeOption()
      }
    },
    selected: {
      handler(value: FileTypeOption) {
        this.$emit('selectedFileType', value)
      },
      deep: true
    }
  },
  methods: {},
  mounted() {
    this.selected = new FileTypeOption()
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
