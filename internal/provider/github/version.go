package github

type PackageInfo struct {
	Metadata PackageInfoMetadata `json:"metadata"`
}

type PackageInfoMetadata struct {
	Container PackageInfoContainer `json:"container"`
}

type PackageInfoContainer struct {
	Tags []string `json:"tags"`
}
