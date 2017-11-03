package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ricocheting/logparse/ngparser"
)

func main() {
	f, err := os.Open("junk/access.log")
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

/*
use boltdb to save those stats
run the parser -> get the top 200 -> save that

check and flag the day of the first record
if any record rolls into the next day, write what we have currently processed to the database, clear the stats in memory, then continue parsing the log

If bucket for next day exists, close out "yesterday" and add the tally to the monthly/yearly totals
*/
