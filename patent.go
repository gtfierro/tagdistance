package main

import (
	"sync"
)

//var Patents = make(map[string](*Patent))
var Patents = [](*Patent){}
var patentChannel = make(chan *Patent)
var patentwg sync.WaitGroup

type Patent struct {
	number int
	tags   []int
}

func makePatent(number int, taglist []int) *Patent {
	p := &Patent{number: number, tags: taglist}
	return p
}

func (p *Patent) jaccardDistance(c *Patent) float64 {
	var intersection, union float64
	intersection = 0
	union = float64(len(p.tags))
	found := false
	for _, ctag := range c.tags {
		for _, ptag := range p.tags {
			if ctag == ptag {
				intersection += 1
				found = true
				break
			}
		}
		if !found {
			union += 1
			found = false
		}
	}
	return intersection / union
}
