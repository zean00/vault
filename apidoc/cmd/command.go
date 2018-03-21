package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/apidoc/apidoc"
	"github.com/hashicorp/vault/builtin/logical/aws"
	"github.com/hashicorp/vault/vault"
)

func Run() int {
	doc := apidoc.NewDoc()

	// we can choose to build different things at this point
	buildDoc(doc)

	// we can choose to render different things at this point
	oapi, err := apidoc.NewOAPIRenderer(2)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return 1
	}
	oapi.Render(doc)

	return 0
}

// buildDoc is a sample of how to populate a Document with content
// from backends or other sources.
func buildDoc(doc *apidoc.Document) {
	// Load the /sys backend, and then append the separate manual paths
	apidoc.LoadBackend("sys", vault.Backend(), doc)
	doc.AddPath("sys", vault.ManualPaths()...)

	// Load another backend to show how separate mounts could be presented.
	// This will be in a separate, "aws" group in the output OAPI.
	apidoc.LoadBackend("aws", aws.Backend().Backend, doc)
}
