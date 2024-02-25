package system

import (
	"os"
	"strings"

	"github.com/compspec/compspec-go/pkg/plugin"
)

const (
	memoryInfoFile = "/proc/meminfo"
)

// getMemoryInformation parses /proc/meminfo to get node memory metadata
func getMemoryInformation() (plugin.PluginSection, error) {
	info := plugin.PluginSection{}

	raw, err := os.ReadFile(memoryInfoFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")

	// We need custom parsing, the sections per processor are split by newlines
	for _, line := range lines {

		// I don't see any empty lines, etc.
		line = strings.Trim(line, " ")
		parts := strings.Split(line, ":")

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Replace parens with underscore. Leave camel case for the rest...
		key = strings.ReplaceAll(key, "(", "_")
		key = strings.ToLower(strings.ReplaceAll(key, ")", ""))
		info[key] = value
	}
	return info, nil
}
