package model

import "gorm.io/gorm"

type UserInfo struct {
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Email    string `json:"email" gorm:"column:email"`
	Id       int32  `json:"id" gorm:"column:id"`
	CreateAt string `json:"create_at" gorm:"column:create_time"`
}

func (*UserInfo) TableName() string {
	return "user_info"
}

func (u *UserInfo) CheckUserNameAndPass(db *gorm.DB, username string, pass string) error {

	err := db.Where("username = ? and password=?", username, pass).First(&u).Error

	if err != nil {
		return err
	}

	return nil
}
