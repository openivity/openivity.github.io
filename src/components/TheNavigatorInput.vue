<template>
  <!-- input -->
  <div class="navigator">
    <div class="navigator-input mx-auto">
      <input
        class="form-control form-control-sm"
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
  max-width: 320px;
}
</style>
