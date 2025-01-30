package model

import "gorm.io/gorm"

type UserLog struct {
	Id         int    `json:"id" gorm:"column:id"`
	ActionName string `json:"actionName" gorm:"column:action_name"`
	ActionUser string `json:"actionUser" gorm:"column:action_user"`
	CreateAt   int64  `json:"createAt" gorm:"column:create_time"`
}

func (*UserLog) TableName() string {
	return "user_log"
}

func (u *UserLog) Top5(db *gorm.DB) []UserLog {
	var res []UserLog
	db.Order("create_time desc").Limit(5).Find(&res)
	return res
}
