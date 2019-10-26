package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const (
	downloadFolder = "/mnt/scratch/shared/SG_KIRILL/controls"
	tableFile      = "encode.files.txt"
)

func main() {
	file, err := os.Open("/path/to/file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := strings.Split(scanner.Text(), "\t")
		if data[0] != "Accession" {

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
