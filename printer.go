package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ricocheting/logparse/internal"
	"github.com/ricocheting/logparse/storage"
)

type PageMonth struct {
	Hits        internal.StatMonth
	IPS         internal.StatMonth
	Pages       internal.StatMonth
	Extensions  internal.StatMonth
	StatusCodes internal.StatMonth
	DateCreated string
}

type PageYear struct { //[YYYY][MM]
	Hits        internal.StatTotal
	Pages       internal.StatTotal
	DateCreated string
}

func main() {
	templateFolder := flag.String("templates", "templates/", "Template folder. Include trailing slash")
	outFolder := flag.String("out", "http/", "Output folder. Include trailing slash")
	flag.Parse()

	store := storage.NewStore(filepath.Join("data", "db"))
	if err := store.Open(); err != nil {
		panic("Error opening storage (db possibly still open by another process): " + err.Error())
	}
	page := PageMonth{}

	hits, _ := store.ListBaseNumber(internal.HitsBucket)
	ips, _ := store.ListBaseNumber(internal.IPSBucket)
	pages, _ := store.ListPages(internal.ExtensionsBucket)
	extensions, _ := store.ListBaseStats(internal.ExtensionsBucket)
	statusCodes, _ := store.ListBaseStats(internal.StatusCodesBucket)
	page.DateCreated = time.Now().Format("Mon Jan _2 15:04:05 2006")

	//fmt.Printf("Hits: %+v\n", page.Extensions)

	pageYear := PageYear{
		Hits:        hits,
		Pages:       pages,
		DateCreated: page.DateCreated,
	}

	fmap := template.FuncMap{
		"formatDate":       internal.FormatShortDate,
		"formatCommas":     internal.FormatCommas,
		"formatShortHand":  internal.FormatShortHand,
		"formatStatusCode": internal.FormatStatusCodeName,
		"formatMonth":      internal.FormatMonth,
		"formatShortMonth": internal.FormatShortMonth,
		"formatMonthYear":  internal.FormatMonthYear,
		"pathDirectory":    internal.PathDirectory,
		"pathFilename":     internal.PathFilename,
	}

	tIndex := template.Must(template.New("index.html").Funcs(fmap).ParseFiles(*templateFolder + "index.html"))
	tMonth := template.Must(template.New("month.html").Funcs(fmap).ParseFiles(*templateFolder + "month.html"))

	// Write the month files
	for year, yearData := range hits.Years {
		// create year folders
		pathname := *outFolder + internal.PathDirectory(year) + "/"
		os.MkdirAll(pathname, os.ModePerm)

		for month, monthData := range yearData.Months {
			page.Hits = *monthData
			page.IPS = *ips.Years[year].Months[month]
			page.Pages = *pages.Years[year].Months[month]
			page.Extensions = *extensions.Years[year].Months[month]
			page.StatusCodes = *statusCodes.Years[year].Months[month]

			filename := internal.PathFilename(month)

			// Write the file
			file, err := os.Create(pathname + filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			if err := tMonth.ExecuteTemplate(file, "month.html", page); err != nil {
				fmt.Println(err)
			}
		}
	}

	// Write the main index year file
	file, err := os.Create(*outFolder + "index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := tIndex.ExecuteTemplate(file, "index.html", pageYear); err != nil {
		fmt.Println(err)
	}

}
