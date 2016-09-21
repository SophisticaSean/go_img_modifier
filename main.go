package main

import (
	"image"
	"os"
	"strconv"
	"sync"

	"github.com/disintegration/imaging"
)

func process(img_path string, name string, blur_amount int, wg *sync.WaitGroup) {
	img, _ := imaging.Open(img_path)
	bigImg := img
	max := img.Bounds().Max
	var min image.Point
	min.X = max.X / blur_amount
	min.Y = max.Y / blur_amount

	smallImg := imaging.Resize(img, min.X, min.Y, imaging.NearestNeighbor)
	bigImg = imaging.Resize(smallImg, max.X, max.Y, imaging.NearestNeighbor)

	file_name := name + strconv.Itoa(blur_amount) + ".png"
	imaging.Save(bigImg, file_name)
	wg.Done()
}

func multi_process(args []string) {
	var wg sync.WaitGroup
	img_path := args[1]
	img_name := args[2]
	blur_amount, _ := strconv.Atoi(args[3])
	lock_img, _ := imaging.Open(args[4])
	monitor_count, _ := strconv.Atoi(args[5])

	for i := 2; i <= blur_amount; i++ {
		wg.Add(1)
		go process(img_path, img_name, i, &wg)
	}
	wg.Wait()

	final_file := img_name + strconv.Itoa(blur_amount) + ".png"
	final_file_name := img_name + strconv.Itoa(blur_amount+1) + ".png"
	img, _ := imaging.Open(final_file)
	if monitor_count == 1 {
		lockImg := imaging.OverlayCenter(img, lock_img, 1.0)
		imaging.Save(lockImg, final_file_name)
	} else if monitor_count == 3 {
		var (
			mon1         image.Point
			mon2         image.Point
			mon3         image.Point
			lockImgPoint image.Point
		)

		lockImgPoint.X = lock_img.Bounds().Max.X / 2
		lockImgPoint.Y = lock_img.Bounds().Max.Y / 2

		mon1.X = (1920 / 2) - lockImgPoint.X
		mon1.Y = (1200 / 2) - lockImgPoint.Y
		mon2.X = mon1.X + 1920
		mon2.Y = mon1.Y
		mon3.X = mon1.X
		mon3.Y = mon1.Y + 1200

		finalImg := imaging.Overlay(img, lock_img, mon1, 1.0)
		finalImg = imaging.Overlay(finalImg, lock_img, mon2, 1.0)
		finalImg = imaging.Overlay(finalImg, lock_img, mon3, 1.0)
		imaging.Save(finalImg, final_file_name)
	}
}

func main() {
	multi_process(os.Args)
	//img, _ := imaging.Open(os.Args[1])
	//bigImg := img
	//max := img.Bounds().Max
	//var min image.Point
	//var mid image.Point
	//blur_amount, _ := strconv.Atoi(os.Args[3])
	//min.X = max.X / blur_amount
	//min.Y = max.Y / blur_amount
	//mid.X = max.X / 2
	//mid.Y = max.Y / 2

	//if os.Args[4] == "blur" {
	//bigImg = imaging.Blur(img, float64(blur_amount))
	//} else {
	//smallImg := imaging.Resize(img, min.X, min.Y, imaging.NearestNeighbor)
	//bigImg = imaging.Resize(smallImg, max.X, max.Y, imaging.NearestNeighbor)
	//}
	//if len(os.Args) == 6 {
	//lock, _ := imaging.Open(os.Args[5])
	//bigImg = imaging.OverlayCenter(bigImg, lock, 1.0)
	//}
	//imaging.Save(bigImg, os.Args[2])
}
