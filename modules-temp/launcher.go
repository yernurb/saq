package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	//	Split the entire command up using ' -' as the delimeter
	parts := strings.Split(`./main 0 video-0.avi`, " ")

	//	The first part is the command, the rest are the args:
	head := parts[0]
	args := parts[1:len(parts)]

	fileIndex := 0

	for {
		//	Format the command
		fileName := "video-" + strconv.Itoa(fileIndex) + ".avi"
		args[1] = fileName
		fileIndex++
		cmd := exec.Command(head, args...)
		//	Sanity check -- just print out the detected args:
		for _, arg := range cmd.Args {
			log.Println(arg)
		}

		//	Sanity check -- capture stdout and stderr:
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		//	Run the command
		cmd.Run()

		//	Output our results
		fmt.Printf("Result: %v / %v", out.String(), stderr.String())
	}
}
