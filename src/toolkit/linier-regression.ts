export class Point {
  X: number = 0
  Y: number = 0

  constructor(x: number, y: number) {
    this.X = x
    this.Y = y
  }
}

export class LinierRegression {
  private a: number = 0
  private b: number = 0

  train(points: Point[]) {
    let sumX: number = 0
    let sumY: number = 0
    let sumXY: number = 0
    let sumX2: number = 0

    points.forEach((point) => {
      sumX += point.X
      sumY += point.Y
      sumXY += point.X * point.Y
      sumX2 += point.X * point.X
    })

    const n = points.length

    // Formula:
    // b = n * Σxy - (Σx) * (Σy) / n * Σx² - (Σx)²
    // a = Σy - b(Σx) / n
    this.b = (n * sumXY - sumX * sumY) / (n * sumX2 - sumX * sumX)
    this.a = (sumY - this.b * sumX) / n
  }

  // y = a + bx
  predictY(x: number): number {
    return this.a + x * this.b
  }
}
