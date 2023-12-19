<script setup lang="ts">
import { ActivityFile } from '@/spec/activity'
import { FileType, Manufacturer, Product, ToolMode } from '@/spec/activity-service'
import type { PropType } from 'vue'
</script>
<template>
  <div>
    <label>Target Device</label>
    <div v-if="selectedFileType == FileType.FIT">
      <!-- FIT-->
      <div class="col-12 pt-1">
        <v-select
          label="label"
          placeholder="Please select a device"
          :selectable="isSelectable"
          :clearable="false"
          :options="dataSource"
          v-model="selected"
        >
          <template #option="{ label, heading }">
            <strong v-if="heading">{{ label }}</strong>
            <span class="ps-3" v-else>{{ label }}</span>
          </template>
        </v-select>
      </div>
      <div class="col-12">
        <div class="row pt-2" v-show="isManualInputProduct">
          <p v-show="!isManufacturerHasProduct">
            Currently, we don't have a device mapping for
            <strong>"{{ selected.manufacturerName }}"</strong>. We require a valid
            <i>product_id</i> to generate a FIT file for the targeted device. Without a valid
            <i>product_id</i>, the resulting FIT file may not be recognized by other platforms.
          </p>
          <div class="col">
            <label for="product" class="form-label sub-label"
              >Product ID <span class="d-inline color-mandatory fs-6">*</span></label
            >
            <input
              id="product"
              type="number"
              step="1"
              class="form-control form-control-sm"
              placeholder="-- Please input product_id --"
              title="product_id"
              aria-label="product_id"
              aria-describedby="product_id"
              v-model="selected.productId"
              @keypress="isNumber($event)"
            />
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <!-- GPX and TCX -->
      <div class="col-12 pt-1">
        <input
          class="form-control form-control-sm"
          v-model="deviceName"
          placeholder="-- Please input device name --"
          :disabled="selectedFileType == FileType.Unsupported"
        />
      </div>
    </div>
  </div>
</template>
<script lang="ts">
export class DeviceOption {
  label: string = ''
  heading?: boolean = false
  manufacturerId?: number | undefined = undefined
  manufacturerName?: string = ''
  isManualInputProduct?: boolean = false
  productId?: number | undefined = undefined
  productName?: string | undefined = undefined

  constructor(data?: DeviceOption) {
    this.label = data?.label ?? '-- Please select a device --'
    this.heading = data?.heading ?? false
    this.manufacturerId = data?.manufacturerId ?? undefined
    this.manufacturerName = data?.manufacturerName ?? ''
    this.isManualInputProduct = data?.isManualInputProduct ?? false
    this.productId = data?.productId ?? undefined
    this.productName = data?.productName ?? undefined
  }
}

class ManufacturerMap {
  manufacturer: Manufacturer = new Manufacturer()
  productMap: Map<number, Product> = new Map()

  constructor(data?: ManufacturerMap) {
    this.manufacturer = data?.manufacturer ?? new Manufacturer()
    this.productMap = data?.productMap ?? new Map()
  }
}

