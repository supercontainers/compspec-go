package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/supercontainers/compspec-go/cmd/compspec/extract"
	"github.com/supercontainers/compspec-go/pkg/types"
)

// Command names
const (
	VersionCommand = "version"
	ExtractCommand = "extract"
)

// I know this text is terrible, just having fun for now
var (
	Usage = `              
┏┏┓┏┳┓┏┓┏┏┓┏┓┏
┗┗┛┛┗┗┣┛┛┣┛┗ ┗
	  ┛  ┛    

Usage:
  comspec version
  comspec extract
`
)

func RunVersion() {
	fmt.Printf("⭐️ compspec version %s\n", types.Version)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println(Usage)
		os.Exit(1)
	}

	// use cobra / pflags instead?
	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case VersionCommand:
		RunVersion()
		break
	case ExtractCommand:
		err := extract.Run(cmdArgs)
		if err != nil {
			log.Fatalf("Issue with extraction: %s", err)
		}
		break
	default:
		log.Fatalf("😱️ Invalid command: %s", cmd)
	}
}
