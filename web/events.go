package mmr

import (
	"os"
	"sort"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/de"
	"github.com/blevesearch/bleve/mapping"
	"github.com/pilu/traffic"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type (
	EventService struct {
		database Database
		tableame string
		index    bleve.Index
	}
)

func NewEventService(database Database, tablename, indexDir string) (*EventService, error) {

	err := database.Table(tablename).DropIndices()
	if err != nil {
		return nil, errors.Wrap(err, "error dropping indices of event service database")
	}

	err = database.Table(tablename).EnsureIndices([][]string{
		{"start"},
		{"organizerid", "start"},
		{"recurrency", "recurrencyend"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating indices of event service database")
	}

	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = de.AnalyzerName

	indexFile := indexDir + string(os.PathSeparator) + "events.bleve"
	err = os.RemoveAll(indexFile)
	if err != nil {
		return nil, errors.Wrapf(err, "error removing full text index file: %s", indexFile)
	}

	index, err := bleve.New(indexFile, mapping)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating full text index file: %s", indexFile)
	}

	go func() {

		var events []*Event
		err = database.Table(tablename).Find(bson.M{}, &events, "start")
		if err != nil {
			err = errors.Wrap(err, "error loading events for indexing")
			traffic.Logger().Print(err.Error())
			return
		}

		batch := index.NewBatch()
		for _, event := range events {
			err := batch.Index(event.Id.Hex(), bson.M{"title": event.Title, "location": event.Addr.Name})
			if err != nil {
				err = errors.Wrapf(err, "error adding event %s to full text index batch", event.Id.String())
				traffic.Logger().Print(err.Error())
				return
			}
		}
		err = index.Batch(batch)
		if err != nil {
			err = errors.Wrap(err, "error writing events to full text index")
			traffic.Logger().Print(err.Error())
			return
		}
	}()

	return &EventService{database, tablename, index}, nil
}

func (events *EventService) table() Table {

	return events.database.Table(events.tableame)
}

func (events *EventService) Cities() ([]string, error) {

	var cities []string
	err := events.table().Distinct(bson.M{}, "addr.city", &cities)

	if err != nil {
		return nil, errors.Wrap(err, "error loading cities")
	}

	return cities, nil
}

func (events *EventService) Locations() ([]string, error) {

	var locations []string
	err := events.table().Distinct(bson.M{}, "addr.name", &locations)

	if err != nil {
		return nil, errors.Wrap(err, "error loading locations")
	}

	return locations, nil
}

type SearchHit struct {
	Name string `json:"name"`
	Url  string `json:"url,omitempty"`
}

func (events *EventService) SearchText(query string) ([]SearchHit, error) {

	if isEmpty(query) {
		return nil, nil
	}

	resultMap := make(map[string]SearchHit)
	tokenStream, err := events.index.Mapping().(*mapping.IndexMappingImpl).AnalyzeText("de", []byte(query))
	if err != nil {
		return nil, errors.Wrapf(err, "error tokenizing full text query: %s", query)
	}

	for _, token := range tokenStream {
		fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewPrefixQuery(string(token.Term)), 1000, 0, false)
		fullTextSearch.IncludeLocations = true
		results, err := events.index.Search(fullTextSearch)
		if err != nil {
			return nil, errors.Wrapf(err, "error searching for term %s in event full text index", string(token.Term))
		} else {
			for _, hit := range results.Hits {
				event, err := events.Load(bson.ObjectIdHex(hit.ID))
				if err != nil {
					return nil, errors.Wrapf(err, "error loading event by full text search result: %s", hit.ID)
				}
				for field := range hit.Locations {
					if field == "title" {
						resultMap[event.Title] = SearchHit{Name: event.Title, Url: event.Url()}
					} else if field == "location" {
						resultMap[event.Addr.Name] = SearchHit{Name: event.Addr.Name}
					}
				}
			}
		}
	}

	i := 0
	result := make([]SearchHit, len(resultMap))
	for _, hit := range resultMap {
		result[i] = hit
		i++
	}

	return result, nil
}

func (events *EventService) Load(id bson.ObjectId) (*Event, error) {

	var event Event

	err := events.table().LoadById(id, &event)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading event: %s", id.String())
	}

	return &event, nil
}

