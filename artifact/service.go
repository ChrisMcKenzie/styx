package artifact

type Artifact struct {
	// ...
}

type ArtifactService interface {
	UploadArtifact(*Artifact) error
}
