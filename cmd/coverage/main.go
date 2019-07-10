package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/vault/logical/framework" // TODO wrong dependency I think? anything without Vault?
	"io/ioutil"
	"os"
	"strings"
)

// TODO if an endpoint's only method is LIST, omit it from the comparison
// TODO actually, make this an api endpoint in the plugin? or put a list in the readme or linked there?
// TODO ensure the paths in the inventory match something in openapi but exclude enterprise endpoints
// TODO also make sure they're only covered once and have a match?
// TODO make an exception for GenericPath and UnknownPath
var pathToOpenAPIDoc = flag.String("openapi-doc", "", "path/to/openapi.json")

// This tool is used for generating a coverage reports regarding
// how much of the Vault API can be consumed with the Terraform
// Vault provider.
func main() {
	flag.Parse()
	if pathToOpenAPIDoc == nil || *pathToOpenAPIDoc == "" {
		fmt.Println("'openapi-doc' is required")
		os.Exit(1)
	}
	doc, err := ioutil.ReadFile(*pathToOpenAPIDoc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	oasDoc := &framework.OASDocument{}
	if err := json.NewDecoder(bytes.NewBuffer(doc)).Decode(oasDoc); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for path := range oasDoc.Paths {
		if strings.Contains(path, "pki") && !strings.Contains(path, "auth") && strings.Contains(path, "sign") {
			fmt.Println(path)
		}
	}

	// Compare
	// Output
}
