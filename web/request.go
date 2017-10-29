package mmr

import (
	"encoding/json"

	"github.com/pilu/traffic"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type (
	Request struct {
		*traffic.Request
	}
)

func (r *Request) ReadJson(v interface{}) error {

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return errors.Wrap(err, "error reading request body JSON")
	}

	return nil
}

func (r *Request) ReadEmailAndPwd() (*emailAndPwd, error) {

	var form emailAndPwd
	err := r.ReadJson(&form)
	if err != nil {
		return nil, errors.Wrap(err, "error reading email and password from JSON request")
	}

	return &form, nil
}

func (r *Request) ReadSendMail() (*sendMail, error) {

	var form sendMail
	err := r.ReadJson(&form)
	if err != nil {
		return nil, errors.Wrap(err, "error reading send mail command from JSON request")
	}

	return &form, nil
}

func (r *Request) ReadUser() (*User, error) {

	var data User
	user := &data
	err := r.ReadJson(user)
	if err != nil {
		return nil, errors.Wrap(err, "error reading user from JSON request")
	}

	return user, nil
}

func (r *Request) ReadSessionId() (bson.ObjectId, error) {

	cookie, err := r.Cookie("SESSIONID")
	if err != nil || !bson.IsObjectIdHex(cookie.Value) {
		return "", errors.New("Failed to read session id.")
	}

	return bson.ObjectIdHex(cookie.Value), nil
}
