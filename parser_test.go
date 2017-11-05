package main

// NOTE: you can't run more than one of these at a time as the Bolt database will lock each other. so use
// # go test -run=TestFilenameParse

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

// RAW dump the users in the database
func TestFilenameParse(t *testing.T) {
	fmt.Printf("We're Running\n")
	basename := "/ruby/pokemon/268.shtml"
	ext := filepath.Ext(basename)

	fmt.Printf("extention: " + ext + "\n")
	fmt.Printf("Last: " + filepath.Dir(basename) + "\n")

}

// RAW dump the users in the database
func TestCPU(t *testing.T) {
	fmt.Printf("number of CPU: %+v\n", runtime.NumCPU())
}
