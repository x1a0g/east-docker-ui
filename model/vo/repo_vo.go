package vo

import "east-docker-ui/model"

type RepoListVo struct {
	Total    int64              `json:"total"`
	RepoList []model.DockerRepo `json:"repoList"`
}

type RepoDownVo struct {
	Id string `json:"id"`
}
