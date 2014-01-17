package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

var filewg sync.WaitGroup

/**
  Given a filename and a channel, reads each line of the file, removes the newline
  and sends the line as a byte slice into the channel. This function is called by
  tags.makeTags
*/
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
		channel <- line[:len(line)-1]
	}
	filewg.Wait()
}
