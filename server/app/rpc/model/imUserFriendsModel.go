package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
)

var _ ImUserFriendsModel = (*customImUserFriendsModel)(nil)

type (
	// ImUserFriendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customImUserFriendsModel.
	ImUserFriendsModel interface {
		imUserFriendsModel
		customImUserFriendsLogicModel
	}

	customImUserFriendsModel struct {
		*defaultImUserFriendsModel
	}

	customImUserFriendsLogicModel interface {
	}
)

// NewImUserFriendsModel returns a model for the database table.
func NewImUserFriendsModel(conn *gorm.DB, c cache.CacheConf) ImUserFriendsModel {
	return &customImUserFriendsModel{
		defaultImUserFriendsModel: newImUserFriendsModel(conn, c),
	}
}

func (m *defaultImUserFriendsModel) customCacheKeys(data *ImUserFriends) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}
