package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	downloadDirectory = "/mnt/scratch/shared/SG_KIRILL/"
	inputTable        = "encode.files.txt"
	taskCapacity      = 10
	separator         = "\t"
)

func main() {

	filenamesChannel := make(chan string, taskCapacity)
	freeResources := make(chan struct{}, taskCapacity)

	file, err := os.Open(inputTable)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	go func() {
		for {
			select {
			case <-freeResources:
				filename := <-filenamesChannel
				// download
				go func(fname string) {
					cmd := exec.Command(
						"wget",
						fname,
						"-P",
						downloadDirectory,
					)
					log.Println(cmd)
					err := cmd.Run()
					if err != nil {
						log.Print(err)
					}
					freeResources <- struct{}{}

				}(filename)
			}
		}
	}()

	for i := 0; i < taskCapacity; i++ {
		freeResources <- struct{}{} // fill it initially
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), separator)
		if line[0] != "Accession" {
			filenamesChannel <- line[6]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// read table
	// store tasks
	// download everything
}
