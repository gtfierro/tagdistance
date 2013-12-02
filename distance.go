package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var commit sync.WaitGroup

func calculateDistances() {
	patentwg.Add(50000)
	for i := 0; i < 50000; i++ {
		go func() {
			for p := range patentChannel {
				for _, c := range Patents {
					distance := p.jaccardDistance(c)
					if p.number == c.number {
						continue
					}
					if distance <= .7 {
						fmt.Println(PatentMap[p.number], PatentMap[c.number], distance)
					}
				}
			}
			patentwg.Done()
		}()
	}
	for _, p := range Patents {
		patentChannel <- p
	}
	close(patentChannel)
	patentwg.Wait()
}

func calculateExternalDistances(filename string) {
	externalFileChannel := make(chan []byte)
	externalPatentChannel := make(chan *Patent)
	externalPatentMap := make(map[int]string)
	patentwg.Add(5)
	fileindex := int32(0)
	for i := 0; i < 5; i++ {
		go func() {
			fileindex := atomic.AddInt32(&fileindex, 1)
			filename := "out/" + strconv.Itoa(int(fileindex))
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()
			w := bufio.NewWriter(file)
			for p := range externalPatentChannel {
				for _, c := range Patents {
					distance := p.jaccardDistance(c)
					if p.number == c.number {
						continue
					}
					fmt.Fprint(w, externalPatentMap[p.number], PatentMap[c.number], distance)
				}
				w.Flush()
			}
			patentwg.Done()
		}()
	}
	go readFile(filename, externalFileChannel)
	linecount := 0
	for line := range externalFileChannel {
		linecount += 1
		parsed := strings.Split(string(line), ",")
		number := parsed[0]
		tagline := parsed[1]
		tags := strings.Split(tagline, " ")
		taglist := make([]int, len(tags), len(tags))
		for i, tag := range tags {
			taglist[i] = Tags[tag]
		}
		externalPatentChannel <- makePatent(linecount, taglist)
		externalPatentMap[linecount] = number
	}
	close(externalPatentChannel)
	patentwg.Wait()
}
