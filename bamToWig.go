package main

import (
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
	taskCapacity    = 10
)

//OSReadDir ...
func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

var wg sync.WaitGroup

func main() {
	mu := &sync.Mutex{}
	seen := make(map[string]bool, 0)

	filenamesChannel := make(chan string, taskCapacity)
	freeResources := make(chan struct{}, taskCapacity)

	files, _ := OSReadDir(signalDirectory)
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
					fmt.Println(basename)
					mu.Lock()
					if _, ok := seen[basename]; !ok {
						seen[basename] = true
						mu.Unlock()
						input := filepath.Join(
							signalDirectory,
							fmt.Sprintf("%s.bam", basename),
						)
						output := filepath.Join(
							signalDirectory,
							fmt.Sprintf("%s.wig", basename),
						)
						cmd := exec.Command(
							"bamToWig",
							input,
							output,
						)
						log.Println(cmd)
						err := cmd.Run()
						if err != nil {
							log.Print(err)
						}
						// if err == nil {
						// 	cmd := exec.Command(
						// 		"rm",
						// 		input,
						// fmt.Sprintf("%.bam.bai", basename),
						// fmt.Sprintf("%_depth.txt", basename),
						// 	)
						// 	err := cmd.Run()
						// 	if err != nil {
						// 		log.Print(err)
						// 	}
						// }
					} else {
						mu.Unlock()
					}

					freeResources <- struct{}{}
				}
			}
		}()
	}

	for _, file := range files {
		filenamesChannel <- file
	}
	wg.Wait()
}
