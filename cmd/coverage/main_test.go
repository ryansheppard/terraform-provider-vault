package main

import (
	"flag"
	"testing"
)

func Test_Main(t *testing.T) {
	pathToDoc := "/home/tbex/go/src/github.com/terraform-providers/terraform-provider-vault/testdata/openapi.json"
	if err := flag.Set("openapi-doc", pathToDoc); err != nil {
		t.Fatal(err)
	}
	main()
}
