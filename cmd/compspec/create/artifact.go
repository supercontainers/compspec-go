package create

import (
	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/plugins/creators/artifact"
)

// Artifact will create a compatibility artifact based on a request in YAML
// TODO likely want to refactor this into a proper create plugin
func Artifact(specname string, fields []string, saveto string, allowFail bool) error {

	// assemble options for node creator
	creator, err := artifact.NewPlugin()
	if err != nil {
		return err
	}

	options := plugin.PluginOptions{
		StrOpts: map[string]string{
			"specname": specname,
			"saveto":   saveto,
		},
		BoolOpts: map[string]bool{
			"allowFail": allowFail,
		},
		ListOpts: map[string][]string{
			"fields": fields,
		},
	}
	return creator.Create(options)
}
