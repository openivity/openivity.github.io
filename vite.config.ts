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
      includeAssets: ['favicon.ico', 'wasm/activity-service.wasm'],
      workbox: {
        runtimeCaching: [
          {
            urlPattern: new RegExp('/assets/.*\\.(svg|ttf|woff2)$'), // cache assets that don't automatically included.
            handler: 'CacheFirst'
          },
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
        additionalData: `@import "~bootstrap/scss/bootstrap";`
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
