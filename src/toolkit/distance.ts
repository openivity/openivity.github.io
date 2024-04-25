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

export function distanceToHuman(n?: number | null, fractionDigits?: number): string {
  if (!n) return '0 m'

  fractionDigits = fractionDigits ?? 0
  const fraction: number = (fractionDigits ?? 0) <= 0 ? 1 : 10 * fractionDigits

  if (n < 1000) {
    return `${(Math.round(n * fraction) / fraction).toLocaleString()} m`
  }
  return `${(Math.round((n / 1000) * fraction) / fraction).toLocaleString()} km`
}
