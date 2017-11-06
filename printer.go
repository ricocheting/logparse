package main

import (
	"fmt"
	"path/filepath"

	"github.com/ricocheting/logparse/storage"
)

func main() {

	store := storage.NewStore(filepath.Join("../logparse/data", "db"))
	if err := store.Open(); err != nil {
		panic("Error opening storage (db possibly still open by another process): " + err.Error())
	}
	//fmt.Println(p.IPsCount())
	hits, _ := store.ListHits()

	fmt.Printf("Hits: %+v\n", hits)

}
