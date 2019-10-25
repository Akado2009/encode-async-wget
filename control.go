package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	baseDownloadURL   = "https://www.encodeproject.org%s"
	baseURL           = "https://www.encodeproject.org%s?format=json"
	downloadDirectory = "/mnt/scratch/shared/SG_KIRILL/control"
	inputTable        = "encode.files.txt"
	taskCapacity      = 10
	separator         = "\t"
)

var wg sync.WaitGroup

func getLink(urlPart string) string {
	fDesc := &FileDescription{}
	response, err := http.Get(
		fmt.Sprintf(
			baseURL,
			urlPart,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, fDesc)
	if fDesc.FileType == "bam" {
		return fDesc.Href
	}
	return ""
}
func main() {
	mu := &sync.Mutex{}
	seen := make(map[string]bool, 0)

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
					urlPart := <-filenamesChannel

					link := getLink(urlPart)
					if link != "" {
						mu.Lock()
						if _, ok := seen[link]; !ok {
							seen[link] = true
							mu.Unlock()
							finalLink := fmt.Sprintf(baseDownloadURL, link)
							cmd := exec.Command(
								"wget",
								finalLink,
								"-P",
								downloadDirectory,
							)
							log.Println(cmd)
							err := cmd.Run()
							if err != nil {
								log.Print(err)
							}
						}
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
			controls := strings.Split(line[8], ",")
			for _, c := range controls {
				filenamesChannel <- strings.TrimSpace(c)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
