package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ricocheting/logparse/ngparser"
)

func main() {
	logfile := flag.String("log", "", "Log file to process")
	domain := flag.String("domain", "", "Domain for log file. Used for ignoring records. Use format \"example.com\"")
	flag.Parse()

	t1 := time.Now()

	// show what log we're running
	if *domain != "" {
		fmt.Print(*domain + " ") //prepend to date line below
	}
	fmt.Println(t1.Format("2006-01-02 15:04:05"))

	f, err := os.Open(*logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := ngparser.New(*domain)
	p.Parse(f, nil)

	/*p.Parse(f, func(r *ngparser.Record) {
		// if you wanna do some processing to the record
	})*/
	fmt.Printf("Hits: %+v\n", p.Count())
	//fmt.Printf("Top Files: %+v\n", p.Stats(ngparser.Pages, 1000))
	fmt.Printf("Unique Files: %+v\n", p.StatsCount(ngparser.Pages))
	//fmt.Printf("Unique Extentions: %+v\n", p.StatsCount(ngparser.Extensions))
	//fmt.Printf("All Extentions: %+v\n", p.Stats(ngparser.Extensions, 0))

	//fmt.Printf("Unique StatusCodes: %+v\n", p.StatsCount(ngparser.StatusCodes))
	//fmt.Printf("All StatusCodes: %+v\n", p.Stats(ngparser.StatusCodes, 0))

	//fmt.Println(p.IPsCount())

	t2 := time.Now()
	fmt.Printf("Time Run: %d seconds\n\n", int(t2.Sub(t1).Seconds()))
}
