<template>
  <!-- input -->
  <div class="navigator">
    <div class="navigator-input" style="text-align: center">
      <input
        class="input"
        type="file"
        :id="id ?? 'fileInput'"
        accept=".fit"
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
    isWebAssemblySupported: Boolean
  }
}
</script>

<style>
.input {
  color: var(--color-title);
  display: inline-block;
  padding: 15px 20px;
  background-color: var(--color-background-mute);
  color: var(--color-text);
  border-radius: 5px;
  cursor: pointer;
  margin: 0 0 5px 0;
}
</style>
