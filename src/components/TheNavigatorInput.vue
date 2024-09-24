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

<template>
  <!-- input -->
  <div class="navigator">
    <div class="navigator-input mx-auto">
      <input
        class="form-control form-control-sm"
        style="font-size: 1em"
        type="file"
        :id="id"
        multiple
        ref="input"
        v-bind:disabled="!isWebAssemblySupported"
      />
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    id: {
      type: String,
      default: 'fileInput'
    },
    isActivityFileReady: Boolean,
    isWebAssemblySupported: Boolean
  },
  mounted() {
    const isSafari = /^((?!chrome|android).)*safari/i.test(navigator.userAgent)
    if (!isSafari)
      // NOTE: Safari on iOS has specific behavior when it comes to the accept attribute on file input fields.
      //       It doesn't have built-in support for handling certain file types, and it might not recognize or handle the .fit extension as expected.
      this.$nextTick(
        () => ((document.getElementById(this.id) as HTMLInputElement).accept = '.fit, .gpx, .tcx')
      )
  }
}
</script>
<style>
.navigator-input {
  text-align: center;
  max-width: 360px;
}
</style>
