# Comspec in Go

![img/compspec.png](img/compspec.png)

This is a prototype compatibility checking tool. Right now our aim is to use in the context of
[these build matrices](https://github.com/rse-ops/lammps-matrix) for LAMMPS and these prototype [specifications](https://github.com/supercontainers/compspec) that are based off of [Proposal C](https://github.com/opencontainers/wg-image-compatibility/pull/8) of the Compatibility Working Group. This is experimental because all of that is subject (and likely) to change.

## Design

The design is based on the prototype from that pull request, shown below.

![img/proposal-c-plugin-design.png](img/proposal-c-plugin-design.png)

Specifically, I'll try to do simple interfaces (in [plugins](plugins)) to:

 - **Extract** Scan a system to collect compatibility metadata / attributes for a named set of extractions
   - this maybe isn't great practice, but if I build the containers I trust them
   - I can also do something simple like run this tool in a pod-> container small cluster
 - load in an artifact spec (json) from a URL
 - compare with a request for a specific system (with any level of detail desired)
   - "basic" might just be providing the architecture
   - "descriptive" might include other variables

The extraction step is also important, but likely this would happen at build time (maybe by another tool).
For now (since this is a prototype) I'm going to just manually do it.

## Usage

### Build

Build the `compspec` binary with:

```bash
make
```

This generates the `bin/compspec` that you can use:

```bash
./bin/compspec
```
```console
┏┏┓┏┳┓┏┓┏┏┓┏┓┏
┗┗┛┛┗┗┣┛┛┣┛┗ ┗
          ┛  ┛    

[sub]Command required
usage: compspec <Command> [-h|--help] [-n|--name "<value>" [-n|--name "<value>"
                ...]]

                Compatibility checking for container images

Commands:

  version  See the version of compspec
  extract  Run one or more extractors

Arguments:

  -h  --help  Print help information
  -n  --name  One or more specific extractor plugin names
```


### Version

```bash
$ ./bin/compspec version
```
```console
⭐️ compspec version 0.1.0-draft
```

I know, the star should not be there. Fight me.

### List

The list command lists each extractor, and sections available for it.

```bash
$ ./bin/compspec list
```
```console
 Compatibility Plugins              
        TYPE      NAME    SECTION   
        extrator  kernel  boot      
        extrator  kernel  config    
        extrator  kernel  modules   
        extrator  system  cpu       
        extrator  system  processor 
 TOTAL  5                           
```

Note that we will eventually add a description column - it's not really warranted yet!

### Extract

If you want to extract metadata to your local machine, you can use extract! Either just run all extractors and dump to the terminal:

```bash
# Not recommend, it's a lot!
./bin/compspec extract
```

Or target a specific one:

```bash
./bin/compspec extract --name kernel
```

Better, save to json file:

```bash
./bin/compspec extract --name kernel -o test-kernel.json
```

This has a better structure for inspecting easily (only the top of the file is shown):

<details>

<summary>Kernel JSON output</summary>

```json
{
  "extractors": {
    "kernel": {
      "sections": {
        "boot": {
          "BOOT_IMAGE": "/boot/vmlinuz-6.1.0-1028-oem",
          "quiet": "",
          "ro": "",
          "root": "UUID",
          "splash": "",
          "vt.handoff": "7"
        },
        "config": {
          "CONFIG_104_QUAD_8": "m",
          "CONFIG_60XX_WDT": "m",
          "CONFIG_64BIT": "y",
          "CONFIG_6LOWPAN": "m"
        }
      }
    }
  }
}
```

</details>

An extractor is made up of sections, and you can ask for parsing just a specific one. 

```bash
./bin/compspec extract --name kernel[config]
```
```console
⭐️ Running extract...
 --Result for kernel
 -- Section boot
   root: UUID
   ro: 
   quiet: 
   splash: 
   vt.handoff: 7
   BOOT_IMAGE: /boot/vmlinuz-6.1.0-1028-oem
Extraction has run!
```

To ask for more than one, it's a comma separated list.

```bash
./bin/compspec extract --name kernel[config,boot]
```

The ordering of your list is honored.

## Developer

Note that there is a [developer environment](.devcontainer) that provides a consistent version of Go, etc.
However, it won't work with all extractors.  Note that for any command that uses a plugin (e.g., `extract` and `check`)


### Limitations

 - I'm starting with just Linux. I know there are those "other" platforms, but if it doesn't run on HPC or Kubernetes easily I'm not super interested (ahem, Mac and Windows)!
 - not all extractors work in containers (e.g., kernel needs to be on the host)

## TODO

 - likely we want a common configuration file to take an extraction -> check recipe
 - need to develop check plugin family
 - todo thinking around manifest.yaml that has listing of images / artifacts

### Extractors wanted / needed

A `*` indicates required for the work / prototype I want to do

 - power usage data [valorium](https://ipo.llnl.gov/sites/default/files/2023-08/Final_variorum-rnd-100-award.pdf)
 - * architecture [archspec](https://github.com/archspec)
 - * MPI existence / variants
 - * operating system stuff
 - ... please add more!


## Thanks and Previous Art

- I learned about kernel parsing from [mfranczy/compat](https://github.com/mfranczy/compat)

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614