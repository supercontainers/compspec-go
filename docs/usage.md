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
  list     List plugins and known sections
  create   Create a compatibility artifact for the current host according to a
            definition
  match    Match a manifest of container images / artifact pairs against a set
            of host fields

Arguments:

  -h  --help  Print help information
  -n  --name  One or more specific plugins to target names
```

## Version

```bash
$ ./bin/compspec version
```
```console
⭐️ compspec version 0.1.1-draft
```

I know, the star should not be there. Fight me.

## List

The list command lists plugins (extractors and creators), and sections available for extractors.

```bash
$ ./bin/compspec list
```
```console
 Compatibility Plugins                                     
                            TYPE       NAME      SECTION   
 creation plugins                                          
                            creator    artifact            
                            creator    cluster             
-----------------------------------------------------------
 generic kernel extractor                                  
                            extractor  kernel    boot      
                            extractor  kernel    config    
                            extractor  kernel    modules   
-----------------------------------------------------------
 generic system extractor                                  
                            extractor  system    processor 
                            extractor  system    os        
                            extractor  system    arch      
                            extractor  system    memory    
                            extractor  system    cpu       
-----------------------------------------------------------
 generic library extractor                                 
                            extractor  library   mpi       
-----------------------------------------------------------
 node feature discovery                                    
                            extractor  nfd       cpu       
                            extractor  nfd       kernel    
                            extractor  nfd       local     
                            extractor  nfd       memory    
                            extractor  nfd       network   
                            extractor  nfd       pci       
                            extractor  nfd       storage   
                            extractor  nfd       system    
                            extractor  nfd       usb       
 TOTAL                                 6         20        
```

Note that we will eventually add a description column - it's not really warranted yet!

## Create

The create command handles two kinds of creation (sub-commands). Each of these is currently linked to a creation plugin.

 - **artifact**: create a compatibility artifact to describe an environment or application
 - **nodes** create a json graph format summary of nodes (a directory with one or more extracted metadata JSON files with node metadata)

The artifact case is described here. For the node case, you can read about it in the [rainbow scheduler](rainbow) documentation.

### Artifact

The create artifact command is how you take a compatibility request, or a YAML file that has a mapping between the extractors defined by this tool and your compatibility metadata namespace, and generate an artifact. The artifact typically will be a JSON dump of key value pairs, scoped under different namespaces, that you might push to a registry to live alongside a container image, and with the intention to eventually use it to check compatiility against a new system. To run create we can use the example in the top level repository:

```bash
./bin/compspec create artifact --in ./examples/lammps-experiment.yaml
```

Note that you'll see some errors about fields not being found! This is because we've implemented this for the fields to be added custom, on the command line.
The idea here is that you can add custom metadata fields during your build, which can be easier than adding for automated extraction. Let's add them now.

```bash
# a stands for "append" and it can write a new field or overwrite an existing one
./bin/compspec create artifact --in ./examples/lammps-experiment.yaml -a custom.gpu.available=yes
```
```console
{
  "version": "0.0.0",
  "kind": "CompatibilitySpec",
  "metadata": {
    "name": "lammps-prototype",
    "schemas": {
      "io.archspec": "https://raw.githubusercontent.com/supercontainers/compspec/main/archspec/compspec.json",
      "org.supercontainers": "https://raw.githubusercontent.com/supercontainers/compspec/main/supercontainers/compspec.json"
    }
  },
  "compatibilities": [
    {
      "name": "org.supercontainers",
      "version": "0.0.0",
      "attributes": {
        "hardware.gpu.available": "yes",
        "mpi.implementation": "mpich",
        "mpi.version": "4.1.1",
        "os.name": "Ubuntu 22.04.3 LTS",
        "os.release": "22.04.3",
        "os.vendor": "ubuntu",
        "os.version": "22.04"
      }
    },
    {
      "name": "io.archspec",
      "version": "0.0.0",
      "attributes": {
        "cpu.model": "13th Gen Intel(R) Core(TM) i5-1335U",
        "cpu.target": "amd64",
        "cpu.vendor": "GenuineIntel"
      }
    }
  ]
}
```

Awesome! That, as simple as it is, is our compatibility artifact. I ran the command on my host just now, but run for a container image during
a build will generate it for that context. We would want to save this to file:

```bash
./bin/compspec create artifact --in ./examples/lammps-experiment.yaml -a custom.gpu.available=yes -o ./examples/generated-compatibility-spec.json
```

And that's it! We would next (likely during CI) push this compatibility artifact to a URI that is likely (TBA) linked to the image.
For now we will manually remember the pairing, at least until the compatibility working group figures out the final design!

## Match

Match is the command you would use to check a potential host against one or more existing artifacts and find matches.
For a small experiment of using create against a set of containers and then testing how to do a match, we are going to place content
in [examples/check-lammps](examples/check-lammps). Note that we generated the actual compatibility spec and pushed with oras before running the example here! Following that, we might use the manifest in that directory to run a match. Note that since most use cases aren't checking the images in the manifest list against the host running the command, we instead
provide the parameters about the expected runtime host to them.

### Check Artifacts

You can do a check to ensure that all your artifacts exist.

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --check-artifacts
```
```console
Checking artifacts complete. There were 0 artifacts missing.
```

