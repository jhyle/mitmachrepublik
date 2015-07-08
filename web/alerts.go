package mmr

import (
	"labix.org/v2/mgo/bson"
	"time"
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

func (alerts *AlertService) Load(id bson.ObjectId) (*Alert, error) {

	var alert Alert
	err := alerts.table().LoadById(id, &alert)
	return &alert, err
}

func (alerts *AlertService) FindAlerts(weekday time.Weekday) ([]Alert, error) {

	var result []Alert
	err := alerts.table().Find(bson.M{"weekdays": weekday}, &result, "id")
	return result, err
}

func (alerts *AlertService) Store(alert *Alert) error {

	_, err := alerts.table().UpsertById(alert.Id, alert)
	return err
}

func (alerts *AlertService) Delete(id bson.ObjectId) error {

	return alerts.table().DeleteById(id)
}
