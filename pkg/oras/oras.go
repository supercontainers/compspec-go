package oras

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/supercontainers/compspec-go/pkg/types"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
	"sigs.k8s.io/yaml"
)

// Oras will provide an interface to retrieve an artifact, specifically
// a compatibillity spec artifact media type
// LoadArtifact retrieves the artifact from the url string
// and returns based on the media type
func LoadArtifact(uri string, mediaType string) (types.CompatibilityRequest, error) {
	request := types.CompatibilityRequest{}
	ctx := context.Background()
	repo, err := remote.NewRepository(uri)
	if err != nil {
		return request, err
	}

	// Disable plain http for now, assume using production registry
	repo.PlainHTTP = false

	// Split reference into tag, or assume latest
	tag := "latest"
	if strings.Contains(uri, ":") {
		parts := strings.Split(uri, ":")
		tag = parts[1]
		uri = parts[0]
	}

	// Fetch manifest for the tag
	desc, readCloser, err := repo.FetchReference(ctx, tag)
	if err != nil {
		return request, err
	}
	defer readCloser.Close()

	// Read the pulled content
	manifestBytes, err := content.ReadAll(readCloser, desc)
	if err != nil {
		return request, err
	}

	// Going to be a big wild here and not check the mnaifest media type.
	// We'd want to find oras, but no reason it can't be pushed another way...
	// unmarshall it
	var manifest oci.Manifest
	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		return request, err
	}

	// Loop through layers and find the media type we are looking for
	for _, layer := range manifest.Layers {

		// Skip layers that are not the compatibility spec... we seek
		if layer.MediaType != mediaType {
			continue
		}

		// Get the descriptor for the digest we want
		desc, err := repo.Blobs().Resolve(ctx, string(layer.Digest))
		if err != nil {
			return request, err
		}

		// Download using the descriptor
		readCloser, err := repo.Fetch(ctx, desc)
		if err != nil {
			return request, err
		}
		defer readCloser.Close()

		// Read the descriptor into bytes
		vr := content.NewVerifyReader(readCloser, desc)
		buffer := bytes.NewBuffer(nil)
		_, err = io.Copy(buffer, vr)
		if err != nil {
			return request, err
		}

		// note: users should not trust the the read content until Verify returns nil
		if err := vr.Verify(); err != nil {
			return request, err
		}

		// Convert this into our Compatibility Request
		// reading from the buffer into bytes proper
		readContent, err := io.ReadAll(buffer)
		if err != nil {
			return request, err
		}

		err = yaml.Unmarshal(readContent, &request)
		if err != nil {
			return request, err
		}
	}
	return request, nil
}
