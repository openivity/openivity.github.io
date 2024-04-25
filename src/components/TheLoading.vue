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
  <div class="spinner-container">
    <span v-bind:style="styles" class="spinner spinner--rotate-diamond">
      <div v-bind:style="diamondStyle" class="diamond"></div>
      <div v-bind:style="diamondStyle" class="diamond"></div>
      <div v-bind:style="diamondStyle" class="diamond"></div>
    </span>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    size: {
      default: '70px'
    }
    // color: {
    //   default: '#41b883'
    // }
  },
  computed: {
    diamondStyle() {
      let size = parseInt(this.size)
      return {
        width: size / 4 + 'px',
        height: size / 4 + 'px'
        // '--bg-color': this.color
      }
    },
    styles() {
      let size = parseInt(this.size)
      return {
        width: this.size,
        height: size / 4 + 'px'
      }
    }
  }
}
</script>

<style lang="scss" scoped>
// $accent: #41b883;
$duration: 1500ms;
$timing: linear;

.spinner-container {
  position: fixed;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  z-index: 10000;

  .spinner {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    * {
      line-height: 0;
      box-sizing: border-box;
    }
    .diamond {
      position: absolute;
      left: 0;
      top: 0;
      border-radius: 2px;
      background: var(--green-text);
      transform: translateX(-50%) rotate(45deg) scale(0);
      animation: diamonds $duration $timing infinite;
      @for $i from 1 through 4 {
        &:nth-child(#{$i}) {
          animation-delay: calc($duration * -1 / 1.5) * $i;
        }
      }
    }
  }
}
@keyframes diamonds {
  50% {
    left: 50%;
    transform: translateX(-50%) rotate(45deg) scale(1);
  }
  100% {
    left: 100%;
    transform: translateX(-50%) rotate(45deg) scale(0);
  }
}
</style>
