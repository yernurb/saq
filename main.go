package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

// Adds "0" in front of single digit val, for example "7" becomes "07". Assumes 0 <= val <= 60
func addZeroTime(val int) string {
	if val < 10 {
		return "0" + strconv.Itoa(val)
	}
	return strconv.Itoa(val)
}

// Return current time as: "{Weekday} / {Year} {Month} {Day} / {Hour}:{Minute}:{Second}"
func textifyTime() string {
	year, month, day := time.Now().Date()
	weekday := time.Now().Weekday()
	hour, min, sec := time.Now().Clock()
	hourText := addZeroTime(hour)
	minText := addZeroTime(min)
	secText := addZeroTime(sec)
	return weekday.String() + " / " + strconv.Itoa(year) + " " + month.String() + " " + strconv.Itoa(day) + " / " + hourText + ":" + minText + ":" + secText
}

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

	webcam.Set(gocv.VideoCaptureFrameWidth, 640)
	webcam.Set(gocv.VideoCaptureFrameHeight, 480)

	// Prepare image container matrix and an array of images
	img := gocv.NewMat()
	defer img.Close()
	var imgArray [maxFrames]gocv.Mat

	if ok := webcam.Read(&img); !ok {
		fmt.Printf("Cannot read device %v\n", deviceID)
		return
	}
	fmt.Println(img.Cols(), img.Rows())

	// Check FPS

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

		currentTime := textifyTime()
		fmt.Println(currentTime)
		pt := image.Pt(30, 30)
		gocv.PutText(&img, currentTime, pt, gocv.FontHersheySimplex, 0.6, color.RGBA{255, 0, 0, 0}, 2)
		imgArray[i] = img.Clone()
	}
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
