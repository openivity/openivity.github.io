<!-- Copyright (C) 2023 Openivity

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. -->

<script setup lang="ts">
import { Session } from '@/spec/activity'
import { Sport, ToolMode } from '@/spec/activity-service'
import type { PropType } from 'vue'
</script>
<template>
  <div>
    <label>Change Sport</label>
    <div class="col-12 pb-1" v-for="(ses, index) in sessions" :key="index">
      <label class="form-label sub-label">Session {{ index + 1 }}'s sport:</label>
      <v-select
        label="name"
        placeholder="Please select a sport"
        :clearable="false"
        :options="sports"
        v-model="sessionSports[index]"
        :disabled="toolMode == ToolMode.Unknown"
      >
      </v-select>
    </div>
  </div>
</template>
<script lang="ts">
export default {
  props: {
    toolMode: { type: Number as PropType<ToolMode>, required: true },
    sessions: { type: Array<Session>, required: true },
    sports: { type: Array<Sport>, required: true }
  },
  data() {
    return {
      sessionSports: Array<Sport | null>()
    }
  },
  computed: {
    sportMap(): Map<string, Sport> {
      const map = new Map()
      for (let i = 0; i < this.sports.length; i++) {
        map.set(this.sports[i].name, this.sports[i])
      }
      return map
    }
  },
  watch: {
    sessions: {
      handler(values) {
        this.updateSessionSports(values)
      }
    },
    sessionSports: {
      handler(values: Sport[]) {
        this.$emit(
          'sessionSports',
          values.flatMap((s) => s.name)
        )
      },
      deep: true
    }
  },
  methods: {
    updateSessionSports(sessions: Session[]) {
      for (let i = 0; i < sessions.length; i++) {
        const ses = sessions[i]
        this.sessionSports[i] = this.sportMap.get(ses.sport) ?? this.sportMap.get('Generic')!
      }
    }
  },
  mounted() {
    this.updateSessionSports(this.sessions)
  }
}
</script>
<style scoped>
@import '@/assets/tools.scss';
</style>
