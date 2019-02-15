package main

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()
	fmt.Println(img.Size())
	var imgArray [250]gocv.Mat

	frameNum := 0
	maxFrameNum := 10
	start := time.Now()
	i := 0
	for {
		if i >= 250 {
			break
		}
		frameNum++
		if frameNum > maxFrameNum {
			frameNum = 0
			elapsed := time.Since(start)
			start = time.Now()
			fmt.Println("FPS:", float64(maxFrameNum)/elapsed.Seconds())
		}
		webcam.Read(&img)
		imgArray[i] = img.Clone()
		window.IMShow(img)
		window.WaitKey(1)
		i++
	}

	for i, im := range imgArray {
		fmt.Println(i, im.Size())
	}
}
