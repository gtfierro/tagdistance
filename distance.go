package main

import (
	"fmt"
	"strings"
	"sync"
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
	patentwg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for p := range externalPatentChannel {
				fmt.Println(p.number, externalPatentMap[p.number])
				for _, c := range Patents {
					distance := p.jaccardDistance(c)
					fmt.Println(externalPatentMap[p.number], PatentMap[c.number], distance)
				}
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

//TODO: dont' do all distances twice! Check for key existance
//func calculateDistances() {
//	for i := 0; i < 500; i++ {
//		patentwg.Add(1)
//		go func() {
//			conn := pool.Get()
//			defer conn.Close()
//			defer patentwg.Done()
//			for p := range patentChannel {
//				fmt.Println(p.number)
//				var distance float64
//				for _, c := range Patents {
//					res, _ := redis.Bool(conn.Do("HEXISTS", c.number, p.number))
//					if res {
//						distance, _ = redis.Float64(conn.Do("HGET", c.number, p.number))
//					} else {
//						distance = p.jaccardDistance(c)
//					}
//					conn.Do("HSET", p.number, c.number, distance)
//				}
//			}
//		}()
//	}
//	patentwg.Wait()
//}

//func getDistance(p1, p2 string) float64 {
//	r := pool.Get()
//	num1 := PatentMap[p1]
//	num2 := PatentMap[p2]
//	res, err := redis.Float64(r.Do("HGET", num1, num2))
//	if err != nil {
//		fmt.Println(err)
//	}
//	return res
//}
