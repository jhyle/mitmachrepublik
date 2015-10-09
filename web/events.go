package mmr

import (
	"github.com/blevesearch/bleve"
	"github.com/pilu/traffic"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
	"sort"
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

	go func() {

		var events []*Event
		err = database.Table(eventTableName).Find(bson.M{}, &events, "start")
		if err != nil {
			traffic.Logger().Print(err.Error())
			return
		}

		batch := eventIndex.NewBatch()
		for _, event := range events {
			err := batch.Index(event.Id.Hex(), bson.M{"title": event.Title, "location": event.Addr.Name})
			if err != nil {
				traffic.Logger().Print(err.Error())
				return
			}
		}
		err = eventIndex.Batch(batch)
		if err != nil {
			traffic.Logger().Print(err.Error())
			return
		}
	}()

	os.RemoveAll(indexDir + string(os.PathSeparator) + "dates.bleve")
	dateIndex, err := bleve.New(indexDir+string(os.PathSeparator)+"dates.bleve", mapping)
	if err != nil {
		return nil, err
	}

	go func() {

		var dates []*Date
		err := database.Table(dateTableName).Find(bson.M{"start": bson.M{"$gte": time.Now()}}, &dates, "start")
		if err != nil {
			traffic.Logger().Print(err.Error())
			return
		}

		batch := dateIndex.NewBatch()
		for _, date := range dates {
			err := batch.Index(date.Id.Hex(), bson.M{"title": date.Title, "location": date.Addr.Name})
			if err != nil {
				traffic.Logger().Print(err.Error())
				return
			}
		}

		err = dateIndex.Batch(batch)
		if err != nil {
			traffic.Logger().Print(err.Error())
			return
		}
	}()

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
		fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewMatchPhraseQuery(search), 1000, 0, false)
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

