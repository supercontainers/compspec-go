package kernel

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// A module holds a version and parameters, and helper functions to derive them.
type Module struct {
	Name       string
	Path       string
	Version    string
	Parameters map[string]string
}

// NewModule prepares a module for a path
func NewModule(name, defaultVersion string) *Module {

	// Get the name, and then we can read the fullpath
	fullPath := filepath.Join(kernelModules, name)
	module := Module{Name: name, Path: fullPath}
	module.setVersion(defaultVersion)
	return &module
}

// Derive the version of the module
func (m *Module) setVersion(version string) {
	v, err := os.ReadFile(filepath.Join(m.Path, "version"))

	// If we don't have an error, use the derived version
	if err == nil {
		version = string(v)
	}
	m.Version = strings.TrimSpace(version)
}

// Key is the prefix key for the higher level modules metadata
func (m *Module) Key() string {
	return fmt.Sprintf("module.%s", m.Name)
}

// parameterPath will return the full path to parameters
func (m *Module) parameterPath() string {
	return filepath.Join(m.Path, "parameters")
}

// setParameters parses the module root / parameters directory
func (m *Module) SetParameters() error {

	// parameters are in the module directory here
	params, err := os.ReadDir(m.parameterPath())
	if err != nil {

		// OK if no parameters
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Make exactly what we need
	data := make(map[string]string, len(params))

	for _, param := range params {
		name := param.Name()
		content, err := m.readParameterFile(name)
		if err != nil {
			return err
		}
		data[name] = content
	}
	m.Parameters = data
	return nil
}

// readParameterFile reads a parameter file and returns the value
func (m *Module) readParameterFile(name string) (string, error) {

	content, err := os.ReadFile(filepath.Join(m.parameterPath(), name))
	if err != nil {
		var pathErr *fs.PathError

		// TODO how often does this happen?
		if errors.As(err, &pathErr) && pathErr.Err == unix.EPERM || pathErr.Err == unix.EACCES {
			fmt.Printf("Error: cannot read parameter path because of EPERM and EACCES: %s\n", err)
		}
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
