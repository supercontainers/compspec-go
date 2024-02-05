package system

import (
	"fmt"
	"os/exec"

	"github.com/supercontainers/compspec-go/pkg/extractor"
	"github.com/supercontainers/compspec-go/pkg/utils"
)

const (
	linkerAMD64 = "/lib64/ld-linux-x86-64.so.2"
	linkeri386  = "/lib/ld-linux.so.2"
	linkerARM64 = "/lib/ld-linux-aarch64.so.1"
)

var (
	linkerPaths = map[string]string{"amd64": linkerAMD64, "i386": linkeri386, "arm64": linkerARM64}
)

// getOSArch determines arch based on the ld linux path
func getOsArch() (string, error) {

	// Detect OS architecture based on presence of this file
	for arch, path := range linkerPaths {
		exists, err := utils.PathExists(path)
		if err != nil {
			return "", err
		}
		if exists {
			return arch, nil
		}
	}
	return "", fmt.Errorf("cannot find architecture based on linker file")
}

// getArchInformation gets architecture information
func getArchInformation() (extractor.ExtractorSection, error) {
	info := extractor.ExtractorSection{}

	// Read in architectures
	arch, err := getOsArch()
	if err != nil {
		return info, err
	}
	info["name"] = arch

	// Try to run arch command to get more details, OK if we don't have it
	path, err := exec.LookPath("arch")
	if err == nil {
		output, err := utils.RunCommand([]string{path})
		if err == nil {
			info["arch"] = output
		}
	}
	return info, nil
}
