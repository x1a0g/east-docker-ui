package dto

type CreateRepoDto struct {
	RepoName     string `json:"repoName"`
	RepoAddr     string `json:"repoAddr"`
	RepoUsername string `json:"repoUsername"`
	RepoPassword string `json:"repoPassword"`
	RepoDesc     string `json:"repoDesc"`
}

type SearchRepoDto struct {
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type UpdateRepoDto struct {
	Id           string `json:"id"`
	RepoName     string `json:"repoName"`
	RepoAddr     string `json:"repoAddr"`
	RepoUsername string `json:"repoUsername"`
	RepoPassword string `json:"repoPassword"`
	RepoDesc     string `json:"repoDesc"`
}
