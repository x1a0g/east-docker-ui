package dto

type SearchImageDto struct {
	Keyword string `json:"keyword"`
}

type CreateImageDto struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	RepoId  string `json:"repoId"`
}

type ImportImageDto struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	FileId  string `json:"fileId"`
}
