package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	rootDir = "huge_dir"
)

func main() {
	args := os.Args[1:]

	// iteration argument
	if len(args) < 1 {
		log.Fatal("Missing iteration argument")
	}
	iteration, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// parseJSON argument
	parseJSON := false
	if len(args) >= 2 && args[1] == "--parse-json" {
		parseJSON = true
	}

	// Delete the huge_dir if it already exists
	if _, err := os.Stat(rootDir); err == nil {
		err = os.RemoveAll(rootDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create the huge_dir
	err = os.Mkdir(rootDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Populate child paths
	paths := make([]string, iteration, iteration)
	for i := 0; i < iteration; i++ {
		paths[i] = filepath.Join(rootDir, strconv.Itoa(i)+".json")
	}

	// Setup the content for each file
	content, err := ioutil.ReadFile("../common/bench_data.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create the wait group for waiting goroutines
	var writeWg sync.WaitGroup
	writeWg.Add(iteration)

	// Benchmarking write
	log.Printf("Creating %v files...", iteration)

	start := time.Now()
	for _, path := range paths {
		go func(path string) {
			ioutil.WriteFile(path, content, 0755)
			writeWg.Done()
		}(path)
	}
	writeWg.Wait()
	elapsed := time.Now().Sub(start)
	log.Println(elapsed)

	var readWg sync.WaitGroup
	readWg.Add(iteration)

	// Benchmarking read
	if parseJSON {
		log.Printf("Reading and parsing %v files...", iteration)
	} else {
		log.Printf("Reading %v files...", iteration)
	}

	start = time.Now()
	for _, path := range paths {
		go func(path string) {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			if parseJSON {
				var d interface{}
				err := json.Unmarshal(bytes, &d)
				if err != nil {
					log.Fatal(err)
				}
			}

			readWg.Done()
		}(path)
	}
	readWg.Wait()
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)

	log.Print("Completed")
}
