import '@/assets/wasm/wasm_exec.js'

const go = new Go()

const wasmUrl = new URL('/wasm/fitsvc.wasm', import.meta.url)
export const fitsvc = WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject).then(
  (wasm) => {
    go.run(wasm.instance)
  }
)
