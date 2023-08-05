package mapper

import (
	"github.com/google/wire"
	"github.com/xh-polaris/platform-comment/biz/infrastructure/config"
	"github.com/zeromicro/go-zero/core/stores/mon"
)

const HistoryCollectionName = "history"

var _ HistoryModel = (*customHistoryModel)(nil)

type (
	// HistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHistoryModel.
	HistoryModel interface {
		historyModel
	}

	customHistoryModel struct {
		*defaultHistoryModel
	}
)

// NewHistoryModel returns a model for the mongo.
func NewHistoryModel(config *config.Config) HistoryModel {
	conn := mon.MustNewModel(config.Mongo.URL, config.Mongo.DB, HistoryCollectionName)
	return &customHistoryModel{
		defaultHistoryModel: newDefaultHistoryModel(conn),
	}
}

var HistorySet = wire.NewSet(
	NewHistoryModel,
)
