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
	Directories internal.StatMonth
	StatusCodes internal.StatMonth
	DateCreated string
	Domain      string
}

type PageYear struct {
	Hits        internal.StatTotal
	Pages       internal.StatTotal
	DateCreated string
	Domain      string
}

type PageError struct {
	Errors      []internal.StatErrorPage
	DateCreated string
	Domain      string
}

func main() {
	templateFolder := flag.String("templates", "templates/", "Template folder. Include trailing slash")
	outFolder := flag.String("out", "http/", "Output folder. Include trailing slash")
	domain := flag.String("domain", internal.DefaultDomain, "Domain for log file. Used for ignoring records. Use format \"example.com\"")
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
	directories, _ := store.ListBaseStats(internal.DirectoriesBucket)
	statusCodes, _ := store.ListBaseStats(internal.StatusCodesBucket)
	errors, _ := store.ListErrors(internal.ErrorsBucket)
	page.DateCreated = time.Now().Format("Mon Jan _2 15:04:05 2006")
	page.Domain = *domain

	//fmt.Printf("Hits: %+v\n", page.Extensions)

	pageYear := PageYear{
		Hits:        hits,
		Pages:       pages,
		DateCreated: page.DateCreated,
		Domain:      *domain,
	}

	pageError := PageError{
		Errors:      errors,
		DateCreated: page.DateCreated,
		Domain:      *domain,
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
		"sortedTopTwenty":  internal.TopTwentyStatsForDay,
	}

	tIndex := template.Must(template.New("index.html").Funcs(fmap).ParseFiles(*templateFolder + "index.html"))
	tMonth := template.Must(template.New("month.html").Funcs(fmap).ParseFiles(*templateFolder + "month.html"))
	tErrors := template.Must(template.New("errors.html").Funcs(fmap).ParseFiles(*templateFolder + "errors.html"))

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

			//if *directories.Years[year] && *directories.Years[year].Months[month] {
			_, ok := directories.Years[year]

			if ok {
				_, ok = directories.Years[year].Months[month]

				if ok {
					page.Directories = *directories.Years[year].Months[month]
				} else {
					page.Directories = internal.StatMonth{}
				}
			} else {
				page.Directories = internal.StatMonth{}
			}

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

	// Write the error file
	file, err := os.Create(*outFolder + "errors.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := tErrors.ExecuteTemplate(file, "errors.html", pageError); err != nil {
		fmt.Println(err)
	}

	/*for page, missing := range errors.Page {

		// walk through the new data
		for key, value := range missing {
			fmt.Printf("%d %s %s\n", value, page, key)
		}

	}*/

	// Write the main index year file
	file, err = os.Create(*outFolder + "index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := tIndex.ExecuteTemplate(file, "index.html", pageYear); err != nil {
		fmt.Println(err)
	}

}
