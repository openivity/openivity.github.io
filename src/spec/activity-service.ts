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
  encodeTook: number
  serializationTook: number
  totalElapsed: number
  fileName: string
  fileType: string
  fileBytes: Uint8Array

  constructor(data?: any) {
    const casted = data as EncodeResult
    this.err = casted?.err
    this.encodeTook = casted?.encodeTook
    this.serializationTook = casted?.serializationTook
    this.totalElapsed = casted?.totalElapsed
    this.fileName = casted?.fileName
    this.fileType = casted?.fileType
    this.fileBytes = casted?.fileBytes
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

export class Product {
  id: number = 0
  name: string = ''

  constructor(data?: any) {
    const casted = data as Product
    this.id = casted?.id
    this.name = casted?.name
  }
}
