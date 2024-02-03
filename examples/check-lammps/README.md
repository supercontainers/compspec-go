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
# TODO it looks like the arch is coming from my host, this would need to be run at build time alongside
# the machine it was built on. We'd also want to make sure it's documented this is the case
cmd=". /etc/profile && /tmp/data/generate-artifact.sh no /tmp/data/specs/compspec-intel-mpi-rocky-9-amd64.json"
docker run -v $PWD:/tmp/data -it ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64 /bin/bash -c "$cmd"

# This generates ./specs/compspec-intel-mpi-rocky-9-amd64.json, let's push to a registry with oras
oras push ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64-compspec --artifact-type application/org.supercontainers.compspec ./specs/compspec-intel-mpi-rocky-9-amd64.json:application/org.supercontainers.compspec
```

Here is how we might see it:

```bash
oras blob fetch --output - ghcr.io/rse-ops/lammps-matrix:intel-mpi-rocky-9-amd64-compspec@sha256:376ea8d492aa8e8db312cecc34bbc729d18fc2b30e891deb2ffdffa38c7db3a5
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
        "hardware.gpu.available": "no",
        "mpi.implementation": "intel-mpi",
        "mpi.version": "2021.8",
        "os.name": "Rocky Linux 9.3 (Blue Onyx)",
        "os.release": "9.3",
        "os.vendor": "rocky",
        "os.version": "9.3"
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

Let's run a quick script that will generate this for a few images so we can better prototype our matching (check) command:

```bash
./extract.sh
```

Note that arch reflects our host environment, which is an issue we need to figure out. This
also needs to be able to build for the right host arch. Next we will capture the URIs of these together in a manifest and put into our compspec tool.