package mmr

import (
	"github.com/blevesearch/bleve"
	"github.com/pilu/traffic"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

type (
	EventService struct {
		database       Database
		eventTableName string
		dateTableName  string
		eventIndex     bleve.Index
		dateIndex      bleve.Index
	}
)

func NewEventService(database Database, eventTableName, dateTableName, indexDir string) (*EventService, error) {

	err := database.Table(eventTableName).DropIndices()
	if err != nil {
		return nil, err
	}

	err = database.Table(eventTableName).EnsureIndices([][]string{
		{"organizerid", "start"},
		{"recurrency"},
	})
	if err != nil {
		return nil, err
	}

	err = database.Table(dateTableName).DropIndices()
	if err != nil {
		return nil, err
	}

	err = database.Table(dateTableName).EnsureIndices([][]string{
		{"start", "addr.city", "addr.pcode"},
		{"eventid", "start"},
		{"organizerid", "start"},
	})

	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = "de"

	os.RemoveAll(indexDir + string(os.PathSeparator) + "events.bleve")
	eventIndex, err := bleve.New(indexDir+string(os.PathSeparator)+"events.bleve", mapping)
	if err != nil {
		return nil, err
	}

	var events []*Event
	err = database.Table(eventTableName).Find(bson.M{}, &events, "start")
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		err := eventIndex.Index(event.Id.Hex(), bson.M{"title": event.Title})
		if err != nil {
			return nil, err
		}
	}

	os.RemoveAll(indexDir + string(os.PathSeparator) + "dates.bleve")
	dateIndex, err := bleve.New(indexDir+string(os.PathSeparator)+"dates.bleve", mapping)
	if err != nil {
		return nil, err
	}

	var dates []*Date
	err = database.Table(dateTableName).Find(bson.M{"start": bson.M{"$gte": time.Now()}}, &dates, "start")
	if err != nil {
		return nil, err
	}

	for _, date := range dates {
		err := dateIndex.Index(date.Id.Hex(), bson.M{"title": date.Title, "location": date.Addr.Name})
		if err != nil {
			return nil, err
		}
	}

	return &EventService{database, eventTableName, dateTableName, eventIndex, dateIndex}, nil
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

func (events *EventService) Locations() ([]string, error) {

	var locations []string
	err := events.eventTable().Distinct(bson.M{}, "addr.name", &locations)
	return locations, err
}

func (events *EventService) Dates(query string) ([]string, error) {

	if isEmpty(query) {
		return make([]string, 0), nil
	}

	dates := make(map[string]string)
	tokenStream, err := events.dateIndex.Mapping().AnalyzeText("de", []byte(query))
	if err != nil {
		return nil, err
	}

	for _, token := range tokenStream {
		fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewPrefixQuery(string(token.Term)), 1000, 0, false)
		results, err := events.dateIndex.Search(fullTextSearch)
		if err != nil {
			return nil, err
		} else {
			for _, hit := range results.Hits {
				date, err := events.LoadDate(bson.ObjectIdHex(hit.ID))
				if err != nil {
					return nil, err
				}
				for field := range hit.Locations {
					if field == "title" {
						dates[date.Title] = date.Title
					} else if field == "location" {
						dates[date.Addr.Name] = date.Addr.Name
					}
				}
			}
		}
	}
	
	i := 0
	result := make([]string, len(dates))
	for date := range dates {
		result[i] = date
		i++
	}

	return result, nil
}

func (events *EventService) buildQuery(search, place string, dates [][]time.Time, targetIds, categoryIds []int, withImagesOnly bool) bson.M {

	query := make([]bson.M, 0, 6)

	if !isEmpty(search) {
		fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewMatchQuery(search), 1000, 0, false)
		results, err := events.dateIndex.Search(fullTextSearch)
		if err != nil {
			traffic.Logger().Print(err.Error())
		} else {
			ids := make([]bson.ObjectId, results.Hits.Len())
			for i, hit := range results.Hits {
				ids[i] = bson.ObjectIdHex(hit.ID)
			}
			query = append(query, bson.M{"_id": bson.M{"$in": ids}})
		}
	}

	if !isEmpty(place) {
		postcodes := Postcodes(place)
		placesQuery := make([]bson.M, len(postcodes)+1)
		for i, postcode := range postcodes {
			placesQuery[i] = bson.M{"addr.pcode": postcode}
		}
		placesQuery[len(postcodes)] = bson.M{"addr.city": place}
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
			datesQuery[i] = bson.M{"start": rangeQuery}
		}
		query = append(query, bson.M{"$or": datesQuery})
	}

	if len(targetIds) > 0 && targetIds[0] != 0 {
		targetsQuery := make([]bson.M, len(targetIds))
		for i, targetId := range targetIds {
			targetsQuery[i] = bson.M{"targets": targetId}
		}
		query = append(query, bson.M{"$or": targetsQuery})
	}

	if len(categoryIds) > 0 && categoryIds[0] != 0 {
		categoriesQuery := make([]bson.M, len(categoryIds))
		for i, categoryId := range categoryIds {
			categoriesQuery[i] = bson.M{"categories": categoryId}
		}
		query = append(query, bson.M{"$or": categoriesQuery})
	}

	if withImagesOnly {
		query = append(query, bson.M{"image": bson.M{"$exists": true, "$ne": ""}})
	}

	return bson.M{"$and": query}
}

