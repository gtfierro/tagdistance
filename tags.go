package main

import (
	_ "fmt"
	"strings"
)

var Tags = make(map[string]int)
var PatentMap = make(map[int]string)

func makeTags(filename string) {
	index := 0
	linecount := 0
	fileChannel := make(chan []byte)
	go readFile(filename, fileChannel)
	for line := range fileChannel {
		linecount += 1
		parsed := strings.Split(string(line), ",")
		number := parsed[0]
		tagline := parsed[1]
		tags := strings.Split(tagline, " ")
		taglist := make([]int, len(tags), len(tags))
		for i, tag := range tags {
			if Tags[tag] == 0 {
				Tags[tag] = index
				index += 1
			}
			taglist[i] = Tags[tag]
		}
		Patents = append(Patents, makePatent(linecount, taglist))
		//Patents[number] = makePatent(linecount, taglist)
		PatentMap[linecount] = number
	}
	commit.Add(len(Patents))
}
