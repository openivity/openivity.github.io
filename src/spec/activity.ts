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

export const SPORT_GENERIC = 'Generic'
export const UNKNOWN = 'Unknown'

export class ActivityFile {
  creator: Creator = new Creator()
  timezone: number = 0
  sessions: Session[] = []

  constructor(json?: any) {
    const casted = json as ActivityFile
    this.creator = casted?.creator
    this.timezone = casted?.timezone
    this.sessions = casted?.sessions
  }
}

export class Creator {
  name: string = UNKNOWN
  manufacturer: number = 0
  product: number = 0
  timeCreated: string | null = null
}

export class Session {
  timestamp: string = ''
  startTime: string = ''
  endTime: string = ''
  sport: string = SPORT_GENERIC
  totalMovingTime: number | null = null
  totalElapsedTime: number | null = null
  totalDistance: number | null = null
  totalAscent: number | null = null
  totalDescent: number | null = null
  totalCycles: number | null = null
  totalCalories: number | null = null
  avgSpeed: number | null = null
  maxSpeed: number | null = null
  avgHeartRate: number | null = null
  maxHeartRate: number | null = null
  avgCadence: number | null = null
  maxCadence: number | null = null
  avgPower: number | null = null
  maxPower: number | null = null
  avgTemperature: number | null = null
  maxTemperature: number | null = null
  avgAltitude: number | null = null
  maxAltitude: number | null = null
  avgPace: number | null = null
  avgElapsedPace: number | null = null

  timezone: number = 0
  workoutType: WorkoutType = WorkoutType.Moving
  laps: Lap[] = []
  records: Record[] = []

  // additional info
  timeCreated: string | null = null
  creatorName: string = UNKNOWN
}

export enum WorkoutType {
  Moving = 0,
  Stationary
}

export class Lap {
  timestamp: string | null = null
  totalMovingTime: number | null = null
  totalElapsedTime: number | null = null
  totalDistance: number | null = null
  totalAscent: number | null = null
  totalDescent: number | null = null
  totalCalories: number | null = null
  avgSpeed: number | null = null
  maxSpeed: number | null = null
  avgHeartRate: number | null = null
  maxHeartRate: number | null = null
  avgCadence: number | null = null
  maxCadence: number | null = null
  avgPower: number | null = null
  maxPower: number | null = null
  avgPace: number | null = null
  avgElapsedPace: number | null = null
}

export class Record {
  timestamp: string | null = null
  positionLat: number | null = null
  positionLong: number | null = null
  distance: number | null = null
  speed: number | null = null
  altitude: number | null = null
  cadence: number | null = null
  heartRate: number | null = null
  power: number | null = null
  temperature: number | null = null
  grade: number = 0
  pace: number | null = null

  constructor(data?: any) {
    const casted = data as Record
    this.positionLat = casted?.positionLat
    this.positionLong = casted?.positionLong
    this.altitude = casted?.altitude
    this.cadence = casted?.cadence
    this.distance = casted?.distance
    this.heartRate = casted?.heartRate
    this.speed = casted?.speed
    this.timestamp = casted?.timestamp
    this.power = casted?.power
    this.temperature = casted?.temperature
    this.grade = casted?.grade
    this.pace = casted?.pace
  }
}
