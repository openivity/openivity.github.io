export class ActivityFile {
  fileId: FileId
  activity: Activity
  sessions: Array<Session>
  records: Array<Record>

  constructor(undefined)
}

export class Activity {
  event: string
  eventType: string
  numSession: number
  timestamp: Date
  type: string
}

export class FileId {
  manufacturer: string
  product: number
  timeCreated: Date
}

export class Session {
  sport: string
  subSport: string
  totalMovingTime: number
  totalElapsedTime: number
  totalTimerTime: number
  totalDistance: number
  totalAscent: number
  totalDescent: number
  totalCycles: number
  totalCalories: number
  avgSpeed: number
  maxSpeed: number
  avgHeartRate: number
  maxHeartRate: number
  avgCadence: number
  maxCadence: number
  avgPower: number
  maxPower: number
  avgTemperature: number
  maxTemperature: number
  avgAltitude: number
  maxAltitude: number
}

export class Record {
  altitude: number
  cadence: number
  distance: number
  heartRate: number
  positionLat: number
  positionLong: number
  speed: number
  timestamp: Date
}
