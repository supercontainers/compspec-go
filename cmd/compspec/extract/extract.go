package extract

import (
	"fmt"
	"os"
	"runtime"

	p "github.com/supercontainers/compspec-go/plugins"
)

// Run will run an extraction of host metadata
func Run(filename string, pluginNames []string) error {
	fmt.Printf("⭐️ Running extract...\n")

	// Womp womp, we only support linux! There is no other way.
	operatingSystem := runtime.GOOS
	if operatingSystem != "linux" {
		return fmt.Errorf("🤓️ Sorry, we only support linux.")
	}

	// parse [section,...,section] into named plugins and sections
	// return plugins
	plugins, err := p.GetPlugins(pluginNames)
	if err != nil {
		return err
	}

	// Extract data for all plugins
	result, err := plugins.Extract()

	// If a filename is provided, save to json
	if filename != "" {

		// This returns an array of bytes
		b, err := result.ToJson()
		if err != nil {
			return fmt.Errorf("There was an issue marshalling to JSON: %s\n", err)
		}
		err = os.WriteFile(filename, b, 0644)
		if err != nil {
			return err
		}
	} else {
		result.Print()
	}

	fmt.Println("Extraction has run!")
	return nil
}
