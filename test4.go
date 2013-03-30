package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
  var num_prev, num int = -1, -1
  return func() int {
    if num_prev == -1 {
      if num == -1 {
        num = 0
        return num
      } else {
        num_prev = 0
        num = 1
        return num
      }
    }
    current := num + num_prev
    num_prev = num
    num = current

    return current
  }
}

func main() {
  f := fibonacci()
  for i := 0; i < 10; i++ {
    fmt.Println(f())
  }
}

