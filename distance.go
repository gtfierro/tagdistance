package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync"
)

var wg sync.WaitGroup

var Distances = make(map[string](map[string]float64))

////TODO: dont' do all distances twice! Check for key existance
//func calculateDistances() {
//    for p := range patentChannel {
//        patentwg.Add(1)
//        go func() {
//            conn := pool.Get()
//            defer conn.Close()
//            defer patentwg.Done()
//            for _, c := range Patents {
//                distance := p.jaccardDistance(c)
//                conn.Send("HSET", p.number, c.number, distance)
//            }
//            conn.Flush()
//            _, err := conn.Receive()
//            if err != nil {
//              panic(err)
//            }
//        }()
//    }
//    patentwg.Wait()
//}
//

func calculateDistances2() {
	for i := 0; i < 1000; i++ {
		patentwg.Add(1)
		go func() {
			defer patentwg.Done()
			for p := range patentChannel {
				distancemap := make(map[string]float64)
				for _, c := range Patents {
					distancemap[c.number] = p.jaccardDistance(c)
				}
				Distances[p.number] = distancemap
			}
		}()
	}
    patentwg.Wait()
}

//TODO: dont' do all distances twice! Check for key existance
func calculateDistances() {
	for i := 0; i < 1000; i++ {
		patentwg.Add(1)
		go func() {
			conn := pool.Get()
			defer conn.Close()
			defer patentwg.Done()
			for p := range patentChannel {
				fmt.Println(p.number)
				var distance float64
				for _, c := range Patents {
					res, _ := redis.Bool(conn.Do("HEXISTS", c.number, p.number))
					if res {
						distance, _ = redis.Float64(conn.Do("HGET", c.number, p.number))
					} else {
						distance = p.jaccardDistance(c)
					}
					conn.Do("HSET", p.number, c.number, distance)
				}
			}
		}()
	}
	patentwg.Wait()
}

func getDistance(p1, p2 string) float64 {
	r := pool.Get()
	res, err := redis.Float64(r.Do("HGET", p1, p2))
	if err != nil {
		fmt.Println(err)
	}
	return res
}
