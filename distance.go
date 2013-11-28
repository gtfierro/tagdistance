package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync"
)

var wg sync.WaitGroup
var commit sync.WaitGroup

func calculateDistances() {
	patentwg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for p := range patentChannel {
				distancemap := make(map[int]float64)
				for _, c := range Patents {
					distance := p.jaccardDistance(c)
					if distance >= .5 {
						distancemap[c.number] = distance
					}
				}
				go commitDistanceMap(p.number, distancemap)
			}
			patentwg.Done()
		}()
	}
	for _, p := range Patents {
		patentChannel <- p
	}
	close(patentChannel)
	patentwg.Wait()
	commit.Wait()
}

func commitDistanceMap(p int, dm map[int]float64) {
	conn := pool.Get()
	defer conn.Close()
	defer commit.Done()
	for num, dist := range dm {
		conn.Do("HSET", p, num, dist)
	}
	fmt.Println(p)
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

func getDistance(p1, p2 string) float64 {
	r := pool.Get()
	num1 := PatentMap[p1]
	num2 := PatentMap[p2]
	res, err := redis.Float64(r.Do("HGET", num1, num2))
	if err != nil {
		fmt.Println(err)
	}
	return res
}
