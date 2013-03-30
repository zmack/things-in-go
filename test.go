package main

import (
  "fmt"
  "math"
)

func Square(value float64) float64 {
  var z, previous, iterations = 10.0, 0.0, 0

  for math.Abs(previous - z) > 0.0001 {
    previous = z
    z = z - (z*z - value)/(2*z)
    iterations += 1
  }

  fmt.Printf("%d iterations\n", iterations)

  return z
}

func main() {
  fmt.Println(Square(133))
}
