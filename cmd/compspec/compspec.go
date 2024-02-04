package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/supercontainers/compspec-go/cmd/compspec/create"
	"github.com/supercontainers/compspec-go/cmd/compspec/extract"
	"github.com/supercontainers/compspec-go/cmd/compspec/list"
	"github.com/supercontainers/compspec-go/cmd/compspec/match"
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
	matchCmd := parser.NewCommand("match", "Match a manifest of container images / artifact pairs against a set of host fields")

	// Shared arguments (likely this will break into check and extract, shared for now)
	pluginNames := parser.StringList("n", "name", &argparse.Options{Help: "One or more specific plugins to target names"})

	// Extract arguments
	filename := extractCmd.String("o", "out", &argparse.Options{Help: "Save extraction to json file"})
	allowFail := extractCmd.Flag("f", "allow-fail", &argparse.Options{Help: "Allow any specific extractor to fail (and continue extraction)"})

	// Match arguments
	matchFields := matchCmd.StringList("m", "match", &argparse.Options{Help: "One or more key value pairs to match"})
	manifestFile := matchCmd.String("i", "in", &argparse.Options{Required: true, Help: "Input manifest list yaml that contains pairs of images and artifacts"})
	printMapping := matchCmd.Flag("p", "print", &argparse.Options{Help: "Print mapping of images to attributes only."})
	printGraph := matchCmd.Flag("g", "print-graph", &argparse.Options{Help: "Print schema graph"})
	checkArtifacts := matchCmd.Flag("c", "check-artifacts", &argparse.Options{Help: "Check that all artifacts exist"})
	allowFailMatch := matchCmd.Flag("f", "allow-fail", &argparse.Options{Help: "Allow an artifact to be missing (and not included)"})

	// Create arguments
	options := createCmd.StringList("a", "append", &argparse.Options{Help: "Append one or more custom metadata fields to append"})
	specname := createCmd.String("i", "in", &argparse.Options{Required: true, Help: "Input yaml that contains spec for creation"})
	specfile := createCmd.String("o", "out", &argparse.Options{Help: "Save compatibility json artifact to this file"})
	mediaType := createCmd.String("m", "media-type", &argparse.Options{Help: "The expected media-type for the compatibility artifact"})
	allowFailCreate := createCmd.Flag("f", "allow-fail", &argparse.Options{Help: "Allow any specific extractor to fail (and continue extraction)"})

	// Now parse the arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(Header)
		fmt.Println(parser.Usage(err))
		return
	}

	if extractCmd.Happened() {
		err := extract.Run(*filename, *pluginNames, *allowFail)
		if err != nil {
			log.Fatalf("Issue with extraction: %s\n", err)
		}
	} else if createCmd.Happened() {
		err := create.Run(*specname, *options, *specfile, *allowFailCreate)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else if matchCmd.Happened() {
		err := match.Run(*manifestFile, *matchFields, *mediaType, *printMapping, *printGraph, *allowFailMatch, *checkArtifacts)
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
