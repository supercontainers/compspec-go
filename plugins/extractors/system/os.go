// Derived from:
// https://github.com/zcalusic/sysinfo/blob/30169cfb37112a562cbf9133494a323764ad852c/os.go
// under an MIT license

package system

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/supercontainers/compspec-go/pkg/extractor"
	"github.com/supercontainers/compspec-go/pkg/utils"
)

const (
	linkerAMD64       = "/lib64/ld-linux-x86-64.so.2"
	linkeri386        = "/lib/ld-linux.so.2"
	osReleaseFile     = "/etc/os-release"
	versionDebianFile = "/etc/debian_version"
	versionCentosFile = "/etc/centos-release"
	versionRHELFile   = "/etc/redhat-release"
	versionRockyFile  = "/etc/rocky-release"
)

var (
	regexName    = regexp.MustCompile(`^PRETTY_NAME=(.*)$`)
	regexID      = regexp.MustCompile(`^ID=(.*)$`)
	regexVersion = regexp.MustCompile(`^VERSION_ID=(.*)$`)
	regexUbuntu  = regexp.MustCompile(`[\( ]([\d\.]+)`)
	regexCentos  = regexp.MustCompile(`^CentOS( Linux)? release ([\d\.]+)`)
	regexRocky   = regexp.MustCompile(`^Rocky( Linux)? release ([\d\.]+)`)
	regexRHEL    = regexp.MustCompile(`[\( ]([\d\.]+)`)
	linkerPaths  = map[string]string{"amd64": linkerAMD64, "i386": linkeri386}
)

// getOSArch determines arch based on the ld linux path
func getOsArch() (string, error) {

	// Detect OS architecture based on presence of this file
	for arch, path := range linkerPaths {
		exists, err := utils.FileExists(path)
		if err != nil {
			return "", err
		}
		if exists {
			return arch, nil
		}
	}
	return "", fmt.Errorf("cannot find architecture based on linker file")
}

// addOsRelease adds in metadata fields from the os release file

func addOsRelease(info *extractor.ExtractorSection) error {

	// Determine OS release by reading this file
	f, err := os.Open(osReleaseFile)
	if err != nil {
		return fmt.Errorf("cannot find %s to determine OS metadata and release", osReleaseFile)
	}
	defer f.Close()

	// Use regular expression to find match for name
	s := bufio.NewScanner(f)
	for s.Scan() {
		text := s.Text()

		// Name
		match := regexName.FindStringSubmatch(text)
		if match != nil {
			(*info)["arch.os.name"] = strings.Trim(match[1], `"`)
			continue
		}

		// ID for os
		match = regexID.FindStringSubmatch(text)
		if match != nil {
			(*info)["arch.os.vendor"] = strings.Trim(match[1], `"`)
			continue
		}

		match = regexVersion.FindStringSubmatch(text)
		if match != nil {
			(*info)["arch.os.version"] = strings.Trim(match[1], `"`)
		}
	}
	return nil
}

// readOsRelease looks for different os version files to read
func readOsRelease(prettyName string, vendor string) (string, error) {

	switch vendor {
	case "debian":
		raw, err := os.ReadFile(versionDebianFile)
		if err != nil {
			return "", err
		}
		return string(raw), nil
	case "ubuntu":
		match := regexUbuntu.FindStringSubmatch(prettyName)
		if match != nil {
			return match[1], nil
		}
	case "centos":
		raw, err := os.ReadFile(versionCentosFile)
		if err != nil {
			return "", err
		}
		match := regexCentos.FindStringSubmatch(string(raw))
		if match != nil {
			return match[2], nil
		}
	case "rocky":
		raw, err := os.ReadFile(versionRockyFile)
		if err != nil {
			return "", err
		}
		match := regexRocky.FindStringSubmatch(string(raw))
		if match != nil {
			// Rocky Linux release 9.3 (Blue Onyx)
			parts := strings.Split(match[2], " ")
			return parts[0], nil
		}
	case "rhel":
		raw, err := os.ReadFile(versionRHELFile)
		if err != nil {
			return "", err
		}
		match := regexRHEL.FindStringSubmatch(string(raw))
		if match != nil {
			return match[1], nil
		}
		match = regexRHEL.FindStringSubmatch(prettyName)
		if match != nil {
			return match[1], nil
		}
	}
	return "", fmt.Errorf("cannot find os release")
}

// getOSInformation gets operating system level metadata
func getOsInformation() (extractor.ExtractorSection, error) {
	info := extractor.ExtractorSection{}

	// Read in architectures
	arch, err := getOsArch()
	if err != nil {
		return info, err
	}
	info["arch.name"] = arch

	// Read in metadata (version, id, name) from os release file
	err = addOsRelease(&info)
	if err != nil {
		return info, err
	}

	// Read in the os release metadata
	osRelease, err := readOsRelease(info["arch.os.name"], info["arch.os.vendor"])
	if err != nil {
		return info, err
	}
	info["arch.os.release"] = osRelease
	return info, nil
}
