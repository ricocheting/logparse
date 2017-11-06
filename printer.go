package main

import (
	"html/template"
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

	var t = template.Must(template.ParseFiles("templates/main.html"))

	t.Execute(os.Stdout, page)

}
