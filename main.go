package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/yernurb/saq/tracker"
	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tsavevideo [camera ID] [video file]")
		return
	}

	deviceID := os.Args[1]
	saveFile := os.Args[2]
	const maxFrames = 250

	saq := new(tracker.Saq)
	// Prepare camera module
	saq.Init(0)

	// Prepare image container matrix and an array of images
	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok {
		fmt.Printf("Cannot read device %v\n", deviceID)
		return
	}
	fmt.Println(img.Cols(), img.Rows())

	start := time.Now()
	// Main loop
	for i := 0; i < maxFrames; i++ {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}

		if img.Empty() {
			continue
		}

		currentTime := tracker.TextifyTime()
		fmt.Println(currentTime)
		pt := image.Pt(30, 30)
		gocv.PutText(&img, currentTime, pt, gocv.FontHersheySimplex, 0.6, color.RGBA{255, 0, 0, 0}, 2)
		imgArray[i] = img.Clone()
	}
	// Check FPS
	elapsed := time.Since(start)
	FPS := float64(maxFrames) / elapsed.Seconds()
	fmt.Println("Average FPS:", FPS)
	// Prepare video writer object
	writer, err := gocv.VideoWriterFile(saveFile, "MP42", FPS, img.Cols(), img.Rows(), true)
	if err != nil {
		fmt.Printf("error opening video writer device: %v\n", saveFile)
		return
	}
	defer writer.Close()

	for i := 0; i < maxFrames; i++ {
		writer.Write(imgArray[i])
	}
}
