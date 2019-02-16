package tracker

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

// AddZeroTime adds "0" in front of single digit val, for example "7" becomes "07". Assumes 0 <= val <= 60
func AddZeroTime(val int) string {
	if val < 10 {
		return "0" + strconv.Itoa(val)
	}
	return strconv.Itoa(val)
}

// TextifyTime returns current time as: "{Weekday} / {Year} {Month} {Day} / {Hour}:{Minute}:{Second}"
func TextifyTime() string {
	year, month, day := time.Now().Date()
	weekday := time.Now().Weekday()
	hour, min, sec := time.Now().Clock()
	hourText := AddZeroTime(hour)
	minText := AddZeroTime(min)
	secText := AddZeroTime(sec)
	return weekday.String() + " / " + strconv.Itoa(year) + " " + month.String() + " " + strconv.Itoa(day) + " / " + hourText + ":" + minText + ":" + secText
}

// Saq structure with corresponding methods for video capture
type Saq struct {
	camera        *gocv.VideoCapture
	writer        *gocv.VideoWriter
	frameBuffer1  [maxFrames]gocv.Mat
	frameBuffer2  [maxFrames]gocv.Mat
	currentBuffer *[maxFrames]gocv.Mat
	frame         gocv.Mat
	currentIndex  int32
	err           error
}

// SwapBuffer swaps the frames buffer
func (t *Saq) SwapBuffer() {
	if t.currentBuffer == &t.frameBuffer1 {
		t.currentBuffer = &t.frameBuffer2
	} else {
		t.currentBuffer = &t.frameBuffer1
	}
}

// StartCapture starts the main infinite loop, capturing and inserting each into buffers
func (t *Saq) StartCapture(deviceID int) {
	// Preparing the camera module
	t.camera, t.err = gocv.OpenVideoCapture(deviceID)
	if t.err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer t.camera.Close()

	// Preparing frame container
	t.frame = gocv.NewMat()
	defer t.frame.Close()

	// Test read from the camera module
	t.camera.Read(&t.frame)
	fmt.Println(t.frame.Cols(), t.frame.Rows())

	// Setting the camera resolution
	t.camera.Set(gocv.VideoCaptureFrameWidth, frameWidth)
	t.camera.Set(gocv.VideoCaptureFrameHeight, frameHeight)

	// Setting the initial frame buffer pointer
	t.currentBuffer = &t.frameBuffer1
	t.currentIndex = 0

	// Starting the main loop
	start := time.Now()
	for {
		t.camera.Read(&t.frame)
		if t.frame.Empty() {
			continue
		}
		currentTime := TextifyTime()
		//fmt.Println(currentTime)
		pt := image.Pt(30, 30)
		gocv.PutText(&t.frame, currentTime, pt, gocv.FontHersheySimplex, 0.6, color.RGBA{255, 0, 0, 0}, 2)
		t.currentBuffer[t.currentIndex] = t.frame.Clone()

		// If current buffer is full swap it and write the collected frames into file
		t.currentIndex++
		if t.currentIndex >= maxFrames {
			t.currentIndex = 0
			t.SwapBuffer()
			elapsed := time.Since(start)
			FPS := float64(maxFrames) / elapsed.Seconds()
			fmt.Println("FPS:", FPS)
			start = time.Now()
		}
	}
}
