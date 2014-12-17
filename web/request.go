package mmr

import (
	"encoding/json"
	"errors"
	"github.com/pilu/traffic"
	"labix.org/v2/mgo/bson"
)

type (
	Request struct {
		*traffic.Request
	}
)

func (r *Request) ReadJson(v interface{}) error {

	return json.NewDecoder(r.Body).Decode(v)
}

func (r *Request) ReadEmailAndPwd() (*emailAndPwd, error) {

	var form emailAndPwd
	err := r.ReadJson(&form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *Request) ReadUser() (*User, error) {

	var data User
	user := &data
	err := r.ReadJson(user)
	return user, err
}

func (r *Request) ReadSessionId() (bson.ObjectId, error) {

	cookie, err := r.Cookie("SESSIONID")
	if err != nil || !bson.IsObjectIdHex(cookie.Value) {
		return "", errors.New("Failed to read session id.")
	}

	return bson.ObjectIdHex(cookie.Value), nil
}
