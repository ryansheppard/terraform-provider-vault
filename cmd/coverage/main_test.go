package main

import (
	"flag"
	"testing"
)

func Test_Main(t *testing.T) {
	pathToDoc := "../../testdata/openapi.json"
	if err := flag.Set("openapi-doc", pathToDoc); err != nil {
		t.Fatal(err)
	}
	main()
}
