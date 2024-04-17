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

import { SPORT_GENERIC, type Session } from '@/spec/activity'

// Check session has pace
export function sessionHasPace(session: Session): boolean {
  switch (session.sport) {
    case 'Hiking':
    case 'Walking':
    case 'Running':
    case 'Swimming':
    case 'Transition':
    case SPORT_GENERIC:
      for (let j = 0; j < session.records.length; j++) {
        const rec = session.records[j]
        if (rec.pace != null) return true
      }
  }
  return false
}
