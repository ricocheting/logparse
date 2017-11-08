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

type Page struct {
	Hits        []internal.Stat
	IPS         []internal.Stat
	Pages       []internal.Stat
	Extensions  internal.StatCollection
	StatusCodes internal.StatCollection
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
	page := Page{}

	page.Hits, _ = store.ListBaseNumber(internal.HitsBucket)
	page.IPS, _ = store.ListBaseNumber(internal.IPSBucket)
	page.Pages, _ = store.ListPages(internal.ExtensionsBucket)
	page.Extensions, _ = store.ListBaseStats(internal.ExtensionsBucket)
	page.StatusCodes, _ = store.ListBaseStats(internal.StatusCodesBucket)
	page.DateCreated = time.Now().Format("Mon Jan _2 15:04:05 2006")

	//fmt.Printf("Hits: %+v\n", page.Extensions)

	fmap := template.FuncMap{
		"formatDate":       internal.FormatShortDate,
		"formatCommas":     internal.FormatCommas,
		"formatShortHand":  internal.FormatShortHand,
		"formatStatusCode": internal.FormatStatusCodeName,
	}
	t := template.Must(template.New("index.html").Funcs(fmap).ParseFiles(*templateFolder + "index.html"))

	// Write the file
	file, err := os.Create(*outFolder + "index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := t.ExecuteTemplate(file, "index.html", page); err != nil {
		fmt.Println(err)
	}

}
