<script setup lang="ts">
import { FileType, ToolMode } from '@/spec/activity-service'
import type { PropType } from 'vue'
</script>
<template>
  <div>
    <label>Target File Type</label>
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
          FIT is currently the most advanced file format for storing activity data.
        </p>
        <p v-show="selected.value == FileType.GPX">
          <strong>Note: GPX does not support the 'Power' data field.</strong>
        </p>
        <p v-show="selected.value == FileType.TCX">
          <strong>Note: GPX does not support the 'Temperature' data field.</strong>
        </p>
        <p v-show="selected.value != FileType.Unsupported && selected.value != FileType.FIT">
          If your target platform is Strava, we recommend choosing FIT instead. Strava have shifted
          their support to FIT files. If you choose GPX or TCX, the target device name may not map
          correctly with Strava's device mapping database.
        </p>
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
