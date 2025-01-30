package model

type DockerImageInfo struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	RepoTags []string `json:"repoTags"`
	Created  string   `json:"created"`
	Size     string   `json:"size"`
}
