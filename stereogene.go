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

const (
	simultaniousPrograms = 50
	signalDirectory      = "/mnt/scratch/shared/SG_KIRILL/samples"
	samplesTable         = "encode.files.txt"
	resultDir            = "/mnt/scratch/shared/SG_KIRILL/results"
	separator            = "\t"
)

var wg sync.WaitGroup

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
		if strings.Contains(file.Name(), "bgraph") {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	mu := &sync.Mutex{}
	seen := make(map[string]bool, 0)

	runsChannel := make(chan []Sample, simultaniousPrograms)
	runsResources := make(chan struct{}, simultaniousPrograms)

	runsErrors := make(chan ErrorReport, simultaniousPrograms)
	// runsExchange := make(chan []Sample, simultaniousPrograms)

	// runsPythonErrors := make(chan error, simultaniousPrograms)

	preprocessedFiles, err := OSReadDir(signalDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// Runs Stereogene
	for i := 0; i < simultaniousPrograms; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-runsResources:
					pair := <-runsChannel
					fId := pair[0]
					sId := pair[1]
					key := fmt.Sprintf("%s_%s", fId.Accession, sId.Accession)
					// Generate & run
					mu.Lock()
					if _, ok := seen[key]; !ok {
						seen[key] = true
						mu.Unlock()

						// Create dir
						dir := filepath.Join(
							resultDir,
							fmt.Sprintf(
								"%s_%s", fId.Accession, sId.Accession,
							),
						)
						// statFile := fmt.Sprintf(
						// 	"%s_%s_statistics",
						// 	fId.Accession,
						// 	sId.Accession,
						// )
						statFile := filepath.Join(
							dir,
							fmt.Sprintf(
								"%s_%s_statistics",
								fId.Accession,
								sId.Accession,
							),
						)
						// paramFile := fmt.Sprintf(
						// 	"%s_%s_parameters",
						// 	fId.Accession,
						// 	sId.Accession,
						// )
						paramFile := filepath.Join(
							dir,
							fmt.Sprintf(
								"%s_%s_parameters",
								fId.Accession,
								sId.Accession,
							),
						)
						_ = os.MkdirAll(dir, os.ModePerm)
						cmd := exec.Command(
							"StereoGene",
							"chrom=/home/akado2009/human_chrom",
							// window
							fmt.Sprintf("resPath=%s", dir),
							"wSize=300000",
							// kernel
							"kernelSigma=1000",
							// statistics
							// parameters
							fmt.Sprintf(
								"params=%s", paramFile,
							),
							fmt.Sprintf(
								"statistics=%s", statFile,
							),
							fId.Link,
							sId.Link,
						)

						out, err := cmd.CombinedOutput()
						runsErrors <- ErrorReport{
							Err:  err,
							Data: string(out),
						}
						// runsExchange <- pair
					} else {
						mu.Unlock()
					}
					runsResources <- struct{}{}
				}
			}
		}()
	}

	// Processes output & run python
	for i := 0; i < simultaniousPrograms; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case err := <-runsErrors:
					// task := <-runsExchange
					if err.Err != nil {
						log.Fatalf("err in %+v\n", err.Data)
					}
				}

			}
		}()
	}

	// ParsesTable
	sampleSlice := make([]Sample, 0)
	file, err := os.Open(samplesTable)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	for i := 0; i < simultaniousPrograms; i++ {
		runsResources <- struct{}{}
	}
	cointCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), separator)
		if line[0] != "Accession" && stringInSlice(fmt.Sprintf("%s.bgraph", line[0]), preprocessedFiles) {
			cl := line[3]
			cointCount++
			if cl == "." {
				if line[4] != "." {
					sampleSlice = append(sampleSlice, Sample{
						Accession: line[0],
						Dataset:   line[1],
						Tissue:    line[2],
						CellLine:  line[4],
						Link:      filepath.Join(signalDirectory, fmt.Sprintf("%s.bgraph", line[0])),
						Feature:   line[5],
					})
				}
			} else {
				sampleSlice = append(sampleSlice, Sample{
					Accession: line[0],
					Dataset:   line[1],
					Tissue:    line[2],
					CellLine:  cl,
					Link:      filepath.Join(signalDirectory, fmt.Sprintf("%s.bgraph", line[0])),
					Feature:   line[5],
				})
			}
		}
	}
	log.Println("cointCount", cointCount)
	log.Println("len", len(sampleSlice))
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(sampleSlice)-1; i++ {
		for j := i + 1; j < len(sampleSlice); j++ {
			// s := make([]Sample, 2, 2)
			// s[0] = sampleSlice[i]
			// s[1] = sampleSlice[j]
			// runsChannel <- s
			if sampleSlice[i].CellLine == sampleSlice[j].CellLine {
				s := make([]Sample, 2, 2)
				s[0] = sampleSlice[i]
				s[1] = sampleSlice[j]
				runsChannel <- s
				// pairsCount++
			}
			if sampleSlice[i].Feature == sampleSlice[j].Feature {
				s := make([]Sample, 2, 2)
				s[0] = sampleSlice[i]
				s[1] = sampleSlice[j]
				// pairsCount++
			}
		}
	}
	wg.Wait()
}

// Pass pairs ID1_ID2
