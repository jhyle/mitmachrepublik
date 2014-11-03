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

	SearchResult interface {
		SetCount(int)
		SetStart(int)
		GetData() interface{}
		GetSize() int
		GetItem(int) Item
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

func (session *Session) GetId() bson.ObjectId {
	return session.Id
}

func (session *Session) SetId(id bson.ObjectId) {
	session.Id = id
}
