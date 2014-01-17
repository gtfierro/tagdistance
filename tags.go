package main

import (
	_ "fmt"
	"strings"
)

var Tags = make(map[string]int)
var PatentMap = make(map[int]string)

/**
  Accepts filename pointing to CSV file with structure:

    document ID,space separated list of tags/terms

  iterates through each of the lines, and creates:
  * Tags: a map[string]int that maps each unique term to an integer
          Allows us to save memory by representing the list of tags
          for a patent just using integers instead of strings
  * Patents: a slice of patent structs
  * PatentMap: a map[int]string that acts as a lookup table for patent
          document numbers. We can have an integer be a pointer to a patent
          instead of using the longer patent string.
  * adds to the `commit` sync.WaitGroup a number of counters equal to the
    number of patent documents read in the input file

  Tags, Patents, PatentMap are global variables

  Nothing is returned by this function.
*/
func makeTags(filename string) {
	index := 1
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
		PatentMap[linecount] = number
	}
	commit.Add(len(Patents))
}
