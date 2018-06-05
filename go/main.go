package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	rootDir   = "huge_dir"
	iteration = 10000
)

func main() {
	if _, err := os.Stat(rootDir); err == nil {
		log.Fatalf("%v already exists, please delete it and try again", rootDir)
	}

	err := os.Mkdir(rootDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	paths := make([]string, iteration, iteration)
	content := []byte(strings.Repeat("A", 100000))
	for i := 0; i < iteration; i++ {
		paths[i] = filepath.Join(rootDir, strconv.Itoa(i)+".txt")
	}

	var writeWg sync.WaitGroup
	writeWg.Add(iteration)

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

	log.Printf("Reading %v files...", iteration)
	start = time.Now()
	for _, path := range paths {
		go func(path string) {
			ioutil.ReadFile(path)
			readWg.Done()
		}(path)
	}
	readWg.Wait()
	elapsed = time.Now().Sub(start)
	log.Println(elapsed)
}
