package extract

import (
	"fmt"
	"runtime"

	"github.com/supercontainers/compspec-go/plugins/extractors/kernel"
)

// Run will run an extraction of host metadata
func Run(args []string) error {
	fmt.Printf("⭐️ running extract\n")

	// TODO allow to specify extractors, right now we import and run all of them

	// Womp womp, we only support linux! There is no other way.
	os := runtime.GOOS
	if os != "linux" {
		return fmt.Errorf("Sorry, we only support linux.")
	}
	e := kernel.KernelExtractor{}
	result, err := e.Extract(args)
	if err != nil {
		fmt.Printf("There was a kernel extreaction error: %s\n", err)
		return err
	}
	result.Print()
	fmt.Println("Extraction has run!")
	return nil
}
