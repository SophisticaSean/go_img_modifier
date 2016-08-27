package main

import (
	"image"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

func main() {
	img, _ := imaging.Open(os.Args[1])
	bigImg := img
	max := img.Bounds().Max
	var min image.Point
	var mid image.Point
	blur_amount, _ := strconv.Atoi(os.Args[3])
	min.X = max.X / blur_amount
	min.Y = max.Y / blur_amount
	mid.X = max.X / 2
	mid.Y = max.Y / 2

	if (os.Args[4] == "blur") {
		bigImg = imaging.Blur(img, float64(blur_amount))
	} else {
		smallImg := imaging.Resize(img, min.X, min.Y, imaging.NearestNeighbor)
		bigImg = imaging.Resize(smallImg, max.X, max.Y, imaging.NearestNeighbor)
	}
	if (len(os.Args) == 6) {
		lock, _ := imaging.Open(os.Args[5])
		bigImg = imaging.OverlayCenter(bigImg, lock, 1.0)
	}
	imaging.Save(bigImg, os.Args[2])
}
