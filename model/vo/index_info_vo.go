package vo

type IndexInfoVo struct {
	HostName      string `json:"hostName"`
	OsVersion     string `json:"OsVersion"`
	DockerVersion string `json:"dockerVersion"`
	PtVersion     string `json:"PtVersion"`
	Who           string `json:"who"`
}

type StaticCount struct {
	TotalContainers   int `json:"totalContainers"`
	ContainersRunning int `json:"containersRunning"`
	ContainersStopped int `json:"ContainersStopped"`
	ContainersPaused  int `json:"ContainersPaused"`
}
