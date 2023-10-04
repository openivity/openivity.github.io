import '@/assets/wasm/wasm_exec.js'

const go = new Go()

const wasmUrl = new URL('/wasm/fitsvc.wasm', import.meta.url)
const wasm = WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject).then((wasm) => {
  go.run(wasm.instance)
})

onmessage = async (e) => {
  await wasm
  // @ts-ignore
  const result = decode(e.data)
  postMessage(result)
}
