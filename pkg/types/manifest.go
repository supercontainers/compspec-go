package types

// A ManifestList is a manual definition of images and paired artifacts
type ManifestList struct {
	Images []ImagePair `json:"images"`
}

type ImagePair struct {
	Name     string `json:"name"`
	Artifact string `json:"artifact"`
}
