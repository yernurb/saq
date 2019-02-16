package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

const (
	maxFrames   = 250
	frameWidth  = 640
	frameHeight = 480
)

func writeBuffer(fileName string, FPS float64, buffer *[maxFrames]gocv.Mat) {
	writer, err := gocv.VideoWriterFile(fileName, "MP42", FPS, frameWidth, frameHeight, true)
	if err != nil {
		fmt.Printf("error opening video writer device: %v\n", fileName)
		return
	}
	defer writer.Close()
	for i := 0; i < maxFrames; i++ {
		writer.Write(buffer[i])
		buffer[i].Close()
	}
	writer.Close()
}

func main() {
	fmt.Println("Starting capture...")
	deviceID := 0
	camera, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer camera.Close()

	// Preparing double buffer
	var (
		frameBuffer1  [maxFrames]gocv.Mat
		frameBuffer2  [maxFrames]gocv.Mat
		currentBuffer *[maxFrames]gocv.Mat
	)

	// Preparing frame container
	frame := gocv.NewMat()
	defer frame.Close()

	// Test read from the camera module
	camera.Read(&frame)
	fmt.Println(frame.Cols(), frame.Rows())

	// Setting the camera resolution
	camera.Set(gocv.VideoCaptureFrameWidth, frameWidth)
	camera.Set(gocv.VideoCaptureFrameHeight, frameHeight)

	// Setting the initial frame buffer pointer
	currentBuffer = &frameBuffer1
	currentIndex := 0

	// Starting the main loop
	fileIndex := 0
	start := time.Now()
	for {
		camera.Read(&frame)
		if frame.Empty() {
			continue
		}

		pt := image.Pt(30, 30)
		text := time.Now().Format("Monday / _2 January 2006 / 15:04:05")
		//fmt.Println(text)
		gocv.PutText(&frame, text, pt, gocv.FontHersheySimplex, 0.6, color.RGBA{255, 0, 0, 0}, 2)
		if currentBuffer == &frameBuffer1 {
			frameBuffer1[currentIndex] = frame.Clone()
		} else {
			frameBuffer2[currentIndex] = frame.Clone()
		}
		currentIndex++

		// If current buffer is full swap it and write the collected frames into file
		if currentIndex >= maxFrames {
			elapsed := time.Since(start)
			FPS := float64(maxFrames) / elapsed.Seconds()
			fmt.Println("FPS:", FPS)
			fileIndex++
			fileName := "video-" + strconv.Itoa(fileIndex) + ".avi"
			if currentBuffer == &frameBuffer1 {
				go writeBuffer(fileName, FPS, &frameBuffer1)
				currentBuffer = &frameBuffer2
			} else {
				go writeBuffer(fileName, FPS, &frameBuffer2)
				currentBuffer = &frameBuffer1
			}

			currentIndex = 0
			start = time.Now()
		}
	}
}