func (events *EventService) LoadBySource(source, sourceId string) (*Event, error) {

	var event Event
	query := bson.M{"source": source, "sourceid": sourceId}

	err := events.table().Load(query, &event, "_id")
	if err != nil {
		return nil, errors.Wrapf(err, "error loading event by source: %s and sourceId: %s", source, sourceId)
	}

	return &event, nil
}

func (events *EventService) futureQuery() bson.M {

	return bson.M{
		"$or": []bson.M{
			bson.M{"start": bson.M{"$gte": time.Now()}},
			bson.M{
				"$and": []bson.M{
					bson.M{"recurrency": bson.M{"$gt": 0}},
					bson.M{
						"$or": []bson.M{
							bson.M{"recurrencyend": time.Time{}},
							bson.M{"recurrencyend": bson.M{"$gt": time.Now()}},
						},
					},
				},
			},
		},
	}
}

func (events *EventService) buildQuery(search, place string, targetIds, categoryIds []int, withImagesOnly bool) bson.M {

	query := []bson.M{events.futureQuery()}

	if !isEmpty(search) {
		fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewMatchPhraseQuery(search), 1000, 0, false)
		results, err := events.index.Search(fullTextSearch)
		if err != nil {
			err = errors.Wrap(err, "error searching in events full text index")
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

func (events *EventService) search(query, place string, dates [][]time.Time, targetIds, categoryIds []int, withImagesOnly bool) ([]*Event, error) {

	var sourceEvents []*Event
	err := events.table().Find(events.buildQuery(query, place, targetIds, categoryIds, withImagesOnly), &sourceEvents)
	if err != nil {
		return nil, errors.Wrap(err, "error searching events")
	}

	filteredEvents := make([]*Event, 0)
	for _, event := range sourceEvents {
		for _, dateRange := range dates {
			if event.RecurresIn(dateRange[0], dateRange[1]) {
				filteredEvents = append(filteredEvents, event)
				break
			}
		}
	}

	return filteredEvents, nil
}

func (events *EventService) Search(query, place string, dates [][]time.Time, targetIds, categoryIds []int, withImagesOnly bool, page, pageSize int) (*EventSearchResult, error) {

	filteredEvents, err := events.search(query, place, dates, targetIds, categoryIds, withImagesOnly)
	if err != nil {
		return nil, err
	}

	var start time.Time
	if len(dates) > 0 {
		start = dates[0][0]
	} else {
		start = time.Now()
	}
	eventList := NewEventList(filteredEvents, start)
	sort.Sort(eventList)

	var result EventSearchResult
	result.SetCount(len(filteredEvents))
	result.SetStart(page * pageSize)
	result.Events = make([]*Event, 0)

	for i := page * pageSize; i < len(filteredEvents) && i < (page+1)*pageSize; i++ {
		result.Events = append(result.Events, filteredEvents[i])
	}

	return &result, nil
}

func (events *EventService) Count(query, place string, dates [][]time.Time, targetIds, categoryIds []int) (int, error) {

	filteredEvents, err := events.search(query, place, dates, targetIds, categoryIds, false)
	if err != nil {
		return 0, errors.Wrap(err, "error counting events")
	}

	return len(filteredEvents), nil
}

func (events *EventService) SearchFutureEventsOfUser(userId bson.ObjectId, page, pageSize int) (*EventSearchResult, error) {

	query := []bson.M{events.futureQuery()}
	query = append(query, bson.M{"organizerid": userId})

	var eventList []*Event
	err := events.table().Find(bson.M{"$and": query}, &eventList, "start")
	if err != nil {
		return nil, errors.Wrapf(err, "error searching future event ids of user: %s", userId.String())
	}

	var result EventSearchResult

	start := page * pageSize
	if start < len(eventList) {
		end := min(start+pageSize, len(eventList))
		result.Events = make([]*Event, end-start)
		for i := 0; i < end-start; i++ {
			result.Events[i] = eventList[start+i]
		}
	} else {
		result.Events = make([]*Event, 0)
	}
	result.Start = start
	result.Count = len(eventList)

	return &result, nil
}

func (events *EventService) FindEvents() ([]Event, error) {

	var allEvents []Event

	err := events.table().Find(bson.M{}, &allEvents, "start")
	if err != nil {
		return nil, errors.Wrap(err, "error loading all events")
	}

	return allEvents, nil
}

func (events *EventService) FindSimilarEvents(event *Event, count int) ([]*Event, error) {

	query := []bson.M{events.futureQuery()}
	query = append(query, bson.M{"_id": bson.M{"$ne": event.Id}})
	query = append(query, bson.M{"addr.city": event.Addr.City})

	categories := make([]bson.M, 0)
	for _, category := range event.Categories {
		categories = append(categories, bson.M{"categories": category})
	}
	if len(categories) > 0 {
		query = append(query, bson.M{"$or": categories})
	}

	targets := make([]bson.M, 0)
	for _, target := range event.Targets {
		targets = append(targets, bson.M{"targets": target})
	}
	if len(targets) > 0 {
		query = append(query, bson.M{"$or": targets})
	}

	page := 0
	pageSize := 10
	var err error
	var result EventSearchResult
	eventsList := make([]*Event, 0)

	for len(eventsList) < count {
		err = events.table().Search(bson.M{"$and": query}, page*pageSize, pageSize, &result, "start")
		if err != nil || len(result.Events) == 0 {
			break
		}
		for _, event := range result.Events {
			found := false
			for _, have := range eventsList {
				if event.Id == have.Id {
					found = true
					break
				}
			}
			if !found {
				eventsList = append(eventsList, event)
				if len(eventsList) == count {
					break
				}
			}
		}
		page++
	}

	if err != nil {
		return nil, errors.Wrapf(err, "error loading similar events for event: %s", event.Id.String())
	}

	return eventsList, nil
}

func (events *EventService) FindEventsOfUser(userId bson.ObjectId, sort string) ([]Event, error) {

	var result []Event

	err := events.table().Find(bson.M{"organizerid": userId}, &result, sort)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading events of user: %s", userId.String())
	}

	return result, nil
}

func (events *EventService) SearchEventsOfUser(userId bson.ObjectId, search string, page, pageSize int, sort string) (*EventSearchResult, error) {

	var query bson.M
	var result EventSearchResult

	if isEmpty(search) {
		query = bson.M{"organizerid": userId}
	} else {
		descr := []bson.M{bson.M{"organizerid": userId}}
		if !isEmpty(search) {
			fullTextSearch := bleve.NewSearchRequestOptions(bleve.NewMatchPhraseQuery(search), 1000, 0, false)
			results, err := events.index.Search(fullTextSearch)
			if err != nil {
				return nil, errors.Wrapf(err, "error retrieving events of user %s from full text index", userId.String())
			}
			ids := make([]bson.ObjectId, results.Hits.Len())
			for i, hit := range results.Hits {
				ids[i] = bson.ObjectIdHex(hit.ID)
			}
			descr = append(descr, bson.M{"_id": bson.M{"$in": ids}})
		}
		query = bson.M{"$and": descr}
	}

	err := events.table().Search(query, page*pageSize, pageSize, &result, sort)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading events of user: %s", userId.String())
	}

	return &result, nil
}

func (events *EventService) Store(event *Event) error {

	_, err := events.table().UpsertById(event.Id, event)
	if err != nil {
		return errors.Wrapf(err, "error upserting event: %s", event.Id.String())
	}

	go func() {
		events.index.Delete(event.Id.Hex())
		events.index.Index(event.Id.Hex(), bson.M{"title": event.Title, "location": event.Addr.Name})
	}()

	return nil
}

func (events *EventService) DeleteOfUser(userId bson.ObjectId) error {

	eventsOfUser, err := events.FindEventsOfUser(userId, "start")
	if err != nil {
		return errors.Wrapf(err, "error deleting events of user: %s", userId.String())
	}

	for _, event := range eventsOfUser {
		err = events.Delete(event.Id)
		if err != nil {
			return errors.Wrapf(err, "error deleting events of user: %s", userId.String())
		}
	}

	return nil
}

func (events *EventService) Delete(id bson.ObjectId) error {

	err := events.table().DeleteById(id)
	if err != nil {
		return errors.Wrapf(err, "error deleting event: %s", id.String())
	}

	go func() {
		events.index.Delete(id.Hex())
	}()

	return nil
}

func (events *EventService) Stop() error {

	events.index.Close()

	return nil
}
