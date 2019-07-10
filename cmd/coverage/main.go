package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/hashicorp/vault/logical/framework"
	"github.com/terraform-providers/terraform-provider-vault/vault"
)

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
	// This is the path, to whether they've been observed in the OpenAPI doc
	vaultPaths := make(map[string]bool)
	for path := range oasDoc.Paths {
		vaultPaths[path] = false
	}

	for _, desc := range vault.DataSourceRegistry {
		for _, path := range desc.PathInventory {
			if path == vault.GenericPath || path == vault.UnknownPath {
				continue
			}
			seenBefore, isCurrentlyInVault := vaultPaths[path]
			if !isCurrentlyInVault && !desc.EnterpriseOnly {
				fmt.Println(path + " is not currently in Vault")
			}
			if seenBefore {
				fmt.Println(path + " is in the Terraform Vault Provider multiple times")
			}
			vaultPaths[path] = true
		}
	}

	for _, desc := range vault.ResourceRegistry {
		for _, path := range desc.PathInventory {
			if path == vault.GenericPath || path == vault.UnknownPath {
				continue
			}
			seenBefore, isCurrentlyInVault := vaultPaths[path]
			if !isCurrentlyInVault && !desc.EnterpriseOnly {
				fmt.Println(path + " is not currently in Vault")
			}
			if seenBefore {
				fmt.Println(path + " is in the Terraform Vault Provider multiple times")
			}
			vaultPaths[path] = true
		}
	}

	supportedVaultEndpoints := []string{}
	unSupportedVaultEndpoints := []string{}
	for path, seen := range vaultPaths {
		if seen {
			supportedVaultEndpoints = append(supportedVaultEndpoints, path)
		} else {
			unSupportedVaultEndpoints = append(unSupportedVaultEndpoints, path)
		}
	}

	fmt.Println(" ")
	fmt.Printf("%.0f percent coverage\n", float64(len(supportedVaultEndpoints))/float64(len(vaultPaths))*100)
	fmt.Printf("%d of %d vault paths are supported\n", len(supportedVaultEndpoints), len(vaultPaths))
	fmt.Printf("%d of %d vault paths are unsupported\n", len(unSupportedVaultEndpoints), len(vaultPaths))

	fmt.Println(" ")
	fmt.Println("SUPPORTED")
	sort.Strings(supportedVaultEndpoints)
	for _, path := range supportedVaultEndpoints {
		fmt.Println("    " + path)
	}

	fmt.Println(" ")
	fmt.Println("UNSUPPORTED")
	sort.Strings(unSupportedVaultEndpoints)
	for _, path := range unSupportedVaultEndpoints {
		fmt.Println("    " + path)
	}
}
