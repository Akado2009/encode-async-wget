// package main

// import (
// 	"log"
// 	"os/exec"
// )

// type WGetTask struct {
// 	Accession string
// 	Type      string
// 	Filename  string
// 	CellLine  string
// }

// const (
// 	dataDir         = "data"
// 	defaultCapacity = 10000
// )

// //GetChipSeqExperiments ...
// func GetChipSeqExperiments() ([]string, error) {
// 	return []string{"https://yandex.ru/"}, nil
// }

// func main() {

// 	taskChannel := make(chan string, defaultCapacity)
// 	successChannel := make(chan string, defaultCapacity)
// 	errorChannel := make(chan string, defaultCapacity)

// 	// Get all ChipSeq experiments
// 	experiments, err := GetChipSeqExperiments()
// 	if err != nil {
// 		log.Fatalf("failed to get all chipseq exps: %v", err)
// 	}

// 	// In goroutine pass WGetTasks into the channel of size ~10-20
// 	// In main -> check if there is a task - spawn a goroutine downloading it
// 	// If err || success - log

// 	go func(exps []string) {
// 		for j := 0; j < 100; j++ {
// 			for i := 0; i < 100; i++ {
// 				taskChannel <- exps[0]
// 			}
// 		}
// 	}(experiments)

// 	go func() {
// 		for {
// 			task := <-taskChannel
// 			cmd := exec.Command(
// 				"wget",
// 				task,
// 			)

// 			_ = cmd.Run()
// 		}
// 	}()
// 	for {
// 		select {
// 		case success := <-successChannel:
// 			log.Printf("Downloaded: %v", success)
// 		case fail := <-errorChannel:
// 			log.Printf("Failed: %v", fail)
// 		}
// 	}
// }
