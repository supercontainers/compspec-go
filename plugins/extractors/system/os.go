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

	"github.com/compspec/compspec-go/pkg/plugin"
)

const (
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
)

// readOsRelease gets the name, version, and vendor from the os release file
func parseOsRelease() (string, string, string, error) {

	var name, version, vendor string

	// Determine OS release by reading this file
	f, err := os.Open(osReleaseFile)
	if err != nil {
		return name, version, vendor, fmt.Errorf("cannot find %s to determine OS metadata and release", osReleaseFile)
	}
	defer f.Close()

	// Use regular expression to find match for name
	s := bufio.NewScanner(f)
	for s.Scan() {
		text := s.Text()

		// Name
		match := regexName.FindStringSubmatch(text)
		if match != nil {
			name = strings.Trim(match[1], `"`)
			continue
		}

		// ID for os
		match = regexID.FindStringSubmatch(text)
		if match != nil {
			vendor = strings.Trim(match[1], `"`)
			continue
		}

		match = regexVersion.FindStringSubmatch(text)
		if match != nil {
			version = strings.Trim(match[1], `"`)
		}
	}
	return name, version, vendor, nil
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
func getOsInformation() (plugin.PluginSection, error) {
	info := plugin.PluginSection{}

	// Get the name, version, and vendor
	name, version, vendor, err := parseOsRelease()
	if err != nil {
		return info, err
	}
	info["name"] = name
	info["version"] = version
	info["vendor"] = vendor

	// Read in the os release metadata
	osRelease, err := readOsRelease(name, vendor)
	if err != nil {
		return info, err
	}
	info["release"] = osRelease
	return info, nil
}
