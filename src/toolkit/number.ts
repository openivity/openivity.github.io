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

export function avg(...values: (number | null)[]): number | null {
  let valid = 0
  let result = 0

  for (let i = 0; i < values.length; i++) {
    const value = values[i]
    if (!value) continue
    result += value
    valid++
  }

  if (valid == 0) return null

  return result / valid
}

export function max(...values: (number | null)[]): number | null {
  let max = 0
  let valid = 0
  for (let i = 0; i < values.length; i++) {
    const value = values[i]
    if (!value) continue
    if (value > max) max = value
    valid++
  }

  if (valid == 0) return null

  return max
}

export function sum(...values: (number | null)[]): number | null {
  let sum = 0
  let valid = 0
  for (let i = 0; i < values.length; i++) {
    const value = values[i]
    if (value == null) continue
    sum += value
    valid++
  }

  if (valid == 0) return null

  return sum
}
