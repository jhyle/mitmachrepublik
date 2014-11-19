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
		Id         bson.ObjectId `bson:"_id" json:",omitempty"`
		Title      string
		Image      string
		Descr      string
		Web        string
		Categories int64
		Start      time.Time
		End        time.Time
		Addr       Address
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
		"Kinder & Familie": 1 << 0,
		"Jugendliche":      1 << 1,
		"Studenten":        1 << 2,
		"Berufstätige":     1 << 3,
		"Eltern":           1 << 4,
		"Senioren":         1 << 5,
		"Leute treffen":    1 << 6,
		"Sport":            1 << 7,
		"Gärtnern":         1 << 8,
		"Kultur":           1 << 9,
		"Bildung":          1 << 10,
		"Religion":         1 << 11,
		"Umwelt":           1 << 12,
		"Tierschutz":       1 << 13,
		"Demonstrationen":  1 << 14,
		"Soziales":         1 << 15,
		"Ehrenamt":         1 << 16,
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
