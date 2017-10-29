package mmr

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		Load(interface{}, Item, ...string) error
		LoadById(bson.ObjectId, Item) error
		CountById(bson.ObjectId) (int, error)
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

	ErrorNotFound struct {
		collection string
		id         string
	}
)

func (e *ErrorNotFound) Error() string {

	return fmt.Sprintf("could not find %s in %s", e.id, e.collection)
}

func NewMongoDb(mongoUrl string, dbName string) (Database, error) {

	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil, errors.Wrap(err, "error dialing MongoDB")
	}

	return &mongoDb{session, dbName}, nil
}

func (db *mongoDb) Table(name string) Table {

	db.session.Refresh()
	return &mongoTable{db.session.DB(db.name).C(name)}
}

func (db *mongoDb) CreateSession(userId bson.ObjectId) (bson.ObjectId, error) {

	session := Session{bson.NewObjectId(), userId, time.Now()}

	id, err := db.Table("session").UpsertById(session.GetId(), &session)
	if err != nil {
		return "", errors.Wrap(err, "error creating session")
	}

	return id, nil
}

func (db *mongoDb) RemoveOldSessions(olderThan time.Duration) error {

	date := time.Now().Add(-olderThan)

	err := db.Table("session").Delete(bson.M{"contact": bson.M{"$lt": date}})
	if err != nil {
		return errors.Wrap(err, "error deleting old sessions")
	}

	return nil
}

func (db *mongoDb) RemoveSession(sessionId bson.ObjectId) error {

	err := db.Table("session").DeleteById(sessionId)
	if err != nil {
		return errors.Wrapf(err, "error deleting session: %s", sessionId.String())
	}

	return nil
}

func (db *mongoDb) Disconnect() {

	db.session.Close()
}

func (db *mongoDb) LoadUserByEmailAndPassword(email string, password string) (*User, error) {

	var user User
	err := db.session.DB(db.name).C("user").Find(bson.M{"email": email, "pwd": password}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &ErrorNotFound{"user", email}
		} else {
			return nil, errors.Wrap(err, "error checking user by email and password")
		}
	}

	return &user, nil
}

func (db *mongoDb) LoadUserBySessionId(sessionId bson.ObjectId) (*User, error) {

	var session Session
	err := db.Table("session").LoadById(sessionId, &session)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &ErrorNotFound{"session", sessionId.String()}
		} else {
			return nil, errors.Wrapf(err, "error loading session: %s", sessionId.String())
		}
	}

	var user User
	err = db.Table("user").LoadById(session.UserId, &user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &ErrorNotFound{"user", session.UserId.String()}
		} else {
			return nil, errors.Wrapf(err, "error loading user: %s", session.UserId.String())
		}
	}

	session.Contact = time.Now()
	_, err = db.Table("session").UpsertById(sessionId, &session)
	if err != nil {
		return nil, errors.Wrapf(err, "error storing session: %s", sessionId.String())
	}

	return &user, nil
}

func (table *mongoTable) DropIndices() error {

	indices, err := table.collection.Indexes()
	if err != nil {
		return errors.Wrapf(err, "error loading indices of %s", table.collection.Name)
	}

	for _, index := range indices {
		if index.Name != "_id_" {
			err = table.collection.DropIndexName(index.Name)
			if err != nil {
				return errors.Wrapf(err, "error dropping index %s of %s", index.Name, table.collection.Name)
			}
		}
	}

	return nil
}

func (table *mongoTable) EnsureIndex(keys ...string) error {

	err := table.collection.EnsureIndexKey(keys...)
	if err != nil {
		return errors.Wrapf(err, "error creating index %+v in %s", keys, table.collection.Name)
	}

	return nil
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

	err := table.collection.FindId(id).One(item)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &ErrorNotFound{table.collection.Name, id.String()}
		} else {
			return errors.Wrapf(err, "error loading %s from %s", id.String(), table.collection.Name)
		}
	}

	return nil
}

func (table *mongoTable) CountById(id bson.ObjectId) (int, error) {

	cnt, err := table.collection.FindId(id).Count()
	if err != nil {
		return 0, errors.Wrapf(err, "error counting id %s in %s", id.String(), table.collection.Name)
	}

	return cnt, nil
}

func (table *mongoTable) Load(query interface{}, item Item, orderBy ...string) error {

	err := table.collection.Find(query).Sort(orderBy...).One(item)
	if err != nil {
		return errors.Wrapf(err, "error loading item from %s with query: %+v", table.collection.Name, query)
	}

	return nil
}

func (table *mongoTable) Find(query interface{}, result interface{}, orderBy ...string) error {

	err := table.collection.Find(query).Sort(orderBy...).All(result)
	if err != nil {
		return errors.Wrapf(err, "error loading items from %s with query: %+v", table.collection.Name, query)
	}

	return nil
}

func (table *mongoTable) Distinct(query interface{}, field string, result interface{}) error {

	err := table.collection.Find(query).Distinct(field, result)
	if err != nil {
		return errors.Wrapf(err, "error loading distinct items from %s with query: %+v", table.collection.Name, query)
	}

	return nil
}

func (table *mongoTable) Count(query interface{}) (int, error) {

	cnt, err := table.collection.Find(query).Count()
	if err != nil {
		return 0, errors.Wrapf(err, "error counting items from %s with query: %+v", table.collection.Name, query)
	}

	return cnt, nil
}

func (table *mongoTable) Search(query interface{}, skip int, limit int, result SearchResult, orderBy ...string) error {

	find := table.collection.Find(query).Sort(orderBy...)
	count, err := find.Count()
	if err != nil {
		return errors.Wrapf(err, "error counting items from %s with query: %+v", table.collection.Name, query)
	}
	result.SetCount(count)
	result.SetStart(skip)

	data := result.GetData()
	err = find.Skip(skip).Limit(limit).All(data)
	if err != nil {
		return errors.Wrapf(err, "error searching items from %s with query: %+v", table.collection.Name, query)
	}

	return nil
}

func (table *mongoTable) UpsertById(id bson.ObjectId, item Item) (bson.ObjectId, error) {

	info, err := table.collection.UpsertId(id, item)
	if err != nil {
		return "", errors.Wrapf(err, "error upserting item %s in %s", id.String(), table.collection.Name)
	}

	if info.UpsertedId != nil {
		return info.UpsertedId.(bson.ObjectId), nil
	} else {
		return id, nil
	}
}

func (table *mongoTable) Delete(query interface{}) error {

	_, err := table.collection.RemoveAll(query)
	if err != nil {
		return errors.Wrapf(err, "error deleting items with query %+v in %s", query, table.collection.Name)
	}

	return nil
}

func (table *mongoTable) DeleteById(id bson.ObjectId) error {

	err := table.collection.RemoveId(id)
	if err != nil {
		return errors.Wrapf(err, "error deleting item %s in %s", id.String(), table.collection.Name)
	}

	return nil
}
