package mmr

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	Item interface {
		GetId() bson.ObjectId
		SetId(id bson.ObjectId)
	}

	Address struct {
		Name     string
		Street   string
		Pcode    string
		City     string
	}

	User struct {
		Id       bson.ObjectId `bson:"_id" json:",omitempty"`
		Email    string
		Pwd      string
		Image    string
		Descr    string
		Web      string
		Addr     Address
		Approved bool
	}

	Session struct {
		Id      bson.ObjectId `bson:"_id"`
		UserId  bson.ObjectId
		Contact time.Time
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
		Events []Event
	}

	OrganizerSearchResult struct {
		Count      int
		Start      int
		Organizers []User
	}
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

func (user *User) IsAdmin() bool {

	return user.Email == "admin@mitmachrepublik.de"
}

func (addr Address) IsEmpty() bool {

	return isEmpty(addr.City) && isEmpty(addr.Name) && isEmpty(addr.Pcode) && isEmpty(addr.Street)
}

func (event *Event) GetId() bson.ObjectId {
	return event.Id
}

func (event *Event) SetId(id bson.ObjectId) {
	event.Id = id
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
	return &result.Events[i]
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
	return &result.Organizers[i]
}
