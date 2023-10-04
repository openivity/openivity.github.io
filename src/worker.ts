import '@/assets/wasm/wasm_exec.js'

const go = new Go()

const wasmUrl = '/wasm/fitsvc.wasm'
const wasm = WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject).then((wasm) => {
  go.run(wasm.instance)
})

onmessage = async (e) => {
  await wasm
  // @ts-ignore
  const result = decode(e.data)
  postMessage(result)
}
