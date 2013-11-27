package main

import (
	"strings"
)

var Tags = make(map[string]int)

func makeTags(filename string) {
	index := 0
	fileChannel := make(chan []byte)
	go readFile(filename, fileChannel)
	for line := range fileChannel {
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
		Patents[number] = makePatent(number, taglist)
	}
}
