package model

import (
	"time"
)

type User struct {
	Id            int64     `gorm:"column:id;primary_key;not null;auto_increment" json:"id"`                          //
	Nickname      string    `gorm:"column:nickname;not null" json:"nickname"`                                         //昵称
	Phone         string    `gorm:"column:phone;not null" json:"phone"`                                               //手机号
	Password      string    `gorm:"column:password;not null" json:"password"`                                         //密码
	Email         string    `gorm:"column:email;not null" json:"email"`                                               //邮箱
	Status        int       `gorm:"column:status;not null" json:"status"`                                             //状态；0=启用，1、禁用
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:current_timestamp" json:"created_at"`           //创建时间
	LastLoginTime time.Time `gorm:"column:last_login_time;not null;default:current_timestamp" json:"last_login_time"` //最后一次登录时间
	LastLoginIp   string    `gorm:"column:last_login_ip;not null" json:"last_login_ip"`                               //最后一次登录ip地址
	UpdateAt      time.Time `gorm:"column:update_at;not null;default:current_timestamp" json:"update_at"`             //修改时间
	Avatar        string    `gorm:"column:avatar;not null" json:"avatar"`                                             //头像
}

func (m *User) TableName() string {
	return "user"
}
