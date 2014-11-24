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
