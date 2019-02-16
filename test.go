package main

import (
	"fmt"

	"github.com/yernurb/saq/tracker"
)

func main() {
	saq := new(tracker.Saq)
	fmt.Println("Starting capture...")
	saq.StartCapture(0)
}
