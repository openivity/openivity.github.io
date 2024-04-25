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

export class Summary {
  sport: string | null = null
  subSport: string | null = null
  totalMovingTime: number | null = null
  totalElapsedTime: number | null = null
  totalTimerTime: number | null = null
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
}
