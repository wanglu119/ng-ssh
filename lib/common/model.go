package common

import (

)

//Machine
type Machine struct {
	Name     string `json:"name" gorm:"type:varchar(50);unique_index"`
	Host     string `json:"host" gorm:"type:varchar(50)"`
	Ip       string `json:"ip" gorm:"type:varchar(80)"`
	Port     uint   `json:"port" gorm:"type:int(6)"`
	User     string `json:"user" gorm:"type:varchar(20)"`
	Password string `json:"password,omitempty"`
	Key      string `json:"key,omitempty"`
	Type     string `json:"type" gorm:"type:varchar(20)"`
}


// User
type User struct {
	Username string `json:"username,omitempty"`
	Token string `json:"token,omitempty"`
}