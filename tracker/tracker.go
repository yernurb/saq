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

// swapBuffer swaps the frames buffer
func (t *Saq) swapBuffer() {
	if t.currentBuffer == &t.frameBuffer1 {
		t.currentBuffer = &t.frameBuffer2
	} else {
		t.currentBuffer = &t.frameBuffer1
	}
}

func (t *Saq) writeBuffer(fileName string, FPS float64) {
	// Prepare video writer object
	fmt.Println("test1")
	t.writer, t.err = gocv.VideoWriterFile(fileName, "MP42", FPS, frameWidth, frameHeight, true)
	if t.err != nil {
		fmt.Printf("error opening video writer device: %v\n", fileName)
		return
	}
	defer t.writer.Close()
	for i := 0; i <= maxFrames; i++ {
		fmt.Println(i)
		if t.currentBuffer == &t.frameBuffer1 {
			t.writer.Write(t.frameBuffer1[i])
		} else {
			t.writer.Write(t.frameBuffer2[i])
		}
	}
	t.writer.Close()
}

func wrtieBuffer(fileName string, FPS float64, buffer [maxFrames]gocv.Mat) {
	writer, err := gocv.VideoWriterFile(fileName, "MP42", FPS, frameWidth, frameHeight, true)
	if err != nil {
		fmt.Printf("error opening video writer device: %v\n", fileName)
		return
	}
	defer writer.Close()
	for i := 0; i <= maxFrames; i++ {
		fmt.Println(i)
		writer.Write(buffer[i])
	}
	writer.Close()
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
	fileIndex := 0
	start := time.Now()
	for {
		t.camera.Read(&t.frame)
		if t.frame.Empty() {
			continue
		}

		pt := image.Pt(30, 30)
		text := time.Now().Format("Monday / _2 January 2006 / 15:04:05")
		//fmt.Println(text)
		gocv.PutText(&t.frame, text, pt, gocv.FontHersheySimplex, 0.6, color.RGBA{255, 0, 0, 0}, 2)
		if t.currentBuffer == &t.frameBuffer1 {
			t.frameBuffer1[t.currentIndex] = t.frame.Clone()
		} else {
			t.frameBuffer2[t.currentIndex] = t.frame.Clone()
		}

		// If current buffer is full swap it and write the collected frames into file
		t.currentIndex++
		if t.currentIndex >= maxFrames {
			elapsed := time.Since(start)
			FPS := float64(maxFrames) / elapsed.Seconds()
			fmt.Println("FPS:", FPS)
			fileIndex++
			fileName := "video-" + strconv.Itoa(fileIndex) + ".avi"
			if t.currentBuffer == &t.frameBuffer1 {
				go writeBuffer(fileName, FPS, t.frameBuffer1)
			} else {
				go writeBuffer(fileName, FPS, t.frameBuffer2)
			}

			t.currentIndex = 0
			t.swapBuffer()
			start = time.Now()
		}
	}
}