func (events *EventService) LoadBySource(source, sourceId string) (*Event, error) {

	var event Event
	query := bson.M{"source": source, "sourceid": sourceId}
	err := events.eventTable().Load(query, &event, "_id")
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

func (events *EventService) SearchFutureEventsOfUser(userId bson.ObjectId, page, pageSize int) (*EventSearchResult, error) {

	var result EventSearchResult

	var eventIds []bson.ObjectId
	query := bson.M{"$and": []bson.M{bson.M{"organizerid": userId}, bson.M{"start": bson.M{"$gte": time.Now()}}}}
	err := events.dateTable().Distinct(query, "eventid", &eventIds)
	if err != nil {
		return &result, err
	}
	
	var eventsOfOrganizer Events
	err = events.eventTable().Find(bson.M{"_id": bson.M{"$in": eventIds}}, &eventsOfOrganizer, "_id")
	if err != nil {
		return &result, err
	}
	
	for i := 0; i < len(eventsOfOrganizer); i++ {
		var date Date
		query := bson.M{"$and": []bson.M{bson.M{"eventid": eventsOfOrganizer[i].Id}, bson.M{"start": bson.M{"$gte": time.Now()}}}}
		err = events.dateTable().Load(query, &date, "start")
		if err != nil {
			return &result, err
		}
		eventsOfOrganizer[i].Start = date.Start
		eventsOfOrganizer[i].End = date.End
	}
	sort.Sort(eventsOfOrganizer)

	start := page*pageSize
	if start < len(eventsOfOrganizer) {
		end := min(start + pageSize, len(eventsOfOrganizer))
		result.Events = make([]*Event, end - start)
		for i := 0; i < end - start; i++ {
			result.Events[i] = &eventsOfOrganizer[start + i]
		}
	} else {
		result.Events = make([]*Event, 0)
	}
	result.Start = start
	result.Count = len(eventsOfOrganizer)

	return &result, err
}

func (events *EventService) FindEvents() ([]Event, error) {

	var allEvents []Event
	err := events.eventTable().Find(bson.M{}, &allEvents, "start")
	return allEvents, err
}

func (events *EventService) FindNextDates() ([]Date, error) {

	allEvents, err := events.FindEvents()
	if err != nil {
		return nil, err
	}

	dates := make([]Date, 0)
	for _, event := range allEvents {
		datesOfEvent, err := events.FindDatesOfEvent(event.Id, "start")
		if err != nil {
			return nil, err
		}
		if len(datesOfEvent) > 0 {
			dates = append(dates, datesOfEvent[0])
		}
	}
	
	return dates, nil
}

func (events *EventService) FindSimilarDates(date *Date, count int) ([]Date, error) {

	query := []bson.M{bson.M{"eventid": bson.M{"$ne": date.EventId}}, bson.M{"addr.city": date.Addr.City}, bson.M{"start": bson.M{"$gt": date.Start}}}
	
	categories := make([]bson.M, 0)
	for _, category := range date.Categories {
		categories = append(categories, bson.M{"categories": category})
	}
	if len(categories) > 0 {
		query = append(query, bson.M{"$or": categories})
	}

	targets := make([]bson.M, 0)
	for _, target := range date.Targets {
		targets = append(targets, bson.M{"targets": target})
	}
	if len(targets) > 0 {
		query = append(query, bson.M{"$or": targets})
	}

	dates := make([]Date, 0)
	query1 := bson.M{"$and": query}

	page := 0
	pageSize := 10
	var err error
	var result DateSearchResult

	for len(dates) < count {
		err = events.dateTable().Search(query1, page*pageSize, pageSize, &result, "start")
		if err != nil || len(result.Dates) == 0 {
			break
		}
		for _, date := range result.Dates {
			found := false
			for _, have := range dates {
				if date.EventId == have.EventId {
					found = true
					break
				}
			}
			if !found {
				dates = append(dates, *date)
				if len(dates) == count {
					break
				}
			}
		}
		page++
	}
	return dates, err
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

func (events *EventService) SearchEventsOfUser(userId bson.ObjectId, search string, page, pageSize int, sort string) (*EventSearchResult, error) {

	var result EventSearchResult
	var query bson.M
	if isEmpty(search) {
		query = bson.M{"organizerid": userId}
	} else {
		descr := []bson.M{bson.M{"organizerid": userId}}
		if !isEmpty(search) {
			fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewMatchPhraseQuery(search), 1000, 0, false)
			results, err := events.eventIndex.Search(fullTextSearch)
			if err != nil {
				return nil, err
			}
			ids := make([]bson.ObjectId, results.Hits.Len())
			for i, hit := range results.Hits {
				ids[i] = bson.ObjectIdHex(hit.ID)
			}
			descr = append(descr, bson.M{"_id": bson.M{"$in": ids}})
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
		_, err = events.SyncDates(event)
	} else if !publish {
		err = events.DeleteDatesOfEvent(event.Id)
	}

	if err == nil {
		go func() {
			events.eventIndex.Index(event.Id.Hex(), bson.M{"title": event.Title})
		}()
	}

	return err
}

func (events *EventService) StoreDate(date *Date) error {

	_, err := events.dateTable().UpsertById(date.Id, date)
	if err == nil {
		go func() {
			events.dateIndex.Index(date.Id.Hex(), bson.M{"title": date.Title, "location": date.Addr.Name})
		}()
	}
	return err
}

func (events *EventService) SyncDates(event *Event) ([]bson.ObjectId, error) {

	now := time.Now()
	genDates := events.generateDates(event, now)

	var dates []Date
	query := bson.M{"$and": []bson.M{bson.M{"eventid": event.Id}, bson.M{"start": bson.M{"$gte": now}}}}
	err := events.dateTable().Find(query, &dates, "start")
	if err != nil {
		return nil, err
	}

	n := min(len(genDates), len(dates))
	for i := 0; i < n; i++ {
		genDates[i].Id = dates[i].Id
		err := events.StoreDate(&genDates[i])
		if err != nil {
			return nil, err
		}
	}

	newDates := make([]bson.ObjectId, 0)
	for i := n; i < len(genDates); i++ {
		err := events.StoreDate(&genDates[i])
		if err != nil {
			return nil, err
		}
		newDates = append(newDates, genDates[i].Id)
	}

	for i := n; i < len(dates); i++ {
		err := events.DeleteDate(dates[i].Id)
		if err != nil {
			return nil, err
		}
	}

	return newDates, nil
}

func (events *EventService) UpdateRecurrences() ([]bson.ObjectId, error) {

	var result []Event
	err := events.eventTable().Find(bson.M{"recurrency": Weekly}, &result, "start")
	if err != nil {
		return nil, err
	}

	dates := make([]bson.ObjectId, 0)
	for _, event := range result {
		newDates, err := events.SyncDates(&event)
		if err != nil {
			return nil, err
		}
		dates = append(dates, newDates...)
	}

	err = events.eventTable().Find(bson.M{"recurrency": Monthly}, &result, "start")
	if err != nil {
		return nil, err
	}

	for _, event := range result {
		newDates, err := events.SyncDates(&event)
		if err != nil {
			return nil, err
		}
		dates = append(dates, newDates...)
	}

	return dates, nil
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
		go func() {
			events.dateIndex.Delete(id.Hex())
		}()
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
		go func() {
			events.eventIndex.Delete(id.Hex())
		}()
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
