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
    if (!value) continue
    sum += value
    valid++
  }

  if (valid == 0) return null

  return sum
}
