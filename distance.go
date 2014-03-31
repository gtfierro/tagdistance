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

/**
*/
func calculateDistances(concurrency int) {
	patentwg.Add(concurrency)
	os.RemoveAll("out")
	os.Mkdir("out", 0777)
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
            defer w.Flush()
			for {
				p, more := <-patentChannel
				if !more {
					break
				}
				for _, c := range Patents[p.number:] {
					if p.number == c.number {
						continue
					}
					distance := p.jaccardDistance(c)
					if distance <= 2.0 {
						continue
					}
					fmt.Fprintln(w, PatentMap[p.number], ",", PatentMap[c.number], ",", distance)
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
	os.RemoveAll("out")
	os.Mkdir("out", 0777)
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
            defer w.Flush()
			for {
				p, more := <-externalPatentChannel
				if !more {
					break
				}
				for _, c := range Patents {
					if p.number == c.number {
						continue
					}
					distance := p.jaccardDistance(c)
					if distance == 0.0 {
						continue
					}
					_, err = fmt.Fprintln(w, externalPatentMap[p.number], ",", PatentMap[c.number], ",", distance)
                    if err != nil {
                      fmt.Println(err)
                      fmt.Println(w)
                      os.Exit(1)
                    }
				}
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
    fmt.Println("gothere")
	close(externalPatentChannel)
	patentwg.Wait()
}
