package mmr

import (
//	"errors"
//	"labix.org/v2/mgo/bson"
)

type (
	AlertService struct {
		database  Database
		tablename string
	}
)

func NewAlertService(database Database, tablename string) (*AlertService, error) {

	err := database.Table(tablename).DropIndices()
	if err != nil {
		return nil, err
	}

	err = database.Table(tablename).EnsureIndices([][]string{
		{"weekdays"},
		{"email"},
	})
	if err != nil {
		return nil, err
	}

	return &AlertService{database, tablename}, nil
}

func (alerts *AlertService) table() Table {

	return alerts.database.Table(alerts.tablename)
}

func (alerts *AlertService) Store(alert *Alert) error {

	_, err := alerts.table().UpsertById(alert.Id, alert)
	return err
}
