# Design

The compatibility tool is responsible for extracting information about a system, and comparing the host metadata to a contender container image to determine if it is compatible. This is a two step process that includes:

1. Extracting metadata about the container image at build time
2. Extracting metadata about the host at image selection time, and comparing against a set of contender container images to select the best one.

## Definitions

### Extractor

An **extractor** is a core plugin that knows how to retrieve metadata about a host. An extractor is ususally going to be run for two cases:

1. During CI to extract (and save) metadata about a particular build to put in a compatibility artifact.
2. During image selection to extract information about the host to compare to.

Examples extractors could be "library" or "system."

### Section

A **section** is a group of metadata within an extractor. For example, within "library" a section is for "mpi." This allows a user to specify running the `--name library[mpi]` extractor to ask for the mpi section of the library family. Another example is under kernel.
The user might want to ask for more than one group to be extracted and might ask for `--name kernel[boot,config]`. Section basically provides more granularity to an extractor namespace. For the above two examples, the metadata generated would be organized like:

```
library
   mpi.<attribute>
kernel
  config.<attribute>
  boot.<attribute>
```

For the above, right now I am implementing extractors generally, or "wild-westy" in the sense that the namespace is oriented toward the extractor name and sections it owns (e.g., no community namespaces like archspec, spack, opencontainers, etc). This is subject to change depending on the design the working group decides on.

### Creator

A creator is a plugin that is responsible for creating an artifact that includes some extracted metadata. The creator is agnostic to what it it being asked to generate in the sense that it just needs a mapping. The mapping will be from the extractor namespace to the compatibility artifact namespace. For our first prototype, this just means asking for particular extractor attributes to map to a set of annotations that we want to dump into json. To start there should only be one creator plugin needed, however if there are different structures of artifacts needed, I could imagine more. An example creation specification for a prototype experiment where we care about architecture, MPI, and GPU is provided in [examples](examples).

## Overview

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