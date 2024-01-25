package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
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

	// Shared arguments (likely this will break into check and extract, shared for now)
	pluginNames := parser.StringList("n", "name", &argparse.Options{Help: "One or more specific plugins to target names"})

	// Extract arguments
	filename := extractCmd.String("o", "out", &argparse.Options{Help: "Save extraction to json file"})

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
			log.Fatalf("Issue with extraction: %s", err)
		}
	} else if listCmd.Happened() {
		list.Run(*pluginNames)
	} else if versionCmd.Happened() {
		RunVersion()
	} else {
		fmt.Println(Header)
		fmt.Println(parser.Usage(nil))
	}
}
