export function distanceToHuman(
    n?: number | null,
    fractionDigits?: number,
): string {
    if (!n) return '0 m'

    fractionDigits = fractionDigits ?? 0
    const fraction: number = (fractionDigits ?? 0) <= 0 ? 1 : 10 * fractionDigits

    if (n < 1000) {
        return `${(Math.round(n * fraction) / fraction).toLocaleString()} m`
    }
    return `${(Math.round((n / 1000) * fraction) / fraction).toLocaleString()} km`
}