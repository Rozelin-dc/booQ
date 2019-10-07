package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogTableName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "logs", (&Log{}).TableName())
}

func TestCreateLog(t *testing.T) {
	t.Parallel()

	t.Run("failures", func(t *testing.T) {
		assert := assert.New(t)

		log, err := CreateLog(Log{})
		assert.Error(err)
		assert.Empty(log)

		log, err = CreateLog(Log{ItemID: 66, OwnerId: 66})
		assert.Error(err)
		assert.Empty(log)
	})

	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		owner, _ := GetUserByName("traP")
		item, _ := CreateItem(Item{Name: "testItemForCreateLog"})

		log, err := CreateLog(Log{ItemID: item.ID, OwnerId: owner.ID, Type: 0})
		assert.NoError(err)
		assert.NotEmpty(log)
		assert.Equal(owner.ID, log.OwnerId)
		assert.Equal(item.ID, log.ItemID)
	})
}

func TestGetLogsByItemID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		user, err := CreateUser(User{Name: "testGetLogsByItemIDUser"})
		assert.NoError(err)
		owner, err := CreateUser(User{Name: "testGetLogsByItemIDOwner"})
		assert.NoError(err)
		item, err := CreateItem(Item{Name: "testItemForGetLogsByItemID"})
		assert.NoError(err)
		_, err = RegisterOwner(Owner{UserId: owner.ID, Rentalable: true, Count: 1}, item)
		assert.NoError(err)
		log, err := CreateLog(Log{ItemID: item.ID, OwnerId: owner.ID, UserId: user.ID, Type: 0, Count: 0})
		assert.NoError(err)
		logs, err := GetLogsByItemID(item.ID)

		assert.NoError(err)
		assert.NotEmpty(logs)
		assert.Equal(logs[0].ItemID, log.ItemID)
		assert.Equal(logs[0].OwnerId, log.OwnerId)
		assert.Equal(logs[0].Owner.Name, owner.Name)
		assert.Equal(logs[0].User.Name, user.Name)
	})
}

func TestGetLatestLogs(t *testing.T) {
	assert := assert.New(t)
	item, err := CreateItem(Item{Name: "testGetLatestLogItem"})
	assert.NoError(err)
	itemID := item.ID
	user1, err := CreateUser(User{Name: "testGetUserLatestLogUser1"})
	assert.NoError(err)
	user2, err := CreateUser(User{Name: "testGetUserLatestLogUser2"})
	assert.NoError(err)
	ownerUser1, err := GetUserByName("traP")
	assert.NoError(err)
	owner1 := Owner{
		UserId:     ownerUser1.ID,
		Rentalable: true,
		Count:      1,
	}
	ownerUser2, _ := GetUserByName("sienka")
	owner2 := Owner{
		UserId:     ownerUser2.ID,
		Rentalable: true,
		Count:      1,
	}

	t.Run("failures", func(t *testing.T) {
		log, err := GetLatestLog(66, 66)
		assert.Error(err)
		assert.Empty(log)

		log, err = GetLatestLog(itemID, 66)
		assert.Error(err)
		assert.Empty(log)
	})

	t.Run("success", func(t *testing.T) {
		_, err := RegisterOwner(owner1, item)
		assert.NoError(err)
		_, err = RegisterOwner(owner2, item)
		assert.NoError(err)
		_, err = CreateLog(Log{ItemID: itemID, UserId: user1.ID, OwnerId: ownerUser1.ID, Type: 0, Count: 0})
		assert.NoError(err)
		_, err = CreateLog(Log{ItemID: itemID, UserId: user2.ID, OwnerId: ownerUser1.ID, Type: 0, Count: 1})
		assert.NoError(err)
		_, err = CreateLog(Log{ItemID: itemID, UserId: user2.ID, OwnerId: ownerUser2.ID, Type: 0, Count: 0})
		assert.NoError(err)
		_, err = CreateLog(Log{ItemID: itemID, UserId: user1.ID, OwnerId: ownerUser2.ID, Type: 0, Count: 1})
		assert.NoError(err)

		logs, err := GetLatestLogs(itemID)
		assert.NoError(err)
		assert.NotEmpty(logs)
		for _, log := range logs {
			if log.OwnerId == ownerUser1.ID {
				assert.Equal(ownerUser1.Name, log.Owner.Name)
				assert.Equal(user2.ID, log.UserId)
				assert.Equal(user2.Name, log.User.Name)
				assert.Equal(1, log.Count)
			}
			if log.OwnerId == ownerUser2.ID {
				assert.Equal(ownerUser2.Name, log.Owner.Name)
				assert.Equal(user1.ID, log.UserId)
				assert.Equal(user1.Name, log.User.Name)
				assert.Equal(1, log.Count)
			}
		}
	})
}
