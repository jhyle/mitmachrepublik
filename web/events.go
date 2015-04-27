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

func (events *EventService) buildQuery(place string, dates [][]time.Time, categoryIds []int) bson.M {

	query := make([]bson.M, 0, 3)

	if len(place) > 0 {
		postcodes := Postcodes(place)
		placesQuery := make([]bson.M, len(postcodes)+1)
		for i, postcode := range postcodes {
			placesQuery[i] = bson.M{"event.addr.pcode": postcode}
		}
		placesQuery[len(postcodes)] = bson.M{"event.addr.city": place}
		query = append(query, bson.M{"$or": placesQuery})
	}

	if len(dates) > 0 {
		datesQuery := make([]bson.M, len(dates))
		for i, timespan := range dates {
			rangeQuery := make(bson.M)
			rangeQuery["$gte"] = timespan[0]
			if timespan[1] != timespan[0] {
				rangeQuery["$lt"] = timespan[1]
			}
			datesQuery[i] = bson.M{"event.start": rangeQuery}
		}
		query = append(query, bson.M{"$or": datesQuery})
	}

	if len(categoryIds) > 0 && categoryIds[0] != 0 {
		categoriesQuery := make([]bson.M, len(categoryIds))
		for i, categoryId := range categoryIds {
			categoriesQuery[i] = bson.M{"event.categories": categoryId}
		}
		query = append(query, bson.M{"$or": categoriesQuery})
	}

	return bson.M{"$and": query}
}

func (events *EventService) Count(place string, dates [][]time.Time, categoryIds []int) (int, error) {

	return events.dateTable().Count(events.buildQuery(place, dates, categoryIds))
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
	err := events.dateTable().Search(events.buildQuery(place, dates, categoryIds), page*pageSize, pageSize, &result, sort)
	return result, err
}

func (events *EventService) SearchDatesOfUser(userId bson.ObjectId, page, pageSize int, sort string) (*DateSearchResult, error) {

	var result DateSearchResult
	query := bson.M{"$and": []bson.M{bson.M{"event.organizerid": userId}, bson.M{"event.start": bson.M{"$gte": time.Now()}}}}
	err := events.dateTable().Search(query, page*pageSize, pageSize, &result, sort)
	return &result, err
}

func (events *EventService) FindEventsOfUser(userId bson.ObjectId, sort string) ([]Event, error) {

	var result []Event
	err := events.eventTable().Find(bson.M{"organizerid": userId}, &result, sort)
	return result, err
}

func (events *EventService) SearchEventsOfUser(userId bson.ObjectId, page, pageSize int, sort string) (*EventSearchResult, error) {

	var result EventSearchResult
	err := events.eventTable().Search(bson.M{"organizerid": userId}, page*pageSize, pageSize, &result, sort)
	return &result, err
}

func (events *EventService) generateDates(event *Event, now time.Time) []Date {

	dates := make([]Date, 0)

	var date Date
	date.Id = bson.NewObjectId()
	date.OrganizerId = event.OrganizerId
	date.EventId = event.Id
	date.Title = event.Title
	date.Start = event.Start
	date.End = event.End
	date.Image = event.Image
	date.Categories = event.Categories
	date.Descr = event.Descr
	date.Rsvp = event.Rsvp
	date.Web = event.Web
	date.Addr = event.Addr

	if !date.Start.Before(now) {
		dates = append(dates, date)
	}

	return dates
}

func (events *EventService) Store(event *Event, syncDates bool) error {

	_, err := events.eventTable().UpsertById(event.Id, event)

	if err == nil && syncDates {
		err = events.SyncDates(event)
	}

	return err
}

func (events *EventService) SyncDates(event *Event) error {

	now := time.Now()
	genDates := events.generateDates(event, now)

	var dates []Date
	query := bson.M{"$and": []bson.M{bson.M{"eventid": event.Id}, bson.M{"event.start": bson.M{"$gte": now}}}}
	err := events.dateTable().Find(query, &dates, "start")
	if err != nil {
		return err
	}

	n := min(len(genDates), len(dates))
	for i := 0; i < n; i++ {
		genDates[i].Id = dates[i].Id
		_, err := events.dateTable().UpsertById(genDates[i].Id, &genDates[i])
		if err != nil {
			return err
		}
	}

	for i := n; i < len(genDates); i++ {
		_, err := events.dateTable().UpsertById(genDates[i].Id, &genDates[i])
		if err != nil {
			return err
		}
	}

	for i := n; i < len(dates); i++ {
		err := events.dateTable().DeleteById(dates[i].Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (events *EventService) DeleteDatesOfUser(userId bson.ObjectId) error {

	query := bson.M{"event.organizerid": userId}
	return events.dateTable().Delete(query)
}

func (events *EventService) DeleteOfUser(userId bson.ObjectId) error {

	err := events.DeleteDatesOfUser(userId)
	if err == nil {
		query := bson.M{"organizerid": userId}
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
