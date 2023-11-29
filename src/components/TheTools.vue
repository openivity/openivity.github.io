<template>
  <div class="row m-0 px-3 text-start">
    <label>We have few useful tools to edit your Activity Files.</label>
    <p>
      Some tools may be disabled depending on the number of activities or sessions you have open.
      Again, the entire process is performed locally, so there's no need to worry about privacy.
      <span class="emoji">😉</span>
    </p>
    <div class="pt-3">
      <label>Please select a tool</label>
      <div class="col-12 pt-1">
        <select
          class="form-select form-select-sm"
          name="encodeMode"
          id="encodeMode"
          v-model="encodeMode"
        >
          <option :value="EncodeMode.Unknown" selected disabled>-- Select a tool --</option>
          <option :value="EncodeMode.Edit">Edit Relevant Data</option>
          <option
            :value="EncodeMode.Combine"
            :disabled="activities.length < 2"
            :title="
              activities.length < 2
                ? 'You have only one activity opened, please open multiple activites to be able to use this feature.'
                : ''
            "
          >
            Combine Multiple Activities into One
          </option>
          <option
            :value="EncodeMode.SplitPerSession"
            :disabled="sessions.length < 2"
            :title="
              sessions.length < 2
                ? 'You have only one session in the opened activity, please open multiple activities or open an activity that have multiple sessions to be able to use this feature.'
                : ''
            "
          >
            Split Activities Per Session
          </option>
        </select>
        <div class="pt-1">
          <p v-show="encodeMode == EncodeMode.Edit">
            We will edit relevant data for every input activities. This changes will apply to your
            entire activities, like a Bulk Edit, if you want to edit one activity, please open only
            one at a time.
          </p>
          <p v-show="encodeMode == EncodeMode.Combine">
            We will combine multiple activities into one continuous activity file. Two sequential
            sessions of the same sport will be merged into one session. If the sport is different,
            it will be placed in separate sessions. This process will continue until all sessions
            are combined.
          </p>
          <p v-show="encodeMode == EncodeMode.SplitPerSession">
            We will create new Activity File for every Sessions in all activities.
          </p>
        </div>
      </div>
    </div>
    <div class="pt-3">
      <label>Target device</label>
      <div class="col-12 pt-1">
        <select
          class="form-select form-select-sm"
          name="manufacturers"
          id="manufacturers"
          v-model="selectedManufacturerProduct"
          :disabled="encodeMode == EncodeMode.Unknown"
        >
          <option disabled selected :value="[]">-- Select a device --</option>
          <optgroup
            v-for="(manufacturer, mIndex) in manufacturers"
            v-bind:key="mIndex"
            :label="manufacturer.name"
          >
            <option
              v-for="(product, pIndex) in manufacturer.products"
              v-bind:key="pIndex"
              :value="[manufacturer.id, manufacturer.name, product.id]"
            >
              {{ manufacturer.name }} {{ product.name }}
            </option>
            <option :value="[manufacturer.id, manufacturer.name, NaN]">
              {{ manufacturer.name }} &lt;manual input product_id&gt;
            </option>
          </optgroup>
        </select>
      </div>
      <div class="col-12">
        <div class="row pt-2" v-show="isManualInputProduct">
          <p v-show="!isManufacturerHasProduct">
            Currently, we don't have a device mapping for <strong>"{{ manufacturerName }}"</strong>.
            We require a valid <i>product_id</i> to generate a FIT file for the targeted device.
            Without a valid <i>product_id</i>, the resulting FIT file may not be recognized by other
            platforms.
          </p>
          <div class="col">
            <label for="product" class="form-label sub-label">Product ID</label>
            <input
              id="product"
              type="number"
              step="1"
              class="form-control form-control-sm"
              placeholder="-- Please input product_id --"
              title="product_id"
              aria-label="product_id"
              aria-describedby="product_id"
              v-model="productId"
              @keypress="isNumber($event)"
            />
          </div>
        </div>
      </div>
    </div>
    <div class="pt-3">
      <label>Change Sport</label>
      <div class="col-12 pb-1" v-for="(ses, index) in sessions" :key="index">
        <label class="form-label sub-label">Session {{ index + 1 }}'s sport:</label>
        <select
          class="form-select form-select-sm"
          name="sports"
          id="sports"
          v-model="sessionSports[index]"
          :disabled="encodeMode == EncodeMode.Unknown"
        >
          <option disabled selected :value="null">-- Please select a sport --</option>
          <option v-for="(sport, index) in sports" :key="index" :value="sport.name">
            {{ sport.name }}
          </option>
        </select>
      </div>
    </div>
    <div class="pt-3">
      <div
        class="row m-0"
        style="cursor: pointer"
        data-bs-toggle="collapse"
        data-bs-target="#trimTarget"
        aria-expanded="false"
        aria-controls="trimTarget"
      >
        <div class="text-start p-0">
          <label class="pe-1">Trim Trackpoints</label>
          <i class="fa-regular fa-circle-question" title="Show or Hide Help Text"></i>
        </div>
      </div>
      <div class="collapse show" id="trimTarget">
        <p>
          Select this option if you want to trim certain trackpoints in your activities. For
          instance, if you finish cycling and transition to another mode of transportation without
          turning off your cyclocomputer, this feature allows you to remove the unwanted
          trackpoints.
        </p>
      </div>
      <div class="col-12 pt-2">
        <select
          class="form-select form-select-sm"
          name="trim"
          id="trim"
          v-model="isTrimSelected"
          :disabled="encodeMode == EncodeMode.Unknown"
        >
          <option :value="false">Do Not Trim Trackpoints</option>
          <option :value="true">Trim Trackpoints</option>
        </select>
      </div>
      <div v-if="isTrimSelected">
        <div class="pt-2" v-for="(marker, index) in trimMarkers" :key="index">
          <div class="row px-1">
            <label class="sub-label">Session {{ index + 1 }}: {{ sessions[index].sport }} </label>
            <div class="col text-start"><p>Distance from the Start</p></div>
            <div class="col text-end">
              {{ ((sessions[index].records[marker.startN].distance ?? 0) / 1000).toFixed(2) }}
              km
            </div>
            <input
              class="form-range"
              type="range"
              :min="0"
              :max="sessions[index].records.length - 1"
              v-model="marker.startN"
            />
          </div>
          <div class="row px-1">
            <div class="col text-start"><p>Distance from the End</p></div>
            <div class="col text-end">
              {{
                (
                  ((sessions[index].records[sessions[index].records.length - 1].distance ?? 0) -
                    (sessions[index].records[marker.endN].distance ?? 0)) /
                  1000
                ).toFixed(2)
              }}
              km
            </div>
            <input
              class="form-range"
              type="range"
              :min="0"
              :max="sessions[index].records.length - 1"
              v-model="marker.endN"
            />
          </div>
        </div>
      </div>
    </div>
    <div class="pt-3">
      <div
        class="row m-0"
        style="cursor: pointer"
        data-bs-toggle="collapse"
        data-bs-target="#concealTarget"
        aria-expanded="false"
        aria-controls="concealTarget"
      >
        <div class="text-start p-0">
          <label class="pe-1">Conceal GPS Positions</label>
          <i class="fa-regular fa-circle-question" title="Show or Hide Help Text"></i>
        </div>
      </div>
      <div class="collapse show" id="concealTarget">
        <p>
          Select this option if you wish to conceal your GPS positions at a specific start and end
          distance while maintaining other trackpoints data. It is useful when you want to share the
          activity on social media, like Strava, and you want to avoid revealing the exact location
          of your home for security reasons and you wouldn't want your information stored on any
          platform's server.
        </p>
      </div>
      <div class="col-12 pt-2">
        <div>
          <select
            class="form-select form-select-sm"
            name="conceal"
            id="conceal"
            v-model="isConcealSelected"
            :disabled="encodeMode == EncodeMode.Unknown"
          >
            <option :value="false">Do Not Conceal My GPS Positions</option>
            <option :value="true">Conceal My GPS positions</option>
          </select>
        </div>
        <div v-if="isConcealSelected">
          <div class="pt-2" v-for="(marker, index) in concealMarkers" :key="index">
            <div class="row px-1">
              <label class="sub-label">Session {{ index + 1 }}: {{ sessions[index].sport }} </label>
              <div class="col text-start"><p>Distance from the Start</p></div>
              <div class="col text-end">
                {{ ((sessions[index].records[marker.startN].distance ?? 0) / 1000).toFixed(2) }}
                km
              </div>
              <input
                class="form-range"
                type="range"
                :min="0"
                :max="sessions[index].records.length - 1"
                v-model="marker.startN"
              />
            </div>
            <div class="row px-1">
              <div class="col text-start"><p>Distance from the End</p></div>
              <div class="col text-end">
                {{
                  (
                    (sessions[index].records[sessions[index].records.length - 1].distance! -
                      sessions[index].records[marker.endN].distance!) /
                    1000
                  ).toFixed(2)
                }}
                km
              </div>
              <input
                class="form-range"
                type="range"
                :min="0"
                :max="sessions[index].records.length - 1"
                v-model="marker.endN"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="pt-3">
      <div
        class="row m-0"
        style="cursor: pointer"
        data-bs-toggle="collapse"
        data-bs-target="#fieldsRemovalTarget"
        aria-expanded="false"
        aria-controls="fieldsRemovalTarget"
      >
        <div class="text-start p-0">
          <label class="pe-1">Remove Fields</label>
          <i class="fa-regular fa-circle-question" title="Show or Hide Help Text"></i>
        </div>
      </div>
      <div class="collapse show" id="fieldsRemovalTarget">
        <p>Select any field you wish to remove from the entire trackpoints.</p>
      </div>
      <div v-for="(item, index) in fieldsRemoverList" :key="index">
        <div class="form-check" v-show="showFieldRemover(item.value)">
          <input
            class="form-check-input"
            type="checkbox"
            :id="item.value"
            :value="item.value"
            :disabled="encodeMode == EncodeMode.Unknown"
            v-model="removeFields"
          />
          <label class="form-check-label" style="color: var(--color-text)" :for="item.value">
            {{ item.label }}
          </label>
        </div>
      </div>
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
import { ActivityFile, Record, Session } from '@/spec/activity'
import {
  EncodeMode,
  EncodeSpecifications,
  FileType,
  Manufacturer,
  Marker,
  Sport
} from '@/spec/activity-service'
import { shallowRef } from 'vue'

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
      manufacturerMap: shallowRef(new Map<number, Manufacturer>()),
      sportMap: shallowRef(new Map<string, Sport>()),
      EncodeMode: EncodeMode,
      encodeMode: EncodeMode.Unknown,
      selectedManufacturerProduct: null as unknown as any[],
      manufacturerName: '',
      manufacturerId: null as unknown as number | null,
      productId: null as unknown as number | null,
      sessionSports: Array<string | null>(),
      isTrimSelected: false,
      isConcealSelected: false,
      trimMarkers: Array<Marker>(),
      concealMarkers: Array<Marker>(),
      fieldsRemoverList: [
        { value: 'cadence', label: 'Cadence' },
        { value: 'heartRate', label: 'Heart Rate' },
        { value: 'power', label: 'Power' },
        { value: 'temperature', label: 'Temperature' }
      ],
      removeFields: []
    }
  },
  computed: {
    isManufacturerHasProduct(): Boolean {
      const m = this.manufacturerMap.get(this.manufacturerId!)
      if (m == undefined) return false
      return m.products.length > 0
    },
    isManualInputProduct(): Boolean {
      if (this.manufacturerId == null) return false
      if (!this.isManufacturerHasProduct) return true
      if (this.productId == null) return true
      return false
    },
    isValidToProceed(): boolean {
      if (this.encodeMode == EncodeMode.Unknown) return false
      if (this.manufacturerId == null) return false
      if (this.productId == null) return false

      for (let i = 0; i < this.sessionSports.length; i++) {
        if (this.sessionSports[i] == null) return false
      }

      if (this.isTrimSelected) {
        for (let i = 0; i < this.trimMarkers.length; i++) {
          const m = this.trimMarkers[i]
          if (m.startN) return true
          if (m.endN != this.sessions[i].records.length - 1) return true
        }
        return false
      }
      if (this.isConcealSelected) {
        for (let i = 0; i < this.concealMarkers.length; i++) {
          const m = this.concealMarkers[i]
          if (m.startN) return true
          if (m.endN != this.sessions[i].records.length - 1) return true
        }
        return false
      }
      return true
    },
    hasCadence(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.cadence != null) return true
        }
      }
      return false
    },
    hasHeartRate(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.heartRate != null) return true
        }
      }
      return false
    },
    hasPower(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.power != null) return true
        }
      }
      return false
    },
    hasTemperature(): boolean {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        for (let j = 0; j < ses.records.length; j++) {
          const rec = ses.records[j]
          if (rec.temperature != null) return true
        }
      }
      return false
    }
  },
  watch: {
    isTrimSelected: {
      handler() {
        this.updateTrimMarker()
      }
    },
    isConcealSelected: {
      handler() {
        this.updateConcealMarkers()
      }
    },
    trimMarkers: {
      handler(concealMarkers: Marker[]) {
        this.limitMarkers(concealMarkers)
      },
      deep: true
    },
    concealMarkers: {
      handler(concealMarkers: Marker[]) {
        this.limitMarkers(concealMarkers)
      },
      deep: true
    },
    selectedManufacturerProduct: {
      handler(values: any[]) {
        const [manufacturerId, manufacturerName, productId] = values
        this.manufacturerId = manufacturerId
        this.manufacturerName = manufacturerName
        this.productId = productId != null ? (!isNaN(productId) ? productId : null) : null
      }
    },
    activities: {
      handler() {
        this.updateManufacturerProductFromSource()
      }
    },
    sessions: {
      handler() {
        this.encodeMode = EncodeMode.Unknown
        this.updateSessionSports()
      }
    },
    encodeMode: {
      handler() {
        this.reset()
      }
    }
  },
  methods: {
    updateManufacturerProductFromSource() {
      for (let i = 0; i < this.activities.length; i++) {
        const act = this.activities[i]
        if (act.creator.manufacturer != null && act.creator.product != null) {
          const manufacturerId = act.creator.manufacturer
          const productId = act.creator.product
          const m = this.manufacturerMap.get(manufacturerId)
          if (m == undefined) continue
          this.selectedManufacturerProduct = [manufacturerId, m.name, productId] as any[]
          return
        }
      }
      this.selectedManufacturerProduct = []
    },
    updateSessionSports() {
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        this.sessionSports[i] = this.sportMap.get(ses.sport)?.name ?? null
      }
    },
    isNumber(e: KeyboardEvent) {
      if (isNaN(parseInt(e.key))) e.preventDefault()
    },
    limitMarkers(markers: Marker[]) {
      for (let i = 0; i < markers.length; i++) {
        const m = markers[i]
        m.startN = parseInt(m.startN as unknown as string)
        m.endN = parseInt(m.endN as unknown as string)

        if (m.startN > m.endN) {
          m.endN = m.startN
        }
      }
    },
    updateTrimMarker() {
      if (this.sessions.length == 0) return
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        this.trimMarkers[i] = new Marker({ startN: 0, endN: ses.records.length - 1 })
      }
    },
    updateConcealMarkers() {
      if (this.sessions.length == 0) return
      this.concealMarkers = new Array<Marker>(this.sessions.length)
      for (let i = 0; i < this.sessions.length; i++) {
        const ses = this.sessions[i]
        this.concealMarkers[i] = new Marker({ startN: 0, endN: ses.records.length - 1 })
      }
    },
    showFieldRemover(value: string): boolean {
      switch (value) {
        case 'cadence':
          return this.hasCadence
        case 'heartRate':
          return this.hasHeartRate
        case 'power':
          return this.hasPower
        case 'temperature':
          return this.hasTemperature
        default:
          return false
      }
    },
    reset() {
      this.isConcealSelected = false
      this.isTrimSelected = false
      this.updateTrimMarker()
      this.updateConcealMarkers()
    },
    proceed() {
      if (!this.isValidToProceed) return

      const spec = new EncodeSpecifications({
        encodeMode: this.encodeMode,
        targetFileType: FileType.FIT, // TODO: implement other type
        manufacturerId: this.manufacturerId as number,
        productId: this.productId as number,
        deviceName: '', // TODO: not required for FIT File Type
        sports: this.sessionSports,
        trimMarkers: this.trimMarkers,
        concealMarkers: this.concealMarkers,
        removeFields: this.removeFields
      })

      this.$emit('encodeSpecifications', spec)
    }
  },
  mounted() {
    for (let i = 0; i < this.manufacturers.length; i++) {
      const m = this.manufacturers[i]
      this.manufacturerMap.set(m.id, m)
    }
    for (let i = 0; i < this.sports.length; i++) {
      const s = this.sports[i]
      this.sportMap.set(s.name, s)
    }
    this.updateManufacturerProductFromSource()
    this.updateSessionSports()
  },
  unmounted() {}
}
</script>
<style scoped>
label {
  color: var(--bs-heading-color);
}

.sub-label {
  font-size: 0.8em;
}

.emoji {
  font-size: 1.2em;
  color: white;
}

p {
  font-size: 0.8em;
  margin-bottom: 0;
}

.rtl {
  direction: rtl;
}

option[title]:hover::after {
  background-color: black!;
  content: attr(title);
  position: absolute;
  top: -100%;
  left: 0;
}
</style>
