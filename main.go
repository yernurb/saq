package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

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

	// Prepare camera module
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	// Prepare image container matrix and an array of images
	img := gocv.NewMat()
	defer img.Close()
	var imgArray [maxFrames]gocv.Mat

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

		currentTime := time.Now().Format("Monday / _2 January 2006 / 15:04:05")
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
