package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/ricocheting/logparse/internal"
	"github.com/ricocheting/logparse/storage"
)

type Page struct {
	Hits        []internal.Stat
	IPS         []internal.Stat
	Extensions  internal.StatCollection
	StatusCodes internal.StatCollection
}

func main() {

	store := storage.NewStore(filepath.Join("data", "db"))
	if err := store.Open(); err != nil {
		panic("Error opening storage (db possibly still open by another process): " + err.Error())
	}
	page := Page{}

	page.Hits, _ = store.ListBaseNumber(internal.HitsBucket)
	page.IPS, _ = store.ListBaseNumber(internal.IPSBucket)
	page.Extensions, _ = store.ListBaseStats(internal.ExtensionsBucket)
	page.StatusCodes, _ = store.ListBaseStats(internal.StatusCodesBucket)

	//fmt.Printf("Hits: %+v\n", page.Extensions)

	fmap := template.FuncMap{
		"formatDate":      internal.FormatShortDate,
		"formatCommas":    internal.FormatCommas,
		"formatShortHand": internal.FormatShortHand,
	}
	t := template.Must(template.New("index.html").Funcs(fmap).ParseFiles("templates/index.html"))

	// Write the file
	file, err := os.Create("logs/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := t.ExecuteTemplate(file, "index.html", page); err != nil {
		fmt.Println(err)
	}

}
