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
