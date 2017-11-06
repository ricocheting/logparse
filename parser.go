package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ricocheting/logparse/ngparser"
)

func main() {
	f, err := os.Open("junk/short.access.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := ngparser.New()
	p.Parse(f, nil)

	/*p.Parse(f, func(r *ngparser.Record) {
		// if you wanna do some processing to the record
	})*/

	fmt.Printf("Hits: %+v\n", p.Count())
	//fmt.Printf("Top Files: %+v\n", p.Stats(ngparser.Pages, 1000))
	fmt.Printf("Unique Files: %+v\n", p.StatsCount(ngparser.Pages))
	fmt.Printf("Unique Extentions: %+v\n", p.StatsCount(ngparser.Extensions))
	fmt.Printf("All Extentions: %+v\n", p.Stats(ngparser.Extensions, 0))

	fmt.Printf("Unique StatusCodes: %+v\n", p.StatsCount(ngparser.StatusCodes))
	fmt.Printf("All StatusCodes: %+v\n", p.Stats(ngparser.StatusCodes, 0))

	//fmt.Println(p.IPsCount())
}
