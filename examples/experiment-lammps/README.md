# Experiment LAMMPS

This is a second stage of [check-lammps](../check-lammps) that has a larger matrix of images to check.
We have already pushed the images and artifacts to registry, represented in [manifests.yaml](manifests.yaml)
and will just test running compspec-go. Let's first check that our artifacts all exist (and our manifests list is OK)

```bash
../../bin/compspec match -i ./manifests.yaml --check-artifacts
```
```console
Checking artifacts complete. There were 0 artifacts missing.
```