export default {
  props: {
    toolMode: { type: Number as PropType<ToolMode>, required: true },
    activities: { type: Array<ActivityFile>, required: true },
    manufacturers: { type: Array<Manufacturer>, required: true },
    selectedFileType: { type: Number as PropType<FileType>, required: true }
  },

  data() {
    return {
      selected: new DeviceOption(),
      deviceName: '' // gpx tcx only
    }
  },
  computed: {
    deviceFromFitFile(): DeviceOption {
      for (let i = 0; i < this.activities.length; i++) {
        const act = this.activities[i]

        if (act.creator.manufacturer != null && act.creator.product != null) {
          const m = this.manufacturerMap.get(act.creator.manufacturer!)
          if (m == undefined) continue

          const p = m.productMap.get(act.creator.product)
          if (p == undefined) continue

          return new DeviceOption({
            label: `${m.manufacturer.name} ${p.name}`,
            manufacturerId: m.manufacturer.id,
            manufacturerName: m.manufacturer.name,
            productId: p.id,
            productName: p.name
          })
        }
      }
      return new DeviceOption()
    },
    deviceNameFromGpxTcxFile(): string {
      for (let i = 0; i < this.activities.length; i++) {
        const act = this.activities[i]
        const device = this.deviceMappingForGpxTcx.get(act.creator.name.toLocaleLowerCase())
        if (device != undefined) {
          return device.label
        } else {
          return act.creator.name
        }
      }
      return ''
    },
    dataSource(): DeviceOption[] {
      const dataSource = new Array<DeviceOption>()
      this.manufacturers.forEach((m) => {
        dataSource.push(
          new DeviceOption({
            label: m.name,
            heading: true,
            manufacturerId: m.id,
            manufacturerName: m.name
          })
        )
        m.products.forEach((p) =>
          dataSource.push(
            new DeviceOption({
              label: `${m.name} ${p.name}`,
              manufacturerId: m.id,
              manufacturerName: m.name,
              productId: p.id,
              productName: p.name
            })
          )
        )
        dataSource.push(
          new DeviceOption({
            label: `${m.name} <manual input>`,
            manufacturerId: m.id,
            manufacturerName: m.name,
            isManualInputProduct: true
          })
        )
      })
      return dataSource
    },
    manufacturerMap(): Map<number, ManufacturerMap> {
      const map = new Map()
      this.manufacturers.forEach((m) => {
        const pmap = new Map()
        m.products.forEach((p) => {
          pmap.set(p.id, p)
        })
        map.set(m.id, new ManufacturerMap({ manufacturer: m, productMap: pmap }))
      })
      return map
    },
    isManufacturerHasProduct(): boolean {
      const m = this.manufacturerMap.get(this.selected.manufacturerId!)
      if (m == undefined) return true
      return m?.productMap.size > 0 ?? false
    },
    isManualInputProduct(): boolean {
      if (this.selected.manufacturerId == undefined) return false
      if (!this.isManufacturerHasProduct) return true
      return this.selected.isManualInputProduct ?? false
    },
    deviceMappingForGpxTcx(): Map<string, DeviceOption> {
      const map = new Map<string, DeviceOption>()
      this.manufacturers.forEach((m) => {
        m.products.forEach((p) => {
          map.set(
            `${m.name} ${p.name}`.toLowerCase(),
            new DeviceOption({
              manufacturerId: m.id,
              manufacturerName: m.name,
              productId: p.id,
              productName: p.name,
              label: `${m.name} ${p.name}`
            })
          )
        })
      })
      return map
    }
  },
  watch: {
    deviceFromFitFile: {
      handler(device: DeviceOption) {
        this.selected = device
      }
    },
    deviceNameFromGpxTcxFile: {
      handler(deviceName: string) {
        this.deviceName = deviceName
      }
    },
    selected: {
      handler(device: DeviceOption) {
        this.updateDeviceNameBySelected(device)
        if (isNaN(parseInt(device?.productId as unknown as string))) device.productId = undefined
        this.$emit('selectedDevice', device)
      }
    },
    deviceName: {
      handler(name: string) {
        this.updateSelectedByDeviceName(name)
      }
    }
  },
  methods: {
    isNumber(e: KeyboardEvent) {
      if (isNaN(parseInt(e.key))) e.preventDefault()
    },
    isSelectable(option: any) {
      return option.heading != true
    },
    updateDeviceNameBySelected(device: DeviceOption) {
      if (device.manufacturerId != undefined) this.deviceName = device.label
    },
    updateSelectedByDeviceName(deviceName: string) {
      const device = this.deviceMappingForGpxTcx.get(deviceName.toLocaleLowerCase())
      if (device != undefined) {
        this.selected = device
      }
    }
  },
  mounted() {
    this.selected = this.deviceFromFitFile
    this.deviceName = this.deviceNameFromGpxTcxFile
    this.updateSelectedByDeviceName(this.deviceName)
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
