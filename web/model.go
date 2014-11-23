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
		Name   string
		Street string
		Pcode  string
		City   string
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
		End         time.Time
		Addr        Address
	}

	SearchResult interface {
		SetCount(int)
		SetStart(int)
		GetData() interface{}
		GetSize() int
		GetItem(int) Item
	}
)

var (
	CategoryMap map[string]int = map[string]int{
		"Kinder & Familie": 1,
		"Jugendliche":      2,
		"Studenten":        3,
		"Berufstätige":     4,
		"Eltern":           5,
		"Senioren":         6,
		"Leute treffen":    7,
		"Sport":            8,
		"Gärtnern":         8,
		"Kultur":           9,
		"Bildung":          10,
		"Religion":         11,
		"Umwelt":           12,
		"Tierschutz":       13,
		"Demonstrationen":  14,
		"Soziales":         15,
		"Ehrenamt":         16,
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
