package vo

import "east-docker-ui/model"

type LoginVo struct {
	Token    string          `json:"token"`
	Userinfo *model.UserInfo `json:"userinfo"`
}
