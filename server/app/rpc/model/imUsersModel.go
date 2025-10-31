package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
)

var _ ImUsersModel = (*customImUsersModel)(nil)

type (
	// ImUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customImUsersModel.
	ImUsersModel interface {
		imUsersModel
		customImUsersLogicModel
	}

	customImUsersModel struct {
		*defaultImUsersModel
	}

	customImUsersLogicModel interface {
	}
)

// NewImUsersModel returns a model for the database table.
func NewImUsersModel(conn *gorm.DB, c cache.CacheConf) ImUsersModel {
	return &customImUsersModel{
		defaultImUsersModel: newImUsersModel(conn, c),
	}
}

func (m *defaultImUsersModel) customCacheKeys(data *ImUsers) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}
