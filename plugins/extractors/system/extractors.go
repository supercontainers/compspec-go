package system

import (
	"fmt"

	// TODO this library is a bit old, but seems to do the job for now
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/supercontainers/compspec-go/pkg/extractor"
)

const (
	CpuInfoFile = "/proc/cpuinfo"
)

// getCPUInformation gets information about the system
// TODO we might want to contribute this to archspec-go instead
func getCPUInformation() (extractor.ExtractorSection, error) {
	info := extractor.ExtractorSection{}

	stat, err := linuxproc.ReadCPUInfo(CpuInfoFile)
	if err != nil {
		return info, fmt.Errorf("cannot read %s: %s", CpuInfoFile, err)
	}

	// Manually set summary metrics
	// We need a good definition for these,
	info["cpu.logical.cpus"] = fmt.Sprintf("%d", stat.NumCPU())
	info["cpu.physical.cpus"] = fmt.Sprintf("%d", stat.NumPhysicalCPU())
	info["cpu.cores"] = fmt.Sprintf("%d", stat.NumCore())
	return info, nil
}

// getProcessorInformation returns details about each processor
func getProcessorInformation() (extractor.ExtractorSection, error) {
	info := extractor.ExtractorSection{}

	stat, err := linuxproc.ReadCPUInfo(CpuInfoFile)
	if err != nil {
		return info, fmt.Errorf("cannot read %s: %s", CpuInfoFile, err)
	}

	// Create features for each processor. Note we might want to separate this into a separate
	for _, s := range stat.Processors {
		uid := fmt.Sprintf("processor.%d.", s.CoreId)
		info[uid+"cachesize"] = fmt.Sprintf("%d", s.CacheSize)
		info[uid+"cores"] = fmt.Sprintf("%d", s.Cores)
		info[uid+"model"] = fmt.Sprintf("%d", s.Model)
		info[uid+"physicalid"] = fmt.Sprintf("%d", s.PhysicalId)
		info[uid+"flags"] = fmt.Sprintf("%s", s.Flags)
		info[uid+"mhz"] = fmt.Sprintf("%f", s.MHz)
		info[uid+"model"] = s.ModelName
		info[uid+"vendor"] = s.VendorId
	}
	return info, nil
}