func (events *EventService) Count(query, place string, dates [][]time.Time, targetIds, categoryIds []int) (int, error) {

	return events.dateTable().Count(events.buildQuery(query, place, dates, targetIds, categoryIds, false))
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

func (events *EventService) SearchDates(query, place string, dates [][]time.Time, targetIds, categoryIds []int, withImagesOnly bool, page, pageSize int, sort string) (DateSearchResult, error) {

	var result DateSearchResult
	err := events.dateTable().Search(events.buildQuery(query, place, dates, targetIds, categoryIds, withImagesOnly), page*pageSize, pageSize, &result, sort)
	return result, err
}

func (events *EventService) SearchDatesOfUser(userId bson.ObjectId, page, pageSize int, sort string) (*DateSearchResult, error) {

	var result DateSearchResult
	query := bson.M{"$and": []bson.M{bson.M{"organizerid": userId}, bson.M{"start": bson.M{"$gte": time.Now()}}}}
	err := events.dateTable().Search(query, page*pageSize, pageSize, &result, sort)
	return &result, err
}

func (events *EventService) FindDatesOfEvent(eventId bson.ObjectId, sort string) ([]Date, error) {

	var result []Date
	err := events.dateTable().Find(bson.M{"$and": []bson.M{bson.M{"eventid": eventId}, bson.M{"start": bson.M{"$gte": time.Now()}}}}, &result, sort)
	return result, err
}

func (events *EventService) FindEventsOfUser(userId bson.ObjectId, sort string) ([]Event, error) {

	var result []Event
	err := events.eventTable().Find(bson.M{"organizerid": userId}, &result, sort)
	return result, err
}

func (events *EventService) SearchEventsOfUser(userId bson.ObjectId, search, location string, page, pageSize int, sort string) (*EventSearchResult, error) {

	var result EventSearchResult
	var query bson.M
	if isEmpty(search) && isEmpty(location) {
		query = bson.M{"organizerid": userId}
	} else {
		descr := []bson.M{bson.M{"organizerid": userId}}
		if !isEmpty(search) {
			fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewMatchQuery(search), 1000, 0, false)
			results, err := events.eventIndex.Search(fullTextSearch)
			if err != nil {
				return nil, err
			}
			if results.Total > 0 {
				ids := make([]bson.ObjectId, results.Hits.Len())
				for i, hit := range results.Hits {
					ids[i] = bson.ObjectIdHex(hit.ID)
				}
				descr = append(descr, bson.M{"_id": bson.M{"$in": ids}})
			} else {
				return &result, err
			}
		}
		if !isEmpty(location) {
			descr = append(descr, bson.M{"addr.name": location})
		}
		query = bson.M{"$and": descr}
	}
	err := events.eventTable().Search(query, page*pageSize, pageSize, &result, sort)
	return &result, err
}

func (events *EventService) generateDates(event *Event, now time.Time) []Date {

	dates := make([]Date, 0)

	var date Date
	date.OrganizerId = event.OrganizerId
	date.EventId = event.Id
	date.Title = event.Title
	date.Start = event.Start.Local()
	date.End = event.End.Local()
	date.Image = event.Image
	date.ImageCredit = event.ImageCredit
	date.Targets = event.Targets
	date.Categories = event.Categories
	date.Descr = event.Descr
	date.Rsvp = event.Rsvp
	date.Web = event.Web
	date.Addr = event.Addr

	if !date.Start.Before(now) {
		date.Id = bson.NewObjectId()
		dates = append(dates, date)
	}

	if event.Recurrency != NoRecurrence {
		timeHorizon := event.RecurrencyEnd
		maxTimeHorizon := now.Add(366 * 24 * time.Hour)
		if timeHorizon.IsZero() || timeHorizon.After(maxTimeHorizon) {
			timeHorizon = maxTimeHorizon
		}

		var dateDuration time.Duration = 0
		if !date.End.IsZero() {
			dateDuration = date.End.Sub(date.Start)
		}

		year, week := date.Start.ISOWeek()
		var startOfFirstWeek int = -3
		for {
			firstWeek := time.Date(year, time.January, startOfFirstWeek, 6, 0, 0, 0, time.Local)
			testYear, _ := firstWeek.ISOWeek()
			if testYear == year {
				break
			} else {
				startOfFirstWeek++
			}
		}

		month := int(date.Start.Month())
		hour := date.Start.Hour()
		minute := date.Start.Minute()

		if event.Recurrency == Weekly && len(event.Weekly.Weekdays) > 0 {
			for date.Start.Before(timeHorizon) {
				week += event.Weekly.Interval
				weekDate := time.Date(year, time.January, (7*(week-1))+startOfFirstWeek, 6, 0, 0, 0, time.Local)
				for _, weekday := range event.Weekly.Weekdays {
					day := weekDate
					for day.Weekday() != weekday {
						day = day.Add(24 * time.Hour)
					}
					date.Start = time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, time.Local)
					if !date.Start.Before(now) && date.Start.Before(timeHorizon) {
						if dateDuration != 0 {
							date.End = date.Start.Add(dateDuration)
						}
						date.Id = bson.NewObjectId()
						dates = append(dates, date)
					}
				}
			}
		}

		if event.Recurrency == Monthly && event.Monthly.Interval > 0 {
			for date.Start.Before(timeHorizon) {
				month += event.Monthly.Interval
				monthDate := time.Date(year, time.Month(month), 1, 6, 0, 0, 0, time.Local)
				day := monthDate
				days := make([]time.Time, 0)
				for day.Month() == monthDate.Month() {
					if day.Weekday() == event.Monthly.Weekday {
						days = append(days, day)
					}
					day = day.Add(24 * time.Hour)
				}
				if event.Monthly.Week == LastWeek || int(event.Monthly.Week) >= len(days) {
					day = days[len(days)-1]
				} else {
					day = days[event.Monthly.Week]
				}
				date.Start = time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, time.Local)
				if !date.Start.Before(now) && date.Start.Before(timeHorizon) {
					if dateDuration != 0 {
						date.End = date.Start.Add(dateDuration)
					}
					date.Id = bson.NewObjectId()
					dates = append(dates, date)
				}
			}
		}
	}

	return dates
}

