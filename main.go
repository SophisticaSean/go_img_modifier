package main

import (
	"image"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

func main() {
	img, _ := imaging.Open(os.Args[1])
	max := img.Bounds().Max
	var min image.Point
	blur_amount, _ := strconv.Atoi(os.Args[3])
	min.X = max.X / blur_amount
	min.Y = max.Y / blur_amount

	//smallImg := imaging.Blur(img, 20)
	smallImg := imaging.Resize(img, min.X, min.Y, imaging.NearestNeighbor)
	bigImg := imaging.Resize(smallImg, max.X, max.Y, imaging.NearestNeighbor)
	imaging.Save(bigImg, os.Args[2])
}
