package main

import (
	"flag"
	"gosloc/gosloc"
	"log"
)

func main() {
	// flags
	dirPtr := flag.String("dir", ".", "dir path")
	savePtr := flag.String("save", "", "file name to save")
	dispPtr := flag.Bool("disp", true, "to display or not on stdout")
	extsPtr := flag.String("ext", "", "',' separated extensions to allow")
	flag.Parse()

	gsloc := gosloc.GoSLOC{}
	err := gsloc.Read(*dirPtr)
	if err != nil {
		log.Fatal(err)
	}

	err = gsloc.Process(*extsPtr)
	if err != nil {
		log.Fatal(err)
	}

	err = gsloc.SaveOrDisplay(*savePtr, *dispPtr)
	if err != nil {
		log.Println(err)
	}
}
