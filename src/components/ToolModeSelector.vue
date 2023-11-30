<script setup lang="ts">
import type { ActivityFile, Session } from '@/spec/activity'
import { ToolMode } from '@/spec/activity-service'
</script>
<template>
  <div>
    <label>Please select a tool</label>
    <div class="col-12 pt-1">
      <v-select
        label="label"
        placeholder="-- Please select a tool first --"
        :selectable="isSelectable"
        :clearable="false"
        :searchable="false"
        :options="dataSource"
        v-model="selected"
      >
        <template #option="{ label, value, title, selectable }">
          <span v-if="!selectable" :title="title">{{ label }}</span>
          <span v-if="selectable && value != ToolMode.Unknown">{{ label }}</span>
        </template>
      </v-select>
      <div class="pt-1">
        <p v-show="selected?.value == ToolMode.Edit">
          We will edit relevant data for every input activities. This changes will apply to your
          entire activities, like a Bulk Edit, if you want to edit one activity, please open only
          one at a time.
        </p>
        <p v-show="selected?.value == ToolMode.Combine">
          We will combine multiple activities into one continuous activity file. Two sequential
          sessions of the same sport will be merged into one session. If the sport is different, it
          will be placed in separate sessions. This process will continue until all sessions are
          combined.
        </p>
        <p v-show="selected?.value == ToolMode.SplitPerSession">
          We will create new Activity File for every Sessions in all activities.
        </p>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
export default {
  props: {
    activities: { type: Array<ActivityFile>, required: true },
    sessions: { type: Array<Session>, required: true }
  },
  data() {
    return {
      selected: null as unknown as any | null
    }
  },
  computed: {
    dataSource(): {}[] {
      const dataSource = [
        { label: 'Edit Relevant Data', value: ToolMode.Edit, selectable: true },
        {
          label: 'Combine Multiple Activities into One',
          value: ToolMode.Combine,
          title:
            'You have only one activity opened, please open multiple activites to be able to use this feature.',
          selectable: this.activities.length > 1
        },
        {
          label: 'Split Activities Per Session',
          value: ToolMode.SplitPerSession,
          title:
            'You have only one session in the opened activity, please open multiple activities or open an activity that have multiple sessions to be able to use this feature.',
          selectable: this.sessions.length > 1
        }
      ]
      return dataSource
    }
  },
  watch: {
    sessions: {
      handler() {
        this.selected = null
      }
    },
    selected: {
      handler(value) {
        this.$emit('toolMode', value?.value ?? ToolMode.Unknown)
      }
    }
  },
  methods: {
    isSelectable(option: any) {
      return option.selectable
    }
  },
  mounted() {
    this.$emit('toolMode', this.selected ?? ToolMode.Unknown)
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
