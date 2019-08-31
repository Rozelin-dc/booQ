package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// User userの構造体
type User struct {
	gorm.Model
	Name        string `gorm:"type:varchar(32);unique;not null;size:50" json:"name"`
	DisplayName string `gorm:"type:varchar(64);not null" json:"displayName"`
	IconFileID  string `gorm:"type:varchar(36);not null" json:"iconFileId"`
	Admin       bool   `gorm:"default:false" json:"admin"`
}

// TableName dbのテーブル名を指定する
func (user *User) TableName() string {
	return "users"
}

// GetUser userを取得する
func GetUser(user User) (User, error) {
	if user.Name == "" {
		return User{}, errors.New("Nameが存在しません")
	}
	res := User{}
	db.Where("name = ?", user.Name).First(&res)
	return res, nil
}

// GetUsers 全userを取得する
func GetUsers() []User {
	res := []User{}
	db.Find(&res)
	return res
}

// GetUserByName userをNameから取得する
func GetUserByName(name string) (User, error) {
	res := User{}
	db.Where("name = ?", name).First(&res)
	if res.Name == "" {
		return User{}, errors.New("Nameが不正です")
	}
	return res, nil
}

// GetUserByID userをIDから取得する
func GetUserByID(id int) (User, error) {
	res := User{}
	db.Where("id = ?", id).First(&res)
	if res.Name == "" {
		return User{}, errors.New("Nameが不正です")
	}
	return res, nil
}

// CreateUser userを作成する
func CreateUser(user User) (User, error) {
	if user.Name == "" {
		return User{}, errors.New("Nameが存在しません")
	}
	db.Create(&user)
	return user, nil
}

// UpdateUser userの情報(表示される名前やアイコン、管理者権限)の変更
func UpdateUser(newUser User) (User, error) {
	res := User{}
	if newUser.Name == "" {
		return User{}, errors.New("Nameが存在しません")
	}
	db.Model(&res).Where("name = ?", newUser.Name).Updates(newUser)
	return res, nil
}
