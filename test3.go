package main

import (
  "strings"
  "code.google.com/p/go-tour/wc"
)

func WordCount(s string) map[string]int {
  words := strings.Fields(s)
  counts := make(map[string]int)

  for _, word := range words {
    value, ok := counts[word]
    if ok {
      counts[word] = value + 1
    } else {
      counts[word] = 1
    }
  }

  return counts
}

func main() {
  wc.Test(WordCount)
}
