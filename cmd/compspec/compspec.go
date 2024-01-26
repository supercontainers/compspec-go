package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/supercontainers/compspec-go/cmd/compspec/create"
	"github.com/supercontainers/compspec-go/cmd/compspec/extract"
	"github.com/supercontainers/compspec-go/cmd/compspec/list"
	"github.com/supercontainers/compspec-go/pkg/types"
)

// I know this text is terrible, just having fun for now
var (
	Header = `              
┏┏┓┏┳┓┏┓┏┏┓┏┓┏
┗┗┛┛┗┗┣┛┛┣┛┗ ┗
	  ┛  ┛    
`
)

func RunVersion() {
	fmt.Printf("⭐️ compspec version %s\n", types.Version)
}

func main() {

	parser := argparse.NewParser("compspec", "Compatibility checking for container images")
	versionCmd := parser.NewCommand("version", "See the version of compspec")
	extractCmd := parser.NewCommand("extract", "Run one or more extractors")
	listCmd := parser.NewCommand("list", "List plugins and known sections")
	createCmd := parser.NewCommand("create", "Create a compatibility artifact for the current host according to a definition")

	// Shared arguments (likely this will break into check and extract, shared for now)
	pluginNames := parser.StringList("n", "name", &argparse.Options{Help: "One or more specific plugins to target names"})

	// Extract arguments
	filename := extractCmd.String("o", "out", &argparse.Options{Help: "Save extraction to json file"})

	// Create arguments
	options := parser.StringList("a", "append", &argparse.Options{Help: "One or more custom metadata fields to append"})
	specname := createCmd.String("i", "in", &argparse.Options{Required: true, Help: "Input yaml that contains spec for creation"})
	specfile := createCmd.String("o", "out", &argparse.Options{Help: "Save compatibility json artifact to this file"})

	// Now parse the arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(Header)
		fmt.Println(parser.Usage(err))
		return
	}

	if extractCmd.Happened() {
		err := extract.Run(*filename, *pluginNames)
		if err != nil {
			log.Fatalf("Issue with extraction: %s\n", err)
		}
	} else if createCmd.Happened() {
		err := create.Run(*specname, *options, *specfile)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else if listCmd.Happened() {
		err := list.Run(*pluginNames)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else if versionCmd.Happened() {
		RunVersion()
	} else {
		fmt.Println(Header)
		fmt.Println(parser.Usage(nil))
	}
}
