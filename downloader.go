package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	downloadDirectory = "/mnt/scratch/shared/SG_KIRILL/samples/"
	inputTable        = "encode.files.txt"
	taskCapacity      = 10
	separator         = "\t"
)

var wg sync.WaitGroup

func main() {

	filenamesChannel := make(chan string, taskCapacity)
	freeResources := make(chan struct{}, taskCapacity)

	file, err := os.Open(inputTable)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < taskCapacity; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-freeResources:
					filename := <-filenamesChannel
					cmd := exec.Command(
						"wget",
						filename,
						"-P",
						downloadDirectory,
					)
					log.Println(cmd)
					err := cmd.Run()
					if err != nil {
						log.Print(err)
					}
					freeResources <- struct{}{}
				}
			}
		}()
	}

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

	wg.Wait()
	// read table
	// store tasks
	// download everything
}
