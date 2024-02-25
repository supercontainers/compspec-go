# Comspec in Go

![img/compspec.png](img/compspec.png)

This is a prototype compatibility checking tool. Right now our aim is to use in the context of
[these build matrices](https://github.com/rse-ops/lammps-matrix) for LAMMPS and these prototype [specifications](https://github.com/supercontainers/compspec) that are based off of [Proposal C](https://github.com/opencontainers/wg-image-compatibility/pull/8) of the Compatibility Working Group. This is experimental because all of that is subject (and likely) to change. This project is under development, and you can see our [docs](docs) for early documentation.

 - ⭐️ [Documentation](docs) ⭐️

### Limitations

 - I'm starting with just Linux. I know there are those "other" platforms, but if it doesn't run on HPC or Kubernetes easily I'm not super interested (ahem, Mac and Windows)!
 - not all extractors work in containers (e.g., kernel needs to be on the host)
 - The node feature discovery source doesn't provide mapping of socket -> cores, nor does it give details about logical vs. physical CPU.
  - We will likely want to add hwloc go bindings, but there is a bug currently. 

Note that for development we are using nfd-source that does not require kubernetes:

```bash
go get -u github.com/converged-computing/nfd-source/source@0.0.1
```

There is some bug with installing the traditional way that I can look into later:

```
github.com/compspec/compspec-go/plugins/extractors/nfd imports
        github.com/converged-computing/nfd-source/source: no matching versions for query "latest"
```

And moving forward we will be working from this WIP branch:

```bash
# This is the develop branch
go get -u github.com/converged-computing/nfd-source/source@20d686e64926b80421637e82fb68e6c5f3f9242a
```

Note that the above is the main branch on February 22, 2024!

## TODO

 - add descriptions to sections
 - make extractor sections into interface
 - metadata namespace and exposure: someone writing a spec to create an artifact needs to know the extract namespace (and what is available) for the mapping.
 - tests: matrix that has several different flavors of builds, generating compspec json output to validate generation and correctness
 - likely we want a common configuration file to take an extraction -> check recipe
 - todo thinking around manifest.yaml that has listing of images / artifacts

### Extractors wanted / needed

A `*` indicates required for the work / prototype I want to do

 - power usage data [valorium](https://ipo.llnl.gov/sites/default/files/2023-08/Final_variorum-rnd-100-award.pdf)
 - ... please add more!

## License

This repository contains code derived from [sysinfo](https://github.com/zcalusic/sysinfo/tree/30169cfb37112a562cbf9133494a323764ad852c)
that was released also under an [MIT License](.github/LICENSE-SYSINFO). The library in question exposed needed functionality under a private
interface and required sudo for extra functionality that we did not need.

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614