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
	taskCapacity    = 10
	listFile        = "/home/akado2009/controls.tab"
)

var wg sync.WaitGroup

func RunAndWaitForCommand(cmd *exec.Cmd) error {
	var err error
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

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
						tmpWig := filepath.Join(
							signalDirectory,
							fmt.Sprintf("%s.wig", basename),
						)
						bedOutput := filepath.Join(
							signalDirectory,
							fmt.Sprintf("%s.bed", basename),
						)
						// bedtools bamtobed -i  /mnt/scratch/shared/SG_KIRILL/control/ENCFF255TUT.bam > check.wig
						bedtoolsCMD := exec.Command(
							"bedtools",
							"bamtobed",
							"-i",
							input,
							">",
							tmpWig,
						)
						log.Println("Shit")
						log.Println(bedtoolsCMD)
						if err := RunAndWaitForCommand(bedtoolsCMD); err != nil {
							log.Printf("Error is: %v\n", err)
						}
						// awk -F'\t' '{print $1,$2,$3,$5}' check.wig > new.bed
						awkCMD := exec.Command(
							"awk",
							"'{print $1,$2,$3,$5}'",
							tmpWig,
							">",
							bedOutput,
						)
						log.Println(awkCMD)
						if err := RunAndWaitForCommand(bedtoolsCMD); err != nil {
							log.Printf("Error is: %v\n", err)
						} else {
							// rm -rf check.wig
							rmCMD := exec.Command(
								"rm",
								"-rf",
								tmpWig,
							)
							if err := RunAndWaitForCommand(rmCMD); err != nil {
								log.Printf("Error is: %v\n", err)
							}

							freeResources <- struct{}{}
						}
					} else {
						mu.Unlock()
						freeResources <- struct{}{}
					}
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
