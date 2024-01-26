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
        extractor  kernel   boot      
        extractor  kernel   config    
        extractor  kernel   modules   
        extractor  system   cpu       
        extractor  system   processor 
        extractor  system   os        
        extractor  system   arch      
        extractor  library  mpi       
 TOTAL  8                            
```

Note that we will eventually add a description column - it's not really warranted yet!

## Create

The create command is how you take a compatibility request, or a YAML file that has a mapping between the extractors defined by this tool and your compatibility metadata namespace, and generate an artifact. The artifact typically will be a JSON dump of key value pairs, scoped under different namespaces, that you might push to a registry to live alongside a container image, and with the intention to eventually use it to check compatiility against a new system. To run create
we can use the example in the top level repository:

```bash
./bin/compspec create --in ./examples/lammps-experiment.yaml
```

Note that you'll see some errors about fields not being found! This is because we've implemented this for the fields to be added custom, on the command line.
The idea here is that you can add custom metadata fields during your build, which can be easier than adding for automated extraction. Let's add them now.

```bash
# a stands for "append" and it can write a new field or overwrite an existing one
./bin/compspec create --in ./examples/lammps-experiment.yaml -a custom.gpu.available=yes
```
```console
{
  "version": "0.0.0",
  "kind": "CompatibilitySpec",
  "metadata": {
    "name": "lammps-prototype",
    "jsonSchema": "https://raw.githubusercontent.com/supercontainers/compspec/main/supercontainers/compspec.json"
  },
  "compatibilities": [
    {
      "name": "org.supercontainers.mpi",
      "version": "0.0.0",
      "annotations": {
        "implementation": "mpich",
        "version": "4.1.1"
      }
    },
    {
      "name": "org.supercontainers.hardware.gpu",
      "version": "0.0.0",
      "annotations": {
        "available": "yes"
      }
    },
    {
      "name": "io.archspec.cpu",
      "version": "0.0.0",
      "annotations": {
        "model": "13th Gen Intel(R) Core(TM) i5-1335U",
        "target": "amd64",
        "vendor": "GenuineIntel"
      }
    }
  ]
}
```

Awesome! That, as simple as it is, is our compatibility artifact. I ran the command on my host just now, but run for a container image during
a build will generate it for that context. We would want to save this to file:

```bash
./bin/compspec create --in ./examples/lammps-experiment.yaml -a custom.gpu.available=yes -o ./examples/generated-compatibility-spec.json
```

And that's it! We would next (likely during CI) push this compatibility artifact to a URI that is likely (TBA) linked to the image.
For now we will manually remember the pairing, at least until the compatibility working group figures out the final design!

## Check

Check is the command you would use to check a potential host against one or more existing artifacts.
For a small experiment of using create against a set of containers and then testing how to do a check, we are going to place content
in [examples/check-lammps](examples/check-lammps).

## Extract

Extraction has two use cases, and likely you won't be running this manually, but within the context of another command:

1. Extracting metadata about the container image at build time to generate an artifact (done via "create")
2. Extracting metadata about the host at image selection time, and comparing against a set of contender container images to select the best one (done via "check").

However, for the advanced or interested user, you can run extract as a standalone utility to inspect or otherwise save metadata from extractors.
For example, if you want to extract metadata to your local machine, you can use extract! Either just run all extractors and dump to the terminal:

```bash
# Not recommend, it's a lot!
./bin/compspec extract
```

Or use a specific, named extractor. Each extractor is shown below (with example output). The first example (with MPI) demonstrates
the full ability to specify:

1. A named extraction
2. One or more specific sections known to an extractor
3. Saving to json metadata instead of dumping to terminal


### Extractors 

Current Extractors include:

 - Library: library-specific metadata (e.g., mpi)
 - System: system-specific metadata (e.g., processor, cpu, arch, os)
 - Kernel: kernel-speific metadata (e.g., boot, config, modules)


#### Library

The library extractor currently just has one section for "mpi"

```bash
./bin/compspec extract --name library
```
```console
⭐️ Running extract...
 --Result for library
 -- Section mpi
   variant: mpich
   version: 4.1.1
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
          "variant": "mpich",
          "version": "4.1.1"
        }
      }
    }
  }
}
```

That shows the generic structure of an extractor output. The "library" extractor owns a set of groups (sections) each with their own namespaced attributes.

#### System

The system extractor supports three sections

 - cpu: Basic CPU counts and metadata
 - processor: detailed information on every processor
 - os: operating system information
 - arch: architecture

For example:

```bash
./bin/compspec extract --name system[os]
```
```console
⭐️ Running extract...
 --Result for system
 -- Section os
   release: 22.04.3
   name: Ubuntu 22.04.3 LTS
   version: 22.04
   vendor: ubuntu
Extraction has run!
```

Or for arch:

```bash
./bin/compspec extract --name system[arch]
```
```console
⭐️ Running extract...
 --Result for system
 -- Section arch
   name: amd64
Extraction has run!
```

#### Kernel

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