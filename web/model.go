package mmr

import (
	"github.com/kennygrant/sanitize"
	"html/template"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

type (
	Item interface {
		GetId() bson.ObjectId
		SetId(id bson.ObjectId)
	}

	Address struct {
		Name   string
		Street string
		Pcode  string
		City   string
	}

	User struct {
		Id         bson.ObjectId `bson:"_id" json:",omitempty"`
		Name       string
		Email      string
		Pwd        string
		Image      string
		Categories []int
		Descr      string
		Web        string
		Addr       Address
		AGBs       bool
		Approved   bool
	}

	Session struct {
		Id      bson.ObjectId `bson:"_id"`
		UserId  bson.ObjectId
		Contact time.Time
	}

	Recurrence int

	WeeklyRecurrence struct {
		Interval int
		Weekdays []time.Weekday
	}

	WeekOfMonth int

	MonthlyRecurrence struct {
		Week     WeekOfMonth
		Weekday  time.Weekday
		Interval int
	}

	Event struct {
		Id          bson.ObjectId `bson:"_id" json:",omitempty"`
		OrganizerId bson.ObjectId `json:",omitempty"`
		Title       string
		Image       string
		Descr       string
		Web         string
		Categories  []int
		Start       time.Time
		End         time.Time `json:",omitempty"`
		Recurrency  Recurrence
		Weekly      WeeklyRecurrence
		Monthly     MonthlyRecurrence
		Rsvp        bool
		Addr        Address
	}

	Date struct {
		Id          bson.ObjectId
		EventId     bson.ObjectId
		OrganizerId bson.ObjectId
		Title       string
		Image       string
		Descr       string
		Web         string
		Categories  []int
		Start       time.Time
		End         time.Time
		Rsvp        bool
		Addr        Address
	}

	SearchResult interface {
		SetCount(int)
		SetStart(int)
		GetData() interface{}
		GetSize() int
		GetItem(int) Item
	}

	EventSearchResult struct {
		Count  int
		Start  int
		Events []*Event
	}

	DateSearchResult struct {
		Count int
		Start int
		Dates []*Date
	}

	OrganizerSearchResult struct {
		Count      int
		Start      int
		Organizers []*User
	}
)

const (	
	NoRecurrence Recurrence = iota
	Weekly
	Monthly 

	FirstWeek WeekOfMonth = iota
	SecondWeek
	ThirdWeek
	FourthWeek
	LastWeek
)

var (
	CategoryMap map[string]int = map[string]int{
		"alle Kategorien":  0,
		"Kinder & Familie": 1,
		"Jugendliche":      2,
		"Studenten":        3,
		"Berufstätige":     4,
		"Eltern":           5,
		"Senioren":         6,
		"Leute treffen":    7,
		"Sport":            8,
		"Gärtnern":         9,
		"Kultur":           10,
		"Bildung":          11,
		"Religion":         12,
		"Umwelt":           13,
		"Tierschutz":       14,
		"Demonstrationen":  15,
		"Soziales":         16,
		"Ehrenamt":         17,
	}

	CategoryIdMap map[int]string = map[int]string{
		0:  "alle Kategorien",
		1:  "Kinder & Familie",
		2:  "Jugendliche",
		3:  "Studenten",
		4:  "Berufstätige",
		5:  "Eltern",
		6:  "Senioren",
		7:  "Leute treffen",
		8:  "Sport",
		9:  "Gärtnern",
		10: "Kultur",
		11: "Bildung",
		12: "Religion",
		13: "Umweltschutz",
		14: "Tierschutz",
		15: "Demonstrationen",
		16: "Soziales",
		17: "Ehrenamt",
	}

	CategoryIconMap map[int]string = map[int]string{
		0:  "calendar",
		1:  "child",
		2:  "child",
		3:  "group",
		4:  "group",
		5:  "group",
		6:  "group",
		7:  "glass",
		8:  "futbol-o",
		9:  "leaf",
		10: "institution",
		11: "graduation-cap",
		12: "group",
		13: "globe",
		14: "globe",
		15: "globe",
		16: "heart",
		17: "heart",
	}

	CategoryOrder []string = []string{
		"Kinder & Familie",
		"Jugendliche",
		"Studenten",
		"Berufstätige",
		"Eltern",
		"Senioren",
		"Leute treffen",
		"Sport",
		"Gärtnern",
		"Kultur",
		"Bildung",
		"Religion",
		"Umwelt",
		"Tierschutz",
		"Demonstrationen",
		"Soziales",
		"Ehrenamt",
	}

	DateIdMap map[int]string = map[int]string{
		0: "aktuell",
		1: "heute",
		2: "morgen",
		3: "wochenende",
	}
)

func (user *User) GetId() bson.ObjectId {
	return user.Id
}

func (user *User) SetId(id bson.ObjectId) {
	user.Id = id
}

func (user *User) Url() string {

	return "/veranstalter/" + user.Id.Hex() + "/" + user.Name + "/0"
}

func (user *User) HtmlDescription() template.HTML {

	html, _ := sanitize.HTMLAllowing(user.Descr)
	return template.HTML(html)
}

func (user *User) PlainDescription() string {

	return sanitize.HTML(user.Descr)
}

func (user *User) IsAdmin() bool {

	return user.Email == "admin@mitmachrepublik.de"
}

func (addr Address) IsEmpty() bool {

	return isEmpty(addr.City) && isEmpty(addr.Pcode) && isEmpty(addr.Street)
}

func (event *Event) GetId() bson.ObjectId {
	return event.Id
}

func (event *Event) SetId(id bson.ObjectId) {
	event.Id = id
}

func (event *Event) HtmlDescription() template.HTML {

	html, _ := sanitize.HTMLAllowing(event.Descr)
	return template.HTML(html)
}

func (event *Event) PlainDescription() string {

	return sanitize.HTML(event.Descr)
}

func (date *Date) GetId() bson.ObjectId {
	return date.Id
}

func (date *Date) SetId(id bson.ObjectId) {
	date.Id = id
}

func (date *Date) Url() string {

	categoryNames := make([]string, len(date.Categories))
	for i, id := range date.Categories {
		categoryNames[i] = CategoryIdMap[id]
	}

	return "/veranstaltung/" + strings.Join(categoryNames, ",") + "/" + dateFormat(date.Start) + "/" + date.Id.Hex() + "/" + date.Title
}

func (date *Date) HtmlDescription() template.HTML {

	html, _ := sanitize.HTMLAllowing(date.Descr)
	return template.HTML(html)
}

func (date *Date) PlainDescription() string {

	return sanitize.HTML(date.Descr)
}

func (session *Session) GetId() bson.ObjectId {
	return session.Id
}

func (session *Session) SetId(id bson.ObjectId) {
	session.Id = id
}

func (result *EventSearchResult) SetCount(count int) {
	result.Count = count
}

func (result *EventSearchResult) SetStart(start int) {
	result.Start = start
}

func (result *EventSearchResult) GetData() interface{} {
	return &result.Events
}

func (result *EventSearchResult) GetSize() int {
	return len(result.Events)
}

func (result *EventSearchResult) GetItem(i int) Item {
	return result.Events[i]
}

func (result *DateSearchResult) SetCount(count int) {
	result.Count = count
}

func (result *DateSearchResult) SetStart(start int) {
	result.Start = start
}

func (result *DateSearchResult) GetData() interface{} {
	return &result.Dates
}

func (result *DateSearchResult) GetSize() int {
	return len(result.Dates)
}

func (result *DateSearchResult) GetItem(i int) Item {
	return result.Dates[i]
}

func (result *OrganizerSearchResult) SetCount(count int) {
	result.Count = count
}

func (result *OrganizerSearchResult) SetStart(start int) {
	result.Start = start
}

func (result *OrganizerSearchResult) GetData() interface{} {
	return &result.Organizers
}

func (result *OrganizerSearchResult) GetSize() int {
	return len(result.Organizers)
}

func (result *OrganizerSearchResult) GetItem(i int) Item {
	return result.Organizers[i]
}
