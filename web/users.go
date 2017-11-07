package mmr

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type (
	UserService struct {
		database  Database
		tablename string
	}
)

func NewUserService(database Database, tablename string) (*UserService, error) {

	err := database.Table(tablename).DropIndices()
	if err != nil {
		return nil, errors.Wrap(err, "error deleting indices of user service database")
	}

	err = database.Table(tablename).EnsureIndices([][]string{
		{"approved", "categories", "addr.city", "addr.pcode"},
		{"name"},
		{"email"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating indices of user service database")
	}

	return &UserService{database, tablename}, nil
}

func (users *UserService) table() Table {

	return users.database.Table(users.tablename)
}

func (users *UserService) buildQuery(place string, categoryIds []int) bson.M {

	query := make([]bson.M, 0, 3)

	query = append(query, bson.M{"approved": true})

	if len(place) > 0 {
		postcodes := Postcodes(place)
		placesQuery := make([]bson.M, len(postcodes)+1)
		for i, postcode := range postcodes {
			placesQuery[i] = bson.M{"addr.pcode": postcode}
		}
		placesQuery[len(postcodes)] = bson.M{"addr.city": place}
		query = append(query, bson.M{"$or": placesQuery})
	}

	if len(categoryIds) > 0 && categoryIds[0] != 0 {
		categoriesQuery := make([]bson.M, len(categoryIds))
		for i, categoryId := range categoryIds {
			categoriesQuery[i] = bson.M{"categories": categoryId}
		}
		query = append(query, bson.M{"$or": categoriesQuery})
	}

	return bson.M{"$and": query}
}

func (users *UserService) Count(place string, categoryIds []int) (int, error) {

	cnt, err := users.table().Count(users.buildQuery(place, categoryIds))
	if err != nil {
		return 0, errors.Wrap(err, "error counting users")
	}

	return cnt, nil
}

func (users *UserService) Validate(user *User) error {

	var result []User
	err := users.table().Find(bson.M{"email": user.Email}, &result)
	if err != nil {
		return errors.Wrapf(err, "error finding user with email %s", user.Email)
	}

	for i := 0; i < len(result); i++ {
		if result[i].Email == user.Email && result[i].Id != user.Id {
			return errors.New("Email address is already in use.")
		}
	}

	return nil
}

func (users *UserService) Load(id bson.ObjectId) (*User, error) {

	var user User
	err := users.table().LoadById(id, &user)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading user %s", id.String())
	}

	return &user, nil
}

func (users *UserService) LoadByEmail(email string) (*User, error) {

	var user User
	err := users.table().Load(bson.M{"email": email}, &user, "_id")
	if err != nil {
		return nil, errors.Wrapf(err, "error loading user by email: %s", email)
	}

	return &user, nil
}

func (users *UserService) FindUsers() ([]User, error) {

	var result []User
	err := users.table().Find(bson.M{}, &result, "name")
	if err != nil {
		return nil, errors.Wrap(err, "error loading all users")
	}

	return result, nil
}

func (users *UserService) FindApproved() ([]User, error) {

	var result []User
	err := users.table().Find(bson.M{"approved": true}, &result, "name")
	if err != nil {
		return nil, errors.Wrap(err, "error loading approved users")
	}

	return result, nil
}

func (users *UserService) FindForEvents(events []*Event) (map[bson.ObjectId]*User, error) {

	organizers := make(map[bson.ObjectId]*User)
	for _, event := range events {
		if _, found := organizers[event.OrganizerId]; !found {
			user, err := users.Load(event.OrganizerId)
			if err != nil {
				return nil, err
			}
			organizers[event.OrganizerId] = user
		}
	}

	return organizers, nil
}

func (users *UserService) Search(place string, categoryIds []int, page, pageSize int, sort string) (*OrganizerSearchResult, error) {

	var result OrganizerSearchResult
	err := users.table().Search(users.buildQuery(place, categoryIds), page*pageSize, pageSize, &result, "name")
	if err != nil {
		return nil, errors.Wrap(err, "error searching users")
	}

	return &result, nil
}

func (users *UserService) Store(user *User) error {

	_, err := users.table().UpsertById(user.Id, user)
	if err != nil {
		return errors.Wrapf(err, "error storing user: %s", user.Id.String())
	}

	return nil
}

func (users *UserService) Delete(id bson.ObjectId) error {

	err := users.table().DeleteById(id)
	if err != nil {
		return errors.Wrapf(err, "error deleting user: %s", id.String())
	}

	return nil
}
