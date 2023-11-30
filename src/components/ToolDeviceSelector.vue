<script setup lang="ts">
import { ActivityFile } from '@/spec/activity'
import { Manufacturer, Product, ToolMode } from '@/spec/activity-service'
import type { PropType } from 'vue'
</script>
<template>
  <div>
    <label>Target device</label>
    <div class="col-12 pt-1">
      <v-select
        label="label"
        placeholder="Please select a device"
        :selectable="isSelectable"
        :clearable="false"
        :options="dataSource"
        v-model="selected"
        :disabled="toolMode == ToolMode.Unknown"
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
</template>
<script lang="ts">
export class DeviceOption {
  label: string = ''
  heading?: boolean = false
  manufacturerId?: number | null = null
  manufacturerName: string = ''
  isManualInputProduct?: boolean = false
  productId?: number | null = null
  productName?: string | null = null

  constructor(data?: DeviceOption) {
    this.label = data?.label ?? '-- Please select a device --'
    this.heading = data?.heading ?? false
    this.manufacturerId = data?.manufacturerId ?? null
    this.manufacturerName = data?.manufacturerName ?? ''
    this.isManualInputProduct = data?.isManualInputProduct ?? false
    this.productId = data?.productId ?? null
    this.productName = data?.productName ?? null
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
    manufacturers: { type: Array<Manufacturer>, required: true }
  },

  data() {
    return {
      selected: new DeviceOption()
    }
  },
  computed: {
    deviceFromFile(): DeviceOption {
      let device = new DeviceOption()
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
      return device
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
      if (this.selected.manufacturerId == null) return false
      if (!this.isManufacturerHasProduct) return true
      return this.selected.isManualInputProduct ?? false
    }
  },
  watch: {
    deviceFromFile: {
      handler(value) {
        this.selected = value
      }
    },
    selected: {
      handler(value) {
        this.$emit('selectedDevice', value)
      }
    }
  },
  methods: {
    isNumber(e: KeyboardEvent) {
      if (isNaN(parseInt(e.key))) e.preventDefault()
    },
    isSelectable(option: any) {
      return option.heading != true
    }
  },
  mounted() {
    this.selected = this.deviceFromFile
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
