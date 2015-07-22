package mmr

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	Database interface {
		Table(string) Table
		CreateSession(bson.ObjectId) (bson.ObjectId, error)
		RemoveSession(bson.ObjectId) error
		RemoveOldSessions(time.Duration) error
		LoadUserBySessionId(bson.ObjectId) (*User, error)
		LoadUserByEmailAndPassword(string, string) (*User, error)
		Disconnect()
	}

	Table interface {
		DropIndices() error
		EnsureIndex(...string) error
		EnsureIndices([][]string) error
		LoadById(bson.ObjectId, Item) error
		CountById(bson.ObjectId) (int, error)
		CheckForId(bson.ObjectId) error
		Find(interface{}, interface{}, ...string) error
		Search(interface{}, int, int, SearchResult, ...string) error
		Distinct(interface{}, string, interface{}) error
		Count(interface{}) (int, error)
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

func (db *mongoDb) RemoveOldSessions(olderThan time.Duration) error {

	date := time.Now().Add(-olderThan)
	return db.Table("session").Delete(bson.M{"contact": bson.M{"$lt": date}})
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

	session.Contact = time.Now()
	_, err = db.Table("session").UpsertById(sessionId, &session)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (table *mongoTable) DropIndices() error {

	indices, err := table.collection.Indexes()
	if err != nil {
		return err
	}
	
	for _, index := range indices {
		table.collection.DropIndex(index.Key...)
	}
	
	return nil
}

func (table *mongoTable) EnsureIndex(keys ...string) error {

	return table.collection.EnsureIndexKey(keys...)
}

func (table *mongoTable) EnsureIndices(indices [][]string) error {

	for _, index := range indices {
		err := table.EnsureIndex(index...)
		if err != nil {
			return err
		}
	}
	
	return nil
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

func (table *mongoTable) Find(query interface{}, result interface{}, orderBy ...string) error {

	return table.collection.Find(query).Sort(orderBy...).All(result)
}

func (table *mongoTable) Distinct(query interface{}, field string, result interface{}) error {

	return table.collection.Find(query).Distinct(field, result)
}

func (table *mongoTable) Count(query interface{}) (int, error) {

	return table.collection.Find(query).Count()
}

func (table *mongoTable) Search(query interface{}, skip int, limit int, result SearchResult, orderBy ...string) error {

	find := table.collection.Find(query).Sort(orderBy...)
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
