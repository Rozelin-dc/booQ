package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Log logの構造体
type Log struct {
	gorm.Model
	ItemID  uint      `gorm:"type:int;not null" json:"item_id"`
	UserId  uint      `gorm:"type:int;not null" json:"user_id"`
	User    User      `json:"user"`
	OwnerId uint      `gorm:"type:int;not null" json:"owner_id"`
	Owner   User      `json:"owner"`
	Type    int       `gorm:"type:int;not null" json:"type"`
	Purpose string    `json:"purpose"`
	DueDate time.Time `gorm:"type:datetime;" json:"due_date"`
	Count   int       `gorm:"type:int;not null" json:"count"`
}

type RequestPostLogsBody struct {
	OwnerID uint      `json:"owner_id"`
	Type    int       `json:"type"`
	Purpose string    `json:"purpose"`
	DueDate time.Time `json:"due_date"`
	Count   int       `json:"count"`
}

// TableName dbのテーブル名を指定する
func (log *Log) TableName() string {
	return "logs"
}

// CreateLog logを作成する
func CreateLog(log Log) (Log, error) {
	if log.ItemID == 0 {
		return Log{}, errors.New("ItemIDが存在しません")
	}
	_, err := GetItemByID(log.ItemID)
	if err != nil {
		return Log{}, errors.New("Itemが存在しません")
	}
	_, err = GetUserByID(int(log.OwnerId))
	if err != nil {
		return Log{}, errors.New("Ownerが存在しません")
	}
	db.Create(&log)
	return log, nil
}

// GetLatestLog ownerIDからLatestLogを取得する
func GetLatestLog(itemID, ownerID uint) (Log, error) {
	item, err := GetItemByID(itemID)
	if err != nil {
		return Log{}, err
	}
	exist := false
	for _, owner := range item.Owners {
		if owner.UserId == ownerID {
			exist = true
		}
	}
	if !exist {
		return Log{}, errors.New("指定した所有者はそのItemを所有していません")
	}
	log := Log{}
	db.Set("gorm:auto_preload", true).Order("created_at desc").First(&log, "item_id = ? AND owner_id = ?", itemID, ownerID)
	return log, nil
}

// GetLatestLogs 各所有者ごとの最新のログを取得する
func GetLatestLogs(itemID uint) ([]Log, error) {
	item := Item{}
	item.ID = itemID
	db.Set("gorm:auto_preload", true).First(&item).Related(&item.Owners, "Owners")
	if item.Name == "" {
		return []Log{}, errors.New("該当するItemがありません")
	}
	logs := []Log{}
	log := Log{}
	log.ItemID = itemID
	for _, owner := range item.Owners {
		log.OwnerId = owner.ID
		db.Set("gorm:auto_preload", true).Order("created_at desc").First(&log)
		if log.ID == 0 {
			continue
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// GetLogsByItemID itemIDからLogsを取得する
func GetLogsByItemID(itemID uint) ([]Log, error) {
	// 指定のitemIDのitemが存在するかどうかはここで判別つけていません
	logs := []Log{}
	db.Set("gorm:auto_preload", true).Find(&logs, "item_id = ?", itemID)
	return logs, nil
}
