package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/ricocheting/logparse/internal"
	"github.com/ricocheting/logparse/storage"
)

func main() {
	domain := flag.String("domain", "", "Domain for log file. Used for ignoring records. Use format \"example.com\"")
	resetErrors := flag.Bool("errors", false, "Reset error log data") // -errors
	flag.Parse()

	t1 := time.Now()

	// show what log we're running
	if *domain != "" {
		fmt.Print(*domain + " ") //prepend to date line below
	}
	fmt.Println(t1.Format("2006-01-02 15:04:05"))

	// runs the reset then stops processing
	if *resetErrors == true {
		store := storage.NewStore(filepath.Join("data", "db"))
		if err := store.Open(); err != nil {
			panic("Error opening storage (db possibly still open by another process): " + err.Error())
		}

		store.ResetErrors()

		fmt.Printf("Cleared %s bucket\n", internal.ErrorsBucket)
		return
	}

	//t2 := time.Now()
	//fmt.Printf("Time Run: %d seconds\n\n", int(t2.Sub(t1).Seconds()))
}
