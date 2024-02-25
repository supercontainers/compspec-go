# Rainbow Scheduler

The [rainbow scheduler](https://github.com/converged-computing/rainbow) has a registration step that requires a cluster to send over node metadata. The reason is because when a user sends a request for work, the scheduler needs to understand
how to properly assign it. To do that, it needs to be able to see all the resources (clusters) available to it.

![../img/rainbow-scheduler-register.png](../img/rainbow-scheduler-register.png)

For the purposes of compspec here, we care about the registration step. This is what that includes:

## Registration

1. At registration, the cluster also sends over metadata about itself (and the nodes it has). This is going to allow for selection for those nodes. 
1. When submitting a job, the user no longer is giving an exact command, but a command + an image with compatibility metadata. The compatibility metadata (somehow) needs to be used to inform the cluster selection.
1. At selection, the rainbow schdeuler needs to filter down cluster options, and choose a subset.
 - Level 1: Don't ask, just choose the top choice and submit
 - Level 2: Ask the cluster for TBA time or cost, choose based on that.
 - Job is added to that queue.

Specifically, this means two steps for compspec go:

1. A step to ask each node to extract it's own metadata, saved to a directory.
2. A second step to combine those nodes into a graph.

Likely we will take a simple approach to do an extract for one node that captures it's metadata into Json Graph Format (JGF) and then dumps into a shared directory (we might imagine this being run with a flux job)
and then some combination step.

## Example

In the example below, we will extract node level metadata with `compspec extract` and then generate the cluster JGF to send for registration with compspec create-nodes.

### 1. Extract Metadata

Let's first generate faux node metadata for a "cluster" - I will just run an extraction a few times and generate equivalent files :) This isn't such a crazy idea because it emulates nodes that are the same!

```bash
mkdir -p ./docs/rainbow/cluster
compspec extract --name library --name nfd[cpu,memory,network,storage,system] --name system[cpu,processor,arch,memory] --out ./docs/rainbow/cluster/node-1.json
compspec extract --name library --name nfd[cpu,memory,network,storage,system] --name system[cpu,processor,arch,memory] --out ./docs/rainbow/cluster/node-2.json
compspec extract --name library --name nfd[cpu,memory,network,storage,system] --name system[cpu,processor,arch,memory] --out ./docs/rainbow/cluster/node-3.json
```

### 2. Create Nodes

Now we are going to give compspec the directory, and ask it to create nodes. This will be in JSON graph format. This outputs to the terminal:

```bash
compspec create nodes --cluster-name cluster-red --node-dir ./docs/rainbow/cluster/
```