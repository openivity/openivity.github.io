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

import type { ActivityFile } from './activity'

export class DecodeResult {
  err: string | null = null
  activities: Array<ActivityFile>
  decodeTook: number
  serializationTook: number
  totalElapsed: number

  constructor(json?: any) {
    const casted = json as DecodeResult

    this.err = casted?.err
    this.activities = casted?.activities
    this.decodeTook = casted?.decodeTook
    this.serializationTook = casted?.serializationTook
    this.totalElapsed = casted?.totalElapsed
  }
}

export class EncodeResult {
  err: string | null = null
  deserializeInputTook: number
  encodeTook: number
  serializationTook: number
  totalElapsed: number
  fileName: string
  fileType: string
  filesBytes: Uint8Array[]

  constructor(data?: any) {
    const casted = data as EncodeResult
    this.err = casted?.err
    this.deserializeInputTook = casted?.deserializeInputTook
    this.encodeTook = casted?.encodeTook
    this.serializationTook = casted?.serializationTook
    this.totalElapsed = casted?.totalElapsed
    this.fileName = casted?.fileName
    this.fileType = casted?.fileType
    this.filesBytes = casted?.filesBytes
  }
}

export class EncodeSpecifications {
  toolMode: number = 0
  targetFileType: FileType = 0
  manufacturerId: number = 0
  productId: number = 0
  deviceName: string = 'Unknown'
  sports: (string | null)[] = []
  trimMarkers?: Marker[] | null = []
  concealMarkers?: Marker[] | null = []
  removeFields?: string[] | null = []

  constructor(data: EncodeSpecifications) {
    this.toolMode = data.toolMode
    this.targetFileType = data.targetFileType
    this.manufacturerId = data.manufacturerId
    this.productId = data.productId
    this.deviceName = data.deviceName
    this.sports = data.sports
    this.trimMarkers = data.trimMarkers
    this.concealMarkers = data.concealMarkers
    this.removeFields = data.removeFields
  }
}

export enum ToolMode {
  Unknown = 0,
  Edit,
  Combine,
  SplitPerSession
}

export enum FileType {
  Unsupported = 0,
  FIT,
  GPX,
  TCX
}

export class Marker {
  startN: number = 0
  endN: number = 0

  constructor(data?: Marker) {
    this.startN = data?.startN ?? 0
    this.endN = data?.endN ?? 0
  }
}

export class ManufacturerListResult {
  manufacturers: Manufacturer[] = []

  constructor(data?: any) {
    const casted = data as ManufacturerListResult
    this.manufacturers = casted?.manufacturers
  }
}

export class Manufacturer {
  id: number = 0
  name: string = ''
  products: Product[] = []

  constructor(data?: any) {
    const casted = data as Manufacturer
    this.id = casted?.id
    this.name = casted?.name
    this.products = casted?.products
  }
}

export class SportListResult {
  sports: Sport[] = []

  constructor(data?: any) {
    const casted = data as SportListResult
    this.sports = casted?.sports
  }
}

export class Sport {
  id: number = 0
  name: string = ''
  toleranceMovingSpeed: number = 0

  constructor(data?: Sport) {
    this.id = data?.id ?? 0
    this.name = data?.name ?? ''
    this.toleranceMovingSpeed = data?.toleranceMovingSpeed ?? 0
  }
}

export class Product {
  id: number = 0
  name: string = ''

  constructor(data?: any) {
    const casted = data as Product
    this.id = casted?.id
    this.name = casted?.name
  }
}
