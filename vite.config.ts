// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import { resolve } from 'node:path'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VitePWA({
      registerType: 'autoUpdate',
      workbox: {
        globDirectory: 'dist',
        globPatterns: ['**/*.{js,css,html,ico,wasm,svg,ttf,woff2}'],
        maximumFileSizeToCacheInBytes: (1000 * 10) << 10, // 10MB per file
        runtimeCaching: [
          {
            urlPattern: new RegExp('^https://.*\\.openstreetmap.org/.*\\.png$'), // cache osm tiles.
            handler: 'CacheFirst'
          }
        ]
      }
    })
  ],
  css: {
    preprocessorOptions: {
      scss: {
        // additionalData: `@import "src/assets/main.scss";`
      }
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '~bootstrap': resolve(__dirname, 'node_modules/bootstrap')
    }
  },
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        wasmServices: resolve(__dirname, 'src/workers/wasm-services.ts'),
        activityService: resolve(__dirname, 'src/workers/activity-service.ts')
      }
    }
  }
})
