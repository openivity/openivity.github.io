import { activityService } from '@/workers/wasm-services'

onmessage = async (e) => {
  await activityService
  // @ts-ignore
  const result = decode(e.data)
  postMessage(result)
}
