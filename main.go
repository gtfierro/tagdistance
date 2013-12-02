package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	flag.Parse()
    f, err := os.Create("cprof")
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
	runtime.GOMAXPROCS(runtime.NumCPU())
	filename := flag.Arg(0)
	external := flag.Arg(1)
	makeTags(filename)
	fmt.Println(len(Tags))
	fmt.Println(len(Patents))
	if external != "" {
		calculateExternalDistances(external)
	} else {
		calculateDistances()
	}
    f, err = os.Create("mprof")
    if err != nil {
        log.Fatal(err)
    }
    pprof.WriteHeapProfile(f)
    f.Close()
}
