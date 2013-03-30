package main

import "code.google.com/p/go-tour/pic"

func Pic(dx, dy int) [][]uint8 {
  image := make([][]uint8, dy)

  for i := 0; i < dy; i += 1 {
    image[i] = make([]uint8, dx)
    for j := 0; j < dx; j += 1 {
      image[i][j] = uint8(i*j)
    }
  }

  return image
}

func main() {
	pic.Show(Pic)
}
