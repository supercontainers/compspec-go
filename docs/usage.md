# Usage

## Build

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

## Version

```bash
$ ./bin/compspec version
```
```console
⭐️ compspec version 0.1.0-draft
```

I know, the star should not be there. Fight me.

## List

The list command lists each extractor, and sections available for it.

```bash
$ ./bin/compspec list
```
```console
 Compatibility Plugins               
        TYPE      NAME     SECTION   
        extrator  kernel   boot      
        extrator  kernel   config    
        extrator  kernel   modules   
        extrator  system   cpu       
        extrator  system   processor 
        extrator  system   os        
        extrator  library  mpi       
 TOTAL  7                            
```

Note that we will eventually add a description column - it's not really warranted yet!

## Extract

If you want to extract metadata to your local machine, you can use extract! Either just run all extractors and dump to the terminal:

```bash
# Not recommend, it's a lot!
./bin/compspec extract
```

Or use a specific, named extractor. Each extractor is shown below (with example output). The first example (with MPI) demonstrates
the full ability to specify:

1. A named extraction
2. One or more specific sections known to an extractor
3. Saving to json metadata instead of dumping to terminal

### Library

The library extractor currently just has one section for "mpi"

```bash
./bin/compspec extract --name library
```
```console
⭐️ Running extract...
 --Result for library
 -- Section mpi
   mpi.variant: mpich
   mpi.version: 4.0
Extraction has run!
```

This would be the same as selecting the section explicitly

```bash
./bin/compspec extract --name library[mpi]
```

If you have a lot of data that you want to use later, save to a json file.

```bash
./bin/compspec extract --name library -o test-library.json
cat test-library.json
```
```console
{
  "extractors": {
    "library": {
      "sections": {
        "mpi": {
          "mpi.variant": "mpich",
          "mpi.version": "4.0"
        }
      }
    }
  }
}
```

That shows the generic structure of an extractor output. The "library" extractor owns a set of groups (sections) each with their own namespaced attributes.

### System

The system extractor supports three sections

 - cpu: Basic CPU counts and metadata
 - processor: detailed information on every processor
 - os: operating system information

For example:

```bash
./bin/comspec extract --name system[os]
```
```console
⭐️ Running extract...
 --Result for system
 -- Section os
   arch.os.name: Ubuntu 22.04.3 LTS
   arch.os.version: 22.04
   arch.os.vendor: ubuntu
   arch.os.release: 22.04.3
   arch.name: amd64
Extraction has run!
```

### Kernel

Kernel supports three sections:

 - config: The full kernel configuration
 - boot: the command line provided at boot time
 - modules: kernel modules (very large output)!

```bash
./bin/compspec extract --name kernel
```

You can select one section:

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

And here is how you might select a group of sections:

```bash
./bin/compspec extract --name kernel[boot,config]
```

or write to json:

```bash
./bin/compspec extract --name kernel[boot,config] -o test-kernel.json
```

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

The ordering of your list is honored.

## Developer

Note that there is a [developer environment](.devcontainer) that provides a consistent version of Go, etc.
However, it won't work with all extractors.  Note that for any command that uses a plugin (e.g., `extract` and `check`)