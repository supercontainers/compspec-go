# Check LAMMPS

This will be a small experiment to:

1. Generate compatiility artifacts for a few containers
2. Push them to an OCI registry
3. Create (design) a manifest.yaml that lists them side by side
4. Test (and develop) the check tool to choose the best one (and test different modes for doing this)

## 1. Generate Compatibility Artifacts

Let's first choose a subset of amd64 containers. We are doing this primarily because that is my development host and I can't run anything
for a different arch. These are a subset from [lammps-matrix](https://github.com/rse-ops/lammps-matrix/pkgs/container/lammps-matrix).
Note that we have added a simple release workflow that generates [binaries](https://github.com/supercontainers/compspec-go/releases/tag/1-26-2024-2)
that we can easily grab.

```bash
docker pull ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64
```

Let's use the [generate-artifact.sh](generate-artifact.sh) script, and inside the container, to do exactly that. Note that we have to
source the /etc/profile to emulate what would happen on an entry to the container. Let's walk through how we'd do this with one container
first, generating the artifact and [pushing to oras](https://oras.land/docs/how_to_guides/pushing_and_pulling/).

```bash
# create directory for output specs
mkdir -p ./specs

# arguments are the hasGpu command and the path for the artifact (relative to PWD)
cmd=". /etc/profile && /tmp/data/generate-artifact.sh no /tmp/data/specs/compspec-intel-mpi-rocky-9-amd64.json"
docker run -v $PWD:/tmp/data -it ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64 /bin/bash -c "$cmd"

# This generates ./specs/compspec-intel-mpi-rocky-9-amd64.json, let's push to a registry with oras
oras push ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64-compspec --artifact-type application/org.supercontainers.compspec ./specs/compspec-intel-mpi-rocky-9-amd64.json
```

Here is how we might see it:

```bash
oras blob fetch --output - ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64-compspec@sha256:b68136afad3e4340f0dd4e09c5fea7faf12306cb4b0c1de616703b00d6ffef78
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
        "implementation": "intel-mpi",
        "version": "2021.8"
      }
    },
    {
      "name": "org.supercontainers.os",
      "version": "0.0.0",
      "annotations": {
        "name": "Rocky Linux 9.3 (Blue Onyx)",
        "release": "9.3",
        "vendor": "rocky",
        "version": "9.3"
      }
    },
    {
      "name": "org.supercontainers.hardware.gpu",
      "version": "0.0.0",
      "annotations": {
        "available": "no"
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

This is great! Next we will capture the URIs of these together in a manifest and put into our compspec tool. TBA!