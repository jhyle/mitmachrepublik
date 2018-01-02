package mmr

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"
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
		Source        string        `json:",omitempty"`
		SourceId      string        `json:",omitempty"`
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
		Facebook      bool   `bson:"-"`
		FacebookId    string `json:"-"`
		Addr          Address
	}

	EventList struct {
		start  time.Time
		events []*Event
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

	OrganizerSearchResult struct {
		Count      int
		Start      int
		Organizers []*User
	}

	Topic struct {
		Name        string
		Place       string
		TargetIds   []int
		CategoryIds []int
		DateIds     []int
		FrontPage   bool
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

	ADMIN_EMAIL = "mitmachrepublik@gmail.com"
)

const (
	FromNow = iota
	Today
	Tomorrow
	ThisWeek
	NextWeekend
	NextWeek
	TwoWeeks
	AfterTomorrow
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
		FromNow:       "alle ab jetzt",
		Today:         "heute",
		Tomorrow:      "morgen",
		AfterTomorrow: "übermorgen",
		ThisWeek:      "diese Woche",
		NextWeekend:   "am Wochenende",
		NextWeek:      "nächste Woche",
		TwoWeeks:      "14 Tage",
	}

	DateOrder []int = []int{TwoWeeks, Today, Tomorrow, AfterTomorrow, ThisWeek, NextWeekend, NextWeek, FromNow}

	whiteSpace = regexp.MustCompile(`\s+`)

	Topics map[string]Topic = map[string]Topic{
		"babies-und-kleinkinder": Topic{
			Name:        "Babies & Kleinkinder",
			Place:       "Berlin",
			TargetIds:   []int{19, 20},
			CategoryIds: []int{0},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"sport-und-gesundheit": Topic{
			Name:        "Sport & Gesundheit",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{8, 24},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"natur-und-garten": Topic{
			Name:        "Natur & Garten",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{9, 18},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"eltern-und-familien": Topic{
			Name:        "Eltern & Familien",
			Place:       "Berlin",
			TargetIds:   []int{1, 2, 19, 20, 21},
			CategoryIds: []int{0},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"bildung-und-kultur": Topic{
			Name:        "Bildung & Kultur",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{10, 11},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"umwelt-und-tierschutz": Topic{
			Name:        "Umwelt- & Tierschutz",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{13, 14},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"demonstrationen-und-politik": Topic{
			Name:        "Demonstrationen & Politik",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{15, 23},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"soziales-und-ehrenamt": Topic{
			Name:        "Soziales & Ehrenamt",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{16, 17},
			DateIds:     []int{FromNow},
			FrontPage:   true,
		},
		"heute-in-berlin": Topic{
			Name:        "heute",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{0},
			DateIds:     []int{Today},
			FrontPage:   false,
		},
		"morgen-in-berlin": Topic{
			Name:        "morgen",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{0},
			DateIds:     []int{Tomorrow},
			FrontPage:   false,
		},
		"uebermorgen-in-berlin": Topic{
			Name:        "übermorgen",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{0},
			DateIds:     []int{AfterTomorrow},
			FrontPage:   false,
		},
		"am-wochenende-in-berlin": Topic{
			Name:        "das nächste Wochenende",
			Place:       "Berlin",
			TargetIds:   []int{0},
			CategoryIds: []int{0},
			DateIds:     []int{NextWeekend},
			FrontPage:   false,
		},
	}
)

func (user *User) GetId() bson.ObjectId {
	return user.Id
}

func (user *User) SetId(id bson.ObjectId) {
	user.Id = id
}

func (user *User) Url() string {

	return "/veranstalter/" + user.Id.Hex() + "/" + sanitizePath(user.Name) + "/0"
}

func (user *User) HtmlDescription() template.HTML {

	return noescape(sanitizeHtml(user.Descr))
}

func (user *User) PlainDescription() string {

	return whiteSpace.ReplaceAllString(sanitize.HTML(user.Descr), " ")
}

func (user *User) IsAdmin() bool {

	return user.Email == ADMIN_EMAIL
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

func (event *Event) Url() string {

	targetNames := make([]string, len(event.Targets))
	for i, id := range event.Targets {
		targetNames[i] = TargetIdMap[id]
	}
	categoryNames := make([]string, len(event.Categories))
	for i, id := range event.Categories {
		categoryNames[i] = CategoryIdMap[id]
	}

	return "/veranstaltung/" + sanitizePath(citypartName(event.Addr)) + "/" + sanitizePath(strings.Join(targetNames, "-")) + "/" + sanitizePath(strings.Join(categoryNames, "-")) + "/" + event.Id.Hex() + "/" + sanitizePath(event.Title) + ".html"
}

func (event *Event) HtmlDescription() template.HTML {

	return noescape(sanitizeHtml(event.Descr))
}

func (event *Event) PlainDescription() string {

	return whiteSpace.ReplaceAllString(sanitize.HTML(event.Descr), " ")
}

var weekdayShort map[time.Weekday]string = map[time.Weekday]string{
	time.Monday:    "Mo",
	time.Tuesday:   "Di",
	time.Wednesday: "Mi",
	time.Thursday:  "Do",
	time.Friday:    "Fr",
	time.Saturday:  "Sa",
	time.Sunday:    "So",
}

func (event *Event) Recurrence() string {

	text := ""

	if Weekly == event.Recurrency {

		days := make([]string, 0)
		for _, weekday := range event.Weekly.Weekdays {
			days = append(days, weekdayShort[weekday])
		}

		if event.Weekly.Interval == 1 && len(days) < 7 {
			text += "jeden "
		} else if event.Weekly.Interval > 1 {
			text += fmt.Sprintf("alle %d Wochen ", event.Weekly.Interval)
			if len(days) == 1 {
				text += "am "
			}
		}
		if len(days) == 1 {
			text += fmt.Sprintf("%s", weekday[int(event.Weekly.Weekdays[0])])
		} else if len(days) == 2 && event.Weekly.Interval == 1 {
			text += fmt.Sprintf("%s und %s", weekday[int(event.Weekly.Weekdays[0])], weekday[int(event.Weekly.Weekdays[1])])
		} else if len(days) < 7 {
			text += fmt.Sprintf("%s", strConcat(days))
		} else if event.Weekly.Interval == 1 {
			text += "Montag bis Sonntag"
		} else {
			text += "Mo - So"
		}

	} else if Monthly == event.Recurrency {

		day := ""
		if event.Monthly.Interval > 1 {
			text += fmt.Sprintf("jeden %d. Monat am ", event.Monthly.Interval)
			day = weekdayShort[event.Monthly.Weekday]
		} else {
			text += "am "
			day = weekday[int(event.Monthly.Weekday)]
		}

		if event.Monthly.Week != LastWeek {
			text += fmt.Sprintf("%d. %s", (event.Monthly.Week + 1), day)
		} else {
			text += fmt.Sprintf("letzten %s", day)
		}

		if event.Monthly.Interval == 1 {
			text += " im Monat"
		}
	}

	return text
}

func (event *Event) Dates(from, until time.Time) []time.Time {

	dates := make([]time.Time, 0)
	if !from.After(event.Start) && !until.Before(event.Start) {
		dates = append(dates, event.Start)
	}

	date := event.Start
	if Weekly == event.Recurrency {

		eventDays := make(map[time.Weekday]bool)
		for _, day := range event.Weekly.Weekdays {
			eventDays[day] = true
		}

		for date.Before(from) {
			date = date.AddDate(0, 0, 7*event.Weekly.Interval)
		}
		date = date.AddDate(0, 0, -7*event.Weekly.Interval)

		weekday := date.Weekday()
		for (event.RecurrencyEnd.IsZero() || date.Before(event.RecurrencyEnd)) && date.Before(until) {

			if weekday == time.Sunday {
				date = date.AddDate(0, 0, 7*(event.Weekly.Interval-1))
			}

			date = date.AddDate(0, 0, 1)
			weekday = date.Weekday()

			if eventDays[weekday] && !date.Before(from) && !date.After(until) {
				dates = append(dates, date)
			}
		}

	} else if Monthly == event.Recurrency {

		for (event.RecurrencyEnd.IsZero() || date.Before(event.RecurrencyEnd)) && date.Before(until) {

			var day int
			month := date.Month()
			year := date.Year()
			if event.Monthly.Week != LastWeek {
				day = 7*int(event.Monthly.Week) + 1
			} else {
				day = lengthOfMonth(month, year) - 6
			}

			date = time.Date(year, month, day, event.Start.Hour(), event.Start.Minute(), 0, 0, time.Local)
			for i := 0; i < 7; i++ {
				if date.Weekday() == event.Monthly.Weekday {
					if !date.Before(from) && date.After(event.Start) && (event.RecurrencyEnd.IsZero() || date.Before(event.RecurrencyEnd)) && !date.After(until) {
						dates = append(dates, date)
					}
				}
				date = date.AddDate(0, 0, 1)
			}

			date = time.Date(year, month, 1, 6, 0, 0, 0, time.Local).AddDate(0, event.Monthly.Interval, 0)
		}
	}

	return dates
}

func (event *Event) RecurresIn(from, until time.Time) bool {

	if until.IsZero() || until.Equal(from) {
		if !from.After(event.Start) {
			return true
		} else if Weekly == event.Recurrency || Monthly == event.Recurrency {
			return event.RecurrencyEnd.IsZero() || !from.After(event.RecurrencyEnd)
		}
	} else {
		return len(event.Dates(from, until)) > 0
	}

	return false
}

func (event *Event) NextDate(from time.Time) time.Time {

	if !event.Start.After(from) && (Weekly == event.Recurrency || Monthly == event.Recurrency) {
		var until time.Time
		if Weekly == event.Recurrency {
			until = from.AddDate(0, 0, 7*(event.Weekly.Interval+1))
		} else if Monthly == event.Recurrency {
			until = from.AddDate(0, event.Monthly.Interval+1, 0)
		}

		dates := event.Dates(from, until)
		if len(dates) > 0 {
			return dates[0]
		}
	}

	return event.Start
}

func (list EventList) Len() int {
	return len(list.events)
}

func (list EventList) Swap(i, j int) {
	list.events[i], list.events[j] = list.events[j], list.events[i]
}

func (list EventList) Less(i, j int) bool {
	return list.events[i].NextDate(list.start).Before(list.events[j].NextDate(list.start))
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
