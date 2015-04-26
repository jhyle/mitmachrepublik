package mmr

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	EventService struct {
		database       Database
		eventTableName string
		dateTableName  string
	}
)

func NewEventService(database Database, eventTableName, dateTableName string) (*EventService, error) {

	err := database.Table(eventTableName).EnsureIndices("organizerid", "start")
	if err == nil {
		err = database.Table(dateTableName).EnsureIndices("eventid", "organizerid", "start", "categories", "addr.city", "addr.pcode")
	}

	if err != nil {
		return nil, err
	} else {
		return &EventService{database, eventTableName, dateTableName}, nil
	}
}

func (events *EventService) eventTable() Table {

	return events.database.Table(events.eventTableName)
}

func (events *EventService) dateTable() Table {

	return events.database.Table(events.dateTableName)
}

func (events *EventService) Cities() ([]string, error) {

	var cities []string
	err := events.eventTable().Distinct(bson.M{}, "addr.city", &cities)
	return cities, err
}

func (events *EventService) Count(place string, dates [][]time.Time, categoryIds []int) (int, error) {

	query := buildQuery(place, dates, categoryIds)
	return events.dateTable().Count(query)
}

func (events *EventService) Load(id bson.ObjectId) (*Event, error) {

	var event Event
	err := events.eventTable().LoadById(id, &event)
	return &event, err
}

func (events *EventService) LoadDate(id bson.ObjectId) (*Date, error) {

	var date Date
	err := events.dateTable().LoadById(id, &date)
	return &date, err
}

func (events *EventService) SearchDates(place string, dates [][]time.Time, categoryIds []int, page, pageSize int, sort string) (DateSearchResult, error) {

	var result DateSearchResult
	query := buildQuery(place, dates, categoryIds)
	err := events.dateTable().Search(query, page*pageSize, pageSize, &result, sort)
	return result, err
}

func (events *EventService) SearchDatesOfUser(userId bson.ObjectId, page, pageSize int, sort string) (*DateSearchResult, error) {

	var result DateSearchResult
	query := bson.M{"$and": []bson.M{bson.M{"organizerid": userId}, bson.M{"start": bson.M{"$gte": time.Now()}}}}
	err := events.dateTable().Search(query, page*pageSize, pageSize, &result, sort)
	return &result, err
}

func (events *EventService) SearchEventsOfUser(userId bson.ObjectId, page, pageSize int, sort string) (*EventSearchResult, error) {

	var result EventSearchResult
	err := events.eventTable().Search(bson.M{"organizerid": userId}, page*pageSize, pageSize, &result, sort)
	return &result, err
}

func (events *EventService) Store(event *Event) error {

	_, err := events.eventTable().UpsertById(event.Id, event)
	// TODO update or create dates
	return err
}

func (events *EventService) DeleteOfUser(userId bson.ObjectId) error {

	query := bson.M{"organizerid": userId}
	err := events.dateTable().Delete(query)
	if err == nil {
		err = events.eventTable().Delete(query)
	}
	return err
}

func (events *EventService) Delete(id bson.ObjectId) error {

	err := events.dateTable().Delete(bson.M{"eventid": id})
	if err == nil {
		err = events.eventTable().DeleteById(id)
	}
	return err
}
