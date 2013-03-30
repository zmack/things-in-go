package main

import (
  "fmt"
)

func addThings(value []int) {
  for i := 0; i < len(value); i++ {
    value[i] = 10
  }
}

func main() {
  value := make([]int, 10)
  fmt.Println(value)
  addThings(value)
  fmt.Println(value)
}
