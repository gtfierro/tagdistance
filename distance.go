package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var commit sync.WaitGroup

func calculateDistances(concurrency int) {
	patentwg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(i int) {
			filename := "out/" + strconv.Itoa(i)
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()
			w := bufio.NewWriter(file)
			for p := range patentChannel {
				for _, c := range Patents {
					distance := p.jaccardDistance(c)
					if p.number == c.number {
						continue
					}
					if distance <= .7 {
						fmt.Fprintln(w, PatentMap[p.number], ",", PatentMap[c.number], ",", distance)
					}
				}
			}
			patentwg.Done()
		}(i)
	}
	for _, p := range Patents {
		patentChannel <- p
	}
	close(patentChannel)
	patentwg.Wait()
}

func calculateExternalDistances(concurrency int, filename string) {
	externalFileChannel := make(chan []byte)
	externalPatentChannel := make(chan *Patent)
	externalPatentMap := make(map[int]string)
	patentwg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(i int) {
			filename := "out/" + strconv.Itoa(i)
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
					fmt.Fprintln(w, externalPatentMap[p.number], ",", PatentMap[c.number], ",", distance)
				}
				w.Flush()
			}
			patentwg.Done()
		}(i)
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
