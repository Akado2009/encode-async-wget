package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	signalDirectory = "/mnt/scratch/shared/SG_KIRILL/control"
	taskCapacity    = 2
	listFile        = "/home/akado2009/controls.tab"
)

var wg sync.WaitGroup

func main() {
	mu := &sync.Mutex{}
	seen := make(map[string]bool, 0)

	filenamesChannel := make(chan string, taskCapacity)
	freeResources := make(chan struct{}, taskCapacity)

	for i := 0; i < taskCapacity; i++ {
		freeResources <- struct{}{} // fill it initially
	}
	for i := 0; i < taskCapacity; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-freeResources:
					basename := <-filenamesChannel
					basename = filepath.Base(basename)
					basename = strings.TrimSuffix(basename, filepath.Ext(basename))
					mu.Lock()
					if _, ok := seen[basename]; !ok {
						seen[basename] = true
						mu.Unlock()
						input := filepath.Join(
							signalDirectory,
							fmt.Sprintf("%s.bam", basename),
						)
						cmd := exec.Command(
							"bamToWig",
							"-D",
							signalDirectory,
							input,
						)
						log.Println(cmd)
						err := cmd.Run()
						if err != nil {
							log.Print(err)
						}
						if err == nil {
							cmd := exec.Command(
								"rm",
								input,
								fmt.Sprintf(filepath.Join(signalDirectory, "%s.bam.bai", basename)),
								fmt.Sprintf(filepath.Join(signalDirectory, "%s_depth.txt", basename)),
								fmt.Sprintf(filepath.Join(signalDirectory, "%s.bigwig", basename)),
							)
							err := cmd.Run()
							if err != nil {
								log.Print(err)
							}
						}
					} else {
						mu.Unlock()
					}

					freeResources <- struct{}{}
				}
			}
		}()
	}

	file, err := os.Open(listFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		filenamesChannel <- strings.TrimSpace(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
