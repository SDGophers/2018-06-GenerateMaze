package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sdgophers/2018-06-GenerateMaze/maze"
)

var mapName = "maze.map"

const Usage = "Usage: %v mapfilename.map\n"

func usage(err error) {
	fmt.Printf(Usage, os.Args[0])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	args := flag.Args()

	// expect the first argument to be the file to display, otherwise
	if len(args) > 0 {
		mapName = args[0]
	}

	f, err := os.Open(mapName)
	if err != nil {
		usage(err)
		return
	}
	m, err := maze.ReadMap(f)
	if err != nil {
		usage(err)
		return
	}
	fmt.Println(m)
}
