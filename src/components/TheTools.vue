<script setup lang="ts">
import type { Marker } from '@/spec/activity-service'
import ToolDeviceSelector, { DeviceOption } from './ToolDeviceSelector.vue'
import ToolFieldsRemover from './ToolFieldsRemover.vue'
import ToolModeSelector from './ToolModeSelector.vue'
import ToolSportChanger from './ToolSportChanger.vue'
import ToolTrackpointsConcealer from './ToolTrackpointsConcealer.vue'
import ToolTrackpointsTrimmer from './ToolTrackpointsTrimmer.vue'

import { ActivityFile, Record, Session } from '@/spec/activity'
import {
  EncodeSpecifications,
  FileType,
  Manufacturer,
  Sport,
  ToolMode
} from '@/spec/activity-service'
import { toRaw } from 'vue'
</script>
<template>
  <div class="row m-0 px-2 text-start">
    <label>We have few useful tools to edit your Activity Files.</label>
    <p>Some tools may be disabled depending on these factors:</p>
    <ul class="ps-4 m-0" style="font-size: 0.8em">
      <li>You need at least two activity files to combine them.</li>
      <li>You need at least two sessions (in one or multiple activity files) to split.</li>
    </ul>
    <div class="pt-3">
      <ToolModeSelector
        :activities="activities"
        :sessions="sessions"
        v-on:tool-mode="onToolMode"
      ></ToolModeSelector>
    </div>
    <div class="pt-3">
      <ToolDeviceSelector
        :manufacturers="manufacturers"
        :activities="activities"
        :tool-mode="toolMode"
        v-on:selected-device="onSelectedDevice"
      ></ToolDeviceSelector>
    </div>
    <div class="pt-3">
      <ToolSportChanger
        :sessions="sessions"
        :sports="sports"
        :tool-mode="toolMode"
        v-on:session-sports="onSessionSports"
      ></ToolSportChanger>
    </div>
    <div class="pt-3">
      <ToolTrackpointsTrimmer
        :sessions="sessions"
        :tool-mode="toolMode"
        v-on:markers="onTrimMarkers"
      ></ToolTrackpointsTrimmer>
    </div>
    <div class="pt-3">
      <ToolTrackpointsConcealer
        :sessions="sessions"
        :tool-mode="toolMode"
        v-on:markers="onConcealMarkers"
      ></ToolTrackpointsConcealer>
    </div>
    <div class="pt-3">
      <ToolFieldsRemover
        :sessions="sessions"
        :tool-mode="toolMode"
        v-on:selected-fields="onSelectedFields"
      ></ToolFieldsRemover>
    </div>
    <div class="pt-4">
      <div class="row">
        <div>
          <button class="w-100 btn btn-success" @click="proceed" :disabled="!isValidToProceed">
            Proceed
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    manufacturers: { type: Array<Manufacturer>, required: true },
    sports: { type: Array<Sport>, required: true },
    activities: { type: Array<ActivityFile>, required: true },
    sessions: { type: Array<Session>, required: true },
    combinedRecord: Array<Record>
  },
  data() {
    return {
      toolMode: ToolMode.Unknown,
      selectedDevice: new DeviceOption(),
      sessionSports: new Array<string>(),
      trimMarkers: new Array<Marker>(),
      concealMarkers: new Array<Marker>(),
      selectedFieldRemovers: new Array<string>()
    }
  },
  computed: {
    isValidToProceed(): boolean {
      if (this.toolMode == ToolMode.Unknown) return false
      if (this.selectedDevice == null) return false
      if (this.selectedDevice.productId == null) return false
      return true
    }
  },
  watch: {
    sessions: {
      handler() {
        this.toolMode = ToolMode.Unknown
      }
    }
  },
  methods: {
    onToolMode(value: ToolMode) {
      this.toolMode = value
    },
    onSelectedDevice(value: DeviceOption) {
      this.selectedDevice = value
    },
    onSessionSports(value: string[]) {
      this.sessionSports = value
    },
    onTrimMarkers(markers: Marker[]) {
      this.trimMarkers = this.validateMarkers(markers)
    },
    onConcealMarkers(markers: Marker[]) {
      this.concealMarkers = this.validateMarkers(markers)
    },
    onSelectedFields(value: string[]) {
      this.selectedFieldRemovers = value
    },
    validateMarkers(markers: Marker[]): Marker[] {
      // the value from input range is a string representation of the selected number.
      // we need to convert it to number otherwise, it will be treated as string in json.
      markers.forEach((m) => {
        m.startN = parseInt(m.startN as unknown as string)
        m.endN = parseInt(m.endN as unknown as string)
      })
      return markers
    },
    proceed() {
      if (!this.isValidToProceed) return

      const spec = new EncodeSpecifications({
        toolMode: this.toolMode,
        targetFileType: FileType.FIT, // TODO: implement other type
        manufacturerId: this.selectedDevice.manufacturerId!,
        productId: this.selectedDevice.productId!,
        deviceName: '', // TODO: not required for FIT File Type
        sports: toRaw(this.sessionSports),
        trimMarkers: toRaw(this.trimMarkers),
        concealMarkers: toRaw(this.concealMarkers),
        removeFields: toRaw(this.selectedFieldRemovers)
      })

      this.$emit('encodeSpecifications', spec)
    }
  },
  mounted() {},
  unmounted() {}
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
