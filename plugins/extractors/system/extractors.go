package system

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/pkg/utils"
)

const (
	CpuInfoFile = "/proc/cpuinfo"
)

var (
	armVendors = map[string]string{
		"0x41": "ARM",
		"0x42": "Broadcom",
		"0x43": "Cavium",
		"0x44": "DEC",
		"0x46": "Fujitsu",
		"0x48": "HiSilicon",
		"0x49": "Infineon Technologies AG",
		"0x4d": "Motorola",
		"0x4e": "Nvidia",
		"0x50": "APM",
		"0x51": "Qualcomm",
		"0x53": "Samsung",
		"0x56": "Marvell",
		"0x61": "Apple",
		"0x66": "Faraday",
		"0x68": "HXT",
		"0x69": "Intel",
	}
)

// Note: arm is a jerk
// https://www.cnx-software.com/2018/02/14/human-readable-decoding-of-proc-cpuinfo-for-arm-processors/

// A quasi manual mapping. Here is what I see for my x86 system. When we parse, spaces are replaced with _
// and everything is made lowercase
// processor       : 11
// vendor_id       : GenuineIntel
// cpu family      : 6
// model           : 186
// model name      : 13th Gen Intel(R) Core(TM) i5-1335U
// stepping        : 3
// microcode       : 0x4e0e
// cpu MHz         : 1599.992
// cache size      : 12288 KB
// physical id     : 0
// siblings        : 12
// core id         : 15
// cpu cores       : 10
// apicid          : 30
// initial apicid  : 30
// fpu             : yes
// fpu_exception   : yes
// cpuid level     : 32
// wp              : yes
// flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc cpuid aperfmperf tsc_known_freq pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpuid_fault epb ssbd ibrs ibpb stibp ibrs_enhanced tpr_shadow vnmi flexpriority ept vpid ept_ad fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid rdseed adx smap clflushopt clwb intel_pt sha_ni xsaveopt xsavec xgetbv1 xsaves avx_vnni dtherm ida arat pln pts hwp hwp_notify hwp_act_window hwp_epp hwp_pkg_req hfi umip pku ospke waitpkg gfni vaes vpclmulqdq rdpid movdiri movdir64b fsrm md_clear serialize arch_lbr ibt flush_l1d arch_capabilities
// vmx flags       : vnmi preemption_timer posted_intr invvpid ept_x_only ept_ad ept_1gb flexpriority apicv tsc_offset vtpr mtf vapic ept vpid unrestricted_guest vapic_reg vid ple shadow_vmcs ept_mode_based_exec tsc_scaling usr_wait_pause
// bugs            : spectre_v1 spectre_v2 spec_store_bypass swapgs eibrs_pbrsb
// bogomips        : 4992.00
// clflush size    : 64
// cache_alignment : 64
// address sizes   : 39 bits physical, 48 bits virtual
// power management:

// ARM
// https://github.com/ARM-software/abi-aa
// processor       : 14
// BogoMIPS        : 2100.00
// Features        : fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm jscvt fcma lrcpc dcpop sha3 sm3 sm4 asimddp sha512 sve asimdfhm dit uscat ilrcpc flagm ssbs paca pacg dcpodp svei8mm svebf16 i8mm bf16 dgh rng
// CPU implementer : 0x41
// CPU architecture: 8
// CPU variant     : 0x1
// CPU part        : 0xd40
// CPU revision    : 1

// See https://github.com/randombit/cpuinfo/blob/master/ppc/power8 for other processors

// Get architecture is akin to model?
func getCpuArchitecture(p map[string]string) (string, error) {
	return utils.LookupValue(p, "cpu_family", "cpu_architecture")
}

// TODO need to check this mapping is right
func getCpuVariant(p map[string]string) (string, error) {
	return utils.LookupValue(p, "model_name", "cpu_variant")
}

// Get the CPU vendor, currently supports x86 and arm
func getCpuVendor(p map[string]string) (string, error) {

	vendor, err := utils.LookupValue(p, "vendor_id", "cpu_implementer")
	if err != nil {
		return vendor, err
	}
	vendorName, ok := armVendors[vendor]
	if !ok {
		return vendor, nil
	}
	return vendorName, nil

}

// Get the CPU features or flags
func getCpuFeatures(p map[string]string) (string, error) {
	return utils.LookupValue(p, "flags", "features")
}

// getCPUInformation gets information about the system
// TODO this is not used.
func getCPUInformation() (plugin.PluginSection, error) {
	info := plugin.PluginSection{}

	cores := runtime.NumCPU()

	// This is a guess at best
	info["cores"] = fmt.Sprintf("%d", cores)

	//stat, err := linuxproc.ReadCPUInfo(CpuInfoFile)
	//if err != nil {
	//	return info, fmt.Errorf("cannot read %s: %s", CpuInfoFile, err)
	//}

	// Manually set summary metrics
	// We need a good definition for these,
	//info["logical.cpus"] = fmt.Sprintf("%d", stat.NumCPU())
	//info["physical.cpus"] = fmt.Sprintf("%d", stat.NumPhysicalCPU())
	//info["cores"] = fmt.Sprintf("%d", stat.NumCore())
	return info, nil
}

// getProcessorInformation returns details about each processor
func getProcessorInformation() (plugin.PluginSection, error) {
	info := plugin.PluginSection{}

	raw, err := os.ReadFile(CpuInfoFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")

	// Parse raw processors into list first
	processors := []map[string]string{}
	current := map[string]string{}

	// Only PPC cpuinfo has timebase
	ppcFields := map[string]string{}
	isPPC := strings.Contains(string(raw), "timebase")

	// We need custom parsing, the sections per processor are split by newlines
	for _, line := range lines {
		line = strings.Trim(line, " ")
		// We found a new processor, add the last one and continue
		if line == "" && len(current) > 0 {
			processors = append(processors, current)
			current = map[string]string{}
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) > 1 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Replace spaces with _. Not sure if this is a good idea, but I don't like spaces
			key = strings.ToLower(strings.ReplaceAll(key, " ", "_"))

			// Is it a more global ppc field?
			if isPPC && key != "processor" {
				ppcFields[key] = value
				continue
			}

			if value != "" {
				current[key] = value
			}
		}
	}

	// Create common features for each processor, but also allow all fields as is
	for i, p := range processors {
		uid, ok := p["processor"]
		if !ok {
			fmt.Printf("Warning, processor metadata index %d missing 'processor' uid, skipping\n", i)
			continue
		}
		uid = fmt.Sprintf("%s.", uid)

		// Parse cpu vendor - arm has a lookup
		vendor, err := getCpuVendor(p)
		if err != nil {
			return info, err
		}
		info[uid+"normalized.vendor"] = vendor

		// bogompis should be the same after lowercase
		bogomips, ok := p["bogomips"]
		if ok {
			info[uid+"normalized.botomips"] = bogomips
		}

		// features or flags
		features, err := getCpuFeatures(p)
		if err != nil {
			return info, err
		}
		info[uid+"normalized.features"] = features

		family, err := getCpuArchitecture(p)
		if err != nil {
			return info, err
		}
		info[uid+"normalized.family"] = family

		variant, err := getCpuVariant(p)
		if err != nil {
			return info, err
		}
		info[uid+"normalized.model"] = variant
		if isPPC {
			for key, value := range ppcFields {
				info[uid+key] = value
			}
		}
		for key, value := range p {
			info[uid+"raw."+key] = value
		}

	}
	return info, nil
}
