<template>
  <div class="elevation" v-if="elevations">Elevation</div>
</template>

<script lang="ts">
export default {
  props: {
    elevations: {}
  }
}
</script>
<style>
.elevation {
  color: var(--color-title);
  padding: 10px;
}
</style>
