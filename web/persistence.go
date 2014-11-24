package mmr

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	Database interface {
		Table(string) Table
		CreateSession(bson.ObjectId) (bson.ObjectId, error)
		RemoveSession(bson.ObjectId) error
		LoadUserBySessionId(bson.ObjectId) (*User, error)
		LoadUserByEmailAndPassword(string, string) (*User, error)
		Disconnect()
	}

	Table interface {
		LoadById(bson.ObjectId, Item) error
		CountById(bson.ObjectId) (int, error)
		CheckForId(bson.ObjectId) error
		Find(interface{}, interface{}) error
		Search(interface{}, string, int, int, SearchResult) error
		UpsertById(bson.ObjectId, Item) (bson.ObjectId, error)
		DeleteById(bson.ObjectId) error
		Delete(interface{}) error
	}

	mongoDb struct {
		session *mgo.Session
		name    string
	}

	mongoTable struct {
		collection *mgo.Collection
	}
)

func NewMongoDb(mongoUrl string, dbName string) (Database, error) {

	session, err := mgo.Dial(mongoUrl)
	return &mongoDb{session, dbName}, err
}

func (db *mongoDb) Table(name string) Table {

	return &mongoTable{db.session.DB(db.name).C(name)}
}

func (db *mongoDb) CreateSession(userId bson.ObjectId) (bson.ObjectId, error) {

	session := Session{bson.NewObjectId(), userId, time.Now()}
	return db.Table("session").UpsertById(session.GetId(), &session)
}

func (db *mongoDb) RemoveSession(sessionId bson.ObjectId) error {

	return db.Table("session").DeleteById(sessionId)
}

func (db *mongoDb) Disconnect() {

	db.session.Close()
}

func (db *mongoDb) LoadUserByEmailAndPassword(email string, password string) (*User, error) {

	find := db.session.DB(db.name).C("user").Find(bson.M{"email": email, "pwd": password})
	n, err := find.Count()

	if err != nil {
		return nil, err
	}

	if n != 1 {
		return nil, errors.New("E-Mail + Password not found!")
	}

	var user User
	err = find.One(&user)
	return &user, err
}

func (db *mongoDb) LoadUserBySessionId(sessionId bson.ObjectId) (*User, error) {

	var session Session
	err := db.Table("session").LoadById(sessionId, &session)
	if err != nil {
		return nil, err
	}

	var user User
	err = db.Table("user").LoadById(session.UserId, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (table *mongoTable) LoadById(id bson.ObjectId, item Item) error {

	return table.collection.FindId(id).One(item)
}

func (table *mongoTable) CountById(id bson.ObjectId) (int, error) {

	return table.collection.FindId(id).Count()
}

func (table *mongoTable) CheckForId(id bson.ObjectId) error {

	count, err := table.CountById(id)
	if err != nil {
		return err
	}

	if count == 1 {
		return nil
	} else {
		return errors.New(table.collection.Name + " " + id.Hex() + " not found.")
	}
}

func (table *mongoTable) Find(query interface{}, result interface{}) error {

	return table.collection.Find(query).All(result)
}

func (table *mongoTable) Search(query interface{}, sort string, skip int, limit int, result SearchResult) error {

	find := table.collection.Find(query).Sort(sort)
	count, err := find.Count()
	if err != nil {
		return err
	}
	result.SetCount(count)
	result.SetStart(skip)

	data := result.GetData()
	return find.Skip(skip).Limit(limit).All(data)
}

func (table *mongoTable) UpsertById(id bson.ObjectId, item Item) (bson.ObjectId, error) {

	info, err := table.collection.UpsertId(id, item)
	if err != nil {
		return id, err
	}

	if info.UpsertedId != nil {
		return info.UpsertedId.(bson.ObjectId), nil
	} else {
		return id, nil
	}
}

func (table *mongoTable) Delete(query interface{}) error {

	_, err := table.collection.RemoveAll(query)
	return err
}


func (table *mongoTable) DeleteById(id bson.ObjectId) error {

	return table.collection.RemoveId(id)
}
