import '@/assets/wasm/wasm_exec.js'

export const go = new Go()

const wasmUrl = new URL('/wasm/activity-service.wasm', import.meta.url)
export const activityService = WebAssembly.instantiateStreaming(
  fetch(wasmUrl),
  go.importObject
).then((wasm) => {
  go.run(wasm.instance)
})
