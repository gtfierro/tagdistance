package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

var filewg sync.WaitGroup

func readFile(filename string, channel chan []byte) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	filewg.Add(1)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			filewg.Done()
			close(channel)
			break
		}
		channel <- line
	}
	filewg.Wait()
}
