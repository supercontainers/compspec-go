package create

import (
	"strings"

	"github.com/compspec/compspec-go/plugins/creators/artifact"
)

// Artifact will create a compatibility artifact based on a request in YAML
// TODO likely want to refactor this into a proper create plugin
func Artifact(specname string, fields []string, saveto string, allowFail bool) error {

	// This is janky, oh well
	allowFailFlag := "false"
	if allowFail {
		allowFailFlag = "true"
	}

	// assemble options for node creator
	creator, err := artifact.NewPlugin()
	if err != nil {
		return err
	}
	options := map[string]string{
		"specname":  specname,
		"fields":    strings.Join(fields, "||"),
		"saveto":    saveto,
		"allowFail": allowFailFlag,
	}
	return creator.Create(options)
}
