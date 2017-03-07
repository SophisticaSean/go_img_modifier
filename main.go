package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/lwalen/vips"
)

var options = vips.Options{
	Width:        100,
	Height:       100,
	Crop:         false,
	Enlarge:      true,
	Extend:       vips.EXTEND_WHITE,
	Interpolator: vips.BICUBIC,
	Gravity:      vips.CENTRE,
	Quality:      100,
	Format:       vips.PNG,
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	bounds := image.Bounds()
	return bounds.Max.X, bounds.Max.Y
}

func process(args []string) {
	imgString := args[1]
	destString := args[2]
	factor, _ := strconv.Atoi(args[3])
	lockImg, _ := imaging.Open(args[4])
	monitorCount, _ := strconv.Atoi(args[5])

	width, height := getImageDimension(imgString)
	resizeWidth := width / factor
	resizeHeight := height / factor
	options.Width = resizeWidth
	options.Height = resizeHeight

	// resize down
	orig, _ := os.Open(imgString)
	inBuf, _ := ioutil.ReadAll(orig)
	buf, err := vips.Resize(inBuf, options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// resize back up
	options.Width = width
	options.Height = height
	final, err := vips.Resize(buf, options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	finalBuffer := bytes.NewBuffer(final)

	img, _, err := image.Decode(finalBuffer)
	if err != nil {
		panic(err)
	}

	//dest, err := os.Create(destString)
	//if err != nil {
	//panic(err)
	//}
	//defer dest.Close()

	//dest.Write(final)

	//finalFile := img_name + strconv.Itoa(blurAmount) + ".png"
	//img, _ := imaging.Open(finalFile)

	if monitorCount == 1 {
		lockImg := imaging.OverlayCenter(img, lockImg, 1.0)
		imaging.Save(lockImg, destString)
	} else if monitorCount == 3 {
		var (
			mon1         image.Point
			mon2         image.Point
			mon3         image.Point
			lockImgPoint image.Point
		)

		lockImgPoint.X = lockImg.Bounds().Max.X / 2
		lockImgPoint.Y = lockImg.Bounds().Max.Y / 2

		mon1.X = (1920 / 2) - lockImgPoint.X
		mon1.Y = (1200 / 2) - lockImgPoint.Y
		mon2.X = mon1.X + 1920
		mon2.Y = mon1.Y
		mon3.X = mon1.X + (1920 * 2)
		mon3.Y = mon1.Y

		finalImg := imaging.Overlay(img, lockImg, mon1, 1.0)
		finalImg = imaging.Overlay(finalImg, lockImg, mon2, 1.0)
		finalImg = imaging.Overlay(finalImg, lockImg, mon3, 1.0)
		imaging.Save(finalImg, destString)
	}
}

func main() {
	process(os.Args)
}
