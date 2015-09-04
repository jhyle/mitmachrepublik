package mmr

import (
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"regexp"
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
		Id          bson.ObjectId `bson:"_id" json:",omitempty"`
		Name        string
		Email       string
		Pwd         string
		Image       string
		ImageCredit string
		Categories  []int
		Descr       string
		Web         string
		Addr        Address
		AGBs        bool
		Approved    bool
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
		Id            bson.ObjectId `bson:"_id" json:",omitempty"`
		OrganizerId   bson.ObjectId `json:",omitempty"`
		Title         string
		Image         string
		ImageCredit   string
		Descr         string
		Web           string
		Targets       []int
		Categories    []int
		Start         time.Time
		End           time.Time `json:",omitempty"`
		Recurrency    Recurrence
		RecurrencyEnd time.Time         `json:",omitempty"`
		Weekly        WeeklyRecurrence  `json:",omitempty"`
		Monthly       MonthlyRecurrence `json:",omitempty"`
		Rsvp          bool
		Addr          Address
	}

	Date struct {
		Id          bson.ObjectId `bson:"_id"`
		EventId     bson.ObjectId
		OrganizerId bson.ObjectId
		Title       string
		Image       string
		ImageCredit string
		Descr       string
		Web         string
		Targets     []int
		Categories  []int
		Start       time.Time
		End         time.Time
		Rsvp        bool
		Addr        Address
	}

	Alert struct {
		Id         bson.ObjectId `bson:"_id"`
		Name       string
		Email      string
		Query      string
		Place      string
		Targets    []int
		Categories []int
		Dates      []int
		Radius     int
		Weekdays   []time.Weekday
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

const (
	FromNow = iota
	Today
	Tomorrow
	ThisWeek
	NextWeekend
	NextWeek
	TwoWeeks
)

var (
	TargetMap map[string]int = map[string]int{
		"alle Zielgruppen": 0,
		"Familien":         1,
		"Jugendliche":      2,
		"Studenten":        3,
		"Erwachsene":       4,
		"Eltern":           5,
		"Senioren":         6,
		"Kleinkinder":      19,
		"Babies":           20,
		"Kinder":           21,
	}

	TargetIdMap map[int]string = map[int]string{
		0:  "alle Zielgruppen",
		1:  "Familien",
		2:  "Jugendliche",
		3:  "Studenten",
		4:  "Erwachsene",
		5:  "Eltern",
		6:  "Senioren",
		19: "Kleinkinder",
		20: "Babies",
		21: "Kinder",
	}

	TargetOrder []string = []string{
		"Babies",
		"Kleinkinder",
		"Kinder",
		"Jugendliche",
		"Studenten",
		"Erwachsene",
		"Eltern",
		"Familien",
		"Senioren",
	}

	CategoryMap map[string]int = map[string]int{
		"allen Kategorien":          0,
		"Leute treffen":             7,
		"Sport":                     8,
		"Gärtnern":                  9,
		"Kultur":                    10,
		"Bildung":                   11,
		"Religion":                  12,
		"Umweltschutz":              13,
		"Tierschutz":                14,
		"Demonstration":             15,
		"Soziales":                  16,
		"Ehrenamt":                  17,
		"Natur":                     18,
		"Basteln & Spielen":         22,
		"Politik":                   23,
		"Gesundheit & Wohlbefinden": 24,
		"Handwerk & Kreatives":      25,
	}

	CategoryIdMap map[int]string = map[int]string{
		0:  "allen Kategorien",
		7:  "Leute treffen",
		8:  "Sport",
		9:  "Gärtnern",
		10: "Kultur",
		11: "Bildung",
		12: "Religion",
		13: "Umweltschutz",
		14: "Tierschutz",
		15: "Demonstration",
		16: "Soziales",
		17: "Ehrenamt",
		18: "Natur",
		22: "Basteln & Spielen",
		23: "Politik",
		24: "Gesundheit & Wohlbefinden",
		25: "Handwerk & Kreatives",
	}

	CategoryIconMap map[int]string = map[int]string{
		0:  "calendar",
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
		18: "leaf",
		22: "child",
		23: "globe",
		24: "heart",
		25: "paint-brush",
	}

	CategoryOrder []string = []string{
		"Leute treffen",
		"Basteln & Spielen",
		"Handwerk & Kreatives",
		"Gesundheit & Wohlbefinden",
		"Sport",
		"Gärtnern",
		"Natur",
		"Kultur",
		"Bildung",
		"Religion",
		"Umweltschutz",
		"Tierschutz",
		"Politik",
		"Demonstration",
		"Soziales",
		"Ehrenamt",
	}

	DateIdMap map[int]string = map[int]string{
		FromNow:     "alle ab jetzt",
		Today:       "heute",
		Tomorrow:    "morgen",
		ThisWeek:    "diese Woche",
		NextWeekend: "am Wochenende",
		NextWeek:    "nächste Woche",
		TwoWeeks:    "14 Tage",
	}

	DateOrder []int = []int{TwoWeeks, Today, Tomorrow, ThisWeek, NextWeekend, NextWeek, FromNow}

	whiteSpace = regexp.MustCompile(`\s+`)
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

	return noescape(sanitizeHtml(user.Descr))
}

func (user *User) PlainDescription() string {

	return whiteSpace.ReplaceAllString(sanitize.HTML(user.Descr), " ")
}

func (user *User) IsAdmin() bool {

	return user.Email == "mitmachrepublik@gmail.com"
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

	return noescape(sanitizeHtml(event.Descr))
}

func (event *Event) PlainDescription() string {

	return whiteSpace.ReplaceAllString(sanitize.HTML(event.Descr), " ")
}

func (date *Date) GetId() bson.ObjectId {
	return date.Id
}

func (date *Date) SetId(id bson.ObjectId) {
	date.Id = id
}

func (date *Date) Url() string {

	targetNames := make([]string, len(date.Targets))
	for i, id := range date.Targets {
		targetNames[i] = TargetIdMap[id]
	}
	categoryNames := make([]string, len(date.Categories))
	for i, id := range date.Categories {
		categoryNames[i] = CategoryIdMap[id]
	}

	return "/veranstaltung/" + citypartName(date.Addr) + "/" + strings.Join(targetNames, ",") + "/" + strings.Join(categoryNames, ",") + "/" + dateFormat(date.Start) + "/" + date.Id.Hex() + "/" + date.Title
}

func (date *Date) HtmlDescription() template.HTML {

	return noescape(sanitizeHtml(date.Descr))
}

func (date *Date) PlainDescription() string {

	return whiteSpace.ReplaceAllString(sanitize.HTML(date.Descr), " ")
}

func (session *Session) GetId() bson.ObjectId {
	return session.Id
}

func (session *Session) SetId(id bson.ObjectId) {
	session.Id = id
}

func (alert *Alert) GetId() bson.ObjectId {
	return alert.Id
}

func (alert *Alert) SetId(id bson.ObjectId) {
	alert.Id = id
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
