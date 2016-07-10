package main

import (
	"image/color"
)

func GetTriadColor(main color.RGBA) (oneColor color.RGBA, twoColor color.RGBA) {
	oneColor = color.RGBA{main.B, main.R, main.G, 0xFF}
	twoColor = color.RGBA{main.G, main.B, main.R, 0xFF}
	return
}
