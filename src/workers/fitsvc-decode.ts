import { fitsvc } from '@/workers/fitsvc'

onmessage = async (e) => {
  await fitsvc
  // @ts-ignore
  const result = decode(e.data)
  postMessage(result)
}
