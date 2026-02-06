package main

import (
	"encbench/internal/bench"
	"flag"
	"runtime"
)

func main() {
	all := flag.Bool("all", false, "include slow compressions, serializations and benchmarks")
	flag.Parse()

	runtime.GOMAXPROCS(1)
	bench.Run(*all)
}
