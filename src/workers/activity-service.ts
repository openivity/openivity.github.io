import { activityService } from '@/workers/wasm-services'

onmessage = async (e) => {
  await activityService

  const begin = new Date()

  switch (e.data.type) {
    case 'isReady':
      postMessage({ type: e.data.type })
      break
    case 'decode': {
      // @ts-ignore
      const result = decode(e.data.input)
      postMessage({
        type: e.data.type,
        result: result,
        elapsed: new Date().getTime() - begin.getTime()
      })
      break
    }
    case 'encode': {
      // @ts-ignore
      const result = encode(e.data.input)
      postMessage({
        type: e.data.type,
        result: result,
        elapsed: new Date().getTime() - begin.getTime()
      })
      break
    }
    case 'manufacturerList': {
      // @ts-ignore
      const manufacturers = manufacturerList(e.data.input)
      postMessage({
        type: e.data.type,
        result: manufacturers,
        elapsed: new Date().getTime() - begin.getTime()
      })
      break
    }
  }
}