If this cannot be guaranteed, you can also ask to allow failures in finding them (and those images won't be included in the match):

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --allow-fail
```

### Match without Metadata Attributes

When you run match without options, all images are by default matches, hooray!

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml
```

<details>

<summary>Match without metadata attributes matches all images</summary>

```console
No field criteria provided, all images are matches.
 --- Found matches ---
ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64
ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-8-amd64
ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-20.04
ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-22.04
```

</details>

### Match with Metadata Attributes

Of course the purpose of the match is to actually provide metadata to match to! Let's ask for gpu. Since we are now querying across an entire graph (with more than one schema) we need to provide the full URI of the annotation. Let's ask for GPUs first!

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --match org.supercontainers.hardware.gpu.available=yes
```
```console
 --- Found matches ---
ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-20.04
ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-22.04
```

Now - no GPU!


```console
 --- Found matches ---
ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64
ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-8-amd64
```

That is very simple, but should serve our purposes for now.

### Match with randomize or single

You can choose to shuffle the results:

```bash
$ ./bin/compspec match -i ./examples/check-lammps/manifest.yaml --randomize
```
```console
Schema io.archspec is being added to the graph
Schema org.supercontainers is being added to the graph
No field criteria provided, all images are matches.
 --- Found matches ---
ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-20.04
ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64
ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-8-amd64
ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-22.04
```
Or do the same, but ask for only one result (note we might) change this to N results depending on use cases:

```bash
$ ./bin/compspec match -i ./examples/check-lammps/manifest.yaml --randomize --single
```
```console
Schema io.archspec is being added to the graph
Schema org.supercontainers is being added to the graph
No field criteria provided, all images are matches.
 --- Found matches ---
ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64
```

### Match and use cache

If you intend to run a match request many times (using repeated images) it's good practice to use a cache. This means you'll look in the cache before asking a registry, and save after if it doesn't exist. The directory must exist.

```bash
mkdir -p ./cache
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --cache ./cache
```

You can also save the graph to file:

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --cache ./cache --cache-graph ./cache/lammps-experiment.json
```

You can then use that cached graph later (NOTE this is at your own discretion knowing the schemas needed).

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --cache ./cache --cache-graph ./cache/lammps-experiment.json
```

Likely we will make tools to visualize it that can just show JGF!

### Match to print Metadata attributes

If you want to carefully inspect the metadata attributes discovered with your artifact set, add `-p` or `--print`

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --print
```

Below shows just the last image (the output is large)

<details>

<summary>Match to print metadata attributes associated with each image artifact</summary>

```console
-- Mapping for Images
  image: ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64
  nodes:
  -  compspec-root
  -  io.archspec.cpu
  -  io.archspec.cpu.model
  -  io.archspec.cpu.model.13th Gen Intel(R) Core(TM) i5-1335U
  -  io.archspec.cpu.target
  -  io.archspec.cpu.target.amd64
  -  io.archspec.cpu.vendor
  -  io.archspec.cpu.vendor.GenuineIntel
  -  org.supercontainers
  -  org.supercontainers.hardware
  -  org.supercontainers.hardware.gpu
  -  org.supercontainers.hardware.gpu.available
  -  org.supercontainers.hardware.gpu.available.no
  -  org.supercontainers.mpi
  -  org.supercontainers.mpi.implementation
  -  org.supercontainers.mpi.implementation.intel-mpi
  -  org.supercontainers.mpi.version
  -  org.supercontainers.mpi.version.2021.8
  -  org.supercontainers.os
  -  org.supercontainers.os.name
  -  org.supercontainers.os.name.Rocky Linux 9.3 (Blue Onyx)
  -  org.supercontainers.os.release
  -  org.supercontainers.os.release.9.3
  -  org.supercontainers.os.vendor
  -  org.supercontainers.os.vendor.rocky
  -  org.supercontainers.os.version
  -  org.supercontainers.os.version.9.3
  image: ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-8-amd64
  nodes:
  -  compspec-root
  -  io.archspec.cpu
  -  io.archspec.cpu.model
  -  io.archspec.cpu.model.13th Gen Intel(R) Core(TM) i5-1335U
  -  io.archspec.cpu.target
  -  io.archspec.cpu.target.amd64
  -  io.archspec.cpu.vendor
  -  io.archspec.cpu.vendor.GenuineIntel
  -  org.supercontainers
  -  org.supercontainers.hardware
  -  org.supercontainers.hardware.gpu
  -  org.supercontainers.hardware.gpu.available
  -  org.supercontainers.hardware.gpu.available.no
  -  org.supercontainers.mpi
  -  org.supercontainers.mpi.implementation
  -  org.supercontainers.mpi.implementation.intel-mpi
  -  org.supercontainers.mpi.version
  -  org.supercontainers.mpi.version.2021.8
  -  org.supercontainers.os
  -  org.supercontainers.os.name
  -  org.supercontainers.os.name.Rocky Linux 8.9 (Green Obsidian)
  -  org.supercontainers.os.release
  -  org.supercontainers.os.release.8.9
  -  org.supercontainers.os.vendor
  -  org.supercontainers.os.vendor.rocky
  -  org.supercontainers.os.version
  -  org.supercontainers.os.version.8.9
  image: ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-20.04
  nodes:
  -  compspec-root
  -  io.archspec.cpu
  -  io.archspec.cpu.model
  -  io.archspec.cpu.model.13th Gen Intel(R) Core(TM) i5-1335U
  -  io.archspec.cpu.target
  -  io.archspec.cpu.target.amd64
  -  io.archspec.cpu.vendor
  -  io.archspec.cpu.vendor.GenuineIntel
  -  org.supercontainers
  -  org.supercontainers.hardware
  -  org.supercontainers.hardware.gpu
  -  org.supercontainers.hardware.gpu.available
  -  org.supercontainers.hardware.gpu.available.yes
  -  org.supercontainers.mpi
  -  org.supercontainers.mpi.implementation
  -  org.supercontainers.mpi.implementation.OpenMPI
  -  org.supercontainers.mpi.version
  -  org.supercontainers.mpi.version.4.0.3
  -  org.supercontainers.os
  -  org.supercontainers.os.name
  -  org.supercontainers.os.name.Ubuntu 20.04.6 LTS
  -  org.supercontainers.os.release
  -  org.supercontainers.os.release.20.04.6
  -  org.supercontainers.os.vendor
  -  org.supercontainers.os.vendor.ubuntu
  -  org.supercontainers.os.version
  -  org.supercontainers.os.version.20.04
  image: ghcr.io/rse-ops/lammps-matrix:openmpi-ubuntu-gpu-22.04
  nodes:
  -  compspec-root
  -  io.archspec.cpu
  -  io.archspec.cpu.model
  -  io.archspec.cpu.model.13th Gen Intel(R) Core(TM) i5-1335U
  -  io.archspec.cpu.target
  -  io.archspec.cpu.target.amd64
  -  io.archspec.cpu.vendor
  -  io.archspec.cpu.vendor.GenuineIntel
  -  org.supercontainers
  -  org.supercontainers.hardware
  -  org.supercontainers.hardware.gpu
  -  org.supercontainers.hardware.gpu.available
  -  org.supercontainers.hardware.gpu.available.yes
  -  org.supercontainers.mpi
  -  org.supercontainers.mpi.implementation
  -  org.supercontainers.mpi.implementation.OpenMPI
  -  org.supercontainers.mpi.version
  -  org.supercontainers.mpi.version.4.1.2
  -  org.supercontainers.os
  -  org.supercontainers.os.name
  -  org.supercontainers.os.name.Ubuntu 22.04.3 LTS
  -  org.supercontainers.os.release
  -  org.supercontainers.os.release.22.04.3
  -  org.supercontainers.os.vendor
  -  org.supercontainers.os.vendor.ubuntu
  -  org.supercontainers.os.version
  -  org.supercontainers.os.version.22.04
```

</details>

### Match to show the graph

If you want to look at the schema graph (without images mapped to it) you can do:

```bash
./bin/compspec match -i ./examples/check-lammps/manifest.yaml --print-graph
```

This is in Json Graph Format (JGF). We likely will be developing better visualization tools to show this.

### Match Algorithm

The application and command logic and (very simple) algortithm works as follows.

1. Read in all entries from the list
2. Retrieve their artifacts, look for the "application/org.supercontainers.compspec" layer media type to identify it.
3. Generate a graph from the schemas we find (it is empty, no images mapped to it yet)
4. Map each image into the graph, linking the image identifer to each node in the graph
5. Take a user match request and look for the matching nodes (metadata attributes)
6. Take the intersection of images at all nodes, those are the matches!

Note that the arch represents YOUR host, and if this is run during build time, would be for the build environent. We probably need to think this over more. I need to test this out and ensure that the artifacts reflect the host where they were built, or the one we expect.

## Extract

Extraction has two use cases, and likely you won't be running this manually, but within the context of another command:

1. Extracting metadata about the container image at build time to generate an artifact (done via "create")
2. Extracting metadata about the host at image selection time, and comparing against a set of contender container images to select the best one (done via "match").

However, for the advanced or interested user, you can run extract as a standalone utility to inspect or otherwise save metadata from extractors.
For example, if you want to extract metadata to your local machine, you can use extract! Either just run all extractors and dump to the terminal:

```bash
# Not recommend, it's a lot!
./bin/compspec extract
```

If you want to allow failures (good for development):

```bash
./bin/compspec extract --allow-fail
```


Or use a specific, named extractor. Each extractor is shown below (with example output). The first example (with MPI) demonstrates
the full ability to specify:

1. A named extraction
2. One or more specific sections known to an extractor
3. Saving to json metadata instead of dumping to terminal


### Extractors 

Current Extractors include:

 - Library: library-specific metadata (e.g., mpi)
 - System: system-specific metadata (e.g., processor, cpu, arch, os, memory)
 - Kernel: kernel-speific metadata (e.g., boot, config, modules)
 - Node Feature Discovery: uses the [source](https://github.com/converged-computing/nfd-source) of NFD to derive metadata across many domains (cpu, kernel, local, memory, network, pci, storage, system, usb)

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
 - memory: parses /proc/meminfo and gives results primarily in KB

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
