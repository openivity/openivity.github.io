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
