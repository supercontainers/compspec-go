package library

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/supercontainers/compspec-go/pkg/extractor"
	"github.com/supercontainers/compspec-go/pkg/utils"
)

const (
	MPIRunExec = "mpirun"
)

var (
	regexIntelMPIVersion = regexp.MustCompile(`Version (.*) `)
)

// getMPIInformation returns info on mpi versions and variant
// yes, fairly janky, please improve upon! This is for a prototype
func getMPIInformation() (extractor.ExtractorSection, error) {
	info := extractor.ExtractorSection{}

	// Do we even have mpirun?
	path, err := exec.LookPath(MPIRunExec)
	if err != nil {
		return info, nil
	}

	// Get output from the tool
	command := []string{MPIRunExec, "--version"}
	output, err := utils.RunCommand(command)
	if err != nil {
		return info, err
	}

	// This is really simple - if we find Open MPI in a line, that's the variant
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Open MPI") {
			info["mpi.variant"] = "OpenMPI"
			parts := strings.Split(line, " ")
			info["mpi.version"] = parts[len(parts)-1]
			return info, nil
		}

		// Intel(R) MPI Library for Linux* OS, Version 2021.8 Build 20221129 (id: 339ec755a1)
		if strings.Contains(line, "Intel") {
			info["mpi.variant"] = "intel-mpi"
			match := regexIntelMPIVersion.FindStringSubmatch(line)
			if match != nil {
				parts := strings.Split(match[0], " ")
				info["mpi.version"] = parts[1]
			}
			return info, nil
		}

		// Note that for mpich there is a LOT more metadata
		// Right now I'm assuming if we find Version: it's for Open MPI
		if strings.Contains(line, "Version:") {
			info["mpi.variant"] = "mpich"
			parts := strings.Split(line, " ")
			info["mpi.version"] = parts[len(parts)-1]
			return info, nil
		}

	}

	fmt.Println(output)
	fmt.Printf("%s is available at %s\n", MPIRunExec, path)
	return info, nil
}
