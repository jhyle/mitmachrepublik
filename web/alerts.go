package mmr

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
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
		return nil, errors.Wrap(err, "error dropping indices of alert service database")
	}

	err = database.Table(tablename).EnsureIndices([][]string{
		{"weekdays"},
		{"email"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating indices of alert service database")
	}

	return &AlertService{database, tablename}, nil
}

func (alerts *AlertService) table() Table {

	return alerts.database.Table(alerts.tablename)
}

func (alerts *AlertService) Load(id bson.ObjectId) (*Alert, error) {

	var alert Alert
	err := alerts.table().LoadById(id, &alert)

	if err != nil {
		return nil, errors.Wrapf(err, "error loading alert with id: %s", id.String())
	}

	return &alert, nil
}

func (alerts *AlertService) FindAlerts(weekday time.Weekday) ([]Alert, error) {

	var result []Alert
	err := alerts.table().Find(bson.M{"weekdays": weekday}, &result, "id")

	if err != nil {
		return nil, errors.Wrapf(err, "error loading alerts for weekday: %d", weekday)
	}

	return result, nil
}

func (alerts *AlertService) Store(alert *Alert) error {

	_, err := alerts.table().UpsertById(alert.Id, alert)

	if err != nil {
		return errors.Wrapf(err, "error storing alert: %+v", *alert)
	}

	return nil
}

func (alerts *AlertService) Delete(id bson.ObjectId) error {

	err := alerts.table().DeleteById(id)

	if err != nil {
		return errors.Wrapf(err, "error deleting alert with id: %s", id.String())
	}

	return nil
}
