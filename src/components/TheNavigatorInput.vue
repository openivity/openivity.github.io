<template>
  <!-- input -->
  <div class="navigator">
    <div class="navigator-input mx-auto">
      <input
        class="form-control form-control-sm"
        type="file"
        :id="id ?? 'fileInput'"
        accept=".fit,.gpx"
        multiple
        v-bind:disabled="!isWebAssemblySupported"
      />
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    id: String,
    isActivityFileReady: Boolean,
    isWebAssemblySupported: Boolean
  }
}
</script>
<style>
.navigator-input {
  text-align: center;
  max-width: 320px;
}
</style>
