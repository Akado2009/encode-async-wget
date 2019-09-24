package main

import "log"

func main() {
	if err := LoadConfig(&AppConfig); err != nil {
		log.Fatalf("Error loading a config file: %v", err)
	}

	//read table
	//find needed type (signal)
}
