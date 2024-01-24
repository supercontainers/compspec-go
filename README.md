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

### Limitations

 - I'm starting with just Linux. I know there are those "other" platforms, but if it doesn't run on HPC or Kubernetes easily I'm not super interested (ahem, Mac and Windows)!

## Usage

Note that there is a [developer environment](.devcontainer) that provides a consistent version of Go, etc.
Build all binaries with:

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

Usage:
  comspec version
  comspec extract
```

More usage details will be added soon!

### Version

```bash
$ ./bin/compspec version
⭐️ compspec version 0.1.0-draft
```

### Extract

Coming soon!

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