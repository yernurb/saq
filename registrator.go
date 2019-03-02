package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"

	"github.com/yernurb/saq/modules"
)

var (
	fileIndex int
	pid       int
	written   bool
)

func registerVideo() {
	fileName := "video-" + strconv.Itoa(fileIndex) + ".avi"
	cmd := exec.Command("./main", "0", fileName)
	cmd.Start()
	pid = cmd.Process.Pid
	fmt.Println(pid)
	cmd.Wait()
	data, _ := ioutil.ReadFile(fileName)
	modules.EncryptFile(fileName, data, "paSSwoRd")
}

func statusCheckLoop() {
	for {
		if !written {
			fmt.Println(fileIndex)
			written = true
		}
		if fileIndex == 5 {
			cmd := exec.Command("kill", strconv.Itoa(pid))
			cmd.Run()
		}
	}
}

func main() {
	fileIndex = 0
	go statusCheckLoop()
	for {
		written = false
		registerVideo()
		fileIndex++
	}
	/*
		videoData := modules.DecryptFile("video-0.avi", "paSSwoRd")
		f, _ := os.Create("video.avi")
		defer f.Close()
		f.Write(videoData)
	*/
}
