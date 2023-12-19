import { activityService, go } from '@/workers/wasm-services'

onmessage = async (e) => {
  await activityService

  if (go.exited) return

  const begin = new Date()

  switch (e.data.type) {
    case 'isReady':
      postMessage({ type: e.data.type })
      break
    case 'decode': {
      // @ts-ignore
      const result = decode(e.data.input)
      const resultJson = JSON.parse(result)
      postMessage({
        type: e.data.type,
        result: resultJson,
        elapsed: new Date().getTime() - begin.getTime()
      })
      break
    }
    case 'encode': {
      // @ts-ignore
      const result = encode(e.data.input)
      const resultJson = JSON.parse(result)
      postMessage({
        type: e.data.type,
        result: resultJson,
        elapsed: new Date().getTime() - begin.getTime()
      })
      break
    }
    case 'manufacturerList': {
      // @ts-ignore
      const manufacturers = JSON.parse(manufacturerList(e.data.input))
      postMessage({
        type: e.data.type,
        result: manufacturers,
        elapsed: new Date().getTime() - begin.getTime()
      })
      break
    }
    case 'sportList': {
      // @ts-ignore
      const sports = JSON.parse(sportList(e.data.input))
      postMessage({
        type: e.data.type,
        result: sports,
        elapsed: new Date().getTime() - begin.getTime()
      })
      break
    }
    case 'shutdown':
      // @ts-ignore
      shutdown()
  }
}
