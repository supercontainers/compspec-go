package create

import (
	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/plugins/creators/cluster"
)

// Nodes will read in one or more node extraction metadata files and generate a single nodes JGF graph
// This is intentended for a registration command.
// TODO this should be converted to a creation (converter) plugin
func Nodes(nodesDir, clusterName, nodeOutFile string) error {

	// assemble options for node creator
	creator, err := cluster.NewPlugin()
	if err != nil {
		return err
	}
	options := plugin.PluginOptions{
		StrOpts: map[string]string{
			"nodes-dir":    nodesDir,
			"cluster-name": clusterName,
			"node-outfile": nodeOutFile,
		},
	}
	return creator.Create(options)
}