func (events *EventService) Store(event *Event, publish bool) error {

	_, err := events.eventTable().UpsertById(event.Id, event)

	if err == nil && publish {
		var dates []Date
		dates, err = events.FindDatesOfEvent(event.Id, "id")
		for _, date := range dates {
			if date.OrganizerId != event.OrganizerId {
				date.OrganizerId = event.OrganizerId
				events.StoreDate(&date)
			}
		}
		err = events.SyncDates(event)
	} else if !publish {
		err = events.DeleteDatesOfEvent(event.Id)
	}

	if err == nil {
		err = events.eventIndex.Index(event.Id.Hex(), bson.M{"title": event.Title})
	}

	return err
}

func (events *EventService) StoreDate(date *Date) error {

	_, err := events.dateTable().UpsertById(date.Id, date)
	if err == nil {
		err = events.dateIndex.Index(date.Id.Hex(), bson.M{"title": date.Title, "location": date.Addr.Name})
	}
	return err
}

func (events *EventService) SyncDates(event *Event) error {

	now := time.Now()
	genDates := events.generateDates(event, now)

	var dates []Date
	query := bson.M{"$and": []bson.M{bson.M{"eventid": event.Id}, bson.M{"start": bson.M{"$gte": now}}}}
	err := events.dateTable().Find(query, &dates, "start")
	if err != nil {
		return err
	}

	n := min(len(genDates), len(dates))
	for i := 0; i < n; i++ {
		genDates[i].Id = dates[i].Id
		err := events.StoreDate(&genDates[i])
		if err != nil {
			return err
		}
	}

	for i := n; i < len(genDates); i++ {
		err := events.StoreDate(&genDates[i])
		if err != nil {
			return err
		}
	}

	for i := n; i < len(dates); i++ {
		err := events.DeleteDate(dates[i].Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (events *EventService) UpdateRecurrences() error {

	var result []Event
	err := events.eventTable().Find(bson.M{"recurrency": Weekly}, &result, "start")
	if err != nil {
		return err
	}

	for _, event := range result {
		err = events.SyncDates(&event)
		if err != nil {
			return err
		}
	}

	err = events.eventTable().Find(bson.M{"recurrency": Monthly}, &result, "start")
	if err != nil {
		return err
	}

	for _, event := range result {
		err = events.SyncDates(&event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (events *EventService) DeleteDatesOfEvent(eventId bson.ObjectId) error {

	var dates []Date
	err := events.dateTable().Find(bson.M{"eventid": eventId}, &dates, "start")
	if err != nil {
		return err
	}

	for _, date := range dates {
		err = events.DeleteDate(date.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (events *EventService) DeleteDatesOfUser(userId bson.ObjectId) error {

	var dates []Date
	err := events.dateTable().Find(bson.M{"organizerid": userId}, &dates, "start")
	if err != nil {
		return err
	}

	for _, date := range dates {
		err = events.DeleteDate(date.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (events *EventService) DeleteDate(id bson.ObjectId) error {

	err := events.dateTable().DeleteById(id)
	if err == nil {
		err = events.dateIndex.Delete(id.Hex())
	}
	return err
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

	err := events.DeleteDatesOfEvent(id)
	if err == nil {
		err = events.eventTable().DeleteById(id)
	}
	if err == nil {
		err = events.eventIndex.Delete(id.Hex())
	}
	return err
}

func (events *EventService) Stop() error {

	err := events.eventIndex.Close()
	if err == nil {
		err = events.dateIndex.Close()
	}
	return err
}
